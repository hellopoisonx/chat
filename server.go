package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const TTL = 60 // time to live(s)
type Server struct {
	Ip        string
	Port      int
	OnLineMap map[string]*User
	Message   chan string
	mapLock   sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: map[string]*User{},
		Message:   make(chan string),
	}
}
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s %s\n> ", user.Addr, user.Name, msg)
	s.Message <- sendMsg
}
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, v := range s.OnLineMap {
			if v.Status == 1 {
				continue
			}
			v.C <- msg
		}
		s.mapLock.Unlock()
	}
}
func (s *Server) ReceiveFromUser(user *User, conn net.Conn, isAlive *chan bool) {
	buf := [4096]byte{}
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			user.Offline()
			return
		}
		if err != io.EOF && err != nil {
			log.Printf("Error reading from tcp connection([%s]%s)", user.Addr, user.Name)
			return
		}
		user.HandleMessage(string(buf[:n-1]))
		*isAlive <- true
	}
}
func (s *Server) Handle(conn net.Conn) {
	isAlive := make(chan bool)
	log.Println("Connected")
	user := NewUser(conn, s)
	user.Online()
	go s.ReceiveFromUser(user, conn, &isAlive)
	for {
		timer := time.After(time.Second * TTL)
		select {
		case <-isAlive:
			continue
		case <-timer:
			user.SendToSelf("EOL")
			user.Offline()
		}
	}
}
func startServer(s *Server) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Fatalf("Error listening %v", err)
	}
	defer listener.Close()
	//Main Loop
	go s.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting request, retrying.... %v", err)
			continue
		}
		go s.Handle(conn)
	}
}

func main() {
	server := NewServer("0.0.0.0", 9991)
	startServer(server)
}
