package g

import (
    "fmt"
    "sort"
    "strings"
)

/*
type MetricValue struct {	//!< agent data
    Endpoint  string      `json:"endpoint"`
    Metric    string      `json:"metric"`
    Value     interface{} `json:"value"`
    Step      int64       `json:"step"`
    Type      string      `json:"counterType"`
    Tags      map[string]string `json:"tags"`
    Timestamp int64       `json:"timestamp"`
}
*/

type MetricData struct {      //!< 统一agent,transporter data, 减小内存拷贝
    Endpoint    string				`json:"endpoint"`
    Metric      string				`json:"metric"`
    Value       float64				`json:"value"`
    Step        int64				`json:"step"`
    Type        string				`json:"Type"`
    Tags        map[string]string	`json:"tags"`
    Timestamp   int64				`json:"timestamp"`
}

func (t *MetricData) String() string {
    return fmt.Sprintf("<Endpoint:%s, Metric:%s, Timestamp:%d, Step:%d, Value:%f, Tags:%v, Type:%s>",
        t.Endpoint, t.Metric, t.Timestamp, t.Step, t.Value, t.Tags, t.Type)
}


//!< tags按照key排序, 返回用","拼接而成的tags字符串
func SortedTags(tags map[string]string) string {
    if tags == nil {
        return ""
    }

    size := len(tags)

    if size == 0 {
        return ""
    }

    if size == 1 {
        for k, v := range tags {
            return fmt.Sprintf("%s=%s", k, v)
        }
    }

    keys := make([]string, size)
    i := 0
    for k := range tags {
        keys[i] = k
        i++
    }

    sort.Strings(keys)

    ret := make([]string, size)
    for j, key := range keys {
        ret[j] = fmt.Sprintf("%s=%s", key, tags[key])
    }

    return strings.Join(ret, ",")
}

func (r *MetricData) PK() string {
    if r.Tags == nil || len(r.Tags) == 0 {
        fmt.Sprintf("%s/%s", r.Endpoint, r.Metric)
    }
    return fmt.Sprintf("%s/%s/%s", r.Endpoint, r.Metric, SortedTags(r.Tags))
}

type TransferResp struct {
    Message     string
    Total       int
    Invalid     int
    Latency     int64
}

type JudgeItem struct {
    Endpoint    string
    Metric      string
    Value       float64
    Timestamp   int64
    JudgeType   string
    Tags        map[string]string
}

type GraphItem struct {
    Endpoint	string
    Metric		string
    Tags		map[string]string
    Value		float64
    Timestamp	int64
    DsType		string
    Step		int
    Heartbeat	int
    Min			string
    Max			string
}

type TsdbItem struct {
    Metric      string
    Value       float64
    Timestamp   int64
    Tags        map[string]string
}
