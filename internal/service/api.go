package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SibHelly/TgUrlShorter/internal/models"
)

type UrlShorter interface {
	Create(models.Url) error
	GetAll() ([]*models.Url, error)
	Delete(alias string) error
	Info(alias string) (*models.UrlInfo, error)
}

type UrlService struct {
	baseURL string
	client  *http.Client
}

func NewURLService(apiBaseURL string) *UrlService {
	return &UrlService{
		baseURL: apiBaseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *UrlService) Create(models.Url) error {
	return nil
}

type ApiResponse struct {
	Data   []models.Url `json:"data"`
	Result string       `json:"result"`
}

func (s *UrlService) GetAll() ([]*models.Url, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/my/all", s.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	urls := make([]*models.Url, len(apiResponse.Data))
	for i := range apiResponse.Data {
		urls[i] = &apiResponse.Data[i]
	}

	return urls, nil
}

func (s *UrlService) Delete(alias string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/url/%s", s.baseURL, alias), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

type ApiResponseInfo struct {
	Data   models.UrlInfo `json:"data"`
	Result string         `json:"result"`
}

func (s *UrlService) Info(alias string) (*models.UrlInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/url/%s", s.baseURL, alias), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var apiResponse ApiResponseInfo
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}
	url := &apiResponse.Data
	return url, nil
}
