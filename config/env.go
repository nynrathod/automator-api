package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

type envConfigs struct {
	DbName          string `mapstructure:"DB_NAME"`
	DbHost          string `mapstructure:"DB_HOST"`
	JWTSecrete      string `mapstructure:"JWT_SECRET"`
	SmtpEmail       string `mapstructure:"SMTP_EMAIL"`
	SmtpToken       string `mapstructure:"SMTP_TOKEN"`
	SmtpToEmail     string `mapstructure:"SMTP_TO_EMAIL"`
	VerrifyOtpToken string `mapstructure:"VERIFY_OTP_TOKEN"`
	ENV             string `mapstructure:"ENV"`
}

// Call to load the variables from env
func loadEnvVariables() (config *envConfigs) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName(".env")
	//viper.SetConfigName(".env")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
