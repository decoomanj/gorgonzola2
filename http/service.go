package gorgonzola

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const SERVICE_PORT = 8080

// ServiceServer handles administrative requests
type ServiceServer struct {
	*mux.Router
}

func NewServiceServer() *ServiceServer {
	return &ServiceServer{
		mux.NewRouter(),
	}
}
func (service *ServiceServer) Start() {
	service.StartOn(SERVICE_PORT)
}

func (service *ServiceServer) StartOn(port int) {

	// start the web server
	log.Printf("Service is listening on %d....\n", port)

	// TODO: http://www.hydrogen18.com/blog/stop-listening-http-server-go.html
	if err := http.ListenAndServe(":"+strconv.Itoa(port), service); err != nil {
		log.Fatal("Service ListenAndServe:", err)
	}
}

// Wrap a Handler with AccessLogger and Principal
func (m *ServiceServer) Handle(method string, path string, handler ContextHandler) {
	log.Printf("Adding resource [%s] %s\n", method, path)
	m.HandleFunc(path, Context{
		next: AccessLogger{handler}.ServeCtxHTTP,
	}.ServeHTTP).Methods(method)
}

// Wrap a Handler with AccessLogger and Principal
func (m *ServiceServer) Principal(method string, path string, handler ContextHandler) {
	fmt.Printf("Adding principal resource [%s] %s\n", method, path)
	m.HandleFunc(path, Context{
		next: AccessLogger{Principal{handler}.ServeCtxHTTP}.ServeCtxHTTP,
	}.ServeHTTP).Methods(method)
}

// Handle: Not Allowed Requests
func (ms *ServiceServer) NotAllowed(method string, path string) {
	fmt.Printf("NotAllowed resource [%s] %s\n", method, path)
	MethodNotAllowed := func(w http.ResponseWriter, r *http.Request, c *Context) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	ms.HandleFunc(path, Context{
		next: AccessLogger{MethodNotAllowed}.ServeCtxHTTP,
	}.ServeHTTP).Methods(method)
}
