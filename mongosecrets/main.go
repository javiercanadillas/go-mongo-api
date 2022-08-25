package mongosecrets

import (
	"context"
	"log"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func GetSecret(secretID string, secretVersion string) []byte {
	// GCP project in which to store secrets in Secret Manager.
	projectID := "javiercm-webapp"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	versionPath := strings.Join([]string{
		"projects/",
		projectID,
		"/secrets/",
		secretID,
		"/versions/",
		secretVersion,
	}, "")

	// Get the secret
	getSecretReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: versionPath,
	}

	// Call the API
	result, err := client.AccessSecretVersion(ctx, getSecretReq)
	if err != nil {
		log.Fatalf("Failed to get secret version: %v", err)
	}

	secretContent := result.Payload.Data
	log.Printf("retrieved payload for %s\n", &result.Name)
	return secretContent
}
