package flight

import (
	"context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s Service) Find(ctx context.Context, budget int) (string, error) {
	//TODO: implement once we stop being sponsored
	return "British Airways", nil
}
