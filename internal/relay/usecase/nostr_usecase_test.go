package usecase

import (
	eventMock "github.com/a5347354/rise-workshop/internal/event/mocks"
	deliveryMock "github.com/a5347354/rise-workshop/internal/relay/mocks"

	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_relayUsecase_ReceiveMessage_Event(t *testing.T) {
	Convey("ReceiveMessage Event", t, func() {
		mockEventStore, usecase, mockDelivery := given(t)
		mockEventStore.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockDelivery.EXPECT().Broadcast(gomock.Any(), gomock.Any())
		event := []byte(`["EVENT",{"id":"34c7483fa2c6ad81c3fd97b70d569c3ec32eafb9082aee265cbcecfcef741bcf","pubkey":"826cf76ca45ccf7c24492e873b00c9c6d36caeb0bd478a62a7b9dd5b099af230","created_at":1684074156,"kind":1,"tags":[["d","b07v7s2ic0haospgmeg73i"],["media","https://media.zapstr.live:3118/d91191e30e00444b942c0e82cad470b32af171764c2275bee0bd99377efd4075/naddr1qqtxyvphwcmhxvnfvvcxsct0wdcxwmt9vumnx6gzyrv3ry0rpcqygju59s8g9jk5wzej4ut3wexzyad7uz7ejdm7l4q82qcyqqq856g4xnp7j","http"],["p","d91191e30e00444b942c0e82cad470b32af171764c2275bee0bd99377efd4075","Host"],["p","fa984bd7dbb282f07e16e7ae87b26a2a7b9b90b7246a44771f0cf5ae58018f52","Guest"],["c","Podcast"],["price","402"],["cover","https://s3-us-west-2.amazonaws.com/anchor-generated-image-bank/production/podcast_uploaded_nologo400/36291377/36291377-1673187804611-64b4f8e9f1687.jpg"],["subject","Nostrovia | The Pablo Episode"]],"content":"Nostrovia | The Pablo Episode\n\nhttps://s3-us-west-2.amazona","sig":"a9ef9e63b0f93b42d81fc74c62889c577dfc787a3f0a55fde85c634eafc8b0132703401d698cc0d3305b2109d2890d11cc57c94ce67e9039f132f31b2dec668a"}]`)
		_, err := usecase.ReceiveMessage(context.TODO(), event, nil)
		So(err, ShouldBeNil)
	})
}

func Test_relayUsecase_ReceiveMessage_Req(t *testing.T) {
	Convey("ReceiveMessage Req", t, func() {
		_, usecase, mockDelivery := given(t)
		mockDelivery.EXPECT().Subscribe(gomock.Any(), "publish-check:1", gomock.Any())
		event := []byte(`["REQ","publish-check:1",{"ids":["095ab2e1627b35c90944aa8e3d7ab084691a7a0702d7da9b486e1e6de6aa147c"]}]`)
		_, err := usecase.ReceiveMessage(context.TODO(), event, nil)
		So(err, ShouldBeNil)
	})
}

func Test_relayUsecase_ReceiveMessage_Close(t *testing.T) {
	Convey("ReceiveMessage Close", t, func() {
		_, usecase, mockDelivery := given(t)
		mockDelivery.EXPECT().UnSubscribe(gomock.Any(), "publish-check:3")
		event := []byte(`["CLOSE","publish-check:3"]`)
		_, err := usecase.ReceiveMessage(context.TODO(), event, nil)
		So(err, ShouldBeNil)
	})
}

func given(t *testing.T) (*eventMock.MockStore, *relayUsecase, *deliveryMock.MockNotification) {
	ctrl := gomock.NewController(t)
	mockStore := eventMock.NewMockStore(ctrl)
	mockNotification := deliveryMock.NewMockNotification(ctrl)
	u := &relayUsecase{
		mockStore,
		mockNotification,
	}
	return mockStore, u, mockNotification
}
