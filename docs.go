package main

//This file is generated automatically. Do not try to edit it manually.

var resourceListingJson = `{
    "apiVersion": "1.0.0",
    "swaggerVersion": "2.0", 
    "apis": [
        {
            "path": "/v1",
            "description": "IBM Cloud Terraform Provider API"
        },
        {
            "path": "/v2",
            "description": "IBM Cloud Configuration Discovery API" 
        }
    ],
    "info": {
        "title": "IBM Cloud Provider API",
        "description": "Swagger IBM Cloud Provider API",
        "contact": "sakshiag@in.ibm.com,anilkumar@in.ibm.com,srikar.uppada@in.ibm.com"
    }
}`

var apiDescriptionsJson = map[string]string{"v1": `{
    "apiVersion": "1.0.0",
    "swaggerVersion": "2.0",
    "resourcePath": "/v1",
    "produces": [
        "application/json"
    ],
    "apis": [
        {
            "path": "/v1/configuration",
            "description": "clone the configuration repo",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "ConfHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse",
                    "items": {},
                    "summary": "clone the configuration repo",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "body",
                            "description": "request body",
                            "dataType": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest",
                            "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 400,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}",
            "description": "delete the configuration repo",
            "operations": [
                {
                    "httpMethod": "DELETE",
                    "nickname": "ConfDeleteHandler",
                    "type": "string",
                    "items": {},
                    "summary": "delete the configuration repo",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "Some ID",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}/plan",
            "description": "Execute plan for the configuration.",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "PlanHandler",
                    "type": "",
                    "items": {},
                    "summary": "Execute plan for the configuration.",
                    "parameters": [
                        {
                            "paramType": "header",
                            "name": "SLACK_WEBHOOK_URL",
                            "description": "provide slack webhook url",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "Repo Name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 202,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}/apply",
            "description": "Execute apply for the configuration.",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "ApplyHandler",
                    "type": "",
                    "items": {},
                    "summary": "Execute apply for the configuration.",
                    "parameters": [
                        {
                            "paramType": "header",
                            "name": "SLACK_WEBHOOK_URL",
                            "description": "provide slack webhook url",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "Repo Name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 202,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}/destroy",
            "description": "Execute destroy for the configuration.",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "DestroyHandler",
                    "type": "",
                    "items": {},
                    "summary": "Execute destroy for the configuration.",
                    "parameters": [
                        {
                            "paramType": "header",
                            "name": "SLACK_WEBHOOK_URL",
                            "description": "provide slack webhook url",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "Repo Name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 202,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}/show",
            "description": "Execute show for the configuration.",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "ShowHandler",
                    "type": "",
                    "items": {},
                    "summary": "Execute show for the configuration.",
                    "parameters": [
                        {
                            "paramType": "header",
                            "name": "SLACK_WEBHOOK_URL",
                            "description": "provide slack webhook url",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "Repo Name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 202,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}/{action_name}/{action_id}/log",
            "description": "Get logs for the configuration.",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "LogHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails",
                    "items": {},
                    "summary": "Get logs for the configuration.",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_name",
                            "description": "action name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_id",
                            "description": "action id",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v1/configuration/{repo_name}/{action_name}/{action_id}/status",
            "description": "Get status of the action.",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "StatusHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse",
                    "items": {},
                    "summary": "Get status of the action.",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_name",
                            "description": "action name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_id",
                            "description": "action id",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        }
    ],
    "models": {
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails",
            "required": [
                "id",
                "action"
            ],
            "properties": {
                "action": {
                    "type": "string",
                    "description": "Action Name",
                    "items": {},
                    "format": ""
                },
                "action_id": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "error": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "string",
                    "description": "Name of the configuration",
                    "items": {},
                    "format": ""
                },
                "output": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse",
            "required": [
                "id",
                "action"
            ],
            "properties": {
                "action": {
                    "type": "string",
                    "description": "Action Name",
                    "items": {},
                    "format": ""
                },
                "action_id": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "string",
                    "description": "Name of the configuration",
                    "items": {},
                    "format": ""
                },
                "status": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "timestamp": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest",
            "required": [
                "git_url"
            ],
            "properties": {
                "git_url": {
                    "type": "string",
                    "description": "The git url of your configuraltion",
                    "items": {},
                    "format": ""
                },
                "log_level": {
                    "type": "string",
                    "description": "The log level defing by user.",
                    "items": {},
                    "format": ""
                },
                "variablestore": {
                    "type": "array",
                    "description": "The environments' variable store",
                    "items": {
                        "type":"github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.EnvironmentVariableRequest"
                    },
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse",
            "required": [
                "config_name"
            ],
            "properties": {
                "config_name": {
                    "type": "string",
                    "description": "configuration name",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse",
            "required": [
                "status"
            ],
            "properties": {
                "error": {
                    "type": "string",
                    "description": "Error of the terraform operation.",
                    "items": {},
                    "format": ""
                },
                "status": {
                    "type": "string",
                    "description": "Status of the terraform operation.",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.EnvironmentVariableRequest": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.EnvironmentVariableRequest",
            "properties": {
                "name": {
                    "type": "string",
                    "description": "The variable's name.",
                    "items": {},
                    "format": ""
                },
                "value": {
                    "type": "string",
                    "description": "The variable's value",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.TerraformerConfigRequest": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.TerraformerConfigRequest",
            "required": [
                "git_url"
            ],
            "properties": {
                "git_url": {
                    "type": "string",
                    "description": "The git url of your configuraltion",
                    "items": {},
                    "format": ""
                }
            }
        }
    }
}`}

