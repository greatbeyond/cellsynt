// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Written by David Högborg <d@greatbeyond.se>, 2016
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cellsynt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Client holds username and password, and default values for messages
// The values on the client are default values that can be overridden
// by a message with different value for a field.
type Client struct {
	// Required
	Username string
	Password string

	// Default values, can be overridden by message values.
	OriginatorType     OriginatorType
	Originator         string
	Charset            Charset
	AllowConcat        bool
	DefaultCountryCode string
}

// Response will contain a success flag and the tracking ids that can
// be used for status tracking messages
type Response struct {
	Success     bool
	TrackingIDs []string
}

// NewClient returns a new client instance with some defaults set
// SenderName is the originator, alpha numeric string by default
func NewClient(username, password string, senderName string) *Client {
	return &Client{
		Username:       username,
		Password:       password,
		OriginatorType: OriginatorTypeAlpha,
		Originator:     senderName,
		Charset:        CharsetUTF8,
		AllowConcat:    true,
	}
}

func (c *Client) getParameters(msgParams map[string]string) map[string]string {
	return map[string]string{
		"username":       c.Username,
		"password":       c.Password,
		"originatortype": string(c.OriginatorType),
		"originator":     c.Originator,
		"charset":        string(c.Charset),
		"allowconcat":    ternaryStr(c.AllowConcat, "true", "false"),
	}
}

// SendMessage dispatches a message to the destination
func (c *Client) SendMessage(message Message) (*Response, error) {

	// get the message parameters
	params := message.GetParameters()
	// get the client parameters
	clientParams := c.getParameters(params)

	// place the client params in the msg params, giving priority to the message
	for k, v := range clientParams {
		if _, ok := params[k]; !ok {
			params[k] = v
		}
	}

	parts := []string{}
	for k, v := range params {
		if v != "" {
			parts = append(parts, k+"="+v)
		}
	}
	paramstr := strings.Join(parts, "&")
	body := bytes.NewBuffer([]byte(paramstr))

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", body)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response, err := c.handleResponse(responseData)
	if err != nil {
		log.WithFields(log.Fields{
			"destination": message.Destination(),
			"error":       err.Error(),
			"type":        message.Type(),
		}).Error("error sending message", caller())
		return nil, err
	}

	log.WithFields(log.Fields{
		"type":         message.Type(),
		"destination":  message.Destination(),
		"tracking_ids": response.TrackingIDs,
	}).Info("sent message")

	return response, nil
}

func (c *Client) handleResponse(respBytes []byte) (*Response, error) {

	respStr := string(respBytes)
	if strings.HasPrefix("OK: ", respStr) {
		idstr := strings.TrimPrefix(respStr, "OK: ")
		return &Response{
			Success:     true,
			TrackingIDs: strings.Split(idstr, ","),
		}, nil
	}

	if strings.HasPrefix("Error: ", respStr) {
		return nil, fmt.Errorf("%s", strings.TrimPrefix("Error: ", respStr))
	}

	return nil, fmt.Errorf("response error")
}
