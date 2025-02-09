package entity

import (
	"encoding/json"
	"mime/multipart"
)

type UoloaderDocumentEntity struct {
	Data        multipart.File
	Name        string
	ContentType string
	Size        int64
}

type Message struct {
	ZippedURL     string `json:"zippedURL"`
	TrackingID    string `json:"trackingID"`
	ReceiptHandle string `json:"receiptHandle"`
}

type ErrorMessage struct {
	TrackingID string `json:"trackingID"`
	Message    string `json:"errorMessage"`
}

func ToMessage(data string) Message {
	var message Message

	err := json.Unmarshal([]byte(data), &message)

	if err != nil {
		panic(err)
	}

	return message
}

func ToErrorMessage(data string) ErrorMessage {
	var message ErrorMessage

	err := json.Unmarshal([]byte(data), &message)

	if err != nil {
		panic(err)
	}

	return message
}
