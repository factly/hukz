// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/events": {
            "get": {
                "description": "Get all events",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Show all events",
                "operationId": "get-all-events",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "limit per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/event.paging"
                        }
                    }
                }
            },
            "post": {
                "description": "Create Event",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Create Event",
                "operationId": "add-event",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Event Object",
                        "name": "Event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/event.event"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/events/{event_id}": {
            "get": {
                "description": "Get event by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Show a event by id",
                "operationId": "get-event-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Event ID",
                        "name": "event_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Event"
                        }
                    }
                }
            },
            "put": {
                "description": "Update event by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Update a event by id",
                "operationId": "update-event-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Event ID",
                        "name": "event_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Event Object",
                        "name": "Event",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/event.event"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Event"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete Event by id",
                "tags": [
                    "Events"
                ],
                "summary": "Delete Event by id",
                "operationId": "delete-event-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Event ID",
                        "name": "event_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/webhooks": {
            "get": {
                "description": "Get all webhooks",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Show all webhooks",
                "operationId": "get-all-webhooks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "limit per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhook.paging"
                        }
                    }
                }
            },
            "post": {
                "description": "Create Webhook",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Create Webhook",
                "operationId": "add-webhook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Webhook Object",
                        "name": "Webhook",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhook.webhook"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Webhook"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/webhooks/logs": {
            "get": {
                "description": "Get all webhooks logs",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Show all webhooks logs",
                "operationId": "get-all-webhooks-logs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "limit per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhook.paging"
                        }
                    }
                }
            }
        },
        "/webhooks/{webhook_id}": {
            "get": {
                "description": "Get webhook by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Show a webhook by id",
                "operationId": "get-webhook-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Webhook ID",
                        "name": "webhook_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Webhook"
                        }
                    }
                }
            },
            "put": {
                "description": "Update webhook by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Update a webhook by id",
                "operationId": "update-webhook-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Webhook ID",
                        "name": "webhook_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Webhook Object",
                        "name": "Webhook",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/webhook.webhook"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Webhook"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete webhook by id",
                "tags": [
                    "Webhooks"
                ],
                "summary": "Delete webhook by id",
                "operationId": "delete-webhook-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Webhook ID",
                        "name": "webhook_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "event.event": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "tags": {
                    "type": "string"
                }
            }
        },
        "event.paging": {
            "type": "object",
            "properties": {
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Event"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "model.Event": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by_id": {
                    "type": "integer"
                },
                "deleted_at": {
                    "type": "string"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Webhook"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "tags": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by_id": {
                    "type": "integer"
                }
            }
        },
        "model.Webhook": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by_id": {
                    "type": "integer"
                },
                "deleted_at": {
                    "type": "string"
                },
                "enabled": {
                    "type": "boolean"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Event"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "tags": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by_id": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "webhook.paging": {
            "type": "object",
            "properties": {
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Webhook"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "webhook.webhook": {
            "type": "object",
            "required": [
                "event_ids",
                "url"
            ],
            "properties": {
                "enabled": {
                    "type": "boolean"
                },
                "event_ids": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "tags": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:7790",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Webhooks API",
	Description: "Webhooks Service API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
