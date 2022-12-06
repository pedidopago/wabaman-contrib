import {HostTemplate, ParsedTemplate} from '../lib/index';
import { expect } from 'chai';

const wabaman_template_json = `
{
    "graph_object": {
        "name": "erp_welcome",
        "language": {
            "code": "pt_BR"
        },
        "components": [
            {
                "type": "header",
                "parameters": [
                    {
                        "type": "image",
                        "image": {
                            "link": "https://pedidopago-support.s3.sa-east-1.amazonaws.com/temp/example_image_waba.png"
                        }
                    }
                ]
            }
        ]
    },
    "original": {
        "id": 1,
        "business_id": 1,
        "name": "erp_welcome",
        "category": "transactional",
        "created_at": "2022-10-27T18:31:33Z",
        "languages": [
            {
                "template_ref_id": 1,
                "language_code": "pt_BR",
                "header": {
                    "header_type": "media",
                    "content_example": "https://www.facebook.com/images/fb_icon_325x325.png"
                },
                "body": "Bem-vindo(a) ao Pharmer!",
                "footer": null,
                "buttons_type": "quick_reply",
                "quick_reply_buttons": [
                    {
                        "text": "Solicitar orçamento"
                    },
                    {
                        "text": "Ajuda ou dúvida"
                    }
                ],
                "call_to_action_buttons": null,
                "last_status": null,
                "created_at": "2022-10-27T18:33:43Z"
            }
        ]
    }
}`

const wabaman_obj = JSON.parse(wabaman_template_json) as HostTemplate;

describe("should parse a template", () => {
    it("should parse a template", () => {
        const parsed = new ParsedTemplate(wabaman_obj);
        expect(parsed.body).to.be.equal("Bem-vindo(a) ao Pharmer!");
    });
})

const wabaman_template_json2 = `
{
    "parsed": {
        "template_name": "erp_welcome",
        "language_code": "pt_BR",
        "body": "Bem-vindo(a) ao Pharmer2!"
    }
}`

const wabaman_obj2 = JSON.parse(wabaman_template_json2) as HostTemplate;

describe("should parse a template v2", () => {
    it("should parse a template v2", () => {
        const parsed = new ParsedTemplate(wabaman_obj2);
        expect(parsed.body).to.be.equal("Bem-vindo(a) ao Pharmer2!");
    });
})