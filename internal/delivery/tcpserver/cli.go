package tcpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/Shayan-Ghani/GOAD/config"
	tcphandler "github.com/Shayan-Ghani/GOAD/internal/delivery/tcpserver/handler"
	"github.com/Shayan-Ghani/GOAD/internal/repository"
)

type Server struct {
	repo     repository.Repository
	listener net.Listener
	handler  *tcphandler.RequestHandler
}

func NewServer(repo repository.Repository) (*Server, error) {
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &Server{
		repo:     repo,
		listener: listener,
		handler:  tcphandler.NewRequestHandler(repo),
	}, nil
}

func (s *Server) Start() error {
	defer s.listener.Close()

	fmt.Printf("Server listening on %s ... \n", s.listener.Addr().String())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// defer conn.Close()

	// data := make([]byte, 1024)
	// byteCount, err := conn.Read(data)
	// if err != nil {
	// 	log.Printf("Error reading from connection: %v\n", err)
	// 	return
	// }


	var request tcphandler.CliRequest

	fmt.Printf("request: %v\n", request)

	if err := json.NewDecoder(conn).Decode(&request); err != nil {
		response := tcphandler.CliResponse{
			Err: fmt.Sprintf("invalid request format: %v", err),
		}
		json.NewEncoder(conn).Encode(response)
		return
	}

	fmt.Printf("request: %v\n", request)

	response := s.handler.HandleRequest(request)

	fmt.Printf("response: %v\n", response)

	if err := json.NewEncoder(conn).Encode(response); err != nil {
		log.Printf("Error sending response: %v\n", err)
	}
}

// func main() {

// 	l, err := net.Listen("tcp", config.Addr)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer l.Close()

// 	fmt.Printf("Server listening on %s ... \n", l.Addr().String())

// 	cr := &CliRequest{}

// 	for {
// 		c, err := l.Accept()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		data := make([]byte, 1024)

// 		byteCount, err := c.Read(data)
// 		if err != nil {
// 			fmt.Println(err)

// 			continue
// 		}

// 		fmt.Println(string(data)+ "\n\n")

// 		if err:= json.Unmarshal(data[:byteCount], cr); err != nil {
// 			fmt.Println(err)

// 			continue
// 		}

// 	}

// }

func R() {
	// 	handleItems := func(items interface{}, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		response.Respond(icmd.flags.Format, items)
	// 		return nil
	// 	}

	// 	// handle update
	// 	updates := make(map[string]interface{}, 3)

	// 	flagUpdates := map[string]interface{}{
	// 		"name":        icmd.flags.Name,
	// 		"description": icmd.flags.Description,
	// 	}

	// 	for key, value := range flagUpdates {
	// 		updates[key] = value
	// 	}

}
