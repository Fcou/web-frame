package config

import (
	"testing"

	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
	"github.com/Fcou/web-frame/framework/provider/app"
	"github.com/Fcou/web-frame/framework/provider/env"
	tests "github.com/Fcou/web-frame/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFcouConfig_Normal(t *testing.T) {
	Convey("test fcou config normal case", t, func() {
		basePath := tests.BasePath
		c := framework.NewFcouContainer()
		c.Bind(&app.FcouAppProvider{BaseFolder: basePath})
		c.Bind(&env.FcouEnvProvider{})

		err := c.Bind(&FcouConfigProvider{})
		So(err, ShouldBeNil)

		conf := c.MustMake(contract.ConfigKey).(contract.Config)
		So(conf.GetString("database.mysql.hostname"), ShouldEqual, "127.0.0.1")
		So(conf.GetInt("database.mysql.timeout"), ShouldEqual, 1)
		So(conf.GetFloat64("database.mysql.readtime"), ShouldEqual, 2.3)
		// So(conf.GetString("database.mysql.password"), ShouldEqual, "mypassword")

		maps := conf.GetStringMap("database.mysql")
		So(maps, ShouldContainKey, "hostname")
		So(maps["timeout"], ShouldEqual, 1)

		maps2 := conf.GetStringMapString("databse.mysql")
		So(maps2["timeout"], ShouldEqual, "")

		type Mysql struct {
			Hostname string
			Username string
		}
		ms := &Mysql{}
		err = conf.Load("database.mysql", ms)
		Println(ms)
		So(err, ShouldBeNil)
		So(ms.Hostname, ShouldEqual, "127.0.0.1")
	})
}
