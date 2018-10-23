package cmd

import (
	"fmt"
	"os"

	"github.com/annymsmthd/go-modules-registry/pkg/uploader"

	"github.com/coreos/go-semver/semver"
	"github.com/spf13/cobra"
)

var (
	registryHost   string
	version        string
	moduleLocation string
)

var rootCmd = &cobra.Command{
	Use:   "go-modules-registry-uploader",
	Short: "go-modules-registry is an uploader to put your git module into the registry",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		semversion, err := semver.NewVersion(version)
		if err != nil {
			fmt.Printf("%s is not a valid semver: %v\n", version, err)
			os.Exit(1)
		}

		_, err = os.Stat(moduleLocation)
		if err != nil {
			fmt.Printf("failed checking module location: %v\n", err)
			os.Exit(1)
		}

		loader := uploader.NewUploader(registryHost, moduleLocation, semversion)
		err = loader.Upload()
		if err != nil {
			fmt.Printf("failed uploading: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&registryHost, "registry", "r", "", "The location of the module registry")
	rootCmd.Flags().StringVarP(&version, "version", "v", "", "the version of the module you are uploading. Must be semver")
	rootCmd.Flags().StringVarP(&moduleLocation, "module", "m", "", "The location of the module directory")

	rootCmd.MarkFlagRequired("registry")
	rootCmd.MarkFlagRequired("version")
	rootCmd.MarkFlagRequired("module")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
