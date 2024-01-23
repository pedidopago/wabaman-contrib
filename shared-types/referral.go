package types

import "encoding/json"

// MessageObjectReferral is the object sometimes sent by Meta, attached to a message sent by a client.
type MessageObjectReferral struct {
	// The Meta URL that leads to the ad or post clicked by the customer. Opening this url takes you to the ad viewed by your customer.
	SourceUrl string `json:"source_url"`
	// The type of the ad’s source; ad or post
	SourceType string `json:"source_type"`
	// Meta ID for an ad or a post
	SourceId string `json:"source_id"`
	// Headline used in the ad or post
	Headline string `json:"headline"`
	// Body for the ad or post
	Body string `json:"body"`
	// Media present in the ad or post; image or video
	MediaType string `json:"media_type"`
	// URL of the image, when media_type is an image
	ImageUrl string `json:"image_url"`
	// URL of the video, when media_type is a video
	VideoUrl string `json:"video_url"`
	// URL for the thumbnail, when media_type is a video
	ThumbnailUrl string `json:"thumbnail_url"`
}

func (r MessageObjectReferral) JSONString() string {
	b, _ := json.Marshal(r)

	return string(b)
}
