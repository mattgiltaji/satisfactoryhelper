package cloudfuncs

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// ***** Helpers *****
func getTestClient(ctx context.Context, t *testing.T) (client *firestore.Client) {
	var err error
	googleAuthFileName := "test-statisfactory-helper-auth.json"
	workingDir, err := os.Getwd()
	if err != nil {
		t.Error("Could not determine current directory to load test auth file")
	}
	googleAuthFileLocation := filepath.Join(workingDir, googleAuthFileName)
	client, err = firestore.NewClient(ctx, getEnv("GCP_PROJECT", "bad_proj"), option.WithCredentialsFile(googleAuthFileLocation))
	if err != nil {
		t.Error("Could not connect to test firestore instance")
	}
	return
}

//deleteDocRef is a test helper for deferring the cleanup of created firestore documents.
//Defer a call to this method and it will handle deleting the document reference and checking for errors
func deleteDocRef(ctx context.Context, t *testing.T, doc *firestore.DocumentRef) {
	_, err := doc.Delete(ctx)
	if err != nil {
		t.Errorf("Could not delete document, Error: %v", err)
	}
}
