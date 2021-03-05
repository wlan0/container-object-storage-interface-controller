package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		klog.Fatal(err)
	}
}
