package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/thanhquy1105/vida/server"
)

func main() {
	var (
		// dataDir is the path to the data directory to save key-value data
		dataDir = flag.String("data", "./data", "path to data directory")
		// host and port is the IP and port of this program to listen to consumers connecting to
		host = flag.String("host", "0.0.0.0", "ip to listen")
		port = flag.String("port", "22133", "port to listen")
	)
	flag.Parse()

	log.Println("data directory: ", *dataDir)
	fmt.Println("The number of CPUs that can execute simultaneously: ", runtime.GOMAXPROCS(runtime.NumCPU()))

	server := server.New(*dataDir)

	laddr, err := net.ResolveTCPAddr("tcp", *host+":"+*port)
	if err != nil {
		log.Fatalln(err)
	}

	go server.Serve(laddr)

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-done)

	server.Stop()
}
