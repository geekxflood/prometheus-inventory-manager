package main

import (
	"fmt"
	"log"
	"os"
)

const DefaultPrometheusURL = "http://localhost:9090"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check if output directory exists, if not create it
	os.MkdirAll("output", os.ModePerm)
}

func main() {
	log.Println("Starting prometheus-inventory-exporter")

	// Check if Prometheus URL is set, if not use DefaultPrometheusURL
	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		prometheusURL = DefaultPrometheusURL
	}

	// Get all metrics metadata
	metricsMetadata := GetAllMetricsMetadata(prometheusURL)

	// Get all alerting rules
	alertingRules := GetAllAlertingRules(prometheusURL)

	// Write metrics metadata to CSV
	metricsOutputFilename := "output/metrics.csv"
	err := WriteMetricsMetadataToCSV(metricsMetadata, metricsOutputFilename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metrics metadata written to %s\n", metricsOutputFilename)

	_ = alertingRules

	log.Println("Finished prometheus-inventory-exporter")

}
