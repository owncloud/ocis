package localcache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/loov/hrtime"
)

func TestInitializeMultiple(t *testing.T) {
	var tests = []map[string]interface{}{
		{},
		{
			"capacity": 1024,
		},
		{
			"randomKey":     "RandomValue",
			"notConsidered": true,
			"anotherType":   map[string]string{},
		},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			err := lCache.Initialize(test)
			if err != nil {
				t.Errorf("Initialization failed for params %+v, got error %v", test, err)
			}

			capacitySet := lCache.MaxCap()
			if testCap, ok := test["capacity"]; !ok {
				if capacitySet != 512 {
					t.Errorf("Wrong initialization for params %+v, expected capacity 512, got %d", test, capacitySet)
				}
			} else {
				if capacitySet != testCap {
					t.Errorf("Wrong initialization for params %+v, expected capacity %d, got %d", test, testCap, capacitySet)
				}
			}
		})
	}
}

func storeMultiple(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{}) // 512 cap is enough
			err := lCache.Store(test.key, test.value, test.ttl)
			if err != nil {
				t.Errorf("Storing failed for params %+v, got error %v", test, err)
			}
		})
	}
}

func storeMultipleLimitedCap(t *testing.T) {
	var toBeStoredList = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	lCache := NewLocalCache()
	lCache.Initialize(map[string]interface{}{
		"capacity": 2,
	})

	for _, toBeStored := range toBeStoredList {
		err := lCache.Store(toBeStored.key, toBeStored.value, toBeStored.ttl)
		if err != nil {
			t.Fatalf("Storing failed, got error %v", err)
		}
	}

	var expected = []struct {
		key   string
		value string
	}{
		{"5566", "välòç"},
		{"üñòKey", "98685"},
	}

	expectedLength := len(expected)
	if expectedLength != lCache.MapLen() || expectedLength != lCache.ListLen() {
		t.Fatalf("Storing failed, expected stored length %d, got map length %d and list length %d", expectedLength, lCache.MapLen(), lCache.ListLen())
	}
	currentPos := 0
	traverseCallback := func(key string, value string) {
		expectedResult := expected[currentPos]
		if expectedResult.key != key {
			t.Errorf("Storing failed, wrong key at pos %d, expected %s, got %s", currentPos, expectedResult.key, key)
		}
		if expectedResult.value != value {
			t.Errorf("Storing failed, wrong value at pos %d, expected %s, got %s", currentPos, expectedResult.value, value)
		}
		currentPos++
	}
	lCache.TraverseList(traverseCallback)
}

func TestStore(t *testing.T) {
	t.Run("Multiple", storeMultiple)
	t.Run("MultipleLimitedCap", storeMultipleLimitedCap)
}

func retrieveMiss(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{}) // 512 cap is enough
			lCache.Store("goodKey", "goodValue", 500)
			value, exists, err := lCache.Retrieve(test.key)
			if err != nil {
				t.Errorf("Retrieving failed for params %+v, got error %v", test, err)
			}
			if exists {
				t.Errorf("Retrieving failed for params %+v, value %v exists", test, value)
			}
			if value != "" {
				t.Errorf("Retrieving failed for params %+v, unexpected value %v was returned", test, value)
			}
		})
	}
}

func retrieveMissObsolete(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{
				"capacity": 2,
			})
			lCache.Store(test.key, test.value, test.ttl)
			lCache.Store("goodKey", "goodValue", 500)
			lCache.Store("goodKey2", "goodValue2", 500)
			// test key should have been removed to free space
			value, exists, err := lCache.Retrieve(test.key)
			if err != nil {
				t.Errorf("Retrieving failed for params %+v, got error %v", test, err)
			}
			if exists {
				t.Errorf("Retrieving failed for params %+v, value %v exists", test, value)
			}
			if value != "" {
				t.Errorf("Retrieving failed for params %+v, unexpected value %v was returned", test, value)
			}
		})
	}
}

func retrieveHit(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{}) // 512 cap is enough
			lCache.Store(test.key, test.value, test.ttl)
			value, exists, err := lCache.Retrieve(test.key)
			if err != nil {
				t.Errorf("Retrieving failed for params %+v, got error %v", test, err)
			}
			if !exists {
				t.Errorf("Retrieving failed for params %+v, value %v doesn't exist", test, value)
			}
		})
	}
}

