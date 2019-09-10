package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boilerplate/backend/app/rest/api"
	"github.com/boilerplate/backend/app/rest/auth"
	"github.com/boilerplate/backend/app/store/engine"
	"github.com/boilerplate/backend/app/store/service"
	"github.com/go-pkgz/mongo"
	"github.com/pkg/errors"
)

// ServerCommand with command line flags and env
type ServerCommand struct {
	Store StoreGroup `group:"store" namespace:"store" env-namespace:"STORE"`
	Mongo MongoGroup `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`

	HTTPPort int `long:"httpPort" env:"HTTP_PORT" default:"3000" description:"HTTP port"`

	CommonOpts
}

// StoreGroup defines options group for storage
type StoreGroup struct {
	Type string `long:"type" env:"TYPE" description:"type of storage" choice:"mongo" default:"mongo"`
}

// MongoGroup holds all mongo params, used by store
type MongoGroup struct {
	URL string `long:"url" env:"URL" default:"localhost" description:"mongo url"`
	DB  string `long:"db" env:"DB" default:"backend" description:"mongo database name"`
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
		Authenticator: auth.Authenticator{DataStore: dataStore},
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
	case "mongo":
		mgServer, err := s.makeMongo()

		if err != nil {
			return nil, errors.Wrap(err, "failed to create mongo server")
		}

		conn := mongo.NewConnection(mgServer, s.Mongo.DB, "")

		return engine.NewMongo(conn, 500, 100*time.Millisecond)
	default:
		return nil, errors.Errorf("unsupported store type %s", s.Store.Type)
	}
}

func (s *ServerCommand) makeMongo() (result *mongo.Server, err error) {
	if s.Mongo.URL == "" {
		return nil, errors.New("no mongo url provided")
	}

	if s.Mongo.DB == "" {
		return nil, errors.New("no mongo db provided")
	}

	return mongo.NewServerWithURL(s.Mongo.URL, 10*time.Second)
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
