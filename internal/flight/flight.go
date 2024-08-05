package flight

import (
	"context"
	"time"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s Service) Find(ctx context.Context, budget int) (string, error) {
	//TODO: implement once we stop being sponsored
	time.Sleep(time.Second * 2)
	return "British Airways", nil
}