func retrieveHitExpired(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", -222},
		{"üñòKey", "välòç", -1},
		{"5566", "välòç", -1555},
		{"üñòKey", "98685", -555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{}) // 512 cap is enough
			lCache.Store(test.key, test.value, test.ttl)
			value, exists, err := lCache.Retrieve(test.key)
			if err != nil {
				t.Errorf("Retrieving failed for params %+v, got error %v", test, err)
			}
			if exists {
				t.Errorf("Retrieving failed for params %+v, value %v exists", test, value)
			}
			if value != "" {
				t.Errorf("Retrieving failed for params %+v, unexpected value %v was returned", test, value)
			}
		})
	}
}

func TestRetrieve(t *testing.T) {
	t.Run("Miss", retrieveMiss)
	t.Run("MissObsolete", retrieveMissObsolete)
	t.Run("Hit", retrieveHit)
	t.Run("HitExpired", retrieveHitExpired)
}

func removeHit(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{}) // 512 cap is enough
			lCache.Store(test.key, test.value, test.ttl)
			if err := lCache.Remove(test.key); err != nil {
				t.Errorf("Removing failed for params %+v, got error %v", test, err)
			}
			value, exists, _ := lCache.Retrieve(test.key)
			if exists {
				t.Errorf("Removing failed for params %+v, value %v exists", test, value)
			}
			if value != "" {
				t.Errorf("Removing failed for params %+v, unexpected value %v was returned", test, value)
			}
		})
	}
}

func removeMiss(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 0},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{}) // 512 cap is enough
			if err := lCache.Remove(test.key); err != nil {
				t.Errorf("Removing failed for params %+v, got error %v", test, err)
			}
			value, exists, _ := lCache.Retrieve(test.key)
			if exists {
				t.Errorf("Removing failed for params %+v, value %v exists", test, value)
			}
			if value != "" {
				t.Errorf("Removing failed for params %+v, unexpected value %v was returned", test, value)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	t.Run("Hit", removeHit)
	t.Run("Miss", removeMiss)
}

func concurrentStore(t *testing.T) {
	capacities := []int{2, 35, 100, 10000}
	for _, capacity := range capacities {
		t.Run("Cap"+strconv.Itoa(capacity), func(t *testing.T) {
			concurrentStoreWithCap(t, capacity)
		})
	}
}

func concurrentStoreWithCap(t *testing.T, lCacheCap int) {
	concurrentThreads := []int{10, 50}

	for _, threads := range concurrentThreads {
		t.Run("Threads"+strconv.Itoa(threads), func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{
				"capacity": lCacheCap,
			})

			wg := sync.WaitGroup{}
			wg.Add(threads)

			for i := 0; i < threads; i++ {
				go func(i int) {
					v := strconv.Itoa(i)
					lCache.Store(v, v, 500)
					wg.Done()
				}(i)
			}
			wg.Wait()

			filledSlots := lCacheCap
			if threads < lCacheCap {
				filledSlots = threads
			}

			if lCache.MapLen() != filledSlots || lCache.ListLen() != filledSlots {
				t.Errorf("Unexpected number of keys stored. Expected %d, got map with %d and list with %d", filledSlots, lCache.MapLen(), lCache.ListLen())
			}

			lCache.TraverseList(func(key string, value string) {
				if key != value {
					// key and value should be the same according to what has been stored
					t.Errorf("Wrong value stored under key %s, expected %s, got %s", key, key, value)
				}
			})
			// We can't assume the threads will spawn with the same order each time.
			// If the cache is limited some keys will be removed, but we don't know which ones
			// because the order isn't predictable.
		})
	}
}

func concurrentRetrieve(t *testing.T) {
	capacities := []int{2, 35, 100, 10000}
	for _, capacity := range capacities {
		t.Run("Cap"+strconv.Itoa(capacity), func(t *testing.T) {
			concurrentRetrieveWithCap(t, capacity)
		})
	}
}

