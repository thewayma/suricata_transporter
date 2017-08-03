package g

import (
    "os"
    "log"
    "fmt"
    "runtime"
)

const (
    GAUGE        = "GAUGE"
    COUNTER      = "COUNTER"
    DERIVE       = "DERIVE"
    DEFAULT_STEP = 60
    MIN_STEP     = 30
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
