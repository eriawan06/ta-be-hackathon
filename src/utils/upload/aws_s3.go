package upload

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
)

// S3Info file information
type S3Info struct {
	Endpoint *string
	Key      string
	Secret   string
	Region   string
	Bucket   string
	File     multipart.File
	Filename string
	Filemime string
	Filesize int64
}

// PushS3Buffer ...
func PushS3Buffer(buffer *bytes.Reader, in S3Info) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      &in.Region,
		Credentials: credentials.NewStaticCredentials(in.Key, in.Secret, ""),
		Endpoint:    in.Endpoint,
	})
	if err != nil {
		return err
	}

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(in.Bucket),
		Key:                  aws.String(in.Filename),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 buffer,
		ContentLength:        aws.Int64(in.Filesize),
		ContentType:          aws.String(in.Filemime),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}

	return nil
}

func PushS3(in S3Info) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      &in.Region,
		Credentials: credentials.NewStaticCredentials(in.Key, in.Secret, ""),
		Endpoint:    in.Endpoint,
	})
	if err != nil {
		return err
	}

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(in.Bucket),
		Key:                  aws.String(in.Filename),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 in.File,
		ContentLength:        aws.Int64(in.Filesize),
		ContentType:          aws.String(in.Filemime),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteS3Object ...
func DeleteS3Object(in S3Info) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(in.Region),
		Credentials: credentials.NewStaticCredentials(in.Key, in.Secret, ""),
		Endpoint:    in.Endpoint,
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(in.Bucket),
		Key:    aws.String(in.Filename),
	})
	if err != nil {
		return err
	}

	return nil
}
