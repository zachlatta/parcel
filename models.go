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

	// Arguments is map containing named arguments for the method or response.
	Arguments map[string]interface{}

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
	m.Arguments = arr[1].(map[string]interface{})
	m.ClientID = arr[2].(string)

	return nil
}

// Login provides access to one or multiple accounts.
type Login struct{}

// Account represents an account on the JMAP server. It provides access to
// mail, contacts, and calendars.
//
// A single login may be used to access multiple accounts.
//
// All IDs are only unique within an account. IDs may clash across accounts.
//
// http://jmap.io/spec.html#accounts
type Account struct{}

// MailMessage represents a mail message. A MailMessage is immutable except for
// the boolean `isXXX` status properties and the set of mailboxes it is in.
type MailMessage struct {
	// A unique, immutable ID that does not change if the message changes
	// mailboxes.
	ID string

	// Has the email not yet been read?
	IsUnread bool

	// Has the email been flagged (starred, or pinned)?
	IsFlagged bool

	// Is the email a draft?
	IsDraft bool

	// Has the email been answered (replied to)?
	IsAnswered bool
}
