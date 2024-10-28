package blockchainlite

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	blockchain *Blockchain
	httpServer *http.Server
}

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func NewServer(bcName string) (*Server, error) {
	bc, err := NewBlockchain(bcName)
	if err != nil {
		log.Printf("[ERROR] Error creating blockchain: %v", err)
		return nil, err
	}
	log.Println("[INFO] Blockchain created successfully")
	return &Server{
		blockchain: bc,
		httpServer: &http.Server{},
	}, nil
}

func (s *Server) writeResponse(w http.ResponseWriter, code int, data interface{}, errMsg string) {
	response := Response{
		Code:  code,
		Data:  data,
		Error: errMsg,
	}
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] Error encoding response: %v", err)
	}
}

func (s *Server) AddBlockHandler(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		s.writeResponse(w, http.StatusBadRequest, nil, err.Error())
		log.Printf("[ERROR] Error decoding request body: %v", err)
		return
	}

	if err := s.blockchain.AddBlock(data); err != nil {
		s.writeResponse(w, http.StatusInternalServerError, nil, err.Error())
		log.Printf("[ERROR] Error adding block: %v", err)
		return
	}

	s.writeResponse(w, http.StatusCreated, "Block added successfully", "")
	log.Println("[INFO] Block added successfully")
}

func (s *Server) GetLatestBlockHandler(w http.ResponseWriter, _ *http.Request) {
	block, err := s.blockchain.GetLatestBlock()
	if err != nil {
		s.writeResponse(w, http.StatusInternalServerError, nil, err.Error())
		log.Printf("[ERROR] Error getting latest block: %v", err)
		return
	}
	if block == nil {
		s.writeResponse(w, http.StatusNotFound, nil, "No blocks found")
		log.Println("[WARNING] No blocks found")
		return
	}
	s.writeResponse(w, http.StatusOK, block, "")
	log.Printf("[INFO] Latest block retrieved: %v", block)
}

func (s *Server) GetBlockHistoryHandler(w http.ResponseWriter, _ *http.Request) {
	blocks, err := s.blockchain.GetBlockHistory()
	if err != nil {
		s.writeResponse(w, http.StatusInternalServerError, nil, err.Error())
		log.Printf("[ERROR] Error getting block history: %v", err)
		return
	}
	s.writeResponse(w, http.StatusOK, blocks, "")
	log.Printf("[INFO] Block history retrieved: %d blocks", len(blocks))
}

func (s *Server) Start(addr string) error {
	r := mux.NewRouter()

	r.HandleFunc("/blocks", s.AddBlockHandler).Methods("POST")
	r.HandleFunc("/blocks/latest", s.GetLatestBlockHandler).Methods("GET")
	r.HandleFunc("/blocks/history", s.GetBlockHistoryHandler).Methods("GET")
	s.httpServer.Handler = r
	s.httpServer.Addr = addr
	log.Printf("[INFO] Starting server on %s\n", addr)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[ERROR] ListenAndServe(): %v", err)
		}
	}()

	log.Println("[INFO] Server started successfully")
	return nil
}

func (s *Server) Stop() error {
	log.Println("[INFO] Stopping server...")
	s.blockchain.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Error during server shutdown: %v", err)
		return err
	}
	log.Println("[INFO] Server stopped gracefully")
	return nil
}
