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
	client, err = firestore.NewClient(context.Background(), getenv("GCP_PROJECT", "bad_proj"), option.WithCredentialsFile(googleAuthFileLocation))
	if err != nil {
		t.Error("Could not connect to test firestore instance")
	}
	return
}
