## Overview
This provides functions for context

## Index
- [GetContext[T any](ctx *gin.Context, name string) (T, bool)](#func-GetContext)
- [GetHeader(ctx *gin.Context, header string) string](#func-GetHeader)


### func GetContext

    GetContext[T any](ctx *gin.Context, name string) (T, bool)

This function is used to retrieve any type of value from the `Gin context`. It takes a `Gin context` object `ctx` and a string name as input and returns the value associated with the `name` and a `boolean` indicating whether the value exists in the context.

### func GetHeader

    GetHeader(ctx *gin.Context, header string) string

This function is used to retrieve a `header` value from the `Gin context`. It takes a `Gin context` object `ctx` and a string header as input and returns the value of the header.
