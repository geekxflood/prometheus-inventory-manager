package main

import (
	"os"

	"github.com/charmbracelet/log"
)

const DefaultPrometheusURL = "http://localhost:9090"

var err error

func init() {
	// Check if output directory exists, if not create it
	os.MkdirAll("output", os.ModePerm)
}

func main() {
	log.Info("Starting prometheus-inventory-exporter")

	SetInsecureSSL()

	// Write metrics metadata to CSV
	metricsOutputFilename := "output/metrics.csv"

	// Write alerting rules to CSV
	alertingRulesOutputFilename := "output/alertingRules.csv"

	// Check if Prometheus URL is set, if not use DefaultPrometheusURL
	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		prometheusURL = DefaultPrometheusURL
	}

	// Get all metrics metadata
	metricsMetadata := GetAllMetricsMetadata(prometheusURL)

	// Get all alerting rules
	AlertingRulesResponse := GetAllAlertingRules(prometheusURL)

	err = WriteMetricsMetadataToCSV(metricsMetadata, metricsOutputFilename)
	if err != nil {
		log.Fatal("Write CSV MetricsMetada error", err)
	}
	log.Info("Metrics metadata written to", "path", metricsOutputFilename)

	err = WriteAlertingRulesToCSV(AlertingRulesResponse, alertingRulesOutputFilename)
	if err != nil {
		log.Fatal("Write CSV Rules error", err)
	}
	log.Info("Alerting rules written to", "path", alertingRulesOutputFilename)

	log.Info("Finished prometheus-inventory-exporter")
}
