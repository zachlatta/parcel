package parcel

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

// Capabilities describes the various capabilities of this server.
//
// The spec for this isn't complete and the types of the fields are guessed.
// See the bottom of http://jmap.io/spec.html#accounts for details.
type Capabilities struct {
	DelayedSend          bool  `json:"delayedSend"`
	MaxSendSize          int64 `json:"maxSendSize"`
	SupportsMultiMailbox bool  `json:"supportsMultiMailbox"`
	SupportsThreads      bool  `json:"supportsThreads"`
}

// Account represents an account on the JMAP server. It provides access to
// mail, contacts, and calendars.
//
// A single login may be used to access multiple accounts.
//
// All IDs are only unique within an account. IDs may clash across accounts.
//
// http://jmap.io/spec.html#accounts
type Account struct {
	// ID uniquely identifies the account.
	ID string `json:"id"`

	// Name is a user-friendly string to show when presenting content from this
	// account. e.g. the email address of the account.
	Name string `json:"name"`

	// IsPrimary is whether the account is the primary account. This must be true
	// for exactly one of the accounts returned.
	IsPrimary bool `json:"isPrimary"`

	// IsReadyOnly is true if the user has read-only access to this account. The
	// user may not use the `set` methods with this account.
	IsReadOnly bool `json:"isReadOnly"`

	// HasMail represents whether this account contains mail data. Clients may
	// use the Mail methods with this account.
	HasMail bool `json:"hasMail"`

	// HasContacts represents whether this account contains contact data. Clients
	// may use the Contacts methods with this account.
	HasContacts bool `json:"hasContacts"`

	// HasCalendars represents whether this account contains calendar data.
	// Clients may use the Calendar methods with this account.
	HasCalendars bool `json:"hasCalendars"`

	// Capabilities describes the various capabilities of this server.
	Capabilities Capabilities `json:"capabilities"`
}

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
