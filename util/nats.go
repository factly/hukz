package util

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

// NC nats connection object
var NC *nats.EncodedConn

// ConnectNats connect to nats server
var ConnectNats = func() {
	var err error
	nc, err := nats.Connect(viper.GetString("nats_url"), nats.UserInfo(viper.GetString("nats_user_name"), viper.GetString("nats_user_password")))
	if err != nil {
		log.Fatal(err)
	}
	NC, _ = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
}
