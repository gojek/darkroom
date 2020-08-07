package gcs

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestInterfaces(t *testing.T) {
	c, _ := storage.NewClient(
		context.TODO(),
		option.WithHTTPClient(
			newTestClient(func(req *http.Request) *http.Response {
				if strings.Contains(req.URL.String(), "failed-bucket") {
					return &http.Response{
						StatusCode: 422,
						Body:       ioutil.NopCloser(strings.NewReader("")),
						Header:     make(http.Header),
					}
				}
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("contents")),
					Header:     make(http.Header),
				}
			}),
		),
	)
	basicTests(t, bucketHandle{c.Bucket("bucket")})
	basicFailedTests(t, bucketHandle{c.Bucket("failed-bucket")})
}

func basicTests(t *testing.T, bkt BucketHandle) {
	b := readObject(t, bkt.Object("stiface-test"))
	if got, want := string(b), "contents"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func basicFailedTests(t *testing.T, bkt BucketHandle) {
	b := readFailedObject(t, bkt.Object("stiface-test"))
	if got, want := string(b), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func readObject(t *testing.T, obj ObjectHandle) []byte {
	r, err := obj.NewReader(context.Background())
	if err != nil {
		t.Fatalf("reading %v: %v", obj, err)
	}
	defer r.Close()
	bytesR, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("reading %v: %v", obj, err)
	}
	rr, err := obj.NewRangeReader(context.Background(), 0, -1)
	if err != nil {
		t.Fatalf("reading %v: %v", obj, err)
	}
	defer rr.Close()
	bytesRR, err := ioutil.ReadAll(rr)
	if err != nil {
		t.Fatalf("reading %v: %v", obj, err)
	}
	if bytes.Compare(bytesR, bytesRR) == 0 {
		return bytesR
	}
	return []byte(nil)
}

func readFailedObject(t *testing.T, obj ObjectHandle) []byte {
	r, err := obj.NewReader(context.Background())
	if err == nil {
		defer r.Close()
		t.Fatalf("expected error %v", obj)
	}
	rr, err := obj.NewRangeReader(context.Background(), 0, -1)
	if err == nil {
		defer rr.Close()
		t.Fatalf("expected error %v", obj)
	}
	return []byte(nil)
}
