
## Overview
This provides functions for making API requests, managing headers, etc. It is designed to simplify the process of communicating with external services and handling HTTP responses.


## Index
- [APIRequest(method string, url string, headers map[string]interface{},body map[string]interface{}) (*http.Response, error)](#func-APIRequest)
- [SETHeaders(request http.Request, headers map[string]interface{}) http.Request](#func-SETHeaders)


### func APIRequest

    APIRequest(method string, url string, headers map[string]interface{},body map[string]interface{}) (*http.Response, error)

This function is used to make `API` requests. It takes the `HTTP` method, `URL`, `headers`, and `body` as parameters and returns the `HTTP response` and an error, if any. The function uses the `HTTPClient` variable, which is an instance of the `http.Client` struct, to send the request. The response status is logged using the `log.Infof` function.


### func SETHeaders

    SETHeaders(request http.Request, headers map[string]interface{}) http.Request

This function is used to set the headers for an `HTTP request`. It takes an `http.Request` object and a map of headers as parameters and returns the modified `http.Request` object. The function iterates over the headers map and adds each `key-value` pair to the request's headers using the `request.Header.Add` method