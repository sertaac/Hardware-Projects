/**************************************/
/*                                    */
/*      IPC Server - Go Backend       */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
)

/**************************************************/
/*                                                */
/*            MESSAGE TYPE CONSTANTS              */
/*                                                */
/**************************************************/

const (
	MsgTypeListGames      = "list_games"
	MsgTypeGetGame        = "get_game"
	MsgTypeLaunchGame     = "launch_game"
	MsgTypeGetCategories  = "get_categories"
	MsgTypeGetPlatforms   = "get_platforms"
	MsgTypeGetFavorites   = "get_favorites"
	MsgTypeToggleFavorite = "toggle_favorite"
	MsgTypeGetRecent      = "get_recent"
	MsgTypeScan           = "scan"
	MsgTypeAddScanPath    = "add_scan_path"
	MsgTypeStatus         = "status"
	MsgTypeError          = "error"
	MsgTypeSuccess        = "success"
)

/**************************************************/
/*                                                */
/*           REQUEST / RESPONSE STRUCTS           */
/*                                                */
/**************************************************/

type Request struct {
	Type    string          `json:"type"`
	ID      string          `json:"id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Response struct {
	Type    string      `json:"type"`
	ID      string      `json:"id,omitempty"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type GameListPayload struct {
	Platform string `json:"platform,omitempty"`
	Category string `json:"category,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type ScanPathPayload struct {
	Path string `json:"path"`
}

/**************************************************/
/*                                                */
/*              IPC SERVER STRUCT                 */
/*                                                */
/**************************************************/

type IPCServer struct {
	listener net.Listener
	clients  map[net.Conn]bool
	mu       sync.RWMutex
	port     int
	running  bool
	handler  func(req Request) Response
}

/**************************************************/
/*                                                */
/*            SERVER CONSTRUCTOR                  */
/*                                                */
/**************************************************/

func NewIPCServer(port int) *IPCServer {
	return &IPCServer{
		clients: make(map[net.Conn]bool),
		port:    port,
		running: false,
	}
}

func (s *IPCServer) SetHandler(handler func(req Request) Response) {
	s.handler = handler
}

/**************************************************/
/*                                                */
/*             START / STOP SERVER                */
/*                                                */
/**************************************************/

func (s *IPCServer) Start() error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.listener = listener
	s.running = true
	fmt.Printf("IPC Server started on %s\n", addr)
	go s.acceptConnections()
	return nil
}

func (s *IPCServer) Stop() {
	s.running = false
	s.mu.Lock()
	for conn := range s.clients {
		conn.Close()
	}
	s.clients = make(map[net.Conn]bool)
	s.mu.Unlock()

	if s.listener != nil {
		s.listener.Close()
	}
	fmt.Println("IPC Server stopped")
}

/**************************************************/
/*                                                */
/*          CONNECTION HANDLING                   */
/*                                                */
/**************************************************/

func (s *IPCServer) acceptConnections() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		s.mu.Lock()
		s.clients[conn] = true
		s.mu.Unlock()

		fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
		go s.handleClient(conn)
	}
}

func (s *IPCServer) handleClient(conn net.Conn) {
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
	}()

	reader := bufio.NewReader(conn)

	for s.running {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(conn, "", fmt.Sprintf("Invalid JSON: %v", err))
			continue
		}

		var resp Response
		if s.handler != nil {
			resp = s.handler(req)
		} else {
			resp = s.defaultHandler(req)
		}

		s.sendResponse(conn, resp)
	}
}

func (s *IPCServer) defaultHandler(req Request) Response {
	return Response{
		Type:    MsgTypeStatus,
		ID:      req.ID,
		Success: true,
		Data:    map[string]string{"status": "ready"},
	}
}

func (s *IPCServer) sendResponse(conn net.Conn, resp Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = conn.Write(append(data, '\n'))
	return err
}

func (s *IPCServer) sendError(conn net.Conn, reqID, message string) {
	resp := Response{Type: MsgTypeError, ID: reqID, Success: false, Error: message}
	s.sendResponse(conn, resp)
}

/**************************************************/
/*                                                */
/*             SERVER UTILITIES                   */
/*                                                */
/**************************************************/

func (s *IPCServer) GetPort() int {
	return s.port
}

func (s *IPCServer) IsRunning() bool {
	return s.running
}

func (s *IPCServer) ClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
