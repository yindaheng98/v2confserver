package v2confserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/yindaheng98/vmessconfig"
	"github.com/yindaheng98/vmessconfig/cmd/args"
	"net/http"
	"time"
)

type V2CmdConfig struct {
	VmessCmdConfig *args.CmdConfig `flag:"-"`
	v2conf         *conf.Config
	Interval       uint   `desc:"Interval for get and ping vmess outbounds"`
	Addr           string `desc:"Address where the server listen on"`
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

func (c *V2CmdConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(c.v2conf, "", " ")
	if err != nil {
		fmt.Printf("V2Config.Response failed: %+v\n", err)
	}
	_, err = fmt.Fprint(w, string(j))
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
