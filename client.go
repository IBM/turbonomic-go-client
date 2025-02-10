// Copyright (c) IBM Corporation
// SPDX-License-Identifier: Apache-2.0

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS-IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package turboclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Default base path of Turbonomic API
const (
	BasePathv3 = "/api/v3"
)

// Parameters for creating a Turbonomic client
type ClientParameters struct {
	Baseurl    string
	Hostname   string
	Username   string
	Password   string
	Skipverify bool
}

// Turbonomic Client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type CommonReqParams struct {
	QueryParameters map[string]string
	Headers         map[string]string
}

// Parameters for making a request to Turbonomic's API
type RequestOptions struct {
	Method          string
	Path            string
	ReqDTO          *bytes.Buffer
	CommonReqParams CommonReqParams
	// QueryParameters  map[string]string
	// Headers          map[string]string
}

// Creates a new instance of the Turbonomic Client
func NewClient(clientParams *ClientParameters) (*Client, error) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if clientParams.Skipverify {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: clientParams.Skipverify}
	}

	var basepath string
	if clientParams.Baseurl == "" {
		basepath = BasePathv3
	} else {
		basepath = clientParams.Baseurl
	}
	jar, _ := cookiejar.New(nil)

	client := &AuthRequest{
		basePath: basepath,
		hostname: clientParams.Hostname,
		username: clientParams.Username,
		password: clientParams.Password,
		httpClient: &http.Client{
			Jar:       jar,
			Timeout:   time.Minute,
			Transport: customTransport,
		},
	}

	return clientAuth(client)

}

// Returns *url.URL with query parameters set
func setParams(baseUrl string, QueryParameters map[string]string) (*url.URL, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	query := fullUrl.Query()
	for k, v := range QueryParameters {
		query.Set(k, v)
	}

	fullUrl.RawQuery = query.Encode()
	return fullUrl, err
}

// Make request to Turbonomic API using http package
func (c *Client) request(reqOpt RequestOptions) ([]byte, error) {

	baseUrl := c.BaseURL + reqOpt.Path
	fullUrl, err := setParams(baseUrl, reqOpt.CommonReqParams.QueryParameters)
	if err != nil {
		return nil, err
	}

	restReq, err := http.NewRequest(reqOpt.Method, fullUrl.String(), reqOpt.ReqDTO)
	if err != nil {
		return nil, err
	}

	restReq.Header.Add("Content-Type", "application/json")

	for k, v := range reqOpt.CommonReqParams.Headers {
		restReq.Header.Set(k, v)
	}

	restResp, err := c.HTTPClient.Do(restReq)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(restResp.Body)
	defer restResp.Body.Close()

	if restResp.StatusCode >= 400 {
		err = errors.New(string(respBody))
		return respBody, err
	}

	if err != nil {
		return nil, err
	}

	return respBody, err
}
