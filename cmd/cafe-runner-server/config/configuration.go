package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server ServerConfiguration
	Data   DatabaseConfiguration
}

func GetConfiguration() *Configuration {
	viper.SetConfigName(`config`)
	viper.SetConfigType(`yml`)
	viper.AddConfigPath(`.`)
	viper.SetDefault(`data.path`, `.`)
	viper.SetDefault(`server.port`, `3133`)
	viper.SetDefault(`server.hostname`, `localhost`)                // without protocol and port. like "cafe-server.organization.org"
	viper.SetDefault(`server.externalurl`, `http://localhost:3133`) // external web url. Port and protocols can be redefined for work with reverse-proxy (or in docker).
	viper.SetDefault(`server.cafe.lowport`, `21000`)
	viper.SetDefault(`server.cafe.highport`, `21010`)

	viper.AutomaticEnv()

	var res Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		log.Printf("Trying to create default config")
		if err := viper.WriteConfigAs(`config.yml`); err != nil {
			log.Fatalf("Error writing default config file, %s", err)
		}

	}
	err := viper.Unmarshal(&res)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &res
}

func (t *Configuration) PrintToLog() {
	log.Printf(`data.path=		%s`, t.Data.Path)
	log.Printf(`server.hostname=	%s`, t.Server.Hostname)
	log.Printf(`server.port=	%d`, t.Server.Port)
	log.Printf(`server.externalurl=	%s`, t.Server.ExternalURL)
	log.Printf(`server.cafe.lowport=	%d`, t.Server.Cafe.LowPort)
	log.Printf(`server.cafe.highport=	%d`, t.Server.Cafe.HighPort)
}
