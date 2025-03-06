package search_test

import (
	"context"
	"sync/atomic"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	searchMocks "github.com/owncloud/ocis/v2/services/search/pkg/search/mocks"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/stretchr/testify/mock"
	mEvents "go-micro.dev/v4/events"
)

var _ = DescribeTable("events",
	func(mcks []string, e interface{}, asyncUploads bool) {
		var (
			s     = &searchMocks.Searcher{}
			calls atomic.Int32
		)

		bus, _ := mEvents.NewStream()

		search.HandleEvents(s, bus, log.NewLogger(), &config.Config{
			Events: config.Events{
				AsyncUploads: asyncUploads,
			},
		})

		for _, mck := range mcks {
			s.On(mck, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				calls.Add(1)
			})
		}

		err := events.Publish(context.Background(), bus, e)

		Expect(err).To(BeNil())
		Eventually(func() int {
			return int(calls.Load())
		}, "2s").Should(Equal(len(mcks)))
	},
	Entry("ItemTrashed", []string{"TrashItem", "IndexSpace"}, events.ItemTrashed{}, false),
	Entry("ItemMoved", []string{"MoveItem", "IndexSpace"}, events.ItemMoved{}, false),
	Entry("ItemRestored", []string{"RestoreItem", "IndexSpace"}, events.ItemRestored{}, false),
	Entry("ContainerCreated", []string{"IndexSpace"}, events.ContainerCreated{}, false),
	Entry("FileTouched", []string{"IndexSpace"}, events.FileTouched{}, false),
	Entry("FileVersionRestored", []string{"IndexSpace"}, events.FileVersionRestored{}, false),
	Entry("TagsAdded", []string{"UpsertItem"}, events.TagsAdded{}, false),
	Entry("TagsRemoved", []string{"UpsertItem"}, events.TagsRemoved{}, false),
	Entry("FileUploaded", []string{"IndexSpace"}, events.FileUploaded{}, false),
	Entry("UploadReady", []string{"IndexSpace"}, events.UploadReady{ExecutingUser: &userv1beta1.User{}}, true),
)