func concurrentRetrieveWithCap(t *testing.T, lCacheCap int) {
	concurrentThreads := []int{10, 50}

	for _, threads := range concurrentThreads {
		t.Run("Threads"+strconv.Itoa(threads), func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{
				"capacity": lCacheCap,
			})

			var ttl int64
			for i := 0; i < lCacheCap; i++ {
				v := strconv.Itoa(i)
				ttl = 500
				if i%2 == 1 {
					// if i is odd
					ttl = -500 // ensure it's expired
				}
				lCache.Store(v, v, ttl)
			}

			wg := sync.WaitGroup{}
			wg.Add(threads)

			for i := 0; i < threads; i++ {
				go func(i int) {
					for j := 0; j < lCache.MaxCap(); j++ {
						v := strconv.Itoa(j)
						value, exists, err := lCache.Retrieve(v)
						if err != nil {
							t.Errorf("Thread %d -> error retrieving key %s: %v", i, v, err)
						}
						if j%2 == 0 {
							if !exists {
								t.Errorf("Thread %d -> retrieved expired key %s", i, v)
							}
						} else {
							if exists {
								t.Errorf("Thread %d -> retrieved key that should be expired %s", i, v)
							}
						}
						if exists && value != v {
							// key and value are expected to be the same
							t.Errorf("Thread %d -> wrong value retrieved, expected %s, got %s", i, v, value)
						}
					}
					wg.Done()
				}(i)
			}
			wg.Wait()

			// we should have accessed to all the keys, so the expired ones should have been removed
			filledSlots := (lCacheCap / 2) + (lCacheCap % 2)

			if lCache.MapLen() != filledSlots || lCache.ListLen() != filledSlots {
				t.Errorf("Unexpected number of keys stored. Expected %d, got map with %d and list with %d", filledSlots, lCache.MapLen(), lCache.ListLen())
			}

			lCache.TraverseList(func(key string, value string) {
				if key != value {
					// key and value should be the same according to what has been stored
					t.Errorf("Wrong value stored under key %s, expected %s, got %s", key, key, value)
				}
			})
		})
	}
}

func concurrentRemove(t *testing.T) {
	capacities := []int{2, 35, 100, 10000}
	for _, capacity := range capacities {
		t.Run("Cap"+strconv.Itoa(capacity), func(t *testing.T) {
			concurrentRemoveWithCap(t, capacity)
		})
	}
}

func concurrentRemoveWithCap(t *testing.T, lCacheCap int) {
	concurrentThreads := []int{10, 50}

	for _, threads := range concurrentThreads {
		t.Run("Threads"+strconv.Itoa(threads), func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{
				"capacity": lCacheCap,
			})

			var ttl int64
			for i := 0; i < lCacheCap; i++ {
				v := strconv.Itoa(i)
				ttl = 500
				if i%2 == 1 {
					// if i is odd
					ttl = -500 // ensure it's expired
				}
				lCache.Store(v, v, ttl)
			}

			wg := sync.WaitGroup{}
			wg.Add(threads)

			for i := 0; i < threads; i++ {
				go func(i int) {
					// all threads will try to remove all items at the same time
					shuffledList := rand.Perm(lCache.MaxCap())
					for _, j := range shuffledList {
						v := strconv.Itoa(j)
						err := lCache.Remove(v)
						if err != nil {
							t.Errorf("Thread %d -> error retrieving key %s: %v", i, v, err)
						}
					}
					wg.Done()
				}(i)
			}
			wg.Wait()

			if lCache.MapLen() != 0 || lCache.ListLen() != 0 {
				t.Errorf("Unexpected number of keys stored. Expected 0, got map with %d and list with %d", lCache.MapLen(), lCache.ListLen())
			}

			lCache.TraverseList(func(key string, value string) {
				if key != value {
					// key and value should be the same according to what has been stored
					t.Errorf("Wrong value stored under key %s, expected %s, got %s", key, key, value)
				}
			})
		})
	}
}

func concurrentMix(t *testing.T) {
	capacities := []int{2, 35, 100, 10000}
	for _, capacity := range capacities {
		t.Run("Cap"+strconv.Itoa(capacity), func(t *testing.T) {
			concurrentRemoveWithCap(t, capacity)
		})
	}
}

