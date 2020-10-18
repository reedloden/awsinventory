package awsdata

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/manywho/awsinventory/internal/inventory"
	"github.com/sirupsen/logrus"
)

const (
	// AssetTypeSQSQueue is the value used in the AssetType field when fetching SQS queues
	AssetTypeSQSQueue string = "SQS Queue"

	// ServiceSQS is the key for the SQS service
	ServiceSQS string = "sqs"
)

func (d *AWSData) loadSQSQueues(region string) {
	defer d.wg.Done()

	sqsSvc := d.clients.GetSQSClient(region)

	log := d.log.WithFields(logrus.Fields{
		"region":  region,
		"service": ServiceSQS,
	})

	log.Info("loading data")

	var queueUrls []*string
	done := false
	params := &sqs.ListQueuesInput{}
	for !done {
		out, err := sqsSvc.ListQueues(params)

		if err != nil {
			d.results <- result{Err: err}
			return
		}

		queueUrls = append(queueUrls, out.QueueUrls...)

		if out.NextToken == nil {
			done = true
		} else {
			params.NextToken = out.NextToken
		}
	}

	log.Info("processing data")

	for _, q := range queueUrls {
		out, err := sqsSvc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			QueueUrl: q,
			AttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameApproximateNumberOfMessages),
				aws.String(sqs.QueueAttributeNameApproximateNumberOfMessagesNotVisible),
				aws.String(sqs.QueueAttributeNameQueueArn),
			},
		})
		if err != nil {
			d.results <- result{Err: err}
			return
		}

		d.results <- result{
			Row: inventory.Row{
				UniqueAssetIdentifier: (*q)[strings.LastIndex(aws.StringValue(q), "/")+1:],
				Virtual:               true,
				DNSNameOrURL:          aws.StringValue(q),
				Location:              region,
				AssetType:             AssetTypeSQSQueue,
				Comments:              fmt.Sprintf("%s, %s", aws.StringValue(out.Attributes[sqs.QueueAttributeNameApproximateNumberOfMessages]), aws.StringValue(out.Attributes[sqs.QueueAttributeNameApproximateNumberOfMessagesNotVisible])),
				SerialAssetTagNumber:  aws.StringValue(out.Attributes[sqs.QueueAttributeNameQueueArn]),
			},
		}
	}

	log.Info("finished processing data")
}
