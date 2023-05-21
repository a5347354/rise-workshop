package usecase

import (
	eventMock "github.com/a5347354/rise-workshop/internal/event/mocks"
	clientMock "github.com/a5347354/rise-workshop/pkg/mocks"

	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nbd-wtf/go-nostr"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/sync/errgroup"
)

func Test_clientUsecase_Collect(t *testing.T) {
	Convey("Collect", t, func() {
		mockEventStore, mockNostrClient, usecase := given(t)
		url := "http://localhost"
		mockNostrClient.EXPECT().ConnectURL(gomock.Any(), url).Return(nil)
		eventChan := make(chan *nostr.Event, len(url))
		go func() {
			eventChan <- &nostr.Event{
				ID:      "1111111",
				Kind:    1,
				Content: "Android",
			}
			close(eventChan)
		}()
		mockNostrClient.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(&nostr.Subscription{
			Events: eventChan,
		}, nil)
		mockEventStore.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		var g errgroup.Group
		g.Go(func() error {
			return usecase.Collect(context.TODO(), url)
		})
		err := g.Wait()
		So(err, ShouldBeNil)
	})
}

func given(t *testing.T) (*eventMock.MockStore, *clientMock.MockNostrClient, *clientUsecase) {
	ctrl := gomock.NewController(t)
	mockStore := eventMock.NewMockStore(ctrl)
	mockNostrClient := clientMock.NewMockNostrClient(ctrl)
	u := &clientUsecase{
		client: mockNostrClient,
		eStore: mockStore,
	}
	return mockStore, mockNostrClient, u
}
