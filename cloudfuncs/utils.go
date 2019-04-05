package cloudfuncs

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
)

var client *firestore.Client
var clientOnce sync.Once

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getClient(w http.ResponseWriter) *firestore.Client {
	clientOnce.Do(func() {
		var err error
		client, err = firestore.NewClient(context.Background(), getenv("GCP_PROJECT", "bad_proj"))
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("firestore.NewClient: %v", err)
			return
		}
	})
	return client
}
