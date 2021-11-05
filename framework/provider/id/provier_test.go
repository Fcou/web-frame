package id

import (
	"testing"

	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
	"github.com/Fcou/web-frame/framework/provider/app"
	"github.com/Fcou/web-frame/framework/provider/config"
	"github.com/Fcou/web-frame/framework/provider/env"
	"github.com/Fcou/web-frame/framework/util"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConsoleLog_Normal(t *testing.T) {
	Convey("test fcou console log normal case", t, func() {
		basePath := util.GetExecDirectory()
		c := framework.NewFcouContainer()
		c.Singleton(&app.FcouAppProvider{BasePath: basePath})
		c.Singleton(&env.FcouEnvProvider{})
		c.Singleton(&config.FcouConfigProvider{})

		err := c.Singleton(&FcouIDProvider{})
		So(err, ShouldBeNil)

		idService := c.MustMake(contract.IDKey).(contract.IDService)
		xid := idService.NewID()
		t.Log(xid)
		So(xid, ShouldNotBeEmpty)
	})
}
