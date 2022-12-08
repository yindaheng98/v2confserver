package main

import (
	"context"
	"fmt"
	"github.com/yindaheng98/v2confserver"
	"github.com/yindaheng98/vmessconfig"
	"github.com/yindaheng98/vmessconfig/cmd/args"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
		Addr:           ":80",
		GetVmessList:   "default",
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

	if v2CmdConfig.GetVmessList == "wget" {
		vmessconfig.CustomizeGetVmessList(vmessconfig.WgetGetVmessList)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)

	server := &http.Server{
		Addr:    v2CmdConfig.Addr,
		Handler: v2CmdConfig,
	}
	v2CmdConfig.Start(ctx)
	errCh := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		cancel()
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	case <-ctx.Done():
		shutdownCtx, shutdownCtxCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCtxCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}
	}
}
