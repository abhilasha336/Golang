## Locale Middleware

## Overview
The `Localize` middleware function is designed to handle language localization in a api application. It extracts the language from the request headers, sets it in the context, and passes it to the next middleware or route handler.


### How to Use
To use the `Localize` middleware, follow these steps:

- Import the middleware package in your Go application.
```go
import (
    "gitlab.com/tuneverse/toolkit/middleware"
)
```


```go
	// Add the LocalizationLanguage middleware
	router.Use(middleware.Localize(LocaleOptions{
		HeaderLabel:  "Accept-Language",
		ContextLabel: "lang",
	}))

```