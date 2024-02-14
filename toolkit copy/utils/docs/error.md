## Overview
This module deals with functions necessary for the proper working of localization module.  It deals with three functions say 
GetErrorCodes, ParseFields and GenerateErrorResponse.
## Index
- [ GetErrorCodes(ctx *gin.Context) map[string]any 
    ](#func-GetErrorCodes)

- [ GetErrorCodes(ctx *gin.Context) map[string]any](#func-GetErrorCodes)

- [  ParseFields(ctx *gin.Context, errorType string, fields string, errorsValue map[string]any, endpoint string, method string) (any, bool, float64)
    ](#func-ParseFields)

- [ ParseFields(ctx *gin.Context, errorType string, fields string, errorsValue map[string]any, endpoint string, method string) (any, bool, float64)](#func-ParseFields)

- [ GenerateErrorResponse(message string, errorCode any, errorsValue map[string]any) 
    ](#func-GenerateErrorResponse)

- [ GenerateErrorResponse(message string, errorCode any, errorsValue map[string]any) ](#func-GenerateErrorResponse)



### func GetErrorCodes

    GetErrorCodes(ctx *gin.Context) map[string]any 

    This function is used for getting all errors from context.

### func ParseFields

  ParseFields(ctx *gin.Context, errorType string, fields string, errorsValue map[string]any, endpoint string, method string) (any, bool, float64)
    This function is used to retrieve errors form all errors satisfying the query param parameters like error type, endpoint names, method and field. and returns the error in proper format.
  The general format is as given below
    {
        message   
        errorCode 
        errors    
    }

This function takes `errorType` , which is a string ie., the type of the errors like validation_error, internal_server_error, not_found. `fields` which is a string specifies the fields like address:format|required, `errorsValue`, which is a map of desired error, `key` which is a string used for retrieving the errors, `endpoint` which is a string that specifies the endpoint name and `method` which is a string that specifies the method of that api. This function wiil returns errorCode 


### func GenerateErrorResponse

GenerateErrorResponse(message string, errorCode any, errorsValue map[string]any) any

This function is used for structuring the errors (including error message, error code and errorvalue). That takes `message`which specifies he error message such as validation_error, `error_code` which service specific API field error codes . and list of errors. `field` specifies fieldname and `help` redirects to the error description documentation.
Sample errors
 {
      "field":      field,
	  "message":    messages,
	  "error_code": "A1001",
	  "help":       "https://tuneverse.com/api-error/#errorcode_number",
}
for example:

`{
    "errors": {
        "data": [
            {
                "error_code": "",
                "field": "name",
                "help": "",
                "message": [
                    "This field is required"
                ]
            },
            {
                "error_code": "",
                "field": "email",
                "help": "",
                "message": [
                    "This field is required"
                ]
            }
        ]
    }
}`