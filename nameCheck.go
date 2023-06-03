package main

import (
	"net"
	"strings"
)

func (s *Server) checkName(conn net.Conn, name string) bool {
	if len(name) < 2 {
		conn.Write([]byte("Name should contain more that 2 characters!\n"))
		return false
	}
	for _, r := range name {
		if (r < 48 || r > 122) || (r > 90 && r < 97) {
			conn.Write([]byte("Name shouldn't contain symbols!\n"))
			return false
		}
	}
	for _, client := range s.Clients {
		if client.name == strings.ToLower(name) {
			conn.Write([]byte("Client with such username already exists!\n"))
			return false
		}
	}
	return true
}
