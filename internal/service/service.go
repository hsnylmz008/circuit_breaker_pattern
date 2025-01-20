package service

import (
	"errors"
	"time"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SimulateExternalCall() (interface{}, error) {
	// Simüle edilmiş dış servis çağrısı
	if time.Now().Unix()%2 == 0 {
		return nil, errors.New("service error")
	}
	return "Success!", nil
}
