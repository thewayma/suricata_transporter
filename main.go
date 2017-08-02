package main

import (
	"flag"
	"github.com/thewayma/suricata_transporter/g"
	"github.com/thewayma/suricata_transporter/receiver"
    /*
	"github.com/thewayma/suricata_transporter/sender"
    */
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)
    g.InitLog()

	//sender.Start()
	go receiver.StartRpc()

	select {}
}
