package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/environment"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

const (
	waitTimeout         = 20
	maxNumberOfMessages = 1
)

type QueueManager struct {
	outputQueueURL    string
	outputQueueClient *sqs.Client
	inputQueueURL     string
	inputQueueClient  *sqs.Client
	errorQueueURL     string
	errorQueueClient  *sqs.Client
}

func ConfigQueueManager(environment environment.Environment) QueueManager {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	// Create SQS clients
	inputQueueClient := sqs.NewFromConfig(cfg)
	outputQueueClient := sqs.NewFromConfig(cfg)
	errorQueueClient := sqs.NewFromConfig(cfg)

	return QueueManager{
		inputQueueURL:     environment.VideoProcessingInputQueue,
		outputQueueURL:    environment.VideoProcessedOutputQueue,
		errorQueueURL:     environment.VideoProcessedErrorQueue,
		outputQueueClient: outputQueueClient,
		inputQueueClient:  inputQueueClient,
		errorQueueClient:  errorQueueClient,
	}
}

func (manager *QueueManager) PollMessages(
	chn chan<- *types.Message,
) {
	for {
		// Receive messages from source queue
		result, err := manager.outputQueueClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(manager.outputQueueURL),
			MaxNumberOfMessages: *aws.Int32(maxNumberOfMessages), // Receive only one message at a time
			WaitTimeSeconds:     *aws.Int32(waitTimeout),         // Wait for messages for 20 seconds
		})

		if err != nil {
			log.Fatalf("failed to receive messages, %v", err)
		}

		// Process received messages
		for _, message := range result.Messages {
			// Process the message content here (e.g., transform, enrich, etc.)
			chn <- &message
		}
	}
}

func (manager *QueueManager) PollErrorMessages(
	chn chan<- *types.Message,
) {
	for {
		// Receive messages from source queue
		result, err := manager.errorQueueClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(manager.errorQueueURL),
			MaxNumberOfMessages: *aws.Int32(maxNumberOfMessages), // Receive only one message at a time
			WaitTimeSeconds:     *aws.Int32(waitTimeout),         // Wait for messages for 20 seconds
		})

		if err != nil {
			log.Fatalf("failed to receive error messages, %v", err)
		}

		// Process received messages
		for _, message := range result.Messages {
			// Process the message content here (e.g., transform, enrich, etc.)
			chn <- &message
		}
	}
}

func (manager *QueueManager) WriteMessage(
	videoID *string,
) error {
	// Send processed message to destination queue
	_, err := manager.inputQueueClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(manager.inputQueueURL),
		MessageBody: videoID,
	})

	if err != nil {
		return responses.Wrap("queue manager: error to write message", err)
	}

	fmt.Println("Message processed and sent to destination queue")

	return nil
}

func (manager *QueueManager) DeleteMessage(
	receiptiHandle *string,
) error {
	// Delete message from source queue
	_, err := manager.outputQueueClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(manager.outputQueueURL),
		ReceiptHandle: receiptiHandle,
	})

	if err != nil {
		return responses.Wrap("queue manager: error when deleting message", err)
	}

	return nil
}
