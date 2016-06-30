# cellsynt
Cellsynt API lib in Go

[![wercker status](https://app.wercker.com/status/c215f070f60e91325c5fd6e8bd53c1c5/s "wercker status")](https://app.wercker.com/project/bykey/c215f070f60e91325c5fd6e8bd53c1c5) [![Coverage Status](https://coveralls.io/repos/github/greatbeyond/cellsynt/badge.svg?branch=master)](https://coveralls.io/github/greatbeyond/cellsynt?branch=master) 

## Create a client
A client is used when communicating with Cellsynt.
```
client := cellsynt.NewClient(username, username, senderName),
```

You can create the client directly with some non-default values:
```
client := &cellsynt.Client{
    Username:       username,
    Password:       password,
    OriginatorType: OriginatorTypeAlpha,
    Originator:     senderName,
    Charset:        CharsetUTF8,
    AllowConcat:    true,
}
```

### Send a text
```
textMsg := &cellsynt.TextMessage{
    Destination: &cellsynt.Destination{
        Recipients: []string{"555-123-45"},
    },
    Text:    message.Body,
    Charset: cellsynt.CharsetUTF8,
}
_, err = client.SendMessage(textMsg)
if err != nil {
    ...
}
```

Override client options by including them in the message:
```
textMsg := &cellsynt.TextMessage{
    Destination: &cellsynt.Destination{
        Recipients: []string{"555-123-45"},
    },
    Options: &cellsynt.Options{
        OriginatorType: OriginatorTypeNumeric,
        Originator:     "0703112233",
    },
    Text:    message.Body,
}
_, err = client.SendMessage(textMsg)
```

A tracking ID is returned for each destination
```
response, _ = client.SendMessage(textMsg)
print(response)
```
```
&Response{
    Success: true,
    TrackingIDs: []string{
        "de8c4a032fb45ae65ab9e349a8dc2458",
    },
}
```