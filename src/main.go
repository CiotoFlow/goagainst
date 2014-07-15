package main

import (
	"fmt"
	"flag"
	"./irc"
)

func main() {
	flag.Parse ()

	config, err := LoadConfig(flag.Arg(0))
	if err != nil {
		fmt.Println (err)
		return
	}

	quitChans := make([]chan bool, len(config.Servers));

	i := 0
	for _, server := range config.Servers {
		client := irc.NewIRC(server.Nickname, server.Address,
							 server.Channel, server.UseTls)

		err := client.Connect()
		if (err != nil) {
			fmt.Println(err)
			continue
		}

		go func(idx int) {
			client.Loop()
			if (err != nil) {
				fmt.Println(err)
				client.Disconnect()
			}
			quitChans[idx] <- true
		}(i)

		i++
	}

	for i = 0; i < len(config.Servers); i++ {
		<- quitChans[i]
	}
}
