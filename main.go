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

	"github.com/thanhquy1105/vida/service"
)

var (
	hostAndPort      = flag.String("listen", "0.0.0.0:22133", "ip and port to listen")
	interruptSignals = []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}
)

func main() {
	flag.Parse()
	fmt.Println("The number of CPUs that can execute simultaneously: ", runtime.GOMAXPROCS(runtime.NumCPU()))

	service := service.New()

	laddr, err := net.ResolveTCPAddr("tcp", *hostAndPort)
	if err != nil {
		log.Fatalln(err)
	}

	go service.Serve(laddr)

	done := make(chan os.Signal)
	signal.Notify(done, interruptSignals...)
	log.Println(<-done)

	service.Stop()
}
