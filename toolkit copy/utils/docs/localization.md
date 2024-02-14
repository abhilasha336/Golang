## Overview
This module deals with function for extracting desired errors from all errors and the errorresponse format 

## Index
- [ ExtractFieldErrors(field string, errorsValue bson.M, endpoint string, method string) (map[string]any, bool) 
    ](#func-ExtractFieldErrors)

- [ ErrorResponseFormat(message string, errorCode any, errorsValue map[string]any, key string) any](#func-ErrorResponseFormat)

### func ExtractFieldErrors

    ExtractFieldErrors(field string, errorsValue bson.M, endpoint string, method string) (map[string]any, bool) 

    This function is used to extract desired errors from all errors based on the criteria. It takes `field`, which is a string is of the form "field:value1|value2", `errorsValue`, which is the all errors retrieved from database, `endpoint`, which is a string to specify the endpoint in which the error to be add, `method`, which is a string specifies the method of the api. This function returns the desired error format and a bool value.




### func ErrorResponseFormat

  ErrorResponseFormat(message string, errorCode any, errorsValue map[string]any, key string) any

  This function is used to format the errors in a desired format.
  The general format is as given below
    {
        message   
        errorCode 
        errors    
    }

This function takes `message` , which is a string ie., the message to be display, `errorCode` which is an integer specifies the error code , `errorsValue`, which is a map of desired error, `key` which is a string used for retrieving the errors.




