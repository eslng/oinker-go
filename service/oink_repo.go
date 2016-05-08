package service

import "github.com/eslng/oinker-go/model"

type OinkRepo interface {
	Service
	Create(model.Oink) (model.Oink, error)
	FindByID(id string) (model.Oink, bool, error)
	All() ([]model.Oink, error)
}
