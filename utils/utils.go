package utils

import (
	"io/ioutil"
	"log"
)

// ReadServiceAccountToken returns token of service account inside container
func ReadServiceAccountToken() string {
	content, err := ioutil.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
