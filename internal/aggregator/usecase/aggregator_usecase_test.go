package usecase

import (
	"github.com/a5347354/rise-workshop/internal"
	eventMock "github.com/a5347354/rise-workshop/internal/event/mocks"

	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_aggregatorUsecase_ListEventByKeyword(t *testing.T) {
	Convey("ListEventByKeyword", t, func() {
		mockEventStore, usecase := given(t)
		keyword := "Android"
		mockEventStore.EXPECT().SearchByContent(gomock.Any(), keyword).Return([]internal.Event{
			{
				PK:      1,
				ID:      "11111",
				Kind:    1,
				Content: "Android",
			},
		}, nil)
		events, err := usecase.ListEventByKeyword(context.TODO(), keyword)
		So(err, ShouldBeNil)
		So(events, ShouldHaveLength, 1)
	})
}

func given(t *testing.T) (*eventMock.MockStore, *aggregatorUsecase) {
	ctrl := gomock.NewController(t)
	mockStore := eventMock.NewMockStore(ctrl)
	u := &aggregatorUsecase{
		eStore:      mockStore,
		url:         nil,
		lc:          nil,
		limitClient: nil,
	}
	return mockStore, u
}
