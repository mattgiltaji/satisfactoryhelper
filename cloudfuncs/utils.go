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

func getClient(w http.ResponseWriter) *firestore.Client {
	clientOnce.Do(func() {
		var err error
		client, err = firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("firestore.NewClient: %v", err)
			return
		}
	})
	return client
}
