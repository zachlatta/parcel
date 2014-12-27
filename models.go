package parcel

import (
	"encoding/json"
	"time"
)

// ExchangeMessage represents an individual message for processing by the
// server or returning to the client.
//
// See http://jmap.io/spec.html#the-structure-of-an-exchange
type ExchangeMessage struct {
	// Name specifies the method to be called on the server or the type of
	// response being sent to the client.
	Name string

	// Arguments is map containing named arguments for the method or response.
	Arguments map[string]interface{}

	// ClientID is an arbitary string to be echoed back with the responses
	// emitted bythe method call.
	ClientID string
}

func (m *ExchangeMessage) MarshalJSON() ([]byte, error) {
	arr := make([]interface{}, 3)
	arr[0] = m.Name
	arr[1] = m.Arguments
	arr[2] = m.ClientID

	return json.Marshal(arr)
}

func (m *ExchangeMessage) UnmarshalJSON(j []byte) error {
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
//
// http://jmap.io/spec.html#threads
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

// Emailer represents a sender or recipient of an email.
type Emailer struct {
	// Name represents the name of the sender/recipient. If a name cannot be
	// extracted for an email, this property should be an empty string.
	Name string `json:"name"`

	// Email is the email address of the sender/recipient.
	Email string `json:"email"`
}

// Attachment represents an email attachment.
type Attachment struct {
	// ID identifies an attachment. This is only unique when paired with the
	// message id and has no meaning without reference to that.
	ID string `json:"id"`

	// URL links to the attachment. The HTTP request must be authenticated.
	URL string `json:"url"`

	// Type is the Content-Type of the attachment.
	Type string `json:"type"`

	// Name is the full file name.
	Name string `json:"name"`

	// Size is the size, in bytes, of the attachment when fully decoded (i.e. the
	// number of bytes in the file the user would download).
	Size int64 `json:"size"`

	// IsLine is whether the attachment is referenced by a cid: link from within
	// the HTML body of the message.
	IsInline bool `json:"isInline"`

	// Width is the width (in px) of the image, if the attachment is an image.
	Width int64 `json:"width,omitempty"`

	// Height is the height (in px) of the image, if the attachment is an image.
	Height int64 `json:"height,omitempty"`
}

// Message represents an email message. A message is immutable except for the
// boolean isXXX status properties and the set of mailboxes it is in.
type Message struct {
	// A unique, immutable ID that does not change if the message changes
	// mailboxes.
	ID string `json:"id"`

	// ThreadID is the id fo the thread to which this message belongs.
	ThreadID string `json:"threadId"`

	// MailboxIDs are the ids of the mailboxes the message is in. A message must
	// belong to one or more mailboxes at all times (until it is deleted). This
	// property is mutable.
	MailboxIDs []string `json:"mailboxIds"`

	// InReplyToMessageID is the id of the message this message is a reply to.
	InReplyToMessageID *string `json:"inReplyToMessageId"`

	// IsUnread is whether the email not yet been read. This corresponds to the
	// opposite of the \Seen system flag in IMAP. This property is mutable.
	IsUnread bool `json:"isUnread"`

	// IsFlagged is whether the email been flagged (starred, or pinned). This
	// corresponds to the \Flagged system flag in IMAP. This property is mutable.
	IsFlagged bool `json:"isFlagged"`

	// IsAnswered is whether the email has been replied to. This corresponds to
	// the \Answered system flag in IMAP. This property is mutable.
	IsAnswered bool `json:"isAnswered"`

	// IsDraft is whether the email is a draft. This corresponds to the \Draft
	// system flag in IMAP This property is mutable.
	IsDraft bool `json:"isDraft"`

	// HasAttachment is whether the message has any attachments.
	HasAttachment bool `json:"hasAttachment"`

	// RawURL is a url to download the original RFC2822 message from. The HTTP
	// request must be authenticated.
	RawURL string `json:"rawUrl"`

	// Headers is a map of the header name to (decoded) header value for all
	// headers in the message. For headers that occur multiple times (e.g.
	// Received), the values are concatenated with a single new line (\n)
	// character in between each one.
	Headers map[string]string `json:"headers"`

	// From contains the name/email from the parsed From header of the email. If
	// the email doesn't have a From header, this is nil.
	From *Emailer `json:"from"`

	// To is a slice of emailers representing the parsed To header of the email,
	// in the same order as they appear in the header. If the email doesn't have
	// a To header, this is nil. If the header exists, but does not have any
	// content, this property should be an empty slice.
	To *[]Emailer `json:"to"`

	// CC is a slice of emailers representing the parsed Cc header of the email,
	// in the same order as they appear in the header. If the email doesn't have
	// a Cc header, this is nil. If the header exists, but does not have any
	// content, this property should be an empty slice.
	CC *[]Emailer `json:"cc"`

	// BCC is a slice of emailers representing the parsed Bcc header of the
	// email. If the email doesn't have a Bcc header (which will be true for most
	// emails outside of the Sent mailbox), this is nil. If the header exists,
	// but does not have any content, this property should be an empty slice.
	BCC *[]Emailer `json:"bcc"`

	// ReplyTo is a emailer from the parsed Reply-To header of the email. If the
	// email doesn't have a Reply-To header, this is nil.
	ReplyTo *Emailer `json:"emailer"`

	// Subject is the subject of the message.
	Subject string `json:"subject"`

	// Date is the date the message was sent (or saved, if the message is a
	// draft).
	Date time.Time `json:"date"`

	// Size is the size in bytes of the whole message as counted by the server
	// towards the user's quota.
	Size int64 `json:"size"`

	// Preview is intended to be shown as a preview line on a mailbox listing. It
	// is up to 256 characters of the beginning of a plain text version of the
	// message body.
	Preview string `json:"preview"`

	// TextBody is the plain text body part for the message. If there is only an
	// HTML version of the body, a plain text version will be generated from
	// this.
	TextBody *string `json:"textBody"`

	// HTMLBody is the HTML body part for the message, if present. If there is
	// only a plain text version of the body, an HTML version will be generated
	// from this. Any scripting content, or references to external plugins, must
	// be stripped from the HTML by the server.
	HTMLBody *string `json:"htmlBody"`

	// Attachments is a slice of attachments detailing all the attachments to the
	// message.
	Attachments *[]Attachment `json:"attachments"`

	// AttachedMessages is a map from attachment ids to message objects with less
	// properties and some properties renamed. I'm not really sure what this
	// means and have asked in the mailing list for clarification.
	//
	// See attachedMessages in http://jmap.io/spec.html#messages. Mailing list
	// discussion: https://groups.google.com/forum/#!msg/jmap-discuss/coEBFDY_A7E/vbIhJFYHCogJ
	AttachedMessages *map[string]Message `json:"attachedMessages"`
}

// SearchSnippet represents the relevant portion of a message when searching
// for messages. e.g. the portion of a message that contains a string I
// searched for in its body.
type SearchSnippet struct {
	// MessageID is the id of the message that the snippet applies to.
	MessageID string `json:"messageId"`

	// Subject is the HTML-escaped subject of the message with matching
	// words/phrases wrapped in <b></b> tags, if text from the filter matches the
	// subject. If it does not match, this is nil.
	Subject *string `json:"subject"`

	// Preview is the relevant section of the body (converted to plain text if
	// originally HTML), HTML-escaped, with matching words/phrases wrapped in
	// <b></b> tags, up to 256 characters long, if text from the filter matches
	// the plain-text or HTML body. If it does not match, this is nil.
	Preview *string `json:"preview"`
}

// ContactGroup represents a named set of contacts.
//
// http://jmap.io/spec.html#contact-groups
type ContactGroup struct {
	// ID is the id of the group. This is immutable.
	ID string `json:"id"`

	// Name is the user-visible name for the group, e.g. "Friends". This may be
	// any UTF-8 string of at least 1 character in length and maximum 256 bytes
	// in size. The same name may not be used by two different groups.
	Name string `json:"name"`
}
