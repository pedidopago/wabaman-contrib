
export enum WSMessageType {
    Ping = 0,
    Pong = 1,
    ClientMessage = 2,
    HostMessage = 3,
    ReadByHostReceipt = 4,
    ClientReceipt = 5,
    ContactUpdate = 6,
    MockClientMessages = 230,
    CloseError = 240,
}

export interface WebsocketMessage {
    type: WSMessageType;
    error?: WSError;
    client_message?: ClientMessage;
    host_message?: HostMessage;
}

export enum WSErrorCode {
    InvalidAuth = 1,
    InternalError = 500,
    RedisClosed = 600,
}

export interface WSError {
    code: WSErrorCode;
    message: string;
}

export interface ClientMessage {
    id: number;
    waba_message_id: string;
    waba_from_id: string;
    waba_profile_name: string;
    waba_timestamp: string;
    type: string;
    text?: TextMessage;
    //TODO: continue
}

export interface HostMessage {
    id: number;
    waba_message_id: string;
    phone_id: number;
    host_phone_number: string;
    waba_recipient_id: string;
    waba_timestamp: string;
    type: string;
    text?: TextMessage;
    //TODO: continue
}

export interface TextMessage {
    body: string;
}

export interface HostTemplate {
    graph_object?: FBGraphTemplate;
    original?: TemplateRef;
}

export interface FBGraphTemplate {
    namespace?: string;
    name: string;
    language?: FBGraphLanguageObject;
    components?: FBGraphTemplateComponent[];
}

export interface FBGraphLanguageObject {
    code: string;
}

export enum FBGraphTemplateComponentType {
    Header = "header",
    Body = "body",
    Footer = "footer",
}

export interface FBGraphTemplateComponent {
    type: FBGraphTemplateComponentType;
    parameters?: FBGraphTemplateComponentParameter[];
}

export enum FBGraphTemplateComponentParameterType {
    Currency = "currency",
    DateTime = "date_time",
    Document = "document",
    Image = "image",
    Text = "text",
    Video = "video",
}

export interface FBGraphTemplateComponentParameter {
    type: FBGraphTemplateComponentParameterType;
    text?: string;
    image?: FBGraphMediaObject;
    date_time?: FBGraphDateTimeParameters;
    currency?: FBGraphCurrencyParameters;
    video?: FBGraphMediaObject;
}

export interface FBGraphDateTimeParameters {
    fallback_value: string;
}

export interface FBGraphCurrencyParameters {
    fallback_value: string;
    code: string;
    amount_1000: number; // Amount multiplied by 1000
}

export interface FBGraphMediaObject {
    // Required when type is audio, document, image, sticker, or video and you are not using a link.
    id?: string;
    // Required when type is audio, document, image, sticker, or video and you are not using an uploaded media ID.
    link?: string;
    caption?: string;
    filename?: string;
    provider?: string;
}

export interface TemplateRef {
    id: number;
    business_id: number;
    name: string;
    category: string; // transactional, marketing, disposable_credentials
    created_at: string;
    languages: TemplateRefLanguage[];
}

export interface TemplateRefLanguage {
    language_code: string;
    header?: TplRefHeader;
    body: string;
    footer: string;
    buttons_type: string; // none, call_to_action, quick_reply
    quick_reply_buttons?: TplQuickReplyButton[];
    call_to_action_buttons?: TplCallToActionButton[];
}

export interface TplRefHeader {
    header_type: string; // none, text, media
    content_example: string;
}

export interface TplQuickReplyButton {
    text: string;
}

export interface TplCallToActionButton {
    type: string; // call, url
    text: string;
    url?: TplCallToActionURL;
    call?: TplCallToActionCall;
}

export interface TplCallToActionURL {
    type: string; // static, dynamic
    href: string;
}

export interface TplCallToActionCall {
    cc: string;
    phone: string;
}