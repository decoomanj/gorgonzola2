package gorgonzola

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MicroService struct {
	Admin   *AdminServer
	Service *ServiceServer
	Health  *Health
	Metrics *prometheus.Registry
	name    string
}

type ContextHandler func(http.ResponseWriter, *http.Request, *Context)

// Instantiate a new microservice
func NewMicroService(name string) *MicroService {
	return (&MicroService{
		Admin:   NewAdminServer(),
		Service: NewServiceServer(),
		Health:  NewHealth(),
		Metrics: prometheus.DefaultRegisterer.(*prometheus.Registry),
		name:    name,
	}).RegisterSignalHealth()
}

func (ms *MicroService) RegisterSignalHealth() *MicroService {
	ms.Health.Register(&HealthCheck{
		Name: "signal",
		Handler: func() error {
			if shutdownMode {
				return errors.New("Shutdown mode activated")
			}
			return nil
		},
		Interval: time.Second,
	})
	return ms
}

// Start the administration only
func (ms *MicroService) StartAdmin() {
	ms.StartAdminOn(ADMIN_PORT)
}

// Start the administration only
func (ms *MicroService) StartAdminOn(adminport int) {

	// add health
	ms.Admin.HandleFunc("/health", ms.Health.page)

	// add metrics
	ms.Admin.Handle("/metrics", prometheus.Handler())

	// Start admin server
	go ms.Admin.StartOn(adminport)
}

// Start the service only
func (ms *MicroService) StartService() {

	ms.StartServiceOn(SERVICE_PORT)
}

// Start the service only
func (ms *MicroService) StartServiceOn(servicePort int) {

	go ms.Service.StartOn(servicePort)
}

// Start a microservice with default health page on the given port
func (ms *MicroService) StartOn(servicePort, adminport int) <-chan bool {

	ms.StartAdminOn(adminport)

	ms.StartServiceOn(servicePort)

	return make(chan bool)
}

// Start a microservice with defaults
func (ms *MicroService) Start() <-chan bool {
	return ms.StartOn(SERVICE_PORT, ADMIN_PORT)
}
