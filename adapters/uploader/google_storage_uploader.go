package uploader

import (
	"context"
	"humanlogger/core/ports"
	"io/ioutil"
	"os"
	"path"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GoogleStorageUploader struct {
	BucketName string
}

func (u *GoogleStorageUploader) Upload(s *ports.LoggingSession, packedPath string) error {
	// Read the contents of the packed file into memory
	file, err := os.Open(packedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Create a new Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(".secret_google_storage.json"))
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a new bucket handle
	bucketHandle := client.Bucket(u.BucketName)

	// Create a new object handle and upload the file to Google Cloud Storage
	objectHandle := bucketHandle.Object(path.Base(packedPath))
	writer := objectHandle.NewWriter(ctx)
	defer writer.Close()

	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}
