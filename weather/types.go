package weather

import (
	"fmt"
	"strings"
)

type WeatherReport struct {
	Address string `json:"address"`
	Days    []Day  `json:"days"`
}

type Day struct {
	Datetime   string  `json:"datetime"`
	Temp       float64 `json:"temp"`
	Conditions string  `json:"conditions"`
}

type BusinessInsight struct {
	RiskLevel       string   `json:"risk_level"`
	Summary         string   `json:"summary"`
	Impacts         []string `json:"impacts"`
	Recommendations []string `json:"recommendations"`
}

func (i BusinessInsight) Validate() error {
	switch i.RiskLevel {
	case "low", "moderate", "high", "critical":
	default:
		return fmt.Errorf("AI response has invalid risk_level %q", i.RiskLevel)
	}
	if strings.TrimSpace(i.Summary) == "" {
		return fmt.Errorf("AI response has an empty summary")
	}
	if len(i.Impacts) == 0 || len(i.Recommendations) == 0 {
		return fmt.Errorf("AI response must include impacts and recommendations")
	}
	for _, impact := range i.Impacts {
		if strings.TrimSpace(impact) == "" {
			return fmt.Errorf("AI response contains an empty impact")
		}
	}
	for _, recommendation := range i.Recommendations {
		if strings.TrimSpace(recommendation) == "" {
			return fmt.Errorf("AI response contains an empty recommendation")
		}
	}
	return nil
}
