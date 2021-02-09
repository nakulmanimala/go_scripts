package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	redColor    = "\033[1;33m%s\033[0m"
	yellowColor = "\033[1;31m%s\033[0m"
)

func main() {
	securityGpSet := make(map[string]string)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(sess)
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:PROJECT"),
				Values: []*string{
					aws.String("XXXXXXX"),
				},
			},
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		fmt.Println(err)
	}

	for _, securityGp := range result.SecurityGroups {
		input := &ec2.DescribeNetworkInterfacesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("group-id"),
					Values: []*string{
						aws.String(*securityGp.GroupId),
					},
				},
			},
		}
		result, err := svc.DescribeNetworkInterfaces(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf(yellowColor, *securityGp.GroupId+"\n")

		for _, rSec := range result.NetworkInterfaces {
			for _, gp := range rSec.Groups {
				if _, ok := securityGpSet[*gp.GroupId]; !ok {
					securityGpSet[*gp.GroupId] = *gp.GroupName
				}

			}
		}

		for k, v := range securityGpSet {
			fmt.Printf(redColor, "	|__"+k+"["+v+"]"+"\n")
		}

		for dk := range securityGpSet {
			delete(securityGpSet, dk)
		}
		fmt.Println()
	}

}

