package id

import (
	tests "github.com/Fcou/web-frame/test"
	"testing"

	"github.com/Fcou/web-frame/framework/contract"
	"github.com/Fcou/web-frame/framework/provider/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConsoleLog_Normal(t *testing.T) {
	Convey("test fcou console log normal case", t, func() {
		c := tests.InitBaseContainer()
		c.Bind(&config.HadeConfigProvider{})

		err := c.Bind(&HadeIDProvider{})
		So(err, ShouldBeNil)

		idService := c.MustMake(contract.IDKey).(contract.IDService)
		xid := idService.NewID()
		So(xid, ShouldNotBeEmpty)
	})
}
