package env

import (
	"testing"

	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
	"github.com/Fcou/web-frame/framework/provider/app"
	tests "github.com/Fcou/web-frame/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFcouEnvProvider(t *testing.T) {
	Convey("test fcou env normal case", t, func() {
		basePath := tests.BasePath
		c := framework.NewFcouContainer()
		sp := &app.FcouAppProvider{BaseFolder: basePath}

		err := c.Bind(sp)
		So(err, ShouldBeNil)

		sp2 := &FcouEnvProvider{}
		err = c.Bind(sp2)
		So(err, ShouldBeNil)

		envServ := c.MustMake(contract.EnvKey).(contract.Env)
		So(envServ.AppEnv(), ShouldEqual, "development")
		// So(envServ.Get("DB_HOST"), ShouldEqual, "127.0.0.1")
		// So(envServ.AppDebug(), ShouldBeTrue)
	})
}
