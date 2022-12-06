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

const wabaman_template_json3 = `
{
    "parsed": {
        "template_name": "orcamento",
        "language_code": "pt_BR",
        "header": {
            "header_type": "media",
            "content_example": "example",
            "content": "https://ppv2-development.s3.sa-east-1.amazonaws.com/01EM2CNQFJG4VWS5BMT8B2H6YF/01GK513EJFKGCZ49B1KW4TFY6N/01GKKW56E564W8ZMNG24S7FT5W.jpg",
            "type": "image"
        },
        "body": "Solicitação: *32329*  R$ 1.061,00 no PIX Itens: 6 Expira em 10 dias corridos Clique para aprovar agora",
        "footer": "",
        "buttons_type": "call_to_action",
        "quick_reply_buttons": null,
        "call_to_action_buttons": [
            {
                "type": "url",
                "text": "Frete e pagamento",
                "url": {
                    "type": "dynamic",
                    "href": "https://staging.loja.pedidopago.com.br/s?inquiry=01GKKS950H3R0DXVCZVN844J50"
                }
            },
            {
                "type": "call",
                "text": "Ligar!",
                "call": {
                    "cc": "",
                    "phone": "+5511993903535"
                }
            }
        ]
    }
}`

const wabaman_obj3 = JSON.parse(wabaman_template_json3) as HostTemplate;

const wabaman_template_json4 = `{
    "graph_object": {
      "name": "orcamento",
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
                "link": "https://ppv2-development.s3.sa-east-1.amazonaws.com/01EM2CNQFJG4VWS5BMT8B2H6YF/01GK9RBE8DYZJVSNWCHJ2NHSN1/01GKM83P93V93TZBDZGNP9A9AZ.jpg"
              }
            }
          ]
        },
        {
          "type": "body",
          "parameters": [
            {
              "type": "text",
              "text": "32322"
            },
            {
              "type": "text",
              "text": "R$ 1.669,00 no PIX"
            },
            {
              "type": "text",
              "text": "3"
            }
          ]
        },
        {
          "type": "button",
          "sub_type": "url",
          "parameters": [
            {
              "type": "text",
              "text": "s?inquiry=01GKA2SSYWTJV30ZFVASBFVTTH"
            }
          ],
          "index": 0
        }
      ]
    },
    "original": {
      "id": 0,
      "business_id": 0,
      "name": "orcamento",
      "category": "transactional",
      "created_at": "0001-01-01T00:00:00Z",
      "languages": [
        {
          "template_ref_id": 0,
          "language_code": "pt_BR",
          "header": {
            "header_type": "media",
            "content_example": "https://scontent.whatsapp.net/v/t61.29466-34/316467224_1460837674743763_2615426465772862050_n.jpg?ccb=1-7&_nc_sid=57045b&_nc_ohc=WqcFW--tsmkAX-cRVUS&_nc_ht=scontent.whatsapp.net&edm=AH51TzQEAAAA&oh=01_AdRsOccxSwYUf1cBdCJt03Lu8OeU6o1sVXlNJNVSmACyYw&oe=63B6D157"
          },
          "body": "Solicitação: *{{1}}* 💰 {{2}} 💊 Itens: {{3}} 🗓️ Expira em 10 dias corridos Clique para aprovar agora👇",
          "footer": null,
          "buttons_type": "",
          "quick_reply_buttons": null,
          "call_to_action_buttons": null,
          "last_status": null,
          "created_at": "0001-01-01T00:00:00Z",
          "header_parameter_type": "",
          "header_content_default": null
        }
      ]
    },
    "parsed": {
      "template_name": "orcamento",
      "language_code": "pt_BR",
      "header": {
        "header_type": "media",
        "content_example": "https://scontent.whatsapp.net/v/t61.29466-34/316467224_1460837674743763_2615426465772862050_n.jpg?ccb=1-7&_nc_sid=57045b&_nc_ohc=WqcFW--tsmkAX-cRVUS&_nc_ht=scontent.whatsapp.net&edm=AH51TzQEAAAA&oh=01_AdRsOccxSwYUf1cBdCJt03Lu8OeU6o1sVXlNJNVSmACyYw&oe=63B6D157",
        "content": "https://ppv2-development.s3.sa-east-1.amazonaws.com/01EM2CNQFJG4VWS5BMT8B2H6YF/01GK9RBE8DYZJVSNWCHJ2NHSN1/01GKM83P93V93TZBDZGNP9A9AZ.jpg",
        "type": "image"
      },
      "body": "Solicitação: *32322* 💰 R$ 1.669,00 no PIX 💊 Itens: 3 🗓️ Expira em 10 dias corridos Clique para aprovar agora👇",
      "footer": "",
      "buttons_type": "call_to_action",
      "quick_reply_buttons": null,
      "call_to_action_buttons": [
        {
          "type": "url",
          "text": "Frete e pagamento",
          "url": {
            "type": "dynamic",
            "href": "https://staging.loja.pedidopago.com.br/s?inquiry=01GKA2SSYWTJV30ZFVASBFVTTH"
          }
        },
        {
          "type": "call",
          "text": "Fale conosco",
          "call": {
            "cc": "",
            "phone": "+5511993903535"
          }
        }
      ]
    }
  }`

const wabaman_obj4 = JSON.parse(wabaman_template_json4) as HostTemplate;

describe("should parse a template v2", () => {
    it("should parse a template v2", () => {
        const parsed = new ParsedTemplate(wabaman_obj2);
        expect(parsed.body).to.be.equal("Bem-vindo(a) ao Pharmer2!");
    });
    it("should parse a template v2 with header and buttons", () => {
        const parsed = new ParsedTemplate(wabaman_obj3);
        expect(parsed.call_to_action_buttons[0].text).to.be.equal("Frete e pagamento");
        expect(parsed.call_to_action_buttons[1].text).to.be.equal("Ligar!");
        expect(parsed.header?.content_example).to.be.equal("example");
    });
    it("should parse a template v2 with legacy options", () => {
        const parsed = new ParsedTemplate(wabaman_obj4);
        expect(parsed.call_to_action_buttons[0].text).to.be.equal("Frete e pagamento");
        expect(parsed.call_to_action_buttons[1].text).to.be.equal("Fale conosco");
        expect(parsed.body).to.be.equal("Solicitação: *32322* 💰 R$ 1.669,00 no PIX 💊 Itens: 3 🗓️ Expira em 10 dias corridos Clique para aprovar agora👇");
    });
})