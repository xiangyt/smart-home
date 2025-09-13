package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"smart-home-backend/internal/models"
	"strings"
	"time"
)

type HomeAssistantService struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewHomeAssistantService(baseURL, token string) *HomeAssistantService {
	return &HomeAssistantService{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetStates 获取所有设备状态
func (s *HomeAssistantService) GetStates() ([]models.HADevice, error) {
	url := fmt.Sprintf("%s/api/states", s.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var devices []models.HADevice
	if err := json.Unmarshal(body, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// GetState 获取单个设备状态
func (s *HomeAssistantService) GetState(entityID string) (*models.HADevice, error) {
	url := fmt.Sprintf("%s/api/states/%s", s.baseURL, entityID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var device models.HADevice
	if err := json.Unmarshal(body, &device); err != nil {
		return nil, err
	}

	return &device, nil
}

// CallService 调用Home Assistant服务
func (s *HomeAssistantService) CallService(domain, service string, entityID string, serviceData map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/services/%s/%s", s.baseURL, domain, service)

	payload := models.HAServiceCall{
		Domain:  domain,
		Service: service,
		Target: map[string]interface{}{
			"entity_id": entityID,
		},
		ServiceData: serviceData,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HA API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetHistory 获取设备历史数据
func (s *HomeAssistantService) GetHistory(entityID string, startTime, endTime *time.Time) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/history/period", s.baseURL)

	if startTime != nil {
		url += "/" + startTime.Format(time.RFC3339)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filter_entity_id", entityID)
	if endTime != nil {
		q.Add("end_time", endTime.Format(time.RFC3339))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var history [][]map[string]interface{}
	if err := json.Unmarshal(body, &history); err != nil {
		return nil, err
	}

	if len(history) > 0 {
		return history[0], nil
	}

	return []map[string]interface{}{}, nil
}
