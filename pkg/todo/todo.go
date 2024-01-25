package todo

import (
	"fmt"
	"log"

	"github.com/vidya-ranganathan/mcluster/pkg/apis/cumulonimbus.ai/v1alpha1"
	//"github.com/vidya-ranganathan/mcluster/pkg/work"
)

func Add(spec v1alpha1.MclusterSpec) {
	log.Printf("Specs --> Cluster name %s\n", spec.Name)

	// Build the URL
	url := fmt.Sprintf("http://localhost:8080/cluster/%s", spec.Name)

	// Define the payload to be sent in the request body
	payload := map[string]interface{}{
		"name": spec.Name,
	}

	/*
		// Call the PUT function
		err := work.PutVerb(url, payload)
		if err != nil {
			fmt.Println("Error:", err)
		}
	*/

	fmt.Println(url)
	fmt.Println(payload)
}
