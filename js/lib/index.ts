
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
    template?: HostTemplate;
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
    document?: FBGraphMediaObject;
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

export class ParsedTemplateHeader implements TplRefHeader {
    header_type: string; // none, text, media
    content_example: string;
    content: string;

    constructor(header: TplRefHeader, content: string) {
        this.header_type = header.header_type;
        this.content_example = header.content_example;
        this.content = content;
    }
}

export class ParsedTemplate {
    template_name: string;
    language_code: string;
    header: ParsedTemplateHeader | null;
    body: string;
    footer: string;
    buttons_type: string; // none, call_to_action, quick_reply
    quick_reply_buttons: TplQuickReplyButton[];
    call_to_action_buttons: TplCallToActionButton[];

    constructor(tpl: HostTemplate | null){
        this.template_name = "";
        this.language_code = "";
        this.header = null;
        this.body = "";
        this.footer = "";
        this.buttons_type = "";
        this.quick_reply_buttons = [];
        this.call_to_action_buttons = [];
        if(tpl){
            doTemplate(tpl, this);
        }
    }
}

function doTemplate(tpl: HostTemplate, out: ParsedTemplate) {
    if (!tpl.original) {
        return;
    }
    if (!tpl.graph_object) {
        return;
    }

    let llang = "pt_BR";

    if (tpl.graph_object.language && tpl.graph_object.language.code !== "") {
        llang = tpl.graph_object.language.code;
    }

    tpl.original.languages.forEach((lang) => {
        if (lang.language_code === llang) {
            out.language_code = lang.language_code;
            out.body = doTemplateReplacement(lang.body, findBody(tpl.graph_object));
            out.footer = doTemplateReplacement(lang.footer, findFooter(tpl.graph_object));
            out.buttons_type = lang.buttons_type;
            if (lang.header) {
                let content = "";
                let h = findHeader(tpl.graph_object);
                if (h && h.parameters && h.parameters.length > 0) {
                    const item0 = h.parameters[0];
                    switch(item0.type){
                        case FBGraphTemplateComponentParameterType.Text:
                            content = doTemplateReplacement(item0.text, findHeader(tpl.graph_object));
                            break;
                        case FBGraphTemplateComponentParameterType.Image:
                            content = item0.image?.link || "";
                            break;
                        case FBGraphTemplateComponentParameterType.Video:
                            content = item0.video?.link || "";
                            break;
                        case FBGraphTemplateComponentParameterType.Document:
                            content = item0.document?.link || "";
                            break;
                        default:
                            // other types are not supported in "header"
                            content = "";
                            break;
                    }
                }
                out.header = new ParsedTemplateHeader(lang.header, content);
            }
            if (lang.call_to_action_buttons) {
                out.call_to_action_buttons = lang.call_to_action_buttons;
                //TODO: doTemplateReplacement for children
            }
            if (lang.quick_reply_buttons) {
                out.quick_reply_buttons = lang.quick_reply_buttons;
                //TODO: doTemplateReplacement for children (if possible)
            }
        }
    });
}

function doTemplateReplacement(txt: string | undefined, graph: FBGraphTemplateComponent | undefined): string {
    if (!txt) {
        return "";
    }
    if (!graph) {
        return txt;
    }
    let finalTxt = txt;
    let found = true;
    let tplindex = 0;
    while(found){
        tplindex++;
        if(finalTxt.indexOf(`{{${tplindex}}`) > -1){
            if(graph.parameters && graph.parameters.length >= tplindex){
                finalTxt = textVarReplace(finalTxt, tplindex, graph.parameters[tplindex - 1]);
            } else {
                found = false;
            }
        }else{
            found = false;
        }
    }
    return finalTxt;
}

function findHeader(gobj: FBGraphTemplate | undefined): FBGraphTemplateComponent | undefined {
    if (!gobj) {
        return undefined;
    }
    if (!gobj.components) {
        return undefined;
    }
    return gobj.components.find((c) => c.type === FBGraphTemplateComponentType.Header);
}

function findBody(gobj: FBGraphTemplate | undefined): FBGraphTemplateComponent | undefined {
    if (!gobj) {
        return undefined;
    }
    if (!gobj.components) {
        return undefined;
    }
    return gobj.components.find((c) => c.type === FBGraphTemplateComponentType.Body);
}

function findFooter(gobj: FBGraphTemplate | undefined): FBGraphTemplateComponent | undefined {
    if (!gobj) {
        return undefined;
    }
    if (!gobj.components) {
        return undefined;
    }
    return gobj.components.find((c) => c.type === FBGraphTemplateComponentType.Footer);
}

function textVarReplace(txt: string, index: number, p: FBGraphTemplateComponentParameter | undefined): string {
    if(!p){
        return txt;
    }
    switch(p.type){
        case FBGraphTemplateComponentParameterType.Text:
            return txt.replace(`{{${index}}}`, p.text || "");
        case FBGraphTemplateComponentParameterType.Image:
            return txt.replace(`{{${index}}}`, p.image?.link || "");
        case FBGraphTemplateComponentParameterType.Video:
            return txt.replace(`{{${index}}}`, p.video?.link || "");
        case FBGraphTemplateComponentParameterType.Document:
            return txt.replace(`{{${index}}}`, p.document?.link || "");
        case FBGraphTemplateComponentParameterType.DateTime:
            return txt.replace(`{{${index}}}`, p.date_time?.fallback_value || "");
        case FBGraphTemplateComponentParameterType.Currency:
            return txt.replace(`{{${index}}}`, p.currency?.fallback_value || "");
        default:
            return txt;
    }
}