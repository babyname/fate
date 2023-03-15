package service

import (
	"github.com/babyname/fate"
)

type Service struct {
	fate fate.Fate
}

func New(fate fate.Fate) *Service {
	return &Service{
		fate: fate,
	}
}
