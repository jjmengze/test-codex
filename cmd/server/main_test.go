package main

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"os"
	"testing"
	"time"

	"log-receiver/pkg/auth"

	"github.com/bmizerany/assert"
)

var (
	externalRouterTS *httptest.Server
	externalTSURL    string

	// Tokens for testing
	saoExpiredToken      string
	saoNeverExpiredToken string
	sdsNeverExpiredToken string
	sxxNeverExpiredToken string

	// Test data
	testGzBase64     = "H4sICOgDB18AA3Rlc3QAAwAAAAAAAAAAAA=="
	testGz           []byte
	testPbBase64     = "GgNzYW8iATEqATE6BGd6aXA="
	testPb           []byte
	test1mbGz        []byte
	test30mbGz       []byte
	withFastTrack    = true
	withoutFastTrack = false
)

func TestMain(m *testing.M) {
	// setup
	now := time.Now().UTC().Unix()
	saoExpiredToken = auth.GenIDPJWTToken("sao", "123", "789", "456", now-3600)
	saoNeverExpiredToken = auth.GenIDPJWTToken("sao", "123", "789", "456", now+3600*24)
	sdsNeverExpiredToken = auth.GenIDPJWTToken("sds", "123", "789", "456", now+3600*24)
	sxxNeverExpiredToken = auth.GenIDPJWTToken("sxx", "123", "789", "456", now+3600*24)

	os.Getenv("JWT_PRIVATE_KEY_PATH")

	// Test data
	var err error
	testGz, err = base64.StdEncoding.DecodeString(testGzBase64)
	if err != nil {
		panic(err)
	}
	testPb, err = base64.StdEncoding.DecodeString(testPbBase64)
	if err != nil {
		panic(err)
	}
	test1mbGz, err = makeTestGz(1024 * 1024)
	if err != nil {
		panic(err)
	}
	test30mbGz, err = makeTestGz(1024 * 1024 * 30)
	if err != nil {
		panic(err)
	}

	//externalRouterTS = httptest.NewServer(externalServer())
	defer externalRouterTS.Close()
	externalTSURL = externalRouterTS.URL

	// run testing
	v := m.Run()
	// teardown
	os.Exit(v)
}

func makeTestGz(size int) ([]byte, error) {
	var gzData []byte
	randomData := make([]byte, size)
	_, err := rand.Read(randomData)
	if err != nil {
		return nil, err
	}
	gzData, err = compressToGz(randomData)
	if err != nil {
		return nil, err
	}
	return gzData, nil
}

func compressToGz(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	_, err := gzw.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gzw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestGetHealth(t *testing.T) {
	uri := "/health"
	tests := []struct {
		name               string
		requestBody        io.Reader
		wantHTTPStatusCode int
	}{
		{"OK", nil, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := externalRouterTS.Client().Get(externalTSURL + uri)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.wantHTTPStatusCode, resp.StatusCode)
		})
	}

}

// TODO: This UT need to pass after your implementation
func TestPostActivityLog(t *testing.T) {
	uri := "/api/v2/activity_log/%s"
	method := "POST"
	tests := []struct {
		name               string
		forwarder          string
		productCode        string
		token              string
		contentEncoding    string
		senderID           string
		sourceID           string
		extraInfo          string
		fastTrack          bool
		requestBody        []byte
		wantHTTPStatusCode int
	}{
		{"OK", "", "sao", saoNeverExpiredToken, "gzip", "7601c152-2f85-48de-af99-e8c3da9ae41b", "", "", withoutFastTrack, testGz, http.StatusOK},
		{"OK with x-Fast-Track", "", "sao", saoNeverExpiredToken, "gzip", "7601c152-2f85-48de-af99-e8c3da9ae41b", "", "", withFastTrack, testGz, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(method, fmt.Sprintf(externalTSURL+uri, tt.productCode), bytes.NewReader(tt.requestBody))
			req.Header.Add("Authorization", "Bearer "+tt.token)
			req.Header.Add("Content-Encoding", tt.contentEncoding)
			req.Header.Add("Content-Type", "application/gzip")
			req.Header.Add("X-Sender-ID", tt.senderID)
			req.Header.Add("X-Source-ID", tt.sourceID)
			req.Header.Add("X-Forwarder", tt.forwarder)
			req.Header.Add("X-Extra-Info", tt.extraInfo)

			if tt.fastTrack {
				req.Header.Add("X-Fast-Track", "")
			}

			resp, err := externalRouterTS.Client().Do(req)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.wantHTTPStatusCode, resp.StatusCode)
		})
	}
}

// TODO: This UT need to pass after you implement the checkRequestBody middleware
func TestCheckOver30mb(t *testing.T) {
	uri := "/api/v2/activity_log/%s"
	tests := []struct {
		name               string
		token              string
		requestBody        []byte
		wantHTTPStatusCode int
	}{
		{"no OK over 30mb", saoNeverExpiredToken, test30mbGz, http.StatusRequestEntityTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", fmt.Sprintf(externalTSURL+uri, "sao"), bytes.NewReader(tt.requestBody))
			req.Header.Add("Authorization", "Bearer "+tt.token)
			req.Header.Add("Content-Encoding", "gzip")
			req.Header.Add("Content-Type", "application/gzip")
			resp, err := externalRouterTS.Client().Do(req)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.wantHTTPStatusCode, resp.StatusCode)
		})
	}
}

// TODO: This UT need to pass after you implement the ValidateTokenController middleware
func TestPostValidateToken(t *testing.T) {
	uri := "/api/v2/validate_token/%s"
	method := "POST"
	tests := []struct {
		name               string
		forwarder          string
		productCode        string
		token              string
		wantHTTPStatusCode int
	}{
		{"OK", "", "sao", saoNeverExpiredToken, http.StatusOK},
		{"Key Pair not Matched", "", "sao", sdsNeverExpiredToken, http.StatusUnauthorized},
		{"Product Not Supported", "", "sxx", sxxNeverExpiredToken, http.StatusNotAcceptable},
		//todo not sure what is path not matched? if the path not Matched it would not routing to endpoint....so what is that mean?
		//{"Path and Token Product Code Not Matched", "", "sao", sxxNeverExpiredToken, http.StatusBadRequest},
		{"Expired Token", "", "sao", saoExpiredToken, http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(method, fmt.Sprintf(externalTSURL+uri, tt.productCode), nil)
			req.Header.Add("Authorization", "Bearer "+tt.token)
			req.Header.Add("X-Forwarder", tt.forwarder)
			resp, err := externalRouterTS.Client().Do(req)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.wantHTTPStatusCode, resp.StatusCode)
		})
	}
}
