# Error Localization Middleware
The `ErrorLocalization` middleware is used to retrieve error response data from a localization service and store it in the context for use in downstream handlers. This middleware is useful for applications that need to provide localized error messages to users.


## Usage
To use the `ErrorLocalization` middleware, you must first create an instance of the `ErrorLocaleOptions` struct and pass it to the middleware function. The `ErrorLocaleOptions` struct contains the following fields:

- `Cache`: A cache interface used to store the error response data.
- `CacheExpiration`: The expiration time for the error response data in the cache.
- `CacheKeyLabel`: The cache key label for the error response data.
- `ContextErrorResponse`: The context key label for the error response data.
- `LocalisationServiceURL`: The URL of the localization service.
- `HeaderLanguage`: The header label for the language.


```go
    type ErrorLocaleOptions struct {
        Cache                  cache         `validate:"required"`
        CacheExpiration        time.Duration `validate:"required"`
        CacheKeyLabel          string        `validate:"required"`
        ContextErrorResponse   string
        LocalisationServiceURL string `validate:"required"`
        HeaderLanguage         string
    }
```


To add the `ErrorLocalization` middleware to a Gin router, you can use the following code:

- Import the middleware package in your Go application.
```go
import (
    "gitlab.com/tuneverse/toolkit/middleware"
)
```

```go
    router.Use(middleware.ErrorLocalization(middleware.ErrorLocaleOptions{
        Cache:                  myCache,
        CacheExpiration:        time.Minute,
        CacheKeyLabel:          "error_data",
        ContextErrorResponse:   "error_responses",
        LocalisationServiceURL: "https://my-localization-service.com",
        HeaderLanguage:         "Accept-Language",
    }))

```