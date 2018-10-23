package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/annymsmthd/go-modules-registry/pkg/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	port *int
)

var rootCmd = &cobra.Command{
	Use:   "go-modules-registry",
	Short: "go-modules-registry is a self hosted registry for all your private go module needs",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		storage := viper.GetString("storage")
		settings := &server.Settings{
			FileStorageBasePath: storage,
			Port:                *port,
		}

		server, err := server.NewServer(settings)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ctx, cancel := context.WithCancel(context.Background())
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(server.Run())

		// Wait for SIGINT/SIGTERM
		waiter := make(chan os.Signal, 1)
		signal.Notify(waiter, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-waiter:
		case <-ctx.Done():
		}
		cancel()
		if err := grp.Wait(); err != nil {
			panic(err)
		}
	},
}

func init() {
	port = rootCmd.Flags().IntP("port", "p", 80, "The port to host the server on")
	rootCmd.Flags().StringP("storage", "s", "/tmp/storage", "The storage location for modules")

	viper.BindEnv("storage", "STORAGE_LOCATION")
	viper.BindPFlag("storage", rootCmd.Flag("storage"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
