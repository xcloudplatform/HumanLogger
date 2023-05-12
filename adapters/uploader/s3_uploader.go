package uploader

import (
	"bytes"
	"github.com/ClickerAI/ClickerAI/core/ports"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Uploader struct {
	bucketName string
	region     string
}

func (u *S3Uploader) Upload(s *ports.LoggingSession, packedPath string) error {
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

	// Create a new AWS session with the specified region
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(u.region),
	})
	if err != nil {
		return err
	}

	// Create a new S3 client
	s3Client := s3.New(awsSession)

	// Upload the file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(u.bucketName),
		Key:    aws.String(path.Base(packedPath)),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}
