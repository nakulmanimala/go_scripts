package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func main() {
	fmt.Println("========================")
	costSession, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := costexplorer.New(costSession)
	input := &costexplorer.GetCostAndUsageInput{
		Filter: &costexplorer.Expression{
			Tags: &costexplorer.TagValues{
				Key:    aws.String("PROJECT"),
				Values: aws.StringSlice([]string{"STAG-LENS"}),
			},
		},
		Metrics: aws.StringSlice([]string{"BlendedCost"}),
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String("2020-09-01"),
			End:   aws.String("2020-10-01"),
		},
		Granularity: aws.String("MONTHLY"),
	}
	req, res := svc.GetCostAndUsageRequest(input)
	errr := req.Send()
	if errr != nil {
		panic(errr)
	}
	fmt.Println(res)
	fmt.Println("========================")

}
