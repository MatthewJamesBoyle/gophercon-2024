package hotel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Option struct {
	Name             string `json:"name"`
	PricePerNightGBP int    `json:"price_per_night_gbp"`
	StarRating       int    `json:"rating"`
}

type Response struct {
	HotelOptions []Option `json:"hotel_options"`
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	doer    doer
	baseUrl string
}

func NewService(doer doer, baseUrl string) *Service {
	return &Service{doer: doer, baseUrl: baseUrl}
}

func (s Service) Find(ctx context.Context, budget int) (string, error) {
	if budget == 0 {
		return "", errors.New("budget cannot be 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", s.baseUrl, "hotels"), nil)
	if err != nil {
		return "", err
	}
	res, err := s.doer.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get recomendation: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get recomendation: %s", res.Status)
	}

	var r Response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("failed to decode recomendation: %w", err)
	}
	for _, h := range r.HotelOptions {
		if h.PricePerNightGBP*h.StarRating > 400 {
			return h.Name, nil
		}
	}

	return r.HotelOptions[0].Name, nil
}
