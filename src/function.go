package main

import (
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

func SetInsecureSSL() {

	// Create a new transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Create a new client
	http.DefaultClient = &http.Client{
		Transport: transport,
	}
}

func ApiCaller(url string, method string, body io.Reader, headers map[string]string) ([]byte, int, error) {

	// Create a new request using the provided URL, method, and body
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %v", err)
	}

	// Add the headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading response: %v", err)
	}

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp.StatusCode, fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return respBody, resp.StatusCode, nil
}

func GetAllMetricsMetadata(prometheusURL string) MetricsMetadataResponseType {

	var metricsMetadata MetricsMetadataResponseType

	response, code, err := ApiCaller(prometheusURL+"/api/v1/targets/metadata", "GET", nil, nil)
	if err != nil {
		log.Fatal("Error query targets metadata", err)
	}
	if code != 200 {
		log.Fatal("Error code query http", code)
	}

	// map the response to a metricsMetadata
	err = json.Unmarshal(response, &metricsMetadata)
	if err != nil {
		log.Fatal("Error unmarshall metrics metadata", err)
	}

	return metricsMetadata
}

func GetAllAlertingRules(prometheusURL string) AlertingRulesResponseType {

	var alertingRules AlertingRulesResponseType

	response, code, err := ApiCaller(prometheusURL+"/api/v1/rules", "GET", nil, nil)
	if err != nil {
		log.Fatal("Error query rules", err)
	}
	if code != 200 {
		log.Fatal("Error code query http", code)
	}

	// map the response to a metricsMetadata
	err = json.Unmarshal(response, &alertingRules)
	if err != nil {
		log.Fatal("Error unmarshall metrics metadata", err)
	}

	return alertingRules
}

func WriteMetricsMetadataToCSV(metricsMetadata MetricsMetadataResponseType, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"instance", "job", "metric", "type", "help", "unit"})

	for _, metadata := range metricsMetadata.Data {
		instance := metadata.Target.Instance
		job := metadata.Target.Job
		metric := metadata.Metric
		metricType := metadata.Type
		help := metadata.Help
		unit := metadata.Unit

		writer.Write([]string{instance, job, metric, metricType, help, unit})
	}

	return nil
}

func WriteAlertingRulesToCSV(alertingRules AlertingRulesResponseType, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"instance", "alertname", "query", "summary", "description"})

	// Create the regex pattern to match newline characters
	pattern := regexp.MustCompile(`\n`)

	for _, group := range alertingRules.Data.Groups {
		for _, rule := range group.Rules {
			instance := rule.Labels.Instance
			alertname := rule.Name
			query := rule.Query
			summary := pattern.ReplaceAllString(rule.Annotations.Summary, " ")
			description := pattern.ReplaceAllString(rule.Annotations.Description, " ")

			writer.Write([]string{instance, alertname, query, summary, description})
		}
	}

	return nil
}
