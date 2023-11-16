package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	key := os.Getenv("DO_SPACE_KEY")
	secret := os.Getenv("DO_SPACE_SECRET")
	endpoint := os.Getenv("DO_SPACE_ENDPOINT")
	region := os.Getenv("DO_SPACE_REGION")
	// https://khongfamily.sgp1.digitaloceanspaces.com

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.New(newSession)

	bucket := "khongfamily"

	// List all files in bucket
	listObjectsInput := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	listObjectsOutput, err := s3Client.ListObjects(listObjectsInput)
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listObjectsOutput.Contents {
		log.Println(*object.Key)
		log.Println(aws.StringValue(object.Key))
	}

}

func deleteFileFromBucket(s3Client *s3.S3, bucket string, filepath string) {
	// Delete file from bucket
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	}

	result, err := s3Client.DeleteObject(deleteInput)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}

func uploadFileToBucket(s3Client *s3.S3, bucket string, file *os.File, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// generate uuid for new file
	uuid := uuid.New().String()

	newFilepath := fmt.Sprintf("users/%s.png", uuid)
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(newFilepath),
		Body:   file,
		// public
		ACL: aws.String("public-read"),
	}

	uploadOutput, err := s3Client.PutObject(uploadInput)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(uploadOutput)
}
