package main

import (
	"chat/common"
	"fmt"
	"net"
)

const PRIVATE_CMD = 0
const PUBLIC_CMD = 1

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	Server *Server
	Status uint8 //0: oneline 1: offline
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   common.GetName(),
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		Server: server,
		Status: 0,
	}
	go user.ListenMessage()
	return user
}
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg))
	}
}
func (u *User) Online() {
	s := u.Server
	s.mapLock.Lock()
	if _, ok := s.OnLineMap[u.Addr]; ok {
		u.Name = s.OnLineMap[u.Addr].Name
	}
	s.OnLineMap[u.Addr] = u
	s.mapLock.Unlock()
	s.BroadCast(u, "Online")
}

func (u *User) Offline() {
	s := u.Server
	u.Status = 1
	s.BroadCast(u, "Offline")
}
func (u *User) SendToSelf(msg string) {
	u.conn.Write([]byte(msg + "\n> "))
}
func (u *User) CommandMode(command string, pattern int) {
	s := u.Server
	var msg string
	switch command {
	case "list_all_users":
		s.mapLock.Lock()
		statusMap := map[uint8]string{0: "online", 1: "offline"}
		for _, user := range s.OnLineMap {
			msg = fmt.Sprintf("[%s]%s %s\n%s", user.Addr, user.Name, statusMap[user.Status], msg)
		}
		s.mapLock.Unlock()
	case "exit":
		u.Offline()
		u.conn.Close()
	}
	switch pattern {
	case PUBLIC_CMD:
		s.BroadCast(u, msg)
	case PRIVATE_CMD:
		u.SendToSelf(msg)
	}
}
func (u *User) HandleMessage(msg string) {
	s := u.Server
	if len(msg) < 1 {
		return
	}
	if msg[0] == '@' {
		mode := PUBLIC_CMD
		u.CommandMode(msg[1:], mode)
		return
	}
	if msg[0] == '/' {
		mode := PRIVATE_CMD
		u.CommandMode(msg[1:], mode)
		return
	}
	s.BroadCast(u, msg)
}
