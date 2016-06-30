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
	"net/http"
	"testing"

	t "github.com/greatbeyond/cellsynt/testing"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&CellsyntSuite{})

type CellsyntSuite struct {
	client *Client
	server *t.MockServer
}

func (suite *CellsyntSuite) SetUpSuite(c *C) {
	apiURL = "http://mock.cellsynt.net/sms.php"
}

func (suite *CellsyntSuite) SetUpTest(c *C) {
	suite.server = t.NewMockServer()
	suite.server.SetChecker(c)
	http.DefaultClient = suite.server.HTTPClient

	suite.client = NewClient("username", "password", "sendername")
}

func (suite *CellsyntSuite) TearDownTest(c *C) {
	suite.server.VerifyNoMoreRequests(c)
	suite.server.Close()
	suite.server = nil
}

// -------------------------------------------------------------
// Parameters

func (suite *CellsyntSuite) Test_Client_getParameters_Normal(c *C) {
	params := suite.client.getParameters()
	c.Assert(params, DeepEquals, map[string]string{
		"username":       "username",
		"password":       "password",
		"originatortype": "alpha",
		"originator":     "sendername",
		"charset":        "UTF-8",
		"allowconcat":    "6",
	})
}

func (suite *CellsyntSuite) Test_Client_messageParameters_Normal(c *C) {

	r := &TextMessage{
		Destination: &Destination{
			Recipients: []string{"0046703112233"},
		},
		Text:        "test",
		Charset:     CharsetUTF8,
		AllowConcat: true,
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
	}

	parameters := suite.client.messageParameters(r)
	c.Assert(parameters, Equals, `allowconcat=6&charset=UTF-8&destination=0046703112233&originator=test&originatortype=alpha&password=password&text=test&type=text&username=username`)
}

func (suite *CellsyntSuite) Test_Client_messageParameters_Override(c *C) {

	r := &TextMessage{
		Destination: &Destination{
			Recipients: []string{"0046703112233"},
		},
		Text:        "test",
		Charset:     CharsetISO88591,
		AllowConcat: true,
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
	}

	parameters := suite.client.messageParameters(r)
	c.Assert(parameters, Equals, `allowconcat=6&charset=ISO-8859-1&destination=0046703112233&originator=test&originatortype=alpha&password=password&text=test&type=text&username=username`)
}

// -------------------------------------------------------------
// Response handling

func (suite *CellsyntSuite) Test_Client_handleResponse_OK(c *C) {
	rb := []byte("OK: 92dff27302b754424242fb204620dc18")
	response, err := suite.client.handleResponse(rb)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, &Response{
		Success: true,
		TrackingIDs: []string{
			"92dff27302b754424242fb204620dc18",
		},
	})
}

func (suite *CellsyntSuite) Test_Client_handleResponse_OK_Multiple(c *C) {
	rb := []byte("OK: de8c4a032fb45ae65ab9e349a8dc2458,ed6037d0fe08dd4a4ab5cdcfd5aae653,6a351ae2ef03c3c5e271adcccd140089")
	response, err := suite.client.handleResponse(rb)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, &Response{
		Success: true,
		TrackingIDs: []string{
			"de8c4a032fb45ae65ab9e349a8dc2458",
			"ed6037d0fe08dd4a4ab5cdcfd5aae653",
			"6a351ae2ef03c3c5e271adcccd140089",
		},
	})
}

func (suite *CellsyntSuite) Test_Client_handleResponse_Failure(c *C) {
	rb := []byte("Error: Parameter destination must be set")
	response, err := suite.client.handleResponse(rb)

	c.Assert(err, ErrorMatches, "Parameter destination must be set")
	c.Assert(response, IsNil)
}

func (suite *CellsyntSuite) Test_Client_handleResponse_BadResponse(c *C) {
	rb := []byte("unexpected response")
	response, err := suite.client.handleResponse(rb)

	c.Assert(err, ErrorMatches, "response error: unexpected response")
	c.Assert(response, IsNil)
}

// -------------------------------------------------------------
// Sending

func (suite *CellsyntSuite) Test_Client_SendMessage_MissingDestination_1(c *C) {
	r := &TextMessage{
		Text: "test",
	}
	_, err := suite.client.SendMessage(r)
	c.Assert(err, ErrorMatches, "message has no destination set")
}

func (suite *CellsyntSuite) Test_Client_SendMessage_MissingDestination_2(c *C) {
	r := &TextMessage{
		Destination: &Destination{
			Recipients: []string{},
		},
		Text: "test",
	}
	_, err := suite.client.SendMessage(r)
	c.Assert(err, ErrorMatches, "message has no destination set")
}

func (suite *CellsyntSuite) Test_Client_SendMessage_Normal(c *C) {

	r := &TextMessage{
		Destination: &Destination{
			Recipients: []string{"0046703112233"},
		},
		Text:        "test",
		Charset:     CharsetUTF8,
		AllowConcat: true,
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
	}

	suite.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   "OK: de8c4a032fb45ae65ab9e349a8dc2458",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, Equals, apiURL)
			c.Assert(body, Equals, "allowconcat=6&charset=UTF-8&destination=0046703112233&originator=test&originatortype=alpha&password=password&text=test&type=text&username=username")
		},
	})

	response, err := suite.client.SendMessage(r)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, &Response{
		Success: true,
		TrackingIDs: []string{
			"de8c4a032fb45ae65ab9e349a8dc2458",
		},
	})
}

func (suite *CellsyntSuite) Test_Client_SendMessage_Error(c *C) {
	r := &TextMessage{
		Destination: &Destination{
			Recipients: []string{"0046703112233"},
		},
		Text: "test",
	}

	suite.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   501,
		Body:   "Error: mocked error",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, Equals, apiURL)
			c.Assert(body, Equals, "allowconcat=6&charset=UTF-8&destination=0046703112233&originator=sendername&originatortype=alpha&password=password&text=test&type=text&username=username")
		},
	})

	response, err := suite.client.SendMessage(r)

	c.Assert(err, ErrorMatches, "mocked error")
	c.Assert(response, IsNil)
}
