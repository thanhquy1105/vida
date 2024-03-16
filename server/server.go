package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/thanhquy1105/vida/repository"
)

// Server represents a tcp message queue server
type Server struct {
	dataDir string
	repo    *repository.QueueRepository
	done    chan struct{}
	wg      *sync.WaitGroup
}

// New creates a new tcp message queue server
func New(dataDir string) *Server {
	s := &Server{
		dataDir: dataDir,
		repo:    &repository.QueueRepository{},
		done:    make(chan struct{}),
		wg:      &sync.WaitGroup{},
	}
	s.wg.Add(1)
	return s
}

// Serve starts the tcp message queue server
func (s *Server) Serve(laddr *net.TCPAddr) {
	defer s.wg.Done()

	fmt.Println("start serving...")
	var err error
	s.repo, err = repository.NewRepository(s.dataDir)
	if err != nil {
		log.Fatal(err)
	}

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

// handleConnection is the concurrent function to handle each comsumer
func (s *Server) handleConnection(conn *net.TCPConn) {
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

// Stop calls graceful shutdown for server
func (s *Server) Stop() {
	log.Println("stopping Server and finishing work...")
	close(s.done)
	s.wg.Wait()
}
