package awsdata

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/manywho/awsinventory/internal/inventory"
	"github.com/sirupsen/logrus"
)

const (
	// AssetTypeS3Bucket is the value used in the AssetType field when fetching S3 Buckets
	AssetTypeS3Bucket string = "S3 Bucket"

	// ServiceS3 is the key for the S3 service
	ServiceS3 string = "s3"
)

func (d *AWSData) loadS3Buckets() {
	defer d.wg.Done()

	s3Svc := d.clients.GetS3Client(ValidRegions[0])

	log := d.log.WithFields(logrus.Fields{
		"region":  "global",
		"service": ServiceS3,
	})
	log.Info("loading data")
	out, err := s3Svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		d.results <- result{Err: err}
		return
	}

	log.Info("processing data")
	for _, b := range out.Buckets {
		d.results <- result{
			Row: inventory.Row{
				UniqueAssetIdentifier: aws.StringValue(b.Name),
				AssetType:             AssetTypeS3Bucket,
			},
		}
	}

	log.Info("finished processing data")
}
