package engine

import "github.com/dimebox/cake-chicken/app/store/models"

// Interface combines all store interfaces
type Interface interface {
	Chicken
	Cake
}

// Chicken chicken methods
type Chicken interface {
	AddChicken(username string, prefix string) (models.UserCounter, error)
	FulfillChicken(username string, prefix string) (models.UserCounter, error)
	GetChickenStats(prefix string) ([]models.UserCounter, error)
}

// Cake chicken methods
type Cake interface {
	AddCake(username string, prefix string) (models.UserCounter, error)
	FulfillCake(username string, prefix string) (models.UserCounter, error)
	GetCakeStats(prefix string) ([]models.UserCounter, error)
}
