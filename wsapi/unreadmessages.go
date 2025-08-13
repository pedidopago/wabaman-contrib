package wsapi

type GetUnreadMessagesRequest struct {
	ContactID uint64 `json:"contact_id"`
}

type GetUnreadMessagesResponse struct {
	ContactID      uint64 `json:"contact_id"`
	UnreadMessages uint   `json:"unread_messages"`
}
