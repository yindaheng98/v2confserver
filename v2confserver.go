package v2confserver

import (
	"context"
	"fmt"
	"github.com/yindaheng98/vmessconfig"
	"github.com/yindaheng98/vmessconfig/cmd/args"
	"net/http"
	"time"
)

type V2CmdConfig struct {
	VmessCmdConfig *args.CmdConfig `flag:"-"`
	v2conf         vmessconfig.V2Config
	Interval       uint   `desc:"Interval for get and ping vmess outbounds"`
	Addr           string `desc:"Address where the server listen on"`
	GetVmessList   string `desc:"How to GetVmessList"`
}

func (c *V2CmdConfig) Routine(ctx context.Context) {
	template, err := c.VmessCmdConfig.TemplateConfig.Template()
	if err != nil {
		fmt.Printf("V2Config.Routine failed: %+v\n", err)
		return
	}
	v2conf, err := vmessconfig.VmessConfig(c.VmessCmdConfig.Urls, template, c.VmessCmdConfig.Config, ctx)
	if err != nil {
		fmt.Printf("V2Config.Routine failed: %+v\n", err)
	} else {
		c.v2conf = v2conf
	}
}

func (c *V2CmdConfig) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, string(c.v2conf))
	if err != nil {
		fmt.Printf("V2Config.Response failed: %+v\n", err)
	}
}

func (c *V2CmdConfig) Start(ctx context.Context) {
	c.Routine(ctx)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(c.Interval)):
				c.Routine(ctx)
			}
		}
	}(ctx)
}
