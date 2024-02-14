# Middleware Package
The middleware package contains common `middleware` functions for handling various aspects of HTTP requests in our Go application.

## APIVersionGuard Middleware

### Description
The `APIVersionGuard` middleware is designed to handle API versioning by checking the `Accept-Version` header in incoming HTTP requests. It ensures that the requested API version is supported by the system.


### How to Use
To use the `APIVersionGuard` middleware, follow these steps:

- Import the middleware package in your Go application.
```go
import (
    "gitlab.com/tuneverse/toolkit/middleware"
)
```


- Define your API version lookup logic or use the default c.Param("version") to extract the version from the request URL path.

- Create a `VersionOptions` struct to configure the middleware:

    **VersionParamLookup**: A function that takes a `*gin.Context` and returns the version string. You can provide your custom logic here. If you want to use the default behavior of extracting the version from the URL path, set this field to nil.

    **AcceptedVersions**: A slice of strings representing the API versions accepted by your system.

```go
versionOptions := middleware.VersionOptions{
    VersionParamLookup: nil, // or your custom version lookup function
    AcceptedVersions:   []string{"v1", "v2"},
}

---

versionOptions := middleware.VersionOptions{
    VersionParamLookup: func(c *gin.Context) string {
        return c.Query("version")
    },
    AcceptedVersions:   []string{"v1", "v2"},
}
```

- Apply the middleware to your Gin router using the `APIVersionGuard` function, passing in the `versionOptions` you created.
```go
router := gin.Default()

// Apply APIVersionGuard middleware
router.Use(middleware.APIVersionGuard(versionOptions))
```

- Now, incoming requests will be checked for a valid API version. If the version is missing, not supported, or valid, the middleware will handle it accordingly.