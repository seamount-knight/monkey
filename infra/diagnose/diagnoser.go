package diagnose

import (
	"encoding/json"
	"sync"
	"time"
)

const (
	// StatusOK means the component is K
	StatusOK HealthStatus = "OK"
	// StatusError means there is an error with the component
	StatusError HealthStatus = "ERROR"
)

type HealthStatus string

type Component interface {
	Diagnose() ComponentReport
}

type ComponentReport struct {
	Status     HealthStatus  `json:"status"`
	Name       string        `json:"name"`
	Message    string        `json:"message"`
	Suggestion string        `json:"suggestion"`
	Latency    time.Duration `json:"latency"`
}

// NewReport constructor
func NewReport(component string) *ComponentReport {
	return &ComponentReport{
		Status:  StatusOK,
		Name:    component,
		Message: "ok",
	}
}

func (c *ComponentReport) Check(err error, message, suggestion string) {
	if err != nil {
		c.Status = StatusError
		c.Message = fmt.Sprintf("%s: \"%s\"", message, err.Error())
		c.Suggestion = suggestion
	}
}

func (c *ComponentReport) AddLatency(start time.Time) {
	duration := time.Since(start)
	if c.Latency == time.Duration(0) || c.Latency < duration {
		c.Latency = duration
	}
}

func (c *ComponentReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status     HealthStatus `json:"status"`
		Name       string       `json:"name"`
		Message    string       `json:"message"`
		Suggestion string       `json:"suggestion"`
		Latency    string       `json:"latency"`
	}{
		c.Status,
		c.Name,
		c.Message,
		c.Suggestion,
		c.Latency.String(),
	})
}

type HealthChecker struct {
	Components []Component
}

func New() (*HealthChecker, error) {
	return &HealthChecker{}, nil
}

func (h *HealthChecker) Add(com Component) *HealthChecker {
	if h.Components == nil {
		h.Components = make([]Component, 0, 2)
	}
	h.Components = append(h.Components, com)
	return h
}

func (h *HealthChecker) Check() HealthReport {
	report := &HealthReport{Status: StatusOK}
	if h.Components == nil {
		return *report
	}
	wait := sync.WaitGroup{}
	for _, c := range h.Components {
		wait.Add(1)
		go func(component Component) {
			diagnose := component.Diagnose()
			if diagnose.Status == StatusError && report.Status != StatusError {
				report.Status = StatusError
			}
			report.Add(diagnose)
			wait.Done()
		}(c)
	}
	wait.Wait()
	return *report
}

type HealthReport struct {
	Status  HealthStatus      `json:"status"`
	Details []ComponentReport `json:"details"`
}

func (h *HealthReport) Add(report ComponentReport) {
	if h.Details == nil {
		h.Details = make([]ComponentReport, 0, 3)
	}
	h.Details = append(h.Details, report)
}
