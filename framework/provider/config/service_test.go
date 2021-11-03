package config

import (
	"path/filepath"
	"testing"

	"github.com/Fcou/web-frame/framework/contract"
	tests "github.com/Fcou/web-frame/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFcouConfig_GetInt(t *testing.T) {
	Convey("test fcou env normal case", t, func() {
		basePath := tests.BasePath
		folder := filepath.Join(basePath, "config")
		serv, err := NewFcouConfig(folder, map[string]string{}, contract.EnvDevelopment)
		So(err, ShouldBeNil)
		conf := serv.(*FcouConfig)
		timeout := conf.GetInt("database.mysql.timeout")
		So(timeout, ShouldEqual, 1)
	})
}
