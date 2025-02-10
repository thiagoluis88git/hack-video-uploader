package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// IHttpRequester abstracts HTTP requests into an interface so it can be mocked during
// unit testing.
type IHttpRequester interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	Put(url string, contentLength int64, body io.Reader) (resp *http.Response, err error)
	Delete(url string) (resp *http.Response, err error)
}

// HttpRequester uses the net/http package to make HTTP requests during the scenario.
type HttpRequester struct{}

func (httpReq HttpRequester) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
func (httpReq HttpRequester) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	postRequest, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	postRequest.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(postRequest)
}

func (httpReq HttpRequester) Put(url string, contentLength int64, body io.Reader) (resp *http.Response, err error) {
	putRequest, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	putRequest.ContentLength = contentLength

	putRequest.Header.Add("Content-Type", "video/mp4")

	return http.DefaultClient.Do(putRequest)
}

func (httpReq HttpRequester) Delete(url string) (resp *http.Response, err error) {
	delRequest, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(delRequest)
}

func main() {
	httpRequester := HttpRequester{}
	// open input file
	fi, err := os.Open("video5.mp4")
	if err != nil {
		panic(err)
	}

	fiStat, err := fi.Stat()
	if err != nil {
		panic(err)
	}

	putResponse, err := httpRequester.Put("https://hack-video-processing-bucket.s3.amazonaws.com/97b5b252-70b6-473e-a440-a6530908dedd?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAZAJMBW6KEWJGQFEK%2F20250210%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250210T003822Z&X-Amz-Expires=900&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEJn%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLWVhc3QtMSJIMEYCIQDDBXBNWZlUZcCPKpSsohnZmCHuEl71vSvQVhliWxLS5wIhANYAZ1QybvBEyOk%2BZ%2BRw%2BCuGx0H%2BFuLkb4VmPaWsiV3HKsIFCLL%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQAhoMNjE5MTA0NTQ4NzU2IgygqKCIhtYT2%2FUj%2BboqlgU7PDqGDu7r%2B0%2FShee8nvXF4YlEhXMMMTfHxIgT5iv8LSPpKw7qdjUsyfeM%2FTQoUSAMQMt%2B5jmXFNbrGujNNhrtxH1h7fy8ZuEjAP6SkprNWsoHWHXaMjWE3CUfazaLx%2BMNY6fRBpKw5CCwIeJV38G%2FVHgDPSi9AV0pcK3nmlpG26R4KBdVq6ZzkiXDzwgHetE4TNIeUY8Dp36c48Zl8dbEvsRDkoBHV3VIQdsfZtsXquQvPe7BUau1c7J3%2FQay0s%2F3f2%2FS9pxTdf%2FHYbXc9AeDKr2lhwXVdgM3B7aBiZirw5wNm7tenyHoMIvuvjPwot2r39hmXYYjXoZebCRf1DScJ4mBOHZMKGyKi1O%2B44xb0rpcQLKwvHN4LPVMPXXRRcRQFspt8YVBwRFOc8eHz2HJsoR%2FJ21AvWKjV7Bb16FPqrK6glEb4SWLA7KN3GPY6leSAcsNK5GXB2sKzjFcrpjpIYW2kf3XqzngOHrAPCOr7vi%2BbAlEv8%2BRBGealWXWPrBnyN%2Bosq1VESLH%2BcTZuHM%2BNJyxK03zdsCb4ncJOML%2F4vJF2DK%2F6xB1MT%2FXmHKcflun2DA%2BU%2B9SGKsA7W5lYCZnZNs8Cc6H2P8qEJ7RdWLJnrhCA2EvnkAlxB15qu6voMeUjx9r%2FxAY8f3owmYnEK9qnmV8fyRi7v0qCZ8owv9PwxC00BDj%2BkgN6YPcnxK8pJIuhidAjvGelgQBkTqvbPyofsF3Q%2Bwq3BXjBfW6d807bhPyw%2FgvgF0FYa979mJgsRKh5fTpDQDXbenxgOTFM%2BWev2YdrDFYUcWp%2BNFtS1heHhXiN3F11U2WfuiNr%2F8MRwOU9vY0lLj0pnpl62HSoGZb%2BLZxJ4uB6yIjKd8%2FOK6MRPLbWhSuijCclKW9BjqwARpFkoP%2F2TVdWZ56mGiYgPDPEDXGrZSwvJVC31a5%2B3RDs6Fjsy4869mz1PbivgeoxutstlIV%2FgQHEoBnrP4fonZqcIqh9um%2B42rtEUuqslbBzb57qDBlVu056dJw1Q9Jcu3%2FrJDO0OOTNrmO%2FnbBuR1mnOMo0WCakJIFHWqFrCaqOyyb4X7KGgWhFFpRHBTuqEQ65HhH8lG4zRQ%2FgazU4lqSd4TUKBiAKqbhMw1MKK1R&X-Amz-SignedHeaders=host&X-Amz-Signature=821183176840e815e9d1fe21a000e9871f6d0a595c2a571b64669594a13494a9", fiStat.Size(), fi)
	if err != nil {
		panic(err)
	}
	log.Printf("presigned URL returned %v.", putResponse.StatusCode)
	log.Println(strings.Repeat("-", 88))
}
