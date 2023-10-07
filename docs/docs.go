// Code generated by swaggo/swag at 2023-10-06 21:09:56.984913 +0300 EEST m=+1.603952959. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/details": {
            "get": {
                "description": "Check if service is in running state and can access its resources with additional list of resource and its status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get service detailed status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_status.ResponseDetails"
                        }
                    }
                }
            }
        },
        "/events": {
            "post": {
                "description": "Create event",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Create event",
                "parameters": [
                    {
                        "description": "Create request",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_event.createRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_agnesenvaite_events_internal_event.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_agnesenvaite_events_internal_api_error.ListedError"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Check if service is in running state and can access its resources",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get service status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_status.ResponseStatus"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_agnesenvaite_events_internal_api_error.Error": {
            "type": "object",
            "properties": {
                "parameter": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "github_com_agnesenvaite_events_internal_api_error.ListedError": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_agnesenvaite_events_internal_api_error.Error"
                    }
                }
            }
        },
        "github_com_agnesenvaite_events_internal_event.Response": {
            "type": "object",
            "properties": {
                "audio_qualities": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "created_at": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "invitees": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "languages": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "video_qualities": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "internal_api_event.createRequest": {
            "type": "object",
            "properties": {
                "audio_qualities": {
                    "type": "array",
                    "items": {
                        "type": "string",
                        "enum": [
                            "high",
                            "medium",
                            "low"
                        ]
                    }
                },
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "invitees": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "languages": {
                    "type": "array",
                    "items": {
                        "type": "string",
                        "enum": [
                            "english",
                            "lithuanian",
                            "dutch"
                        ]
                    }
                },
                "name": {
                    "type": "string"
                },
                "video_qualities": {
                    "type": "array",
                    "items": {
                        "type": "string",
                        "enum": [
                            "720p",
                            "1080p",
                            "2160p"
                        ]
                    }
                }
            }
        },
        "internal_api_status.ResponseDetails": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "mysql": " OK"
                    }
                }
            }
        },
        "internal_api_status.ResponseStatus": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Events API",
	Description:      "Events API specification",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}