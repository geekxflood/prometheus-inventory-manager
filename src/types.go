package main

import "time"

type MetricsResponseType struct {
	Status string             `json:"status"`
	Data   []MetricsNamesType `json:"data"`
}

type MetricsNamesType string

type MetricsMetadataResponseType struct {
	Status string                `json:"status"`
	Data   []MetricsMetadataType `json:"data"`
}

type MetricsMetadataType struct {
	Target map[string]string `json:"target"`
	Metric string            `json:"metric"`
	Type   string            `json:"type"`
	Help   string            `json:"help"`
	Unit   string            `json:"unit"`
}

type TargetMetadataType struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

type AlertingRulesResponseType struct {
	Data struct {
		Groups []struct {
			Rules []RuleType `json:"rules"`
		} `json:"groups"`
	} `json:"data"`
	Status string `json:"status"`
}

type RuleType struct {
	Alerts       []AlertType       `json:"alerts"`
	Annotations  map[string]string `json:"annotations"`
	Health       string            `json:"health"`
	Labels       map[string]string `json:"labels"`
	Name         string            `json:"name"`
	Query        string            `json:"query"`
	Type         string            `json:"type"`
	Duration     int               `json:"duration"`
	Groups       []string          `json:"groups"`
	LastExecuted int               `json:"lastExecuted"`
}

type AlertType struct {
	ActiveAt    time.Time         `json:"activeAt"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
	State       string            `json:"state"`
	Value       string            `json:"value"`
}
