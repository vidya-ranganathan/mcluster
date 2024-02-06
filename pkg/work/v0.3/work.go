package work

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func PutVerb(url string, payload map[string]interface{}) error {
	fmt.Println("calling PutVerb")

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create a request with the payload
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Make the PUT request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %+v", err)
	}
	defer resp.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	// Read and print the response body if needed
	// responseBody, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("Response Body:", string(responseBody))

	return nil
}

func DeleteVerb(url string, payload map[string]interface{}) error {
	fmt.Println("calling DeleteVerb")

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create a request with the payload
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Make the DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %+v", err)
	}
	defer resp.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	// Read and print the response body if needed
	// responseBody, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("Response Body:", string(responseBody))

	return nil
}
