package cmd

import (
	"flag"
	"github.com/rdsalakhov/game-keys-store/internal/server"
	"github.com/spf13/viper"
	"log"
)

func main() {
	flag.Parse()
	config := &server.Config{}

	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}

	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
