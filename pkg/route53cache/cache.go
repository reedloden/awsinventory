package route53cache

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

type Cache struct {
	records []*route53.ResourceRecordSet
}

func New(records []*route53.ResourceRecordSet) *Cache {
	return &Cache{
		records: records,
	}
}

func (c *Cache) FindRecordsForInstance(i *ec2.Instance) (results []string) {
	for _, r := range c.records {
		switch aws.StringValue(r.Type) {
		case "CNAME":
			for _, record := range r.ResourceRecords {
				if aws.StringValue(i.PrivateDnsName) != "" && strings.Contains(record.String(), aws.StringValue(i.PrivateDnsName)) {
					results = append(results, aws.StringValue(r.Name))
					break
				}
				if aws.StringValue(i.PublicDnsName) != "" && strings.Contains(record.String(), aws.StringValue(i.PublicDnsName)) {
					results = append(results, aws.StringValue(r.Name))
					break
				}
			}
			break
		default:
			for _, record := range r.ResourceRecords {
				if aws.StringValue(i.PrivateIpAddress) != "" && strings.Contains(record.String(), aws.StringValue(i.PrivateIpAddress)) {
					results = append(results, aws.StringValue(r.Name))
					break
				}
				if aws.StringValue(i.PublicIpAddress) != "" && strings.Contains(record.String(), aws.StringValue(i.PublicIpAddress)) {
					results = append(results, aws.StringValue(r.Name))
					break
				}
			}
			break
		}
	}

	return
}
