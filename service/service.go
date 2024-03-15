package service

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Service struct {
	done chan struct{}
	wg   *sync.WaitGroup
}

func New() *Service {
	s := &Service{
		done: make(chan struct{}),
		wg:   &sync.WaitGroup{},
	}
	s.wg.Add(1)
	return s
}

func (s *Service) Serve(laddr *net.TCPAddr) {
	defer s.wg.Done()

	fmt.Println("start serving...")

	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", listener.Addr())

	for {
		select {
		case <-s.done:
			log.Println("stopping listening on", listener.Addr())
			listener.Close()
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Service) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	defer s.wg.Done()

	for {
		select {
		case <-s.done:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}
	}
}

func (s *Service) Stop() {
	log.Println("stopping service and finishing work...")
	close(s.done)
	s.wg.Wait()
}
