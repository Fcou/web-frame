package config

import (
	"path/filepath"
	"testing"

	"github.com/Fcou/web-frame/framework/contract"
	tests "github.com/Fcou/web-frame/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHadeConfig_GetInt(t *testing.T) {
	container := tests.InitBaseContainer()

	Convey("test fcou env normal case", t, func() {
		appService := container.MustMake(contract.AppKey).(contract.App)
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		folder := filepath.Join(appService.ConfigFolder(), envService.AppEnv())

		serv, err := NewHadeConfig(container, folder, map[string]string{})
		So(err, ShouldBeNil)
		conf := serv.(*HadeConfig)
		timeout := conf.GetString("database.default.timeout")
		So(timeout, ShouldEqual, "10s")
	})
}
