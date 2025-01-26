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

func ToMessage(data string) Message {
	var message Message

	err := json.Unmarshal([]byte(data), &message)

	if err != nil {
		panic(err)
	}

	return message
}
