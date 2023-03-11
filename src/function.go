package main

import (
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func SetInsecureSSL() {
	// Create a new transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a new client
	http.DefaultClient = &http.Client{Transport: tr}
}

func ApiCaller(url string, method string, body io.Reader, headers map[string]string) ([]byte, int, error) {

	// Create a new request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	// Add the headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

func GetAllMetricsMetadata(prometheusURL string) MetricsMetadataResponseType {
	var metricsMetadata MetricsMetadataResponseType

	response, code, error := ApiCaller(prometheusURL+"/api/v1/targets/metadata", "GET", nil, nil)
	if error != nil {
		panic(error)
	}
	if code != 200 {
		panic(code)
	}

	// map the response to a metricsMetadata
	err := json.Unmarshal(response, &metricsMetadata)
	if err != nil {
		panic(error)
	}

	return metricsMetadata
}

func GetAllAlertingRules(prometheusURL string) AlertingRulesResponseType {
	var alertingRules AlertingRulesResponseType

	response, code, error := ApiCaller(prometheusURL+"/api/v1/rules", "GET", nil, nil)
	if error != nil {
		panic(error)
	}
	if code != 200 {
		panic(code)
	}

	// map the response to a metricsMetadata
	err := json.Unmarshal(response, &alertingRules)
	if err != nil {
		panic(error)
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
