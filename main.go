package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Message struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

func handleMessage(message *Message, sesClient *ses.SES) error {
	// Prepare the email subject and body
	subject := message.Subject
	body := fmt.Sprintf("Name: %s %s\nEmail: %s\n\nMessage:\n%s",
		message.FirstName, message.LastName, message.Email, message.Message)

	// Send the email via SES
	input := &ses.SendEmailInput{
		Source: aws.String("dev+portfolio@tituscoleman.com"),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String("dev@tituscoleman.com"),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body),
				},
			},
		},
	}

	_, err := sesClient.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}

func handler(event events.SQSEvent) error {
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	// Create an SQS client
	sqsClient := sqs.New(sess)

	// Create an SES clientexport GOOS=linux
	sesClient := ses.New(sess)

	// Process each message in the event
	for _, record := range event.Records {
		message := &Message{}

		// Parse the message JSON
		err := json.Unmarshal([]byte(record.Body), message)
		if err != nil {
			return err
		}

		// Handle the message
		err = handleMessage(message, sesClient)
		if err != nil {
			return err
		}

		// Delete the processed message from the SQS queue
		_, err = sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(record.EventSourceARN),
			ReceiptHandle: aws.String(record.ReceiptHandle),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