func concurrentMixWithCap(t *testing.T, lCacheCap int) {
	concurrentThreads := []int{10, 50}

	for _, threads := range concurrentThreads {
		t.Run("Threads"+strconv.Itoa(threads), func(t *testing.T) {
			lCache := NewLocalCache()
			lCache.Initialize(map[string]interface{}{
				"capacity": lCacheCap,
			})

			wg := sync.WaitGroup{}
			wg.Add(threads)

			for i := 0; i < threads; i++ {
				switch i % 3 {
				case 0:
					// storing thread
					go func(i int) {
						// all threads will try to remove all items at the same time
						shuffledList := rand.Perm(50000)
						for _, j := range shuffledList {
							v := strconv.Itoa(j)
							err := lCache.Store(v, v, 500)
							if err != nil {
								t.Errorf("Thread %d -> error storing key %s: %v", i, v, err)
							}
						}
						wg.Done()
					}(i)
				case 1:
					// retrieving thread
					go func(i int) {
						// all threads will try to remove all items at the same time
						shuffledList := rand.Perm(50000)
						for _, j := range shuffledList {
							v := strconv.Itoa(j)
							value, exists, err := lCache.Retrieve(v)
							if err != nil {
								t.Errorf("Thread %d -> error retrieving key %s: %v", i, v, err)
							}
							if exists && value != v {
								t.Errorf("Thread %d -> wrong value for key %s: expected %s, got %s", i, v, v, value)
							}
						}
						wg.Done()
					}(i)
				case 2:
					// removing thread
					go func(i int) {
						// all threads will try to remove all items at the same time
						shuffledList := rand.Perm(50000)
						for _, j := range shuffledList {
							v := strconv.Itoa(j)
							err := lCache.Remove(v)
							if err != nil {
								t.Errorf("Thread %d -> error removing key %s: %v", i, v, err)
							}
						}
						wg.Done()
					}(i)
				}
			}
			wg.Wait()

			// results can be quite random, just perform sanity checks
			if lCache.MapLen() != lCache.ListLen() {
				t.Errorf("Unexpected number of keys stored. Got map with %d and list with %d", lCache.MapLen(), lCache.ListLen())
			}

			if lCache.MapLen() > lCache.MaxCap() || lCache.ListLen() > lCache.MaxCap() {
				t.Errorf("Unexpected number of keys stored. Max capacity %d exceeded. Got map with %d and list with %d", lCache.MaxCap(), lCache.MapLen(), lCache.ListLen())
			}

			lCache.TraverseList(func(key string, value string) {
				if key != value {
					// key and value should be the same according to what has been stored
					t.Errorf("Wrong value stored under key %s, expected %s, got %s", key, key, value)
				}
			})
		})
	}
}
func TestConcurrent(t *testing.T) {
	t.Run("Store", concurrentStore)
	t.Run("Retrieve", concurrentRetrieve)
	t.Run("Remove", concurrentRemove)
	t.Run("Mix", concurrentMix)
}

func BenchmarkStore(b *testing.B) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var bench *hrtime.Benchmark
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			bench = hrtime.NewBenchmark(b.N)
			i := 0
			for bench.Next() {
				v := strconv.Itoa(i)
				lCache.Store(v, v, 222)
				i++
			}
		})
		fmt.Println(bench.Histogram(10).StringStats())
	}
}

func BenchmarkStoreSame(b *testing.B) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var bench *hrtime.Benchmark
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			bench = hrtime.NewBenchmark(b.N)
			for bench.Next() {
				lCache.Store("key1", "value1", 222)
			}
		})
		fmt.Println(bench.Histogram(10).StringStats())
	}
}

func BenchmarkRetrieveMiss(b *testing.B) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var bench *hrtime.Benchmark
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			for i := 0; i < lCache.MaxCap(); i++ {
				v := strconv.Itoa(-i)
				lCache.Store(v, v, 222)
			}

			bench = hrtime.NewBenchmark(b.N)
			i := 0
			for bench.Next() {
				v := strconv.Itoa(i)
				lCache.Retrieve(v)
				i++
			}
		})
		fmt.Println(bench.Histogram(10).StringStats())
	}
}

