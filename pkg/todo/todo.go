package todo

import (
	"fmt"
	"log"

	"github.com/vidya-ranganathan/mcluster/pkg/apis/cumulonimbus.ai/v1alpha1"
	"github.com/vidya-ranganathan/mcluster/pkg/work"
)

func Add(spec v1alpha1.MclusterSpec) string {
	log.Printf("Specs --> Cluster name %s\n", spec.Name)

	// Build the URL
	url := fmt.Sprintf("http://localhost:8080/cluster/%s", spec.Name)

	// Call the PUT function
	clusterID, err := work.PutVerb(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("PutVerb: Cluster created with ID =%s\n", clusterID)
	return clusterID
}

func Delete(clusterName string) bool {
	log.Printf("Cluster name %s\n", clusterName)

	// Build the URL
	url := fmt.Sprintf("http://localhost:8080/cluster/%s", clusterName)

	// Call the DELETE function
	err := work.DeleteVerb(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return true
}
