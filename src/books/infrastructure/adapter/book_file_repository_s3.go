package adapter

import (
	"bytes"
	"context"
	"log"
	"main/src/books/domain/repository"
	"mime"

	appError "main/utils/error"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BookFileRepositoryS3 struct {
	ctx        context.Context
	client     *s3.Client
	BucketName string
	BucketKey  string
}

func NewBookFileRepositoryS3(ctx context.Context, client *s3.Client, bucketName, bucketKey string) repository.BookFileRepository {
	return &BookFileRepositoryS3{
		ctx:        ctx,
		client:     client,
		BucketName: bucketName,
		BucketKey:  bucketKey,
	}
}

func (r *BookFileRepositoryS3) SaveBookFile(file *bytes.Reader, bucketKey, fileExt string) *appError.Error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(r.BucketName),
		Key:         aws.String(bucketKey),
		Body:        file,
		ContentType: aws.String(mime.TypeByExtension(fileExt)),
	}

	_, errS3 := r.client.PutObject(r.ctx, input)
	if errS3 != nil {
		log.Printf("Error while putting object to S3: %v\n", errS3)
		return appError.NewBadRequestError("Error while putting object to S3")
	}

	log.Printf("Book file creation completed successfully, book: %+v", bucketKey)
	return nil
}


func (r *BookFileRepositoryS3) DeleteBookFile(bucketKey string) (*appError.Error) {
	input := &s3.DeleteObjectInput{
        Bucket: aws.String(r.BucketName),
        Key:    aws.String(bucketKey),
    }

	_, err := r.client.DeleteObject(context.TODO(), input)
    if err != nil {
        log.Printf("unable to delete object, %v", err)
		return appError.NewBadRequestError("Error while putting object to S3")
    }

    log.Printf("Object %s deleted successfully from %s\n", bucketKey, r.BucketName)
	return nil
}