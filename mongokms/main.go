// Package mongokms implements a Google Cloud KMS client to decrypt encrypted secrets.
package mongokms

import (
	"context"
	"hash/crc32"
	"log"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Decrypts a given ciphertext using symmetric encryption
// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
// ciphertext := []byte("...")  // result of a symmetric encryption call
func DecriptSymmetric(name string, ciphertext []byte) string {

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatalf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Optional, but recommended: Compute ciphertext's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	ciphertextCRC32C := crc32c(ciphertext)

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:             name,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	// Call the API.
	result, err := client.Decrypt(ctx, req)
	if err != nil {
		log.Fatalf("failed to decrypt ciphertext: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		log.Fatalf("decrypt: response corrupted in-transit")
	}

	return string(result.Plaintext)
}
