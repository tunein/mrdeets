package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tunein/dshareiff-playground/mrDeets/metadata"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/tunein/streaming-common/go-config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {

	opts := metadata.Ec2Options{}
	if err := config.Resolve(&opts); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(opts.Region),
	})
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	params := &ec2.DescribeInstancesInput{
		Filters: createFilters(&opts),
	}
	client := ec2.New(s)
	resp, _ := client.DescribeInstances(params)

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {

			p := metadata.NewMinimalEC2Instance(instance)

			if opts.Out != nil {
				results := createOutFilter(&opts, p)
				for _, result := range results {
					fmt.Println(result)
				}
			} else {
				b, _ := json.Marshal(p)
				q, _ := prettyPrintJSON(b)
				fmt.Println(string(q))
			}

		}
	}

}
