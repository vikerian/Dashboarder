package main

import (
	"fmt"
	"log/slog"
	"os"
)

/* using logging as singleton */
var sl *slog.Logger

/* pre run initialization */
func init() {
	sl = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

/* main function */
func main() {
	sl.Info(fmt.Sprintf("Service Dashboarder starting up from bin: %s...", os.Args[0]))
	defer sl.Info("Service Dashboarder quitting on command...")

}
