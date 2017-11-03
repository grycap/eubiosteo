package pergamo

import (
	"fmt"
	"net/http"
	"pergamo/pergamo/state"
	"pergamo/scheduler"
	"strconv"
)

type Server struct {
	config     *Config
	stateStore *state.StateStore
	scheduler  *scheduler.Scheduler
	httpServer *HttpServer
}

func (s *Server) InitializeSchema() error {
	return s.stateStore.InitializeSchema()
}

// Listen es para escuchar todo lo necesario
func (s *Server) Listen() {
	s.httpServer.listen()
}

func (s *Server) RunSPA() {
	if s.config.SPAPort == 0 {
		return
	}

	fmt.Printf("Up and running spa on port %d\n", s.config.SPAPort)

	fs := http.FileServer(http.Dir("./frontend/public"))
	http.Handle("/", fs)

	http.ListenAndServe(":"+strconv.Itoa(s.config.SPAPort), nil)
}

func createDispatcher(config *Config) (scheduler.Dispatch, error) {
	if config.Dispatcher == "local" {
		return scheduler.NewLocalDispatcher(config.SchedulerPath)
	}

	if config.Dispatcher == "slurm" {
		return scheduler.NewSlurmDispatcher(config.SchedulerPath)
	}

	return nil, fmt.Errorf("Dispatcher with name %s does not exists. Use 'local' or 'slurm'", config.Dispatcher)
}

func NewServer(config *Config) (*Server, error) {
	stateStore, err := state.NewStateStore(config.OnedataPath, config.DBDriver, config.DBPath)
	if err != nil {
		return nil, err
	}

	dispatcher, err := createDispatcher(config)
	if err != nil {
		return nil, err
	}

	scheduler, err := scheduler.NewScheduler(config.SchedulerPath, config.OnedataPath, config.APIPort, dispatcher)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		stateStore: stateStore,
		scheduler:  scheduler,
	}

	httpServer, err := NewHttpServer(server, config)
	if err != nil {
		return nil, fmt.Errorf("Http server error %s", err.Error())
	}

	go server.RunSPA()

	server.httpServer = httpServer
	return server, nil
}
