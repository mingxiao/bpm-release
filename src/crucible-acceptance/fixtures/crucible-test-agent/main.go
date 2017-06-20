package main

import (
	"crucible-acceptance/fixtures/crucible-test-agent/handlers"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var port = flag.Int("port", -1, "port the server listens on")

func main() {
	flag.Parse()
	if *port == -1 {
		log.Fatal("no explicit port specified")
	}

	crucibleVar := os.Getenv("CRUCIBLE")
	if crucibleVar == "" {
		log.Fatal("Crucible environment variable not set")
	}

	fmt.Println("Test Agent Started - STDOUT")
	log.Println("Test Agent Started - STDERR")

	http.HandleFunc("/", handlers.Hello)
	http.HandleFunc("/hostname", handlers.Hostname)
	http.HandleFunc("/mounts", handlers.Mounts)
	http.HandleFunc("/processes", handlers.Processes)
	http.HandleFunc("/var-vcap", handlers.VarVcap)
	http.HandleFunc("/var-vcap-data", handlers.VarVcapData)
	http.HandleFunc("/var-vcap-jobs", handlers.VarVcapJobs)
	http.HandleFunc("/whoami", handlers.Whoami)

	errChan := make(chan error)
	signals := make(chan os.Signal)

	signal.Notify(signals)

	go func() {
		errChan <- http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			log.Fatal(err)
		}
	case sig := <-signals:
		log.Fatalf("Signalled: %#v", sig)
	}
}
