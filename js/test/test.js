"use strict";
var expect = require("chai").expect;
var index = require("../dist/index");

const exampleHeaderImage = "https://www.facebook.com/images/fb_icon_325x325.png";

const host_template1 = {
    graph_object: {
        name: "hello_tpl",
        language: {
            code: "pt_BR",
        },
        components: [
            {
                "type": "header",
                "parameters": [
                    {
                        "type": "image",
                        "image": {
                            "link": exampleHeaderImage,
                        }
                    }
                ],
            },
            {
                "type": "body",
                "parameters": [
                    {
                        "type": "text",
                        "text": "Mr Jones"
                    }
                ],
            }
        ],
    },
    original: {
        id: 100,
        business_id: 200,
        name: "hello_tpl",
        category: "transactional",
        created_at: "2002-10-02T10:00:00-05:00",
        languages: [
            {
                language_code: "pt_BR",
                header: {
                    header_type: "image",
                    content_example: "http://google.com/image.png",
                },
                body: "Olá {{1}}",
                footer: "",
                buttons_type: "none",
            }
        ]
    }
}

describe("test template parse", () => {
    const ptpl = new index.ParsedTemplate(host_template1);
    it("should parse template", () => {
        expect(ptpl.header.content).to.equal(exampleHeaderImage);
        expect(ptpl.body).to.equal("Olá Mr Jones");
    });
});

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

const wabaman_obj = JSON.parse(wabaman_template_json);

describe("test real life template parse", () => {
    const ptpl = new index.ParsedTemplate(wabaman_obj);
    it("should parse template", () => {
        expect(ptpl.body).to.equal("Bem-vindo(a) ao Pharmer!");
        expect(ptpl.quick_reply_buttons[0].text).to.equal("Solicitar orçamento");
    });
});