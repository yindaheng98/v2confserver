package main

import (
	"context"
	"fmt"
	"github.com/yindaheng98/v2confserver"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import "github.com/yindaheng98/vmessconfig/cmd/args"

func exit(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%s [balancer|single] -urls https://... -urls https://...\n", os.Args[0])
	args.PrintUsage()
	os.Exit(1)
}

func main() {
	v2CmdConfig := &v2confserver.V2CmdConfig{
		VmessCmdConfig: args.NewCmdConfig(),
		Interval:       1800,
	}
	err := args.AddCmdArgs(v2CmdConfig)
	if err != nil {
		exit(err)
	}
	err = v2CmdConfig.VmessCmdConfig.GenerateCmdArgs()
	if err != nil {
		exit(err)
	}
	err = v2CmdConfig.VmessCmdConfig.ParseCmdArgs(os.Args[1:])
	if err != nil {
		exit(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)

	v2CmdConfig.Start(ctx)
	http.HandleFunc("/", v2CmdConfig.Response)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		return
	}
}
