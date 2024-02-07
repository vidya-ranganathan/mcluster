package work

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PutVerb(url string) (string, error) {
	clusterID := "N/A"

	// Create a simple PUT request without a body
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return clusterID, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Make the PUT request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return clusterID, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is 200 OK
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return clusterID, fmt.Errorf("Error reading response body: %v", err)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return clusterID, fmt.Errorf("Error unmarshaling response JSON: %v", err)
		}

		var ok bool
		// Assuming the server sends the clusterID in the response
		if clusterID, ok = response["clusterID"].(string); !ok {
			return clusterID, fmt.Errorf("clusterID not found in the response")
		}
	} else {
		fmt.Println("Response Status:", resp.Status)
	}

	// returning clusterID now..
	return clusterID, nil
}

func DeleteVerb(url string) error {
	// Create a simple DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Make the DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to delete cluster, response status: %s", resp.Status)
	}

	return nil
}
