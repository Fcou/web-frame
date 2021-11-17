package orm

import (
	"testing"

	"github.com/Fcou/web-frame/framework/contract"
	"github.com/Fcou/web-frame/framework/provider/config"
	tests "github.com/Fcou/web-frame/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFcouConfig_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.FcouConfigProvider{})

	Convey("test config", t, func() {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		config := &contract.DBConfig{}
		err := configService.Load("database.default", config)
		So(err, ShouldBeNil)
	})

	Convey("test default config", t, func() {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		config := &contract.DBConfig{
			ConnMaxIdle: 10,
		}
		err := configService.Load("database.read", config)
		So(err, ShouldBeNil)
		So(config.ConnMaxIdle, ShouldEqual, 10)
	})

	Convey("test base config", t, func() {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		config := &contract.DBConfig{
			ConnMaxOpen: 200,
		}
		err := configService.Load("database", config)
		So(err, ShouldBeNil)
		So(config.ConnMaxOpen, ShouldEqual, 100)
	})

}
