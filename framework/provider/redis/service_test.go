package redis

import (
	"context"
	"github.com/Fcou/web-frame/framework/provider/config"
	"github.com/Fcou/web-frame/framework/provider/log"
	tests "github.com/Fcou/web-frame/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestHadeService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.HadeConfigProvider{})
	container.Bind(&log.HadeLogServiceProvider{})

	Convey("test get client", t, func() {
		fcouRedis, err := NewHadeRedis(container)
		So(err, ShouldBeNil)
		service, ok := fcouRedis.(*HadeRedis)
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
