package tx

import (
    //"fmt"
	"time"
	//"bytes"
	"github.com/toolkits/container/list"
	"github.com/toolkits/concurrent/semaphore"
	."github.com/thewayma/suricata_transporter/g"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

//!< 半同步发送任务, 消费半异步队列
func startSendTasks() {
	cfg := Config()

	// init semaphore
	judgeConcurrent := cfg.Judge.MaxConns
	graphConcurrent := cfg.Graph.MaxConns
	tsdbConcurrent  := cfg.Tsdb.MaxConns

    if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}

	if graphConcurrent < 1 {
		graphConcurrent = 1
	}

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Judge.Cluster {
		queue := JudgeQueues[node]
		go forward2JudgeTask(queue, node, judgeConcurrent)
	}
/*
	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			queue := GraphQueues[node+addr]
			go forward2GraphTask(queue, node, addr, graphConcurrent)
		}
	}

	if cfg.Tsdb.Enabled {
		go forward2TsdbTask(tsdbConcurrent)
	}
    */
}

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeListLimited, node string, concurrent int) {
	batch := Config().Judge.Batch // 一次发送,最多batch条数据
	addr := Config().Judge.Cluster[node]
	sema := semaphore.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems := make([]*JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*JudgeItem)
            //fmt.Printf("%s\n", judgeItems[i])
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, judgeItems []*JudgeItem, count int) {
			defer sema.Release()

			resp := &SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				Log.Error("forward2JudgeTask send judge %s:%s fail: %v", node, addr, err)
				//proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				//proc.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

/*
// Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *list.SafeListLimited, node string, addr string, concurrent int) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	sema := semaphore.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*GraphItem)
		}

		sema.Acquire()
		go func(addr string, graphItems []*GraphItem, count int) {
			defer sema.Release()

			resp := &SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Printf("send to graph %s:%s fail: %v", node, addr, err)
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToGraphCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}

// Tsdb定时任务, 将数据通过api发送到tsdb
func forward2TsdbTask(concurrent int) {
	batch := g.Config().Tsdb.Batch // 一次发送,最多batch条数据
	retry := g.Config().Tsdb.MaxRetry
	sema := semaphore.NewSemaphore(concurrent)

	for {
		items := TsdbQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()

			var tsdbBuffer bytes.Buffer
			for i := 0; i < len(itemList); i++ {
				tsdbItem := itemList[i].(*TsdbItem)
				tsdbBuffer.WriteString(tsdbItem.TsdbString())
				tsdbBuffer.WriteString("\n")
			}

			var err error
			for i := 0; i < retry; i++ {
				err = TsdbConnPoolHelper.Send(tsdbBuffer.Bytes())
				if err == nil {
					proc.SendToTsdbCnt.IncrBy(int64(len(itemList)))
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				proc.SendToTsdbFailCnt.IncrBy(int64(len(itemList)))
				log.Println(err)
				return
			}
		}(items)
	}
}
*/
