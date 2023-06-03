package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Addr     string
	Port     int
	Listener net.Listener
	mu       sync.Mutex
	Clients  map[net.Conn]*client
}

type client struct {
	name   string
	change bool
}

const (
	addr     = "localhost"
	def_port = 8989
)

func NewServer(addr string, port int) (*Server, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(addr), Port: port})
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr:     addr,
		Port:     port,
		Listener: listener,
		Clients:  make(map[net.Conn]*client),
	}, nil
}

func (s *Server) Start() error {
	messageHistory := []string{}
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer s.mu.Unlock()
			defer conn.Close()
			s.mu.Lock()
			s.Clients[conn] = &client{}
			if len(s.Clients) > 10 {
				writeToClient(conn, "Chat is full! Comeback later.")
				delete(s.Clients, conn)
				return
			}
			s.mu.Unlock()
			welcomeGuest(conn)
			buf := make([]byte, 1024)
			current_time := time.Now().Format("01-02-2006 15:04:05")
			for {
				left := false
				writeToClient(conn, "[ENTER YOUR NAME]: ")
				name, err := readFromClient(conn, &buf)
				if err != nil {
					break
				}
				s.mu.Lock()
				if !s.checkName(conn, name) {
					s.mu.Unlock()
					continue
				}
				s.mu.Unlock()
				if !s.Clients[conn].change {
					s.mu.Lock()
					for _, msg := range messageHistory {
						writeToClient(conn, msg)
					}
					s.mu.Unlock()
					s.Broadcast(conn, &messageHistory, name+" has joined the chat!")
				} else {
					s.Broadcast(conn, &messageHistory, s.Clients[conn].name+" has changed name to "+"-> "+name)
				}
				s.mu.Lock()
				s.Clients[conn] = &client{name: name, change: false}
				s.mu.Unlock()
				prefixpasta := fmt.Sprintf("[%s][%s]: ", current_time, name)
				for {
					writeToClient(conn, prefixpasta)
					msg, err := readFromClient(conn, &buf)
					if err != nil {
						left = true
						s.Broadcast(conn, &messageHistory, name+" has left the chat...")
						break
					}
					if msg == "/changeName" {
						s.mu.Lock()
						s.Clients[conn].change = true
						s.mu.Unlock()
						break
					} else if msg == "" {
						continue
					}
					msg = fmt.Sprintf("[%s][%s]: %s", current_time, name, msg)
					s.Broadcast(conn, &messageHistory, msg)
				}
				if left {
					break
				}
			}
			s.mu.Lock()
			delete(s.Clients, conn)
			conn.Close()
		}(conn)
	}
}

func (s *Server) Stop() error {
	s.Listener.Close()
	for client := range s.Clients {
		client.Close()
	}
	return nil
}

func (s *Server) Broadcast(sender net.Conn, msgHistory *[]string, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	*msgHistory = append(*msgHistory, message+"\n")
	current_time := time.Now().Format("01-02-2006 15:04:05")
	save_msg := fmt.Sprintf("{%v} -- {%s}\n", s.Port, message)
	logCatcher(save_msg)
	for conn := range s.Clients {
		if conn == sender || s.Clients[conn].name == "" {
			continue
		}
		prefixpasta := fmt.Sprintf("[%s][%s]: ", current_time, s.Clients[conn].name)
		// writeToClient(conn, "\r")
		writeToClient(conn, "\n"+message+"\n"+prefixpasta)
	}
}

func writeToClient(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error sending response to %s: %v\n", conn.RemoteAddr(), err)
		return
	}
}

func readFromClient(conn net.Conn, buf *[]byte) (string, error) {
	n, err := conn.Read(*buf)
	if err != nil {
		return "", err
	}
	msg := strings.TrimSpace(string((*buf)[:n-1]))
	res := ""
	for _, v := range msg {
		if v == 27 {
			res += "^["
			continue
		}
		res += string(v)

	}
	return res, nil
}

func main() {
	args := os.Args[1:]
	ports := []int{}
	if len(args) == 0 {
		server, err := NewServer(addr, def_port)
		if err != nil {
			fmt.Println(err)
			return
		}
		go server.Start()
		fmt.Println("Listening on the port :", def_port)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		server.Stop()
		return
	}
	for _, arg := range args {
		port, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Invalid argument: ", arg)
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		ports = append(ports, port)
	}
	servers := []*Server{}
	for _, port := range ports {
		server, err := NewServer(addr, port)
		if err != nil {
			fmt.Println(err)
			return
		}
		servers = append(servers, server)
	}
	fmt.Print("Listening on the ports ")
	for _, server := range servers {
		fmt.Print(":", server.Port, " ")
		go server.Start()
	}
	fmt.Println()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	for _, server := range servers {
		server.Stop()
	}
}
