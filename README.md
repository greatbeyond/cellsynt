# cellsynth
Cellsynth API lib in Go

## Create a client
A client is used when communicating with Cellsynth.
```
client := cellsynt.NewClient(username, username, senderName),
```

### Send a text
```
textMsg := &cellsynt.TextMessage{
    BaseMessage: &cellsynt.BaseMessage{
        Destinations: reciptients,
    },
    Text:    message.Body,
    Charset: cellsynt.CharsetUTF8,
}
_, err = client.SendMessage(textMsg)
if err != nil {
    ...
}
```