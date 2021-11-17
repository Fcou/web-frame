package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Fcou/web-frame/framework/provider/config"
	"github.com/Fcou/web-frame/framework/provider/log"
	tests "github.com/Fcou/web-frame/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFcouService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.FcouConfigProvider{})
	container.Bind(&log.FcouLogServiceProvider{})

	Convey("test get client", t, func() {
		fcouRedis, err := NewFcouRedis(container)
		So(err, ShouldBeNil)
		service, ok := fcouRedis.(*FcouRedis)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("redis.write"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		ctx := context.Background()
		err = client.Set(ctx, "foo", "bar", 1*time.Hour).Err()
		So(err, ShouldBeNil)
		val, err := client.Get(ctx, "foo").Result()
		So(err, ShouldBeNil)
		So(val, ShouldEqual, "bar")
		err = client.Del(ctx, "foo").Err()
		So(err, ShouldBeNil)
	})
}
