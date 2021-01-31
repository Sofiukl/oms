package main

import (
	"runtime"

	"github.com/sofiukl/oms/oms-checkout/core"
	"github.com/sofiukl/oms/oms-checkout/dispatcher"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := core.App{}
	app.Initialize()
	dispatcher.StartDispatcher(4)
	app.Run(":" + app.Config.ServerPort)
}
