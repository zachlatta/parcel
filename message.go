package main

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
