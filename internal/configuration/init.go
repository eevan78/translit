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
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("outputDir", "../../output")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("Output directory is\t", configuration.OutputDir)

	// Reading variables without using the model
	fmt.Println("\nReading variables without using the model..")
	fmt.Println("Database is\t", viper.GetString("OutputDir"))

	dictionary.OutputDir = viper.GetString("OutputDir")
}
