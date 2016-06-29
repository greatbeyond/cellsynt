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

import . "gopkg.in/check.v1"

var _ = Suite(&MessageSuite{})

type MessageSuite struct{}

func (suite *MessageSuite) SetUpSuite(c *C) {}

func (suite *MessageSuite) SetUpTest(c *C) {}

func (suite *MessageSuite) TearDownTest(c *C) {}

// -------------------------------------------------------------
// Reciptient

func (suite *MessageSuite) Test_Reciptient_Normal(c *C) {
	r := &Reciptient{
		Destinations: []string{"0703112233"},
		CountryCode:  "46",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination": "0046703112233",
	})
}

func (suite *MessageSuite) Test_Reciptient_Multiple(c *C) {
	r := &Reciptient{
		Destinations: []string{
			"0046703112233",
			"+46703112233",
			"0703112233",
		},
		CountryCode: "46",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination": "0046703112233,0046703112233,0046703112233",
	})
}

// -------------------------------------------------------------
// Options

func (suite *MessageSuite) Test_Options_Normal(c *C) {
	r := &Options{
		OriginatorType: OriginatorTypeAlpha,
		Originator:     "test",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"originatortype": "alpha",
		"originator":     "test",
	})
}

func (suite *MessageSuite) Test_Options_Omitted(c *C) {
	r := &Options{
		Originator: "test",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"originator": "test",
	})
}

// -------------------------------------------------------------
// Text message

func (suite *MessageSuite) Test_TextMessage_Full(c *C) {
	r := &TextMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},
		Text:        "test",
		Charset:     CharsetUTF8,
		AllowConcat: true,
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination":    "0046703112233",
		"type":           "text",
		"text":           "test",
		"charset":        "UTF-8",
		"originatortype": "alpha",
		"originator":     "test",
		"allowconcat":    "6",
	})
}

func (suite *MessageSuite) Test_TextMessage_Minimal(c *C) {
	r := &TextMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},
		Text: "test åäö",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination": "0046703112233",
		"type":        "text",
		"text":        "test+%C3%A5%C3%A4%C3%B6",
	})
}

// -------------------------------------------------------------
// binary emssage

func (suite *MessageSuite) Test_BinaryMessage_Normal(c *C) {
	r := &BinaryMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
		UDH:    []byte("AABBCC001122"),
		Binary: []byte("334455FF"),
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination":    "0046703112233",
		"type":           "binary",
		"originatortype": "alpha",
		"originator":     "test",
		"udh":            "AABBCC001122",
		"data":           "334455FF",
	})
}

func (suite *MessageSuite) Test_BinaryMessage_Minimal(c *C) {
	r := &BinaryMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},

		Binary: []byte("334455FF"),
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination": "0046703112233",
		"type":        "binary",
		"data":        "334455FF",
	})
}

// -------------------------------------------------------------
// flash emssage

func (suite *MessageSuite) Test_FlashMessage_Normal(c *C) {
	r := &FlashMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
		Text:        "test",
		Charset:     CharsetUTF8,
		AllowConcat: true,
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination":    "0046703112233",
		"type":           "flash",
		"originatortype": "alpha",
		"originator":     "test",
		"text":           "test",
		"charset":        "UTF-8",
		"allowconcat":    "6",
	})
}

func (suite *MessageSuite) Test_FlashMessage_Minimal(c *C) {
	r := &FlashMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},
		Text: "test",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination": "0046703112233",
		"type":        "flash",
		"text":        "test",
	})
}

// -------------------------------------------------------------
// unicode emssage

func (suite *MessageSuite) Test_UnicodeMessage_Normal(c *C) {
	r := &UnicodeMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},
		Options: &Options{
			OriginatorType: OriginatorTypeAlpha,
			Originator:     "test",
		},
		Charset:     CharsetUTF8,
		AllowConcat: true,
		Text:        "Ελλάδα",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination":    "0046703112233",
		"type":           "unicode",
		"originatortype": "alpha",
		"originator":     "test",
		"allowconcat":    "6",
		"text":           "%CE%95%CE%BB%CE%BB%CE%AC%CE%B4%CE%B1",
		"charset":        "UTF-8",
	})
}

func (suite *MessageSuite) Test_UnicodeMessage_Minimal(c *C) {
	r := &UnicodeMessage{
		Reciptient: &Reciptient{
			Destinations: []string{"0046703112233"},
		},

		Text: "Ελλάδα",
	}
	c.Assert(r.GetParameters(), DeepEquals, map[string]string{
		"destination": "0046703112233",
		"type":        "unicode",
		"text":        "%CE%95%CE%BB%CE%BB%CE%AC%CE%B4%CE%B1",
	})
}
