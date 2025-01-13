package remote_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"slices"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

const (
	respMsg = `<?xml version="1.0" enconding="UTF-8"?>
	<CompleteMultipartUploadOutput>
		<Location>mockValue</Location>
		<Bucket>mockValue</Bucket>
		<Key>mockValue</key>
		<ETag>mockValue</ETag>
	</CompleteMultipartUploadOutput>
	`
)

type S3Mock struct {
	s3iface.S3API
	files map[string][]byte
	tags  map[string]map[string]string
}

func (m S3Mock) PutObject(in *s3.PutObjectInput) (out *s3.PutObjectOutput, err error) {
	key := path.Join(*in.Bucket, *in.Key)
	m.files[key], err = io.ReadAll(in.Body)

	m.tags[key] = map[string]string{}

	if in.Tagging != nil {
		u, err := url.Parse("/?" + *in.Tagging)
		if err != nil {
			panic(fmt.Errorf("unable to parse AWS S3 Tagging string %q: %w", *in.Tagging, err))
		}

		q := u.Query()
		for k := range q {
			m.tags[key][key] = q.Get(k)
		}
	}

	return &s3.PutObjectOutput{}, nil
}

func emptyList() []string {
	return []string{}
}

func mockS3(ignoreOps []string, status int) *s3.S3 {
	var m sync.Mutex

	partNum := 0
	names := []string{}
	params := []any{}

	svc := s3.New(unit.Session)

	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.UnmarshalError.Clear()
	svc.Handlers.Send.Clear()
	svc.Handlers.Send.PushBack(func(r *request.Request) {
		m.Lock()
		defer m.Unlock()

		if slices.Contains(ignoreOps, r.Operation.Name) {
			names = append(names, r.Operation.Name)
			params = append(params, r.Params)
		}

		r.HTTPResponse = &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewReader([]byte(respMsg))),
		}

		switch data := r.Data.(type) {
		case *s3.CreateMultipartUploadOutput:
			data.UploadId = aws.String("UPLOAD-ID")
		case *s3.UploadPartOutput:
			partNum++
			data.ETag = aws.String(fmt.Sprintf("ETAG%d", partNum))
		case *s3.CompleteMultipartUploadOutput:
			data.Location = aws.String("https://location")
			data.VersionId = aws.String("VERSION-ID")
			data.ETag = aws.String("ETAG")
		case *s3.PutObjectOutput:
			data.VersionId = aws.String("VERSION-ID")
			data.ETag = aws.String("ETAG")
		}
	})

	return svc
}
