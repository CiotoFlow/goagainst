package irc

import (
	"encoding/json"
	"os"
)

type ServerConfig struct {
	Address string
	Nickname string
	Realname string
	AutoJoin []string
	AutoPong bool
	UseTls bool
}

type Config struct {
	Servers []ServerConfig
}

func LoadConfig(filename string) (config Config, err error) {
	config = Config {}
	
	file, err := os.Open (filename)
	if err != nil { return }
	
	dec := json.NewDecoder (file)
	err = dec.Decode(&config)
	if err != nil { return }

	return
}
