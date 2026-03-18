package wsapi

// GetUnreadMessagesRequest is sent by a browser client to query the unread message count
// for a contact.
type GetUnreadMessagesRequest struct {
	ContactID uint64 `json:"contact_id"`
}

// GetUnreadMessagesResponse is the server's reply to a [GetUnreadMessagesRequest].
type GetUnreadMessagesResponse struct {
	ContactID      uint64 `json:"contact_id"`
	UnreadMessages uint   `json:"unread_messages"`
}
