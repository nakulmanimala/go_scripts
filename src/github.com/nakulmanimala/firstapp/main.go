package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Output struct {
	Cost struct {
		CostUnit       string
		TotalCost      string
		TimePeriod     string
		Metrics        string
		CostByServices []ServiceCost
	}
}
type ServiceCost struct {
	ServiceName string `json:"ServiceName"`
	Cost        string `json:"cost"`
}

func main() {
	//fmt.Println("========================")
	costSession, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := costexplorer.New(costSession)
	input := &costexplorer.GetCostAndUsageInput{
		Filter: &costexplorer.Expression{
			Tags: &costexplorer.TagValues{
				Key:    aws.String("PROJECT"),
				Values: aws.StringSlice([]string{"STAG-LENS", "STAG-LENS-HA"}),
			},
		},
		Metrics: aws.StringSlice([]string{"UnblendedCost"}),
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String("2020-09-01"),
			End:   aws.String("2020-10-01"),
		},
		Granularity: aws.String("MONTHLY"),
		GroupBy: []*costexplorer.GroupDefinition{
			&costexplorer.GroupDefinition{
				Key:  aws.String("SERVICE"),
				Type: aws.String("DIMENSION"),
			},
		},
	}
	req, res := svc.GetCostAndUsageRequest(input)
	errr := req.Send()
	if errr != nil {
		panic(errr)
	}
	var output Output
	//fmt.Println(res.ResultsByTime[0].Groups)
	for _, group := range res.ResultsByTime[0].Groups {
		cost := *group.Metrics["UnblendedCost"].Amount
		serviceName := *group.Keys[0]
		serviceCost := ServiceCost{
			Cost:        cost,
			ServiceName: serviceName,
		}
		output.Cost.CostByServices = append(output.Cost.CostByServices, serviceCost)
	}
	output.Cost.CostUnit = "USD"
	output.Cost.TimePeriod = "2020-09-01 to 2020-10-01"
	//fmt.Println(output)
	out, err := json.Marshal(output)
	if err != nil {
		// handle error
	}
	fmt.Println(string(out))
	//fmt.Println("========================")
}
