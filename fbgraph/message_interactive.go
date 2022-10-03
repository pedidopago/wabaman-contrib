package fbgraph

type InteractiveMessageType string

const (
	InteractiveMessageButton      InteractiveMessageType = "button"
	InteractiveMessageList        InteractiveMessageType = "list"
	InteractiveMessageProduct     InteractiveMessageType = "product"
	InteractiveMessageProductList InteractiveMessageType = "product_list"
)

type InteractiveMessageObject struct {
	Type InteractiveMessageType `json:"type"`
	// Required.
	// Action you want the user to perform after reading the message.
	Action *InteractiveMessageAction `json:"action,omitempty"`
	// Optional for type product. Required for other message types.
	// An object with the body of the message.
	//
	// The content of the message. Emojis and markdown are supported.
	// Maximum length: 1024 characters.
	Body *InteractiveTextObject `json:"body,omitempty"`
	// Optional. An object with the footer of the message.
	//
	// The footer content. Emojis, markdown, and links are supported.
	// Maximum length: 60 characters.
	Footer *InteractiveTextObject `json:"footer,omitempty"`
	// Required for type product_list. Optional for other types.
	// Header content displayed on top of a message. You cannot set a header
	// if your interactive object is of product type. See header object for more information.
	Header *InteractiveHeaderObject `json:"header,omitempty"`
}

type InteractiveTextObject struct {
	Text string `json:"text"`
}

type InteractiveHeaderObject struct {
	// Required if type is set to document.
	// Contains the media object for this document.
	Document *MediaObject `json:"document,omitempty"`

	// Required if type is set to image.
	// Contains the media object for this image.
	Image *MediaObject `json:"image,omitempty"`
	// Required if type is set to text.
	// Text for the header. Formatting allows emojis, but not markdown.
	// Maximum length: 60 characters.
	Text string `json:"text,omitempty"`

	// Required.
	// The header type you would like to use. Supported values:
	//     text: Used for List Messages, Reply Buttons, and Multi-Product Messages.
	//     video: Used for Reply Buttons.
	//     image: Used for Reply Buttons.
	//     document: Used for Reply Buttons.
	Type string `json:"type"`

	// Required if type is set to video.
	// Contains the media object for this video.
	Video *MediaObject `json:"video,omitempty"`
}

type MediaObject struct {
	// Required when type is audio, document, image, sticker, or video and you are not using a link.
	ID string `json:"id,omitempty"`
	// Required when type is audio, document, image, sticker, or video and you are not using an uploaded media ID.
	Link string `json:"link,omitempty"`
	// Optional.
	// Describes the specified image or video media.
	// Do not use with audio, document, or sticker media.
	Caption string `json:"caption,omitempty"`
	// Optional.
	// Describes the filename for the specific document. Use only with document media.
	// The extension of the filename will specify what format the document is displayed as in WhatsApp.
	Filename string `json:"filename,omitempty"`
	// Optional. Only used for On-Premises API.
	Provider string `json:"provider,omitempty"`

	InternalMetadata struct {
		WabamanIDHex string `json:"-"`
	} `json:"-"`
}

type InteractiveMessageAction struct {
	// Required for List Messages.
	// Button content. It cannot be an empty string and must be unique within
	// the message. Emojis are supported, markdown is not.
	// Maximum length: 20 characters.
	Button string `json:"button,omitempty"`
	// Required for Reply Buttons.
	// You can have up to 3 buttons. You cannot have leading or trailing spaces when setting the ID.
	Buttons []InteractiveButton `json:"buttons,omitempty"`
	// Required for Single Product Messages and Multi-Product Messages.
	// Unique identifier of the Facebook catalog linked to your WhatsApp Business
	// Account. This ID can be retrieved via the [Meta Commerce Manager](https://business.facebook.com/commerce/).
	CatalogID string `json:"catalog_id,omitempty"`
	// Required for Single Product Messages and Multi-Product Messages.
	// Unique identifier of the product in a catalog.
	//
	// To get this ID go to Meta Commerce Manager and select your Meta Business
	// account. You will see a list of shops connected to your account. Click
	// the shop you want to use. On the left-side panel, click Catalog > Items,
	// and find the item you want to mention. The ID for that item is displayed
	// under the item's name.
	ProductRetailerID string `json:"product_retailer_id,omitempty"`
	// Required for List Messages and Multi-Product Messages.
	// Array of section objects. Minimum of 1, maximum of 10. See section object.
	Sections []InteractiveMessageSection `json:"sections,omitempty"`
}

type InteractiveButton struct {
	Type  string                  `json:"type"` // type: only supported type is reply (for Reply Button)
	Reply *InteractiveReplyButton `json:"reply"`
}

type InteractiveReplyButton struct {
	// Button title. It cannot be an empty string and must be unique within
	// the message. Emojis are supported, markdown is not. Maximum length: 20 characters.
	Title string `json:"title"`
	// Unique identifier for your button. This ID is returned in the webhook when
	// the button is clicked by the user. Maximum length: 256 characters.
	ID string `json:"id"`
}

type InteractiveMessageSection struct {
	// Required if the message has more than one section.
	// Maximum length: 24 characters.
	Title string `json:"title,omitempty"`

	//TODO: product_items at https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages

	// Required for List Messages.
	Rows []InteractiveMessageSectionRow `json:"rows,omitempty"`
}

// Each row must have a title (Maximum length: 24 characters) and an ID
// (Maximum length: 200 characters). You can add a description (Maximum
// length: 72 characters), but it is optional.
type InteractiveMessageSectionRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// "interactive": {
//     "type": "list",
//     "header": {
//       "type": "text",
//       "text": "HEADER_TEXT"
//     },
//     "body": {
//       "text": "BODY_TEXT"
//     },
//     "footer": {
//       "text": "FOOTER_TEXT"
//     },
//     "action": {
//       "button": "BUTTON_TEXT",
//       "sections": [
//         {
//           "title": "SECTION_1_TITLE",
//           "rows": [
//             {
//               "id": "SECTION_1_ROW_1_ID",
//               "title": "SECTION_1_ROW_1_TITLE",
//               "description": "SECTION_1_ROW_1_DESCRIPTION"
//             },
//             {
//               "id": "SECTION_1_ROW_2_ID",
//               "title": "SECTION_1_ROW_2_TITLE",
//               "description": "SECTION_1_ROW_2_DESCRIPTION"
//             }
//           ]
//         },
//         {
//           "title": "SECTION_2_TITLE",
//           "rows": [
//             {
//               "id": "SECTION_2_ROW_1_ID",
//               "title": "SECTION_2_ROW_1_TITLE",
//               "description": "SECTION_2_ROW_1_DESCRIPTION"
//             },
//             {
//               "id": "SECTION_2_ROW_2_ID",
//               "title": "SECTION_2_ROW_2_TITLE",
//               "description": "SECTION_2_ROW_2_DESCRIPTION"
//             }
//           ]
//         }
//       ]
//     }
//   }
