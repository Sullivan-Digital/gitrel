package main

import (
	"gitrel/cmd"
	"gitrel/config"
)

func main() {
	config.InitConfig()
	cmd.Execute()
}
