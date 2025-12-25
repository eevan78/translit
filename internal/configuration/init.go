package configuration

import (
	"flag"
	"fmt"

	"github.com/eevan78/translit/internal/terminal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configuration Configurations
	configVersion string
)

func ConfigInit() {
	if *terminal.ConfigPtr {
		readConfig()
		initVars()
		initFlags()
	}
}

func readConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../configs/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Грешка при читању конфигурационог фајла, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Грешка при раду са конфигурационим фајлом, %v", err)
	}
}

func initVars() {
	defaultVars()

	terminal.OutputDir = viper.GetString("OutputDir")
	configVersion = viper.GetString("Version")
}

func defaultVars() {
	viper.SetDefault("outputDir", "../../output")
	viper.SetDefault("version", "v0.3.0")
}

func initFlags() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	*terminal.C2lPtr = configuration.C2lPtr
	*terminal.L2cPtr = configuration.L2cPtr
	*terminal.HtmlPtr = configuration.HtmlPtr
	*terminal.TextPtr = configuration.TextPtr
	*terminal.InputPathPtr = configuration.InputPathPtr
}
