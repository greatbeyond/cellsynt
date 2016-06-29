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
	"net/url"
	"strings"
)

// Message is a generic interface to all types of messages
type Message interface {
	Destination() string
	Type() string
	GetParameters() map[string]string
}

type Options struct {
	// Optional
	OriginatorType OriginatorType
	Originator     string
}

func (b *Options) GetParameters() map[string]string {
	params := map[string]string{
		"originatortype": string(b.OriginatorType),
		"originator":     b.Originator,
	}
	return clearEmptyParams(params)
}

// Reciptient contains default values that all message types share.
// Theese values can be omitted if you want to use the client default
type Reciptient struct {
	// Required
	Destinations []string

	// Optional, can be included in destination
	CountryCode string
}

// Destination is the Destination address(es) formatted for cellsynt
func (b *Reciptient) Destination() string {

	if b == nil {
		return ""
	}

	phones := []string{}
	for _, phone := range b.Destinations {
		if strings.HasPrefix(phone, "+") {
			phone = "00" + strings.TrimPrefix(phone, "+")
		} else if !strings.HasPrefix(phone, "00") {
			phone = "00" + b.CountryCode + strings.TrimLeft(phone, "0")
		}
		phones = append(phones, phone)
	}

	return strings.Join(phones, ",")
}

func (b *Reciptient) GetParameters() map[string]string {
	params := map[string]string{
		"destination": b.Destination(),
	}
	return clearEmptyParams(params)
}

// TextMessage is used to send a normal text message. Maximum number of characters is 160. Any characters
// specified within character coding GSM 03.38 can be used (e.g. English, Swedish, Norwegian), for
// other languages / alphabets (e.g. Arabic, Japanese, Chinese) please use Unicode.
type TextMessage struct {
	// Required
	Text string

	// Optional
	Charset     Charset
	AllowConcat bool

	*Reciptient
	*Options
}

// Type returns the message type
func (m *TextMessage) Type() string { return "text" }

// GetParameters implements Message interface
func (m *TextMessage) GetParameters() map[string]string {
	params := map[string]string{
		"type":        m.Type(),
		"text":        url.QueryEscape(m.Text),
		"charset":     string(m.Charset),
		"allowconcat": ternaryStr(m.AllowConcat, "6", ""),
	}
	params = mergeParams(params, m.Reciptient.GetParameters())
	if m.Options != nil {
		params = mergeParams(params, m.Options.GetParameters())
	}

	return clearEmptyParams(params)
}

// BinaryMessage  can be used to send settings, bookmarks, visiting cards etc. See relevant phone
// manufacturer's manual for further specifications on message formats available.
// Two parameters are needed to send the binary data: udh and data.
// Parameters can be set individually, it is not mandatory to use both, however, at least one of them
type BinaryMessage struct {
	Binary []byte
	UDH    []byte

	*Reciptient
	*Options
}

// Type returns the message type
func (m *BinaryMessage) Type() string { return "binary" }

// GetParameters implements Message interface
func (m *BinaryMessage) GetParameters() map[string]string {
	params := map[string]string{
		"type": m.Type(),
		"data": string(m.Binary),
		"udh":  string(m.UDH),
	}
	params = mergeParams(params, m.Reciptient.GetParameters())
	if m.Options != nil {
		params = mergeParams(params, m.Options.GetParameters())
	}
	return clearEmptyParams(params)
}

// FlashMessage is used to send a normal text message which is shown directly on screen instead of being saved in
// the inbox on the recipient's mobile phone ("flash message").
// Please note! Support for flash messages can vary depending on mobile phone model, operator
// network and other external factors. If it is not possible to send a mess
type FlashMessage struct {
	// Required
	Text string

	// Optional
	Charset     Charset
	AllowConcat bool

	*Reciptient
	*Options
}

// Type returns the message type
func (m *FlashMessage) Type() string { return "flash" }

// GetParameters implements Message interface
func (m *FlashMessage) GetParameters() map[string]string {
	params := map[string]string{
		"type":        m.Type(),
		"text":        url.QueryEscape(m.Text),
		"charset":     string(m.Charset),
		"allowconcat": ternaryStr(m.AllowConcat, "6", ""),
	}
	params = mergeParams(params, m.Reciptient.GetParameters())
	if m.Options != nil {
		params = mergeParams(params, m.Options.GetParameters())
	}
	return clearEmptyParams(params)
}

// UnicodeMessage can be used if you need to send characters not available within an ordinary
// text message (e.g. Arabic, Japanese). A Unicode SMS can contain maximum 70 characters per
// message (or 67 characters per part for a long SMS).
// Unicode messages are sent by setting parameter type to unicode. When sending Unicode
// messages, the parameter charset can be set to UTF-8 to send data in the UTF-8 character set
// instead which holds characters of most of the world's languages.
type UnicodeMessage struct {
	// Required
	Text string

	// Optional
	Charset     Charset
	AllowConcat bool

	*Reciptient
	*Options
}

// Type returns the message type
func (m *UnicodeMessage) Type() string { return "unicode" }

// GetParameters implements Message interface
func (m *UnicodeMessage) GetParameters() map[string]string {
	params := map[string]string{
		"type":        m.Type(),
		"text":        url.QueryEscape(m.Text),
		"charset":     string(m.Charset),
		"allowconcat": ternaryStr(m.AllowConcat, "6", ""),
	}
	params = mergeParams(params, m.Reciptient.GetParameters())
	if m.Options != nil {
		params = mergeParams(params, m.Options.GetParameters())
	}
	return clearEmptyParams(params)
}
