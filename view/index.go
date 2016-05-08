package view

import (
	"github.com/eslng/oinker-go/model"
)

type Index struct {
	Page
	Oinks   []model.Oink
	IsEmpty bool
}
