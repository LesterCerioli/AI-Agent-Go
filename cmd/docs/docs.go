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
                "description": "Generates code based on prompt and requirements, creates files in specified path",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Agent"],
                "summary": "Generate code using Ollama AI",
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
                            "additionalProperties": {"type": "string"}
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {"type": "string"}
                        }
                    }
                }
            }
        },
        "/agent/context": {
            "get": {
                "description": "Returns information about the current workspace",
                "produces": ["application/json"],
                "tags": ["Agent"],
                "summary": "Get current workspace context",
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
                            "additionalProperties": {"type": "string"}
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Returns the health status of the service",
                "produces": ["application/json"],
                "tags": ["Health"],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {"type": "string"}
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.GenerateCodeRequest": {
            "type": "object",
            "required": ["prompt", "requirements", "project_path"],
            "properties": {
                "prompt": {
                    "type": "string",
                    "example": "Create a REST API for user management"
                },
                "projectName": {
                    "type": "string",
                    "example": "user-api"
                },
                "requirements": {
                    "type": "string",
                    "example": "Go 1.26 with Fiber, PostgreSQL, Docker, Swagger"
                },
                "project_path": {
                    "type": "string",
                    "example": "C:/Users/developer/projects/user-api"
                },
                "language": {
                    "type": "string",
                    "example": "go"
                },
                "description": {
                    "type": "string",
                    "example": "User management API"
                }
            }
        },
        "models.GenerateCodeResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "message": {
                    "type": "string",
                    "example": "Successfully generated 5 files"
                },
                "projectId": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "projectPath": {
                    "type": "string",
                    "example": "C:/Users/developer/projects/user-api"
                },
                "filesCreated": {
                    "type": "array",
                    "items": {"type": "string"},
                    "example": ["main.go", "Dockerfile", "docker-compose.yml"]
                }
            }
        },
        "models.GetCurrentContextResponse": {
            "type": "object",
            "properties": {
                "workspace_path": {
                    "type": "string",
                    "example": "C:/Users/developer/projects/my-project"
                },
                "project_name": {
                    "type": "string",
                    "example": "my-project"
                },
                "files": {
                    "type": "array",
                    "items": {"type": "string"}
                },
                "language": {
                    "type": "string",
                    "example": "go"
                }
            }
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "AI Agent API",
	Description:      "AI Agent with Ollama integration for automatic code generation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
