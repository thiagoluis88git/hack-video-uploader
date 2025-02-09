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
	fi, err := os.Open("video.mp4")
	if err != nil {
		panic(err)
	}

	fiStat, err := fi.Stat()
	if err != nil {
		panic(err)
	}

	putResponse, err := httpRequester.Put("https://hack-video-processing-bucket.s3.amazonaws.com/bc46dfd5-1c12-48f3-a528-06379089f8f3?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAZAJMBW6KBFS43JSR%2F20250209%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250209T210023Z&X-Amz-Expires=900&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEJX%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIArim3ucj4IPlVMU4%2FHvjaBbImlrDmhREwc88U5mnZP4AiEAyLOiUtJZcXMwhNyTlCGrD5x0u3XjJhicJP4vee%2FRzWIqugIIrv%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARACGgw2MTkxMDQ1NDg3NTYiDA%2FdK1GMKj3cwVcs5yqOArzQyi6hp%2B%2FkyiKAPzjH0zTqJW5VEJeuTIfRQfPpCxbbL3IWf%2FvrhbZlqz0nx%2BFUcMbnl0cpTV7MRmrgcNz2kOWmpTappEnzX0jEgzs53odD5mJajrq3OftJQpbZCpqj%2FJaGaS24gmE6lM1jub5BK0XVOf3Sl8HUNWa1SlHeXfVXqDad636OigmYSCwJRAdkgeBhmmYAv3DGFgD0PrTQvEVt2b0%2B1YU9msI7hz3uBjO%2FHCKY3V6rtSCwOFOV0047Rivx0ikW2pRsjSzulbnmFSZNP4BPT0u1Q07gJznVfm7LmLJFvTGTZmIiw2WsPRAr2J3M2lRWUqDhNdU0awpsE6C0aeJh7VsKCWXa9HiAhjD6oqS9BjqdAWddFdVrnQQvnyw48CZguP04xAln%2FFHR8Odhi66bZM%2BJJ606qMix%2BQtRFPFCtS3CJBhJ8MI%2B4cBEy1B423uudlE%2B8yw8ODMiUHeRE4UuX8S244kWmdIJfL%2F%2BCeuCMwdvVr0u6c4yNxqsWaTNaS5zrzPztMApOSY7DdzP0IFdYvMtis2wL7nMOX3wbgfi07bHsWpy87zXRp%2FaprP%2Fl7s%3D&X-Amz-SignedHeaders=host&X-Amz-Signature=1a475e74f5110ee1d34976def86855b145a1b286982b9ba85c9bd33567412802", fiStat.Size(), fi)
	if err != nil {
		panic(err)
	}
	log.Printf("presigned URL returned %v.", putResponse.StatusCode)
	log.Println(strings.Repeat("-", 88))
}
