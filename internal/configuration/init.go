package configuration

import (
	"flag"
	"fmt"

	"github.com/eevan78/translit/internal/dictionary"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configuration Configurations

func ConfigInit() {
	if *dictionary.ConfigPtr == true {
		readConfig()
		initVars()
		inifFlags()
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

	dictionary.OutputDir = viper.GetString("OutputDir")
	dictionary.Version = viper.GetString("Version")
}

func defaultVars() {
	viper.SetDefault("outputDir", "../../output")
	viper.SetDefault("version", "v0.3.0")
}

func inifFlags() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	defaultFlags()
	*dictionary.C2lPtr = configuration.C2lPtr
	*dictionary.L2cPtr = configuration.L2cPtr
	*dictionary.HtmlPtr = configuration.HtmlPtr
	*dictionary.TextPtr = configuration.TextPtr
	*dictionary.InputPathPtr = configuration.InputPathPtr
}

func defaultFlags() {
	configuration.C2lPtr = false
	configuration.L2cPtr = true
	configuration.HtmlPtr = false
	configuration.TextPtr = true
	configuration.InputPathPtr = ""
}
