package main

import "encoding/json"

// Message represents an individual message for processing by the server or
// returning to the client.
//
// See http://jmap.io/spec.html#the-structure-of-an-exchange
type Message struct {
	// Name specifies the method to be called on the server or the type of
	// response being sent to the client.
	Name string

	// Arguments is an object containing named arguments for the method or
	// response.
	Arguments interface{}

	// ClientID is an arbitary string to be echoed back with the responses
	// emitted bythe method call.
	ClientID string
}

func (m *Message) MarshalJSON() ([]byte, error) {
	arr := make([]interface{}, 3)
	arr[0] = m.Name
	arr[1] = m.Arguments
	arr[2] = m.ClientID

	return json.Marshal(arr)
}

func (m *Message) UnmarshalJSON(j []byte) error {
	arr := []interface{}{}
	if err := json.Unmarshal(j, &arr); err != nil {
		return err
	}

	m.Name = arr[0].(string)
	m.Arguments = arr[1]
	m.ClientID = arr[2].(string)

	return nil
}
