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
