package service

import (
	"github.com/dimebox/cake-chicken/app/store/engine"
)

// DataStore wraps engine.Interface with additional methods
type DataStore struct {
	engine.Interface
}