func BenchmarkRetrieveHit(b *testing.B) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var bench *hrtime.Benchmark
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			for i := 0; i < lCache.MaxCap(); i++ {
				v := strconv.Itoa(i)
				lCache.Store(v, v, 222)
			}

			bench = hrtime.NewBenchmark(b.N)
			i := 0
			for bench.Next() {
				v := strconv.Itoa(i % lCache.MaxCap())
				lCache.Retrieve(v)
				i++
			}
		})
		fmt.Println(bench.Histogram(10).StringStats())
	}
}

func concurrentStoreBench(b *testing.B, threads int) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var stopwatch *hrtime.Stopwatch
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			stopwatch = hrtime.NewStopwatch(b.N)
			maxThreads := threads
			if b.N < threads {
				maxThreads = b.N
			}
			for i := 0; i < maxThreads; i++ {
				lapsToDo := b.N / maxThreads
				if i < (b.N % maxThreads) {
					lapsToDo++
				}
				go func(i int, lapsToDo int, sw *hrtime.Stopwatch) {
					for j := 0; j < lapsToDo; j++ {
						lap := sw.Start()
						v := strconv.FormatInt(int64(lap), 10)
						lCache.Store(v, v, 500)
						sw.Stop(lap)
					}
				}(i, lapsToDo, stopwatch)
			}
			stopwatch.Wait()
		})
		fmt.Println(stopwatch.Histogram(10).StringStats())
	}
}

func concurrentRetrieveBench(b *testing.B, threads int) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var stopwatch *hrtime.Stopwatch
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			for i := 0; i < lCache.MaxCap(); i++ {
				v := strconv.Itoa(i)
				lCache.Store(v, v, 500)
			}

			stopwatch = hrtime.NewStopwatch(b.N)
			maxThreads := threads
			if b.N < threads {
				maxThreads = b.N
			}
			for i := 0; i < maxThreads; i++ {
				lapsToDo := b.N / maxThreads
				if i < (b.N % maxThreads) {
					lapsToDo++
				}
				go func(i int, lapsToDo int, sw *hrtime.Stopwatch) {
					for j := 0; j < lapsToDo; j++ {
						lap := sw.Start()
						v := strconv.FormatInt(int64(lap), 10)
						lCache.Retrieve(v)
						sw.Stop(lap)
					}
				}(i, lapsToDo, stopwatch)
			}
			stopwatch.Wait()
		})
		fmt.Println(stopwatch.Histogram(10).StringStats())
	}
}

func concurrentRemoveBench(b *testing.B, threads int) {
	benchTest := map[string]map[string]interface{}{
		"DefCap": {},
		"LimCap": {
			"capacity": 3,
		},
		"BigCap": {
			"capacity": 1 * 1000 * 1000,
		},
	}
	for testname, initParams := range benchTest {
		var stopwatch *hrtime.Stopwatch
		b.Run(testname, func(b *testing.B) {
			lCache := NewLocalCache()
			lCache.Initialize(initParams)

			for i := 0; i < lCache.MaxCap(); i++ {
				v := strconv.Itoa(i)
				lCache.Store(v, v, 500)
			}

			stopwatch = hrtime.NewStopwatch(b.N)
			maxThreads := threads
			if b.N < threads {
				maxThreads = b.N
			}
			for i := 0; i < maxThreads; i++ {
				lapsToDo := b.N / maxThreads
				if i < (b.N % maxThreads) {
					lapsToDo++
				}
				go func(i int, lapsToDo int, sw *hrtime.Stopwatch) {
					for j := 0; j < lapsToDo; j++ {
						lap := sw.Start()
						v := strconv.FormatInt(int64(lap), 10)
						lCache.Remove(v)
						sw.Stop(lap)
					}
				}(i, lapsToDo, stopwatch)
			}
			stopwatch.Wait()
		})
		fmt.Println(stopwatch.Histogram(10).StringStats())
	}
}

func BenchmarkConcurrent(b *testing.B) {
	threads := []int{3, 10, 50}
	for _, nThreads := range threads {
		nt := strconv.Itoa(nThreads)
		b.Run("StoreT"+nt, func(b *testing.B) {
			concurrentStoreBench(b, nThreads)
		})
		b.Run("RetrieveT"+nt, func(b *testing.B) {
			concurrentStoreBench(b, nThreads)
		})
		b.Run("RemoveT"+nt, func(b *testing.B) {
			concurrentStoreBench(b, nThreads)
		})
	}
}
