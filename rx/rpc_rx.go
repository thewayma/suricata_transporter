package rx

import (
	"time"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

    ."github.com/thewayma/suricata_transporter/g"
	"github.com/thewayma/suricata_transporter/tx"
)

type Transfer struct{}

func StartRpc() {
	if !Config().Rpc.Enabled {
		return
	}

	addr := Config().Rpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
	}

	server := rpc.NewServer()
	server.Register(new(Transfer))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (t *Transfer) Update(args []*MetricData, reply *TransferResp) error {
	return RecvMetric(args, reply, "rpc")
}

func RecvMetric(items []*MetricData, reply *TransferResp, from string) error {
	start := time.Now()
	reply.Invalid = 0

    //!< sanity check已前移至agent上
	cfg := Config()
	if cfg.Judge.Enabled {
		tx.Push2JudgeSendQueue(items)
	}

    /*
	if cfg.Graph.Enabled {
		tx.Push2GraphSendQueue(items)
	}

	if cfg.Tsdb.Enabled {
		tx.Push2TsdbSendQueue(items)
	}
    */

	reply.Message = "ok"
	reply.Total   = len(items)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
