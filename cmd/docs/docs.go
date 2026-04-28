package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@aiagent.com"
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
        "/agent/generate": {
            "post": {
                "description": "Generates code based on prompt and creates files in current VS Code workspace",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "Generate code using AI",
                "parameters": [
                    {
                        "description": "Code generation request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.GenerateCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GenerateCodeResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/agent/context": {
            "get": {
                "description": "Returns information about the currently open VS Code workspace",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "Get current VS Code workspace context",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GetCurrentContextResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Returns the health status of the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.GenerateCodeRequest": {
            "type": "object",
            "properties": {
                "prompt": {
                    "type": "string",
                    "example": "Create a REST API for user management"
                },
                "projectName": {
                    "type": "string",
                    "example": "user-api"
                }
            }
        },
        "models.GenerateCodeResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "path": {
                                "type": "string"
                            },
                            "content": {
                                "type": "string"
                            }
                        }
                    }
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "models.GetCurrentContextResponse": {
            "type": "object",
            "properties": {
                "workspaceFolder": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "activeFile": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "AI Agent API",
	Description:      "AI Agent with DeepSeek integration for automatic code generation in VS Code",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
