package receiver

import (
	"time"
	"log"
	"net"
	"strconv"
	"net/rpc"
	"net/rpc/jsonrpc"

    "github.com/thewayma/suricata_transporter/g"
	//"github.com/thewayma/suricata_transporter/sender"
)

type Transfer int

type MetaData struct {
    Metric      string
    Endpoint    string
    Timestamp   int64
    Step        int64
    Value       float64
    CounterType string
    //Tags        map[string]string
}

type MetricValue struct {
    Metric    string
    Endpoint  string
    Timestamp int64
    Step      int64
    Value     interface{}
    Type      string
    //Tags      string
}

type TransferResp struct {
    Message string
    Total   int
    Invalid int
    Latency int64
}

func StartRpc() {
	if !g.Config().Rpc.Enabled {
		return
	}

	addr := g.Config().Rpc.Listen
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

func (t *Transfer) Update(args []*MetricValue, reply *TransferResp) error {
	return RecvMetric(args, reply, "rpc")
}

func RecvMetric(args []*MetricValue, reply *TransferResp, from string) error {
	start := time.Now()
	reply.Invalid = 0

	items := []*MetaData{}
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		if v.Metric == "" || v.Endpoint == "" {
			reply.Invalid += 1
			continue
		}

		if v.Value == "" {
			reply.Invalid += 1
			continue
		}

        fv := &MetaData{
            Metric:      v.Metric,
            Endpoint:    v.Endpoint,
            Timestamp:   v.Timestamp,
            Step:        v.Step,
            CounterType: v.Type,
            //Tags:        v.Tags,
        }

		valid := true
		var vv float64
		var err error

		switch cv := v.Value.(type) {
		case string:
			vv, err = strconv.ParseFloat(cv, 64)
			if err != nil {
				valid = false
			}
		case float64:
			vv = cv
		case int64:
			vv = float64(cv)
		default:
			valid = false
		}

		if !valid {
			reply.Invalid += 1
			continue
		}

		fv.Value = vv
		items = append(items, fv)
	}

/*
	cfg := g.Config()

	if cfg.Judge.Enabled {
		sender.Push2JudgeSendQueue(items)
	}

	if cfg.Graph.Enabled {
		sender.Push2GraphSendQueue(items)
	}

	if cfg.Tsdb.Enabled {
		sender.Push2TsdbSendQueue(items)
	}
    */

	reply.Message = "ok"
	reply.Total   = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
