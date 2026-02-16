package search_test

import (
	"sync/atomic"
	"time"

	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
)

var _ = Describe("SpaceDebouncer", func() {
	var (
		debouncer *search.SpaceDebouncer

		callCount atomic.Int32

		spaceid = &sprovider.StorageSpaceId{
			OpaqueId: "spaceid",
		}
	)

	BeforeEach(func() {
		callCount = atomic.Int32{}
		debouncer = search.NewSpaceDebouncer(50*time.Millisecond, func(id *sprovider.StorageSpaceId) {
			if id.OpaqueId == "spaceid" {
				callCount.Add(1)
			}
		})
	})

	It("debounces", func() {
		debouncer.Debounce(spaceid)
		debouncer.Debounce(spaceid)
		debouncer.Debounce(spaceid)
		Eventually(func() int {
			return int(callCount.Load())
		}, "200ms").Should(Equal(1))
	})

	It("works multiple times", func() {
		debouncer.Debounce(spaceid)
		debouncer.Debounce(spaceid)
		debouncer.Debounce(spaceid)
		time.Sleep(100 * time.Millisecond)

		debouncer.Debounce(spaceid)
		debouncer.Debounce(spaceid)

		Eventually(func() int {
			return int(callCount.Load())
		}, "200ms").Should(Equal(2))
	})

	It("doesn't trigger twice simultaneously", func() {
		debouncer = search.NewSpaceDebouncer(50*time.Millisecond, func(id *sprovider.StorageSpaceId) {
			if id.OpaqueId == "spaceid" {
				callCount.Add(1)
			}
			time.Sleep(300 * time.Millisecond)
		})
		debouncer.Debounce(spaceid)
		time.Sleep(100 * time.Millisecond) // Let it trigger once

		debouncer.Debounce(spaceid)
		time.Sleep(100 * time.Millisecond) // shouldn't trigger as the other run is still in progress
		Expect(int(callCount.Load())).To(Equal(1))

		Eventually(func() int {
			return int(callCount.Load())
		}, "500ms").Should(Equal(2))
	})
})
