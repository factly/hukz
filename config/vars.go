package config

import (
	"log"

	"github.com/spf13/viper"
)

// SetupVars setups all the config variables to run application
func SetupVars() {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("hukz")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("config file not found...")
	}

	if !viper.IsSet("database_host") {
		log.Fatal("please provide database_host config param")
	}

	if !viper.IsSet("database_user") {
		log.Fatal("please provide database_user config param")
	}

	if !viper.IsSet("database_name") {
		log.Fatal("please provide database_name config param")
	}

	if !viper.IsSet("database_password") {
		log.Fatal("please provide database_password config param")
	}

	if !viper.IsSet("database_port") {
		log.Fatal("please provide database_port config param")
	}

	if !viper.IsSet("database_ssl_mode") {
		log.Fatal("please provide database_ssl_mode config param")
	}

	if !viper.IsSet("nats_url") {
		log.Fatal("please provide nats_url config param")
	}

	if !viper.IsSet("nats_user_name") {
		log.Fatal("please provide nats_user_name config param")
	}

	if !viper.IsSet("nats_user_password") {
		log.Fatal("please provide nats_user_password config param")
	}

	if !viper.IsSet("queue_group") {
		log.Fatal("please provide queue_group config param")
	}

}

func DegaToGoogleChat() bool {
	return viper.IsSet("dega_to_google_chat") && viper.GetBool("dega_to_google_chat")
}
