// Copyright © 2018 Zac Blazic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestTerminationHandler_isTerminatingWithFutureDate(t *testing.T) {
	defer gock.Off()

	future := time.Now().Add(time.Minute * 2)

	gock.New("http://169.254.169.254/latest").
		Get("/meta-data/spot/termination-time").
		Reply(http.StatusOK).
		BodyString(future.Format(time.RFC3339))

	handler := NewTerminationHandler(nil, time.Second*30)
	terminating := handler.isTerminating()

	assert.True(t, terminating)
}

func TestTerminationHandler_isTerminatingWithPastDate(t *testing.T) {
	defer gock.Off()

	gock.New("http://169.254.169.254/latest").
		Get("/meta-data/spot/termination-time").
		Reply(http.StatusOK).
		BodyString("2015-01-05T18:02:00Z")

	handler := NewTerminationHandler(nil, time.Second*30)
	terminating := handler.isTerminating()

	assert.False(t, terminating)
}

func TestTerminationHandler_isTerminatingWhenNotTerminating(t *testing.T) {
	defer gock.Off()

	gock.New("http://169.254.169.254/latest").
		Get("/meta-data/spot/termination-time").
		Reply(http.StatusNotFound).
		BodyString("404 - Not Found")

	handler := NewTerminationHandler(nil, time.Second*30)
	terminating := handler.isTerminating()

	assert.False(t, terminating)
}
