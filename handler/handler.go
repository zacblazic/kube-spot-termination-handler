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
	"time"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"

	log "github.com/sirupsen/logrus"
)

// TerminationHandler watches instance metadata for spot termination notices and
// attempts to handle them gracefully.
type TerminationHandler struct {
	client      *ec2metadata.EC2Metadata
	interval    time.Duration
	terminating chan bool
}

// NewTerminationHandler creates a new termination handler.
func NewTerminationHandler(client *ec2metadata.EC2Metadata, interval time.Duration) *TerminationHandler {
	if client == nil {
		session := session.Must(session.NewSession())
		client = ec2metadata.New(session)
	}

	return &TerminationHandler{
		client:      client,
		interval:    interval,
		terminating: make(chan bool),
	}
}

// Start begins the termination handler's main control loop.
func (th *TerminationHandler) Start() {
	log.Info("Starting handler...")

	go th.watch()
	th.handle()
}

func (th *TerminationHandler) watch() {
	log.Info("Starting watcher...")

	for {
		log.Debug("Checking if instance is terminating...")

		if th.isTerminating() {
			log.Info("Instance is terminating, notifying channel")
			th.terminating <- true
			break
		}

		log.Debug("Instance is not terminating.")
		log.Debugf("Sleeping for %v...", th.interval)
		time.Sleep(th.interval)
	}
}

func (th *TerminationHandler) handle() {
	<-th.terminating

	log.Info("Notification received, preparing for termination...")
}

func (th *TerminationHandler) isTerminating() bool {
	value, err := th.client.GetMetadata("spot/termination-time")
	if err != nil {
		return false
	}

	terminationTime, _ := time.Parse(time.RFC3339, value)

	if terminationTime.After(time.Now()) {
		return true
	}

	return false
}
