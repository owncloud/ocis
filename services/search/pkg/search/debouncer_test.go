package search_test

import (
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
)

var _ = Describe("SpaceDebouncer", func() {
	var (
		debouncer *search.SpaceDebouncer
		callCount map[string]int

		userId = &user.UserId{
			OpaqueId: "user",
		}
		spaceid = &sprovider.StorageSpaceId{
			OpaqueId: "spaceid",
		}
	)

	BeforeEach(func() {
		callCount = map[string]int{}
		debouncer = search.NewSpaceDebouncer(50*time.Millisecond, func(id *sprovider.StorageSpaceId, _ *user.UserId) {
			callCount[id.OpaqueId] += 1
		})
	})

	It("debounces", func() {
		debouncer.Debounce(spaceid, userId)
		debouncer.Debounce(spaceid, userId)
		debouncer.Debounce(spaceid, userId)
		Eventually(func() int {
			return callCount["spaceid"]
		}, "200ms").Should(Equal(1))
	})

	It("works multiple times", func() {
		debouncer.Debounce(spaceid, userId)
		debouncer.Debounce(spaceid, userId)
		debouncer.Debounce(spaceid, userId)
		time.Sleep(100 * time.Millisecond)

		debouncer.Debounce(spaceid, userId)
		debouncer.Debounce(spaceid, userId)

		Eventually(func() int {
			return callCount["spaceid"]
		}, "200ms").Should(Equal(2))
	})

	It("doesn't trigger twice simultaneously", func() {
		debouncer = search.NewSpaceDebouncer(50*time.Millisecond, func(id *sprovider.StorageSpaceId, _ *user.UserId) {
			callCount[id.OpaqueId] += 1
			time.Sleep(300 * time.Millisecond)
		})
		debouncer.Debounce(spaceid, userId)
		time.Sleep(100 * time.Millisecond) // Let it trigger once

		debouncer.Debounce(spaceid, userId)
		time.Sleep(100 * time.Millisecond) // shouldn't trigger as the other run is still in progress
		Expect(callCount["spaceid"]).To(Equal(1))

		Eventually(func() int {
			return callCount["spaceid"]
		}, "500ms").Should(Equal(2))
	})
})
