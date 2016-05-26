package engine

import (
	"service"
	"time"
)

func ExampleEngineType_Pause() {

	enginesrv := NewEngine()
	service.StartService(enginesrv, "engine", nil)

	go func() {
		tickChan := time.NewTicker(time.Second * 10).C
		var trick bool = true
		for {
			select {
			case <-tickChan:
				if trick {
					//logger.Debug("engine pause test")
					enginesrv.Pause()
					trick = false
				} else {
					//logger.Debug("engine resume test")
					enginesrv.Resume()
					trick = true
				}
			}
		}
	}()

	//t.Logf("test sigint...")
	//t.Error("sigint test failed")

}
