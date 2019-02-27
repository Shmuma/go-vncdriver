package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/openai/go-vncdriver/gymvnc"
)
import _ "net/http/pprof"

type foo struct {
	bar string
}

func main() {
	f, err := os.Create("/tmp/profile-00.pprof")
    pprof.Lookup("heap").WriteTo(f, 1)
    f.Close()

	gymvnc.ConfigureLogging()

	batch := gymvnc.NewVNCBatch()
	err = batch.Open("conn", gymvnc.VNCSessionConfig{
		Address:          "127.0.0.1:5900",
//		Address:          "3.public-devbox.sci.openai-tech.com:20000",
		Password:         "openai",
		Encoding:         "tight",
		FineQualityLevel: 100,
	})
	if err != nil {
		panic(err)
	}

	f, err = os.Create("/tmp/profile-01.pprof")
    pprof.Lookup("heap").WriteTo(f, 1)
    f.Close()

	start := time.Now()
	updates := 0
	errs := 0
	for i := 0; i < 2000; i++ {
		elapsed := time.Now().Sub(start)
		if elapsed >= time.Duration(1)*time.Second {
			delta := float64(elapsed / time.Second)
			log.Printf("Received: %d, updates=%.2f errs=%.2f", i, float64(updates)/delta, float64(errs)/delta)

			start = time.Now()
			updates = 0
			errs = 0
		}

		batchEvents := map[string][]gymvnc.VNCEvent{
			"conn": []gymvnc.VNCEvent{},
		}
		_, updatesN, errN := batch.Step(batchEvents)
		if errN["conn"] != nil {
			log.Fatalf("error: %+v", errN["conn"])
		}

		updates += len(updatesN["conn"])
		time.Sleep(16 * time.Millisecond)
	}
	f, err = os.Create("/tmp/profile-02.pprof")
    pprof.Lookup("heap").WriteTo(f, 1)
    f.Close()

    batch.Close("conn")

    f, err = os.Create("/tmp/profile-03.pprof")
    pprof.Lookup("heap").WriteTo(f, 1)
    f.Close()


	// f, err := os.Create("/tmp/hi.prof")
	// if err != nil {
	//     log.Fatal("could not create memory profile: ", err)
	// }
	// runtime.GC() // get up-to-date statistics
	// if err := pprof.WriteHeapProfile(f); err != nil {
	//     log.Fatal("could not write memory profile: ", err)
	// }
	// f.Close()
}
