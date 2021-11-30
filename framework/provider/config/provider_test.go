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

func TestHadeConfig_Normal(t *testing.T) {
	Convey("test fcou config normal case", t, func() {
		basePath := tests.BasePath
		c := framework.NewHadeContainer()
		c.Bind(&app.HadeAppProvider{BaseFolder: basePath})
		c.Bind(&env.HadeEnvProvider{})

		err := c.Bind(&HadeConfigProvider{})
		So(err, ShouldBeNil)

		conf := c.MustMake(contract.ConfigKey).(contract.Config)
		So(conf.GetString("database.default.host"), ShouldEqual, "localhost")
		So(conf.GetInt("database.default.port"), ShouldEqual, 3306)
		//So(conf.GetFloat64("database.default.readtime"), ShouldEqual, 2.3)
		// So(conf.GetString("database.mysql.password"), ShouldEqual, "mypassword")

		maps := conf.GetStringMap("database.default")
		So(maps, ShouldContainKey, "host")
		So(maps["host"], ShouldEqual, "localhost")

		maps2 := conf.GetStringMapString("database.default")
		So(maps2["host"], ShouldEqual, "localhost")

		type Mysql struct {
			Host string `yaml:"host"`
		}
		ms := &Mysql{}
		err = conf.Load("database.default", ms)
		So(err, ShouldBeNil)
		So(ms.Host, ShouldEqual, "localhost")
	})
}
