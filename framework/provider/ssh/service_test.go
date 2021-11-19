package ssh

import (
	"github.com/Fcou/web-frame/framework/provider/config"
	"github.com/Fcou/web-frame/framework/provider/log"
	tests "github.com/Fcou/web-frame/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFcouSSHService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.FcouConfigProvider{})
	container.Bind(&log.FcouLogServiceProvider{})

	Convey("test get client", t, func() {
		fcouRedis, err := NewFcouSSH(container)
		So(err, ShouldBeNil)
		service, ok := fcouRedis.(*FcouSSH)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("ssh.web-01"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		session, err := client.NewSession()
		So(err, ShouldBeNil)
		out, err := session.Output("pwd")
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
		session.Close()
	})
}
