package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

// StorageClient is a wrapper for a GCP Session
type StorageClient struct {
	Client  *storage.Client
	bucket  *storage.BucketHandle
	context context.Context
}

// NewStorageClient returns a new GCP client for the GCS API
func NewStorageClient(bucketName string) (*StorageClient, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	bucket := client.Bucket(bucketName)

	if err != nil {
		return nil, err
	}

	return &StorageClient{
		Client:  client,
		bucket:  bucket,
		context: ctx,
	}, nil
}

// AddFileToBucket will upload a single file to GCP, it will require a pre-built gcp session
// and will set file info like content type and encryption on the uploaded file.
func (c *StorageClient) AddFileToBucket(key string, body []byte) error {
	obj := c.bucket.Object(key)
	w := obj.NewWriter(c.context)
	w.Write(body)
	err := w.Close()

	return err
}

// GetFileFromBucket func
func (c *StorageClient) GetFileFromBucket(key string) ([]byte, error) {
	rc, err := c.bucket.Object(key).NewReader(c.context)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

// RenameObject func
func (c *StorageClient) RenameObject(src, dest string) error {
	srcObj := c.bucket.Object(src)
	dstObj := c.bucket.Object(dest)

	if _, err := dstObj.CopierFrom(srcObj).Run(c.context); err != nil {
		return err
	}
	if err := srcObj.Delete(c.context); err != nil {
		return err
	}
	return nil
}
