## Overview
This provides functions which related to the version

## Index
- [GetVersionFromContext(ctx *gin.Context) string](#func-GetVersionFromContext)
- [PrepareVersionName(version string) string](#func-PrepareVersionName)


### func GetVersionFromContext

    GetVersionFromContext(ctx *gin.Context) string

This function is used to retrieve the `API version` from the `Gin context`. It takes a `Gin context` object ctx as input and returns the value of the `Accept-version` header.


### func PrepareVersionName

    PrepareVersionName(version string) string

This function is used to prepare an `API version` name for use in `URLs` or other `contexts` where periods are not allowed. It takes a string version as input and replaces all periods with `underscores`
