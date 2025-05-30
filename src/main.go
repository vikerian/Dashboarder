package main

import (
	"fmt"
	"os"
)

/* using logging from logsingleton */

func main() {
	sl.Info(fmt.Sprintf("Service Dashboarder starting up from bin: %s...", os.Args[0]))
	defer sl.Info("Service Dashboarder quitting on command...")

}
