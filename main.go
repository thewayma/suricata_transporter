package main

import (
	"flag"
	"github.com/thewayma/suricata_transporter/g"
	"github.com/thewayma/suricata_transporter/rx"
	"github.com/thewayma/suricata_transporter/tx"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)
    g.InitLog()

	tx.Start()
	go rx.StartRpc()

	select {}
}
