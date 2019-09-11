package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/petrrusanov/cake-chicken/app/rest/api"
	"github.com/petrrusanov/cake-chicken/app/store/engine"
	"github.com/petrrusanov/cake-chicken/app/store/service"
	"github.com/pkg/errors"
)

// ServerCommand with command line flags and env
type ServerCommand struct {
	Store StoreGroup `group:"store" namespace:"store" env-namespace:"STORE"`
	Bolt BoltGroup `group:"bolt" namespace:"bolt" env-namespace:"BOLT"`

	HTTPPort int `long:"httpPort" env:"HTTP_PORT" default:"3000" description:"HTTP port"`

	CommonOpts
}

// StoreGroup defines options group for storage
type StoreGroup struct {
	Type string `long:"type" env:"TYPE" description:"type of storage" choice:"bolt" default:"bolt"`
}

// BoltGroup holds all bolt params, used by store
type BoltGroup struct {
	Path  string `long:"db" env:"PATH" description:"bolt database path"`
}

type serverApp struct {
	*ServerCommand
	restSrv    *api.Rest
	store      *service.DataStore
	terminated chan struct{}
}

// Execute is the entry point for "server" command
func (s *ServerCommand) Execute(args []string) error {
	log.Printf("[INFO] start server on port: %d", s.HTTPPort)

	ctx, cancel := context.WithCancel(context.Background())

	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("[WARN] interrupt signal")
		cancel()
	}()

	app, err := s.newServerApp()

	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	if err = app.run(ctx); err != nil {
		log.Printf("[INFO] terminated with error %+v", err)
		return err
	}

	log.Printf("[INFO] terminated")

	return nil
}

func (s *ServerCommand) newServerApp() (*serverApp, error) {
	storeEngine, err := s.makeStoreEngine()

	if err != nil {
		return nil, err
	}

	dataStore := &service.DataStore{
		Interface: storeEngine,
	}

	rest := &api.Rest{
		Version:       Revision,
		SharedSecret:  s.SharedSecret,
		DataStore:     dataStore,
	}

	return &serverApp{
		ServerCommand: s,
		restSrv:       rest,
		store:         dataStore,
		terminated:    make(chan struct{}),
	}, nil
}

func (s *ServerCommand) makeStoreEngine() (engine.Interface, error) {
	switch s.Store.Type {
	case "bolt":
		if s.Bolt.Path == "" {
			return nil, errors.New("no bolt path provided")
		}

		return engine.NewBolt(s.Bolt.Path)
	default:
		return nil, errors.Errorf("unsupported store type %s", s.Store.Type)
	}
}

// Run all application objects
func (a *serverApp) run(ctx context.Context) error {
	go func() { // shutdown on context cancellation
		<-ctx.Done()
		a.restSrv.Shutdown()
	}()

	a.restSrv.Run(a.HTTPPort)

	close(a.terminated)

	return nil
}

// Wait for application completion (termination)
func (a *serverApp) Wait() {
	<-a.terminated
}
