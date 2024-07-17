package configuration

import (
	"fmt"

	"github.com/eevan78/translit/internal/dictionary"
	"github.com/spf13/viper"
)

func ConfigInit() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("../../configs/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")

	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Грешка при читању конфигурационог фајла, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("outputDir", "../../output")
	viper.SetDefault("version", "v0.3.0")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Грешка при раду са конфигурационим фајлом, %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("Output directory is\t", configuration.OutputDir)

	// Reading variables without using the model
	fmt.Println("\nReading variables without using the model..")
	fmt.Println("Version is\t", viper.GetString("Version"))
	fmt.Println("OutputDir is\t", viper.GetString("OutputDir"))

	dictionary.OutputDir = viper.GetString("OutputDir")
	dictionary.Version = viper.GetString("Version")
}
