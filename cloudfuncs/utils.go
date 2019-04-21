package cloudfuncs

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var client *firestore.Client
var clientOnce sync.Once

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getClient(w http.ResponseWriter) *firestore.Client {
	clientOnce.Do(func() {
		var err error
		//TODO: see how we encrypt the credential file for CI/CD
		googleAuthFileName := "test-statisfactory-helper-auth.json"
		workingDir, err := os.Getwd()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("Could not determine current directory to load auth file: %v", err)
		}
		googleAuthFileLocation := filepath.Join(workingDir, googleAuthFileName)

		//only try to use the auth file to connect if it exists
		if _, err := os.Stat(googleAuthFileLocation); os.IsNotExist(err) {
			//no auth file, fallback to default credentials
			client, err = firestore.NewClient(context.Background(), getEnv("GCP_PROJECT", "bad_proj"))
		} else {
			client, err = firestore.NewClient(context.Background(), getEnv("GCP_PROJECT", "bad_proj"),
				option.WithCredentialsFile(googleAuthFileLocation))
		}

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("firestore.NewClient: %v", err)
			return
		}
	})
	return client
}