var apiDescriptionsJsonV2 = map[string]string{"v2": `{
    "apiVersion": "1.0.0",
    "swaggerVersion": "3.0",
    "resourcePath": "/v2",
    "produces": [
        "application/json"
    ],
    "apis": [
        {
            "path": "/v2/configuration",
            "description": "clone the configuration repo",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "ConfHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse",
                    "items": {},
                    "summary": "clone the configuration repo",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "body",
                            "description": "request body",
                            "dataType": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest",
                            "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 400,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v2/configuration/{repo_name}/import",
            "description": "Import cloud resources.",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "TerraformerImportHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse",
                    "items": {},
                    "summary": "Import cloud resources.",
                    "parameters": [
                        {
                            "paramType": "query",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "query",
                            "name": "services",
                            "description": "services",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": true,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "query",
                            "name": "command",
                            "description": "Command",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "query",
                            "name": "tags",
                            "description": "Filter parameters.  ex: region:us-east,resource_group:default",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": true,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v2/configuration/{repo_name}/{action_name}/{action_id}/log",
            "description": "Get logs for the import request.",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "LogHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails",
                    "items": {},
                    "summary": "Get logs for the import request.",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_name",
                            "description": "action name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_id",
                            "description": "action id",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v2/configuration/{repo_name}/{action_name}/{action_id}/status",
            "description": "Get status of the import action.",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "StatusHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse",
                    "items": {},
                    "summary": "Get status of the import action.",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_name",
                            "description": "action name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "path",
                            "name": "action_id",
                            "description": "action id",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v2/configuration/{repo_name}",
            "description": "Get terraform tf/state file from repo",
            "operations": [
                {
                    "tags": "discovery",
                    "httpMethod": "GET",
                    "nickname": "TerraformerImportHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse",
                    "items": {},
                    "tags": {
                        "name": "discovery"
                    },
                    "summary": "Download terraform tf/state file from repo.",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "query",
                            "name": "service",
                            "description": "service",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        },
        {
            "path": "/v2/configuration/{repo_name}",
            "description": "Rollback/Delete the tf/state file from repo.",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "TerraformerImportHandler",
                    "type": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse",
                    "items": {},
                    "summary": "Rollback/Delete the tf/state file from repo.",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "repo_name",
                            "description": "repo name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "query",
                            "name": "command",
                            "description": "command",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse"
                        },
                        {
                            "code": 404,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        },
                        {
                            "code": 500,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "string"
                        }
                    ],
                    "produces": [
                        "application/json"
                    ]
                }
            ]
        }
    ],
    "models": {
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionDetails",
            "required": [
                "id",
                "action"
            ],
            "properties": {
                "action": {
                    "type": "string",
                    "description": "Action Name",
                    "items": {},
                    "format": ""
                },
                "action_id": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "error": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "string",
                    "description": "Name of the configuration",
                    "items": {},
                    "format": ""
                },
                "output": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ActionResponse",
            "required": [
                "id",
                "action"
            ],
            "properties": {
                "action": {
                    "type": "string",
                    "description": "Action Name",
                    "items": {},
                    "format": ""
                },
                "action_id": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "string",
                    "description": "Name of the configuration",
                    "items": {},
                    "format": ""
                },
                "status": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "timestamp": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigRequest",
            "required": [
                "git_url"
            ],
            "properties": {
                "git_url": {
                    "type": "string",
                    "description": "The git url of your configuraltion",
                    "items": {},
                    "format": ""
                },
                "config_name": {
                    "type": "string",
                    "description": "The name of config repo",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.ConfigResponse",
            "required": [
                "config_name"
            ],
            "properties": {
                "config_name": {
                    "type": "string",
                    "description": "configuration name",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.StatusResponse",
            "required": [
                "status"
            ],
            "properties": {
                "error": {
                    "type": "string",
                    "description": "Error of the terraform operation.",
                    "items": {},
                    "format": ""
                },
                "status": {
                    "type": "string",
                    "description": "Status of the terraform operation.",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.EnvironmentVariableRequest": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.EnvironmentVariableRequest",
            "properties": {
                "name": {
                    "type": "string",
                    "description": "The variable's name.",
                    "items": {},
                    "format": ""
                },
                "value": {
                    "type": "string",
                    "description": "The variable's value",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.TerraformerConfigRequest": {
            "id": "github.com.terrform-schematics-demo.terraform-provider-ibm-api.utils.TerraformerConfigRequest",
            "required": [
                "git_url"
            ],
            "properties": {
                "git_url": {
                    "type": "string",
                    "description": "The git url of your configuraltion",
                    "items": {},
                    "format": ""
                }
            }
        }
    }
}`}
