package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	// Specify profile for config and region for requests
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
	})

	if err != nil {
		fmt.Println("Error creating session ", err)
		return
	}

	// DeleteBucket function deletes a bucket.
	DeleteBucket(sess, "aws-sample-bucket-1")

	// ListBuckets function lists the buckets in your account.
	buckets, _ := ListBuckets(sess)
	fmt.Println("buckets: ", buckets)

	// CreateBucket Creates an S3 Bucket in the region configured in the shared config
	CreateBucket(sess, "aws-sample-bucket-1")

	// UploadFile function uploads an object to a bucket.
	UploadFile(sess, "iam-store", "serverless.yml")

	// ListObjects function lists the items in a bucket
	objects, _ := ListObjects(sess, "iam-store")
	fmt.Println("objects: ", objects)

	// DownloadObject function downloads an object from a bucket.
	// file, _ := DownloadObject(sess, "iam-store", "serverless.yml", "serverless1.yml")
	// fmt.Println("file: ", file.Name())

	// DeleteObject function deletes an object from a bucket.
	DeleteObject(sess, "iam-store", "serverless.yml")

	// ListObjects function lists the items in a bucket
	objects, _ = ListObjects(sess, "iam-store")
	fmt.Println("objects: ", objects)
}

// ListBuckets function lists the buckets in your account.
func ListBuckets(sess *session.Session) ([]string, error) {
	var buckets []string

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	for _, v := range result.Buckets {
		buckets = append(buckets, *v.Name)
	}

	return buckets, nil
}

// CreateBucket Creates an S3 Bucket in the region configured in the shared config
func CreateBucket(sess *session.Session, bucket string) {
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		fmt.Printf("Unable to create bucket %q, %v", bucket, err)
		return
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		fmt.Printf("Error occurred while waiting for bucket to be created, %v", bucket)
		return
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)
}

// DeleteBucket function deletes a bucket.
func DeleteBucket(sess *session.Session, bucket string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return err
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	return err
}

// ListObjects function lists the items in a bucket
func ListObjects(sess *session.Session, bucket string) ([]string, error) {
	var objects []string

	svc := s3.New(sess)

	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return nil, err
	}

	for _, item := range result.Contents {
		objects = append(objects, *item.Key)
	}

	return objects, nil
}

// UploadFile function uploads an object to a bucket.
func UploadFile(sess *session.Session, bucket string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Unable to open file %v", err)
		return err
	}

	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	return err
}

// DownloadObject function downloads an object from a bucket.
func DownloadObject(sess *session.Session, bucket string, input string, output string) (*os.File, error) {
	file, err := os.Create(output)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(input),
	})

	if err != nil {
		return nil, err
	}

	return file, nil
}

// DeleteObject function deletes an object from a bucket.
func DeleteObject(sess *session.Session, bucket string, filename string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	return err
}
