package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	bucketcontroller "github.com/kubernetes-sigs/container-object-storage-interface-api/controller"
	"github.com/kubernetes-sigs/container-object-storage-interface-controller/pkg/bucketaccessrequest"
	"github.com/kubernetes-sigs/container-object-storage-interface-controller/pkg/bucketrequest"

	"github.com/golang/glog"
)

var cmd = &cobra.Command{
	Use:           "cosi-controller",
	Short:         "central controller for managing bucket* and bucketAccess* API objects",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(c *cobra.Command, args []string) error {
		return run(c.Context(), args)
	},
	DisableFlagsInUseLine: true,
}

var kubeConfig string

func init() {
	viper.AutomaticEnv()

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flag.Set("logtostderr", "true")

	strFlag := func(c *cobra.Command, ptr *string, name string, short string, dfault string, desc string) {
		c.PersistentFlags().
			StringVarP(ptr, name, short, dfault, desc)
	}
	strFlag(cmd, &kubeConfig, "kubeconfig", "", kubeConfig, "path to kubeconfig file")

	hideFlag := func(name string) {
		cmd.PersistentFlags().MarkHidden(name)
	}
	hideFlag("alsologtostderr")
	hideFlag("log_backtrace_at")
	hideFlag("log_dir")
	hideFlag("logtostderr")
	hideFlag("master")
	hideFlag("stderrthreshold")
	hideFlag("vmodule")

	// suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})
	viper.BindPFlags(cmd.PersistentFlags())

}

func run(ctx context.Context, args []string) error {
	ctrl, err := bucketcontroller.NewDefaultObjectStorageController("cosi-controller", "leader-lock", 40)
	if err != nil {
		return err
	}
	ctrl.AddBucketRequestListener(bucketrequest.NewListener())
	ctrl.AddBucketAccessRequestListener(bucketaccessrequest.NewListener())
	return ctrl.Run(ctx)
}
