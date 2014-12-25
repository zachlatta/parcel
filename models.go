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

// Mailbox represents a named set of emails. This is the primary mechanism for
// organizing messages within an account. It is analogous to a folder in IMAP
// or a label in other systems. A mailbox may perform a certain role in the
// system.
//
// http://jmap.io/spec.html#mailboxes
type Mailbox struct {
	// ID is an immutable identifier for the mailbox.
	ID string `json:"id"`

	// Name is the user-visible name for the mailbox, e.g. "Inbox". This may be
	// any UTF-8 string of at least 1 character in length and maximum 256 bytes
	// in size. Mailboxes may have the same name as a sibling mailbox.
	Name string `json:"name"`

	// ParentID is the mailbox id for the parent of this mailbox, or nil if this
	// mailbox is at the top level. Mailboxes form acyclic graphs and, therefore,
	// must not loop.
	ParentID *string `json:"parentId"`

	// Role identifies system mailboxes. This property is immutable.
	//
	// The following values should be used for the relevant mailboxes:
	//
	//     inbox - the mailbox to which new mail is delivered by default, unless
	//             diverted by a rule or spam filter, etc.
	//   archive - messages the user does not need right now, but does not wish
	//             to delete.
	//    drafts - messages the user is currently writing and are not yet sent.
	//    outbox - messages the user has finished writing and wishes to send.
	//      sent - messages the user has sent.
	//     trash - messages the user has deleted.
	//      spam - messages considered spam by the server.
	// templates - drafts which should be used as templates (i.e. used as the
	//             basis for creating new drafts).
	//
	// If the mailbox's role is trash, then it must be treated specially:
	//
	//  * Messages in the Trash are ignored when calculating the unreadThreads
	//    and totalThreads count of other mailboxes.
	//  * Messages not in the Trash are ignored when calculating the
	//    unreadThreads and totalThreads count for the Trash folder.
	Role *string `json:"role"`

	// MustBeOnlyMailbox is whether messages in this mailbox may be in other
	// mailboxes as well.
	MustBeOnlyMailbox bool `json:"mustBeOnlyMailbox"`

	// MayAddMessages represents whether the user may add messages to this
	// mailbox (by either creating a new message or modifying an existing one).
	MayAddMessages bool `json:"mayAddMessages"`

	// MayRemoveMessages represents whether the user may remove messages from
	// this mailbox (by either changing the mailboxes of a message or deleting
	// it).
	MayRemoveMessages bool `json:"mayRemoveMessages"`

	// MayCreateChlid represents whether the user may create a mailbox with this
	// mailbox as its parent.
	MayCreateChild bool `json:"mayCreateChild"`

	// MayRenameMailbox represents whether the user may rename the mailbox or
	// make it as a child of another mailbox.
	MayRenameMailbox bool `json:"mayRenameMailbox"`

	// MayDeleteMailbox represents whether the user may delete the mailbox
	// itself.
	MayDeleteMailbox bool `json:"mayDeleteMailbox"`

	// TotalMessages is the number of messages in this mailbox.
	TotalMessages int64 `json:"totalMessages"`

	// UnreadMessages is the number of messages in this mailbox where the
	// isUnread property of the message is set to true and the isDraft property
	// is false.
	UnreadMessages int64 `json:"unreadMessages"`

	// TotalThreads is the number of threads where at least one message in the
	// thread is in this mailbox.
	TotalThreads int64 `json:"totalThreads"`

	// UnreadThreads is the number of threads where at least one message in the
	// thread has the isUnread property set to true and the isDraft property set
	// to false and at least one message in the thread is in the mailbox. Note,
	// the unread message does not need to be the one in this mailbox.
	UnreadThreads int64 `json:"unreadThreads"`
}

// Thread represents a grouping of replies with the original message. It is
// simply a flat list of messages, ordered by date. Every message belongs to a
// thread, even if it's the only message in the thread.
type Thread struct {
	// ID identifies the thread.
	ID string `json:"id"`

	// MessageIDs are the ids of the messages in the thread, sorted such that:
	//
	// * Any message with isDraft set to true and an inReplyToMessageId property
	//   that corresponds to another message in the thread comes immediately
	//   after that message in the sort order.
	// * Other than that, everything is sorted in date order (the same as the
	//   date property on the Message object), oldest first.
	MessagesIDs []string `json:"messageIds"`
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
