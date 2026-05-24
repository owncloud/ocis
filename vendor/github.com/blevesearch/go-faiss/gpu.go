//  Copyright (c) 2026 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build gpu

package faiss

/*
#include <stddef.h>
#include <faiss/c_api/gpu/StandardGpuResources_c.h>
#include <faiss/c_api/gpu/GpuAutoTune_c.h>
#include <faiss/c_api/gpu/GpuClonerOptions_c.h>
#include <faiss/c_api/gpu/DeviceUtils_c.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	errAccessingGPUDevices = errors.New("error accessing GPU devices")
	errNilIndex            = errors.New("index is nil")
	errNoGPUDevices        = errors.New("no GPU devices available")
)

// memorySpace controls where GPU index data is allocated.
type memorySpace int

const (
	// memorySpaceDevice uses standard GPU memory (cudaMalloc).
	memorySpaceDevice memorySpace = 1
	// memorySpaceUnified uses CUDA managed memory (cudaMallocManaged),
	// allowing the index to exceed GPU memory on Pascal+ (CC 6.0+) GPUs.
	memorySpaceUnified memorySpace = 2
)

const (
	// the minimum amount of free memory that must be available on a GPU to be considered for index cloning.
	minGPUFreeMemory = 512 * 1024 * 1024 // 512 MiB
	// the default memory space to use for GPU indices
	defaultGPUMemoryMode = memorySpaceUnified
)

var (
	gpuCount     int
	loadBalancer *gpuLoadBalancer
)

func init() {
	var err error
	gpuCount, err = numGPUs()
	if err != nil || gpuCount <= 0 {
		gpuCount = 0
	}

	// With exactly one GPU there is nothing to balance; getBestGPUDevice()
	// returns device 0 directly when loadBalancer is nil.
	// TODO: verify if 500 milliseconds is a good interval
	if gpuCount > 1 {
		loadBalancer = newGPULoadBalancer(500 * time.Millisecond)
		go loadBalancer.monitor()
	}
}

// numGPUs returns the number of available GPU devices.
func numGPUs() (int, error) {
	var rv C.int
	c := C.faiss_get_num_gpus(&rv)
	if c != 0 {
		return 0, fmt.Errorf("error getting number of GPUs, err: %v", getLastError())
	}
	return int(rv), nil
}

// gpuLoadBalancer monitors GPU free memory on a fixed interval, keeps a
// memory-sorted list of devices, and hands them out in round-robin order.
// At each interval the list is re-sorted and the round-robin counter resets
// to 0, so the next cycle always starts from the GPU with the most free memory.
type gpuLoadBalancer struct {
	mu            sync.RWMutex
	sortedDevices []int
	idx           atomic.Uint32
	interval      time.Duration
	// scratch buffers reused across refresh calls; only accessed by the monitor goroutine
	freeMemory  []uint64
	scratchDevs []int
}

func newGPULoadBalancer(interval time.Duration) *gpuLoadBalancer {
	lb := &gpuLoadBalancer{
		interval:      interval,
		freeMemory:    make([]uint64, gpuCount),
		scratchDevs:   make([]int, 0, gpuCount),
		sortedDevices: make([]int, 0, gpuCount),
	}
	return lb
}

func (lb *gpuLoadBalancer) monitor() {
	ticker := time.NewTicker(lb.interval)
	defer ticker.Stop()

	// Perform an initial sort before any requests come in.
	lb.refresh()

	for range ticker.C {
		lb.refresh()
	}
}

// refresh queries every GPU for free memory, sorts the device list in descending
// order of free memory, and resets the round-robin counter to 0.
// If all queries fail the sorted list becomes empty, causing nextDevice to error.
func (lb *gpuLoadBalancer) refresh() {
	// Zero freeMemory before querying; failed queries leave their slot as 0,
	// which naturally excludes those devices from selection.
	clear(lb.freeMemory)
	lb.scratchDevs = lb.scratchDevs[:0]

	var wg sync.WaitGroup
	wg.Add(gpuCount)
	for i := 0; i < gpuCount; i++ {
		go func(device int) {
			defer wg.Done()
			var freeBytes C.size_t
			if C.faiss_gpu_free_memory(C.int(device), &freeBytes) == 0 {
				lb.freeMemory[device] = uint64(freeBytes)
			}
		}(i)
	}
	wg.Wait()

	// Only include devices that reported non-zero free memory, and have at least minGPUFreeMemory free.
	for i, mem := range lb.freeMemory {
		if mem > minGPUFreeMemory {
			lb.scratchDevs = append(lb.scratchDevs, i)
		}
	}

	// Shuffle first, then sort descending by free memory to make the
	// sort as "unstable" as possible
	// This is useful to add fairness between GPUs with the same memory
	rand.Shuffle(len(lb.scratchDevs), func(i, j int) {
		lb.scratchDevs[i], lb.scratchDevs[j] = lb.scratchDevs[j], lb.scratchDevs[i]
	})
	// Sort in a descending order by free memory so index 0 is the most appealing GPU.
	sort.Slice(lb.scratchDevs, func(i, j int) bool {
		return lb.freeMemory[lb.scratchDevs[i]] > lb.freeMemory[lb.scratchDevs[j]]
	})

	lb.mu.Lock()
	old := lb.sortedDevices
	lb.sortedDevices = lb.scratchDevs
	lb.scratchDevs = old[:0]
	lb.idx.Store(0)
	lb.mu.Unlock()
}

// nextDevice returns the next GPU device in round-robin order.
// Returns an error if no devices are currently available.
func (lb *gpuLoadBalancer) nextDevice() (int, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	devices := lb.sortedDevices
	n := len(devices)
	if n == 0 {
		return 0, errAccessingGPUDevices
	}

	// atomically allocates the GPU. Minus 1 for zero based index
	idx := lb.idx.Add(1) - 1
	return devices[int(idx%uint32(n))], nil
}

func getBestGPUDevice() (int, error) {
	if gpuCount == 0 {
		return 0, errNoGPUDevices
	}
	// With exactly one GPU there is nothing to balance; always use device 0.
	if loadBalancer == nil {
		return 0, nil
	}
	return loadBalancer.nextDevice()
}

// only expose API used by zapx
type GPUIndexImpl struct {
	idx         *faissIndex
	gpuResource *C.FaissStandardGpuResources
}

func (g *GPUIndexImpl) cPtr() *C.FaissIndex {
	return g.idx.idx
}

func (g *GPUIndexImpl) Train(x []float32) error {
	return g.idx.Train(x)
}

func (g *GPUIndexImpl) Add(x []float32) error {
	return g.idx.Add(x)
}

func (g *GPUIndexImpl) Search(x []float32, k int64) ([]float32, []int64, error) {
	return g.idx.Search(x, k)
}

func (g *GPUIndexImpl) Close() {
	if g.idx != nil {
		g.idx.Close()
		g.idx = nil
	}
	if g.gpuResource != nil {
		C.faiss_StandardGpuResources_free(g.gpuResource)
		g.gpuResource = nil
	}
}

// CloneToGPU transfers a CPU index to the best available GPU based on free memory.
func CloneToGPU(cpuIndex *IndexImpl) (*GPUIndexImpl, error) {
	if cpuIndex == nil {
		return nil, errNilIndex
	}

	// Use the load balancer to select the best GPU device
	device, err := getBestGPUDevice()
	if err != nil {
		return nil, err
	}

	var gpuResource *C.FaissStandardGpuResources
	if code := C.faiss_StandardGpuResources_new(&gpuResource); code != 0 {
		return nil, fmt.Errorf("failed to initialize GPU resources: error code %d, err: %v", code, getLastError())
	}

	// Disable the pre-allocated temp memory pool so that all GPU memory is
	// available for index data; unified memory mode handles intermediate
	// allocations via cudaMalloc/cudaFree on demand.
	if code := C.faiss_StandardGpuResources_noTempMemory(gpuResource); code != 0 {
		C.faiss_StandardGpuResources_free(gpuResource)
		return nil, fmt.Errorf("failed to disable GPU temp memory: error code %d, err: %v", code, getLastError())
	}

	var clonerOpts *C.FaissGpuClonerOptions
	if code := C.faiss_GpuClonerOptions_new(&clonerOpts); code != 0 {
		C.faiss_StandardGpuResources_free(gpuResource)
		return nil, fmt.Errorf("failed to create cloner options: error code %d, err: %v", code, getLastError())
	}
	defer C.faiss_GpuClonerOptions_free(clonerOpts)

	C.faiss_GpuClonerOptions_set_memorySpace(clonerOpts, C.int(defaultGPUMemoryMode))

	var gpuIdx *C.FaissGpuIndex
	code := C.faiss_index_cpu_to_gpu_with_options(
		gpuResource,
		C.int(device),
		cpuIndex.cPtr(),
		clonerOpts,
		&gpuIdx,
	)
	if code != 0 {
		C.faiss_StandardGpuResources_free(gpuResource)
		return nil, fmt.Errorf("failed to transfer index to GPU device %d: error code %d, err: %v", device, code, getLastError())
	}

	idx := &faissIndex{
		idx: (*C.FaissIndex)(unsafe.Pointer(gpuIdx)),
	}

	return &GPUIndexImpl{
		idx:         idx,
		gpuResource: gpuResource,
	}, nil
}

func CloneToCPU(gpuIndex *GPUIndexImpl) (*IndexImpl, error) {
	if gpuIndex == nil {
		return nil, errNilIndex
	}

	var cpuIdx *C.FaissIndex
	code := C.faiss_index_gpu_to_cpu(
		gpuIndex.cPtr(),
		&cpuIdx,
	)
	if code != 0 {
		return nil, fmt.Errorf("failed to transfer index to CPU: %v", getLastError())
	}
	return &IndexImpl{&faissIndex{idx: cpuIdx}}, nil
}
