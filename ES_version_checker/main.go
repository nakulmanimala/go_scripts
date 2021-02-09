package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/ses"
)

type config struct {
	Sender      string
	Recipient1  string
	Recipient2  string
	Subject     string
	Mode        string
	Profile     string
	Region      string
	DomainNames []string
}

func printError(s string, e error) {
	if e != nil {
		fmt.Printf("%s : %s", s, e)
	}
}

func sendMail(sess *session.Session, msg *string, conf *config) {
	svc := ses.New(sess)
	CharSet := "UTF-8"
	subject := conf.Subject
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(conf.Recipient1),
				aws.String(conf.Recipient2),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(*msg),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(conf.Subject),
			},
		},
		Source: aws.String(conf.Sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error(), conf.Recipient1, conf.Recipient2, subject)
		} else {
			fmt.Println(err, conf.Recipient1, conf.Recipient2, subject)
		}
		fmt.Println("error : ", err, " To User  :", conf.Recipient1, conf.Recipient2, " Subject :", subject)
	}
	if result != nil {
		fmt.Println("Message Id: ", *result.MessageId, " To User  :", conf.Recipient1, conf.Recipient2, " Subject :", subject)
		fmt.Println("Email sent to address: ", conf.Recipient1, conf.Recipient2, " Subject :", subject, " Message Id :", *result.MessageId)
	} else {
		fmt.Println("Sending email returned empty result")
	}

}

var (
	sess    *session.Session
	err     error
	message *string
	msgText string = ""
	flag    int
)

func recoverName() {
	if r := recover(); r != nil {
		fmt.Println("\nRecovered from ", r)
	}
}

func checkPanic(i interface{}) {
	if i == nil {
		panic("Panicking!!!!")
	}
}

func main() {

	//Unmarshaling json
	conf := &config{}
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), conf)

	mode := conf.Mode
	defer recoverName()
	// Check wherther it is profile or role
	if mode == "profile" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Profile: conf.Profile,
			Config: aws.Config{
				Region: aws.String(conf.Region),
			},
		})
	} else if mode == "role" {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(conf.Region)},
		)
	} else {
		fmt.Println("Unknown mode")
	}
	printError("Unable to start session :", err)

	es := elasticsearchservice.New(sess)
	for _, domain := range conf.DomainNames {
		input := &elasticsearchservice.DescribeElasticsearchDomainInput{
			DomainName: aws.String(domain),
		}
		r, err := es.DescribeElasticsearchDomain(input)
		printError("Error in describe ES domain :", err)
		checkPanic(r)
		//Checking if there is any update available
		avail := r.DomainStatus.ServiceSoftwareOptions.UpdateAvailable
		if *avail == true {
			msgText = msgText + fmt.Sprintf("Domain: %v\n-----------------------\nupdate Available: %v\nNew Version : %v\nAutomatic Update date: %v\n***********************\n", domain, *r.DomainStatus.ServiceSoftwareOptions.UpdateAvailable, *r.DomainStatus.ServiceSoftwareOptions.NewVersion,*r.DomainStatus.ServiceSoftwareOptions.AutomatedUpdateDate)
			message = &msgText
			flag = flag + 1
		}
	}
	// send email if there any update

	if flag != 0 {
		fmt.Println(*message)
		sendMail(sess, message, conf)
	}

}

