package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"server"`
	CiruitBreaker struct {
		Name            string  `json:"name"`
		MaxRequests     uint32  `json:"maxRequests"`
		Interval        int     `json:"interval"`
		Timeout         int     `json:"timeout"`
		FailureRatio    float64 `json:"failureRatio"`
		MinimumRequests int64   `json:"minimumRequests"`
	} `json:"circuitBreaker"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
