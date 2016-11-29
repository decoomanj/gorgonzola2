package gorgonzola

import (
	"log"
	"net/http"
	"strconv"
)

const ADMIN_PORT = 9090

// AdminServer handles administrative requests
type AdminServer struct {
	*http.ServeMux
}

func NewAdminServer() *AdminServer {
	return &AdminServer{
		http.NewServeMux(),
	}
}

func (admin *AdminServer) Start() {
	admin.StartOn(ADMIN_PORT)
}

func (admin *AdminServer) StartOn(port int) {
	log.Println("Starting administration server")

	// start the web server
	log.Printf("Administrator is listening on %d....\n", port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), admin); err != nil {
		log.Fatal("Administrator ListenAndServe:", err)
	}
}
