// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httputil

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

// NewClient returns a http.Client using the specified http.RoundTripper.
func NewClient(rt http.RoundTripper) *http.Client {
	return &http.Client{Transport: rt}
}

// NewDeadlineClient returns a new http.Client which will time out long running
// requests.
func NewDeadlineClient(timeout time.Duration, proxyURL *url.URL) *http.Client {
	return NewClient(NewDeadlineRoundTripper(timeout, proxyURL))
}

// NewDeadlineRoundTripper returns a new http.RoundTripper which will time out
// long running requests.
func NewDeadlineRoundTripper(timeout time.Duration, proxyURL *url.URL) http.RoundTripper {
	return &http.Transport{
		// Set proxy (if null, then becomes a direct connection)
		Proxy: http.ProxyURL(proxyURL),
		// We need to disable keepalive, because we set a deadline on the
		// underlying connection.
		DisableKeepAlives: true,
		Dial: func(netw, addr string) (c net.Conn, err error) {
			start := time.Now()

			c, err = net.DialTimeout(netw, addr, timeout)

			if err == nil {
				c.SetDeadline(start.Add(timeout))
			}

			return
		},
	}
}

type bearerAuthRoundTripper struct {
	bearerToken string
	rt          http.RoundTripper
}

// NewBearerAuthRoundTripper adds the provided bearer token to a request unless the authorization
// header has already been set.
func NewBearerAuthRoundTripper(bearer string, rt http.RoundTripper) http.RoundTripper {
	return &bearerAuthRoundTripper{bearer, rt}
}

func (rt *bearerAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) == 0 {
		req = cloneRequest(req)
		req.Header.Set("Authorization", "Bearer "+rt.bearerToken)
	}

	return rt.rt.RoundTrip(req)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// Shallow copy of the struct.
	r2 := new(http.Request)
	*r2 = *r
	// Deep copy of the Header.
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
