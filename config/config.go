package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	FetchConfig bool
	RemoteConfig string
	LocalBranchNameConfig string
	RemoteBranchNameConfig string
)

func InitConfig() {
	loadConfig()

	FetchConfig = viper.GetBool("fetch")
	RemoteConfig = viper.GetString("remote")
	LocalBranchNameConfig = viper.GetString("local-branch-name")
	RemoteBranchNameConfig = viper.GetString("remote-branch-name")
}

func loadConfig() {
	viper.SetConfigName(".gitrelrc")
	viper.SetConfigType("env")

	// Look up the directory tree
	dir, err := os.Getwd()
	if err == nil {
		for {
			viper.AddConfigPath(dir)
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	viper.AddConfigPath("$HOME")

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return
	}

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
