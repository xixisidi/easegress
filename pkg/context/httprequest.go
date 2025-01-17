/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package context

import (
	"io"
	"net/http"

	"github.com/tomasen/realip"

	"github.com/megaease/easegress/pkg/util/callbackreader"
	"github.com/megaease/easegress/pkg/util/httpheader"
)

type (
	httpRequest struct {
		std       *http.Request
		header    *httpheader.HTTPHeader
		body      *callbackreader.CallbackReader
		bodyCount int
		metaSize  int
		realIP    string
	}
)

func newHTTPRequest(stdr *http.Request) *httpRequest {
	// Reference: https://golang.org/pkg/net/http/#Request
	//
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	if stdr.Header.Get("Host") == "" {
		stdr.Header.Set("Host", stdr.Host)
	}

	hq := &httpRequest{
		std:    stdr,
		header: httpheader.New(stdr.Header),
		body:   callbackreader.New(stdr.Body),
		realIP: realip.FromRequest(stdr),
	}

	// NOTE: Always count original body, even the body could be changed
	// by SetBody().
	hq.body.OnAfter(func(num int, p []byte, n int, err error) ([]byte, int, error) {
		hq.bodyCount += n
		return p, n, err
	})

	// Reference: https://tools.ietf.org/html/rfc2616#section-5
	//
	// meta length is the length of:
	// w.stdr.Method + " "
	// + stdr.URL.RequestURI() + " "
	// + stdr.Proto + "\r\n",
	// + w.Header().Dump() + "\r\n\r\n"
	//
	// but to improve performance, we won't build this string

	hq.metaSize += len(stdr.Method) + 1
	hq.metaSize += len(stdr.URL.RequestURI()) + 1
	hq.metaSize += len(stdr.Proto) + 2
	hq.metaSize += hq.Header().Length() + 4

	return hq
}

func (r *httpRequest) RealIP() string {
	return r.realIP
}

func (r *httpRequest) Method() string {
	return r.std.Method
}

func (r *httpRequest) SetMethod(method string) {
	r.std.Method = method
}

func (r *httpRequest) Scheme() string {
	if scheme := r.std.URL.Scheme; scheme != "" {
		return scheme
	}

	if scheme := r.std.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}

	if r.std.TLS != nil {
		return "https"
	}

	return "http"
}

func (r *httpRequest) Host() string {
	return r.std.Host
}

func (r *httpRequest) SetHost(host string) {
	r.std.Host = host
}

func (r *httpRequest) Path() string {
	return r.std.URL.Path
}

func (r *httpRequest) SetPath(path string) {
	r.std.URL.Path = path
}

func (r *httpRequest) EscapedPath() string {
	return r.std.URL.EscapedPath()
}

func (r *httpRequest) Query() string {
	return r.std.URL.RawQuery
}

func (r *httpRequest) SetQuery(query string) {
	r.std.URL.RawQuery = query
}

func (r *httpRequest) Fragment() string {
	return r.std.URL.Fragment
}

func (r *httpRequest) Proto() string {
	return r.std.Proto
}

func (r *httpRequest) Header() *httpheader.HTTPHeader {
	return r.header
}

func (r *httpRequest) Cookie(name string) (*http.Cookie, error) {
	return r.std.Cookie(name)
}

func (r *httpRequest) Cookies() []*http.Cookie {
	return r.std.Cookies()
}

func (r *httpRequest) AddCookie(cookie *http.Cookie) {
	r.std.AddCookie(cookie)
}

func (r *httpRequest) Body() io.Reader {
	return r.body
}

func (r *httpRequest) SetBody(reader io.Reader, closePreviousReader bool) {
	r.body.SetReader(reader, closePreviousReader)
}

func (r *httpRequest) Size() uint64 {
	return uint64(r.metaSize + r.bodyCount)
}

func (r *httpRequest) finish() {
	// NOTE: We don't use this line in case of large flow attack.
	// io.Copy(io.Discard, r.std.Body)

	// NOTE: The server will do it for us.
	// r.std.Body.Close()

	r.body.Close()
}

func (r *httpRequest) Std() *http.Request {
	return r.std
}
