package main

import (
	"fmt"
	"flag"
	"./irc"
)

func main() {
	flag.Parse ()

	configFile := "config.json";
	if flag.NArg() > 0 {
		configFile = flag.Arg(0)
	}
	
	config, err := irc.LoadConfig(configFile)
	if err != nil {
		fmt.Println (err)
		return
	}

	quitChans := make([]chan bool, len(config.Servers));

	i := 0
	for _, server := range config.Servers {
		client := irc.NewIRC(server)

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

	for i, _ := range config.Servers {
		<- quitChans[i]
	}
}
