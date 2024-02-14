# Package Version
This package contains a function called `RenderHandler` that is used to execute a given method on an object. The function checks if a version method exists and executes it instead of the given method if it exists.


## Function
### RenderHandler
```go
    func RenderHandler(ctx *gin.Context, object interface{}, method string, args ...interface{})
```

This function takes in four parameters:

- **ctx**: A pointer to a gin.Context object.
- **object**: An interface object.
- **method**: A string representing the method to be executed.
- **args**: A variadic parameter of type `interface{}` representing the arguments to be passed to the method.




## Variables

### ErrAcceptedVersionNotFound
```go
    var ErrAcceptedVersionNotFound = "unable to find the accepted versions in context"
```
This variable is a string representing an error message when the accepted versions are not found in the context

### ErrHeaderVersionNotFound
```go
    var ErrHeaderVersionNotFound = "unable to find the version in context"
```

This variable is a string representing an error message when the version is not found in the context.
