package controllers

import (
	"member/internal/consts"
	constant "member/internal/consts"
	"member/internal/entities"
	"member/internal/usecases"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/core/version"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/models/api"
	"gitlab.com/tuneverse/toolkit/utils"
)

// MemberController handles member-related HTTP requests and routes.
type MemberController struct {
	router   *gin.RouterGroup
	useCases usecases.MemberUseCaseImply
}

// NewMemberController creates a new instance of MemberController.
func NewMemberController(router *gin.RouterGroup, memberUseCase usecases.MemberUseCaseImply) *MemberController {
	return &MemberController{
		router:   router,
		useCases: memberUseCase,
	}
}

// InitRoutes initializes the routes for the MemberController.
//
// It sets up a route for the health check endpoint, where the version is a URL parameter.
// When a GET request is made to this endpoint, it calls the "HealthHandler" function.
//
// Params:
//
//	@member: required - the MemberController instance.

func (member *MemberController) InitRoutes() {
	member.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "HealthHandler")
	})
	member.router.POST("/:version/members/:member_id/billing-address", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "AddBillingAddress")
	})
	member.router.PATCH("/:version/members/:member_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "UpdateMember")
	})
	member.router.PATCH("/:version/members/:member_id/billing-address/:billing_address_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "UpdateBillingAddress")
	})
	member.router.DELETE("/:version/members/:member_id/billing-address/:billing_address_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "DeleteBillingAddress")
	})
	member.router.GET("/:version/members/:member_id/billing-address", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "GetAllBillingAddresses")
	})
	member.router.PATCH("/:version/members/:member_id/change-password", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "ChangePassword")
	})
	member.router.GET("/:version/members/:member_id/reset-password", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "ResetPassword")
	})
	member.router.POST("/:version/members", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "RegisterMember")
	})
	member.router.POST("/:version/members/:member_id/stores", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "AddMemberStores")
	})
	member.router.DELETE("/:version/members/:member_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "DeleteMember")
	})
	member.router.GET("/:version/members/:member_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "ViewMemberProfile")
	})
	member.router.GET("/:version/members", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "ViewMembers")
	})
	member.router.GET("/:version/members/oauth", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "GetBasicMemberDetailsByEmail")
	})
	member.router.POST("/:version/members/:member_id/subscriptions/checkout", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "SubscriptionCheckout")
	})
	member.router.PATCH("/:version/members/:member_id/subscriptions/renewal", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "SubscriptionRenewal")
	})
	member.router.PATCH("/:version/members/:member_id/subscriptions/cancel", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "SubscriptionCancellation")
	})
	member.router.PATCH("/:version/members/:member_id/subscriptions/product-switch", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "SubscriptionProductSwitch")
	})
	member.router.GET("/:version/members/:member_id/subscriptions", func(ctx *gin.Context) {
		version.RenderHandler(ctx, member, "ViewAllSubscriptions")
	})
}

// HealthHandler handles health check requests and responds with the server's health status.
//
// Params:
//
//	@ctx: required - the Gin context for the HTTP request.
func (member *MemberController) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "server running with base version",
	})
}

// AddBillingAddress handles adding a billing address for a member.
// It performs the following steps:
//   - Validates the endpoint and method.
//   - Parses the member_id from the URL parameters.
//   - Binds the JSON request body to the  struct.

//   - Handles validation errors and internal server errors.
//   - Logs relevant information about the process.
//
// Parameters:
//
//	@ctx (*gin.Context): The Gin context for handling the HTTP request.

func (member *MemberController) AddBillingAddress(ctx *gin.Context) {
	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Billing Address failed, endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)

	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Adding Billing Address failed, Failed to fetch error values from context")
		return
	}

	// Extract member_id from the URL parameters (UUID parsing)
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("AddBillingAddress failed, invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
		}
		return
	}
	// Bind the JSON request body to the entities package's BillingAddress struct
	var billingAddress entities.BillingAddress
	if err := ctx.ShouldBindJSON(&billingAddress); err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("AddBillingAddress failed, invalid JSON data: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// Call the use case to add a billing address
	fieldsMap, err := member.useCases.AddBillingAddress(ctx, memberID, billingAddress)
	// Call the use case to add a billing address
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("AddBillingAddress failed!! failed to add billing address: %s", err)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
		}
	}

	// If the fieldsMap has some mappings
	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx.Request.Context()).Errorf("AddBillingAddress failed, failed to add billing address")

		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Billing address added successfully
	logger.Log().WithContext(ctx.Request.Context()).Info("AddBillingAddress: Billing address added successfully")

	// Data added successfully
	ctx.JSON(http.StatusCreated, gin.H{
		"message": consts.SuccessfullyAdded,
	})
}

// AddMemberStores function is employed to incorporate stores into the member store table,
// with each store associated with the partner under which the member is categorized.
func (member *MemberController) AddMemberStores(ctx *gin.Context) {
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()
	var stores []string

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed, endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)

	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Add Member stores failed, Failed to fetch error values from context")
		return
	}

	// Extract member_id from the URL parameters (UUID parsing)
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed, invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
		}
		return
	}

	// Check if payload exists
	if ctx.Request.ContentLength == 0 {

		fieldsMap, err := member.useCases.AddMemberStores(ctx, memberID, stores)

		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed, failed to add stores: %s", err)
			val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
			if hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		}

		// If the fieldsMap has some mappings
		if len(fieldsMap) > 0 {
			fields := utils.FieldMapping(fieldsMap)
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed, failed to add stores")

			val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
			if hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		}

	} else {
		var storeList entities.AddMemberStores
		if err := ctx.ShouldBindJSON(&storeList); err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed, invalid JSON data: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data",
			})
			return
		}

		// Call the use case to add member stores.
		fieldsMap, err := member.useCases.AddMemberStores(ctx, memberID, storeList.Storelist)

		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed,, failed to add stores: %s", err)
			val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
			if hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		}

		// If the fieldsMap has some mappings
		if len(fieldsMap) > 0 {
			fields := utils.FieldMapping(fieldsMap)
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Add Member stores failed, failed to add stores")

			val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
			if hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		}
	}

	logger.Log().WithContext(ctx.Request.Context()).Info("Adding Member stores success: Stores added successfully")

	// Data added successfully
	ctx.JSON(http.StatusCreated, gin.H{
		"message": consts.SuccessfullyAddedMemberStore,
	})
}

// UpdateBillingAddress handles updating a billing address for a member.
// It handles updating a billing address for a member.
// It performs the following steps:
//   - Validates the endpoint and method.
//   - Parses the member_id from the URL parameters.
//   - Binds the JSON request body to the  struct.

//   - Handles validation errors and internal server errors.
//   - Logs relevant information about the process.
//
// Parameters:
//
//	@ctx (*gin.Context): The Gin context for handling the HTTP request.
func (member *MemberController) UpdateBillingAddress(ctx *gin.Context) {
	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)
	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Update billing address failed ,endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}
	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Updating Billing Address failed, Failed to fetch error values from context")
		return
	}

	// Extract member_id from the URL parameters (UUID parsing)
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)
	memberBillingStr := ctx.Param("billing_address_id")
	memberBillingID, err := uuid.Parse(memberBillingStr)
	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Update BillingAddress failed, invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")

		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Bind the JSON request body to the entities package's BillingAddress struct
	var updatedBillingAddress entities.BillingAddress
	if err := ctx.ShouldBindJSON(&updatedBillingAddress); err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Update BillingAddress failed, invalid JSON data: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// Call the use case to update the billing address
	fieldsMap, err := member.useCases.UpdateBillingAddress(ctx, memberID, memberBillingID, updatedBillingAddress)

	if err != nil {
		if err.Error() == "Member does not exist" {
			// Member not exist
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Update BillingAddress failed, member not exist: %s", err.Error())
			val, hasError, errorCode := utils.ParseFields(ctx, consts.NotFoundErr, "", contextError, "", "")
			if !hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		} else {
			// If an error occurs during the use case execution, log the error and return a 500 Internal Server Error response
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Update BillingAddress failed : Internal server error %s", err.Error())
			val, _, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx).Errorf("Update billing address failed, err = %s", fields)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)

		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Billing address updated successfully
	logger.Log().WithContext(ctx.Request.Context()).Info("Update BillingAddress: Billing address updated successfully")

	// Data updated successfully
	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.SuccessUpdated,
	})
}

// ChangePassword handles requests to change a member's password.
//
// This function handles password change requests and expects the following parameters:
//   - member_id (UUID): The unique identifier of the member whose password needs to be changed.
//   - NewPassword (string): The new password to set for the member.
//   - CurrentPassword (string): The current password to validate the change.
//
// It performs the following steps:
//   - Validates the endpoint and method.
//   - Parses the member_id from the URL parameters.
//   - Binds the JSON request body to the passwordChangeRequest struct.
//   - Calls the ChangePassword use case with the memberID and new password.
//   - Handles validation errors and internal server errors.
//   - Logs relevant information about the process.
//
// Parameters:
//
//	@ctx (*gin.Context): The Gin context for handling the HTTP request.
//
// HTTP Request:
//
//	POST /members/:member_id/change-password
//
// JSON Request Body:
//
//	{
//	  "NewPassword": "string",
//	  "CurrentPassword": "string"
//	}
//
// Returns:
//   - If successful, it responds with HTTP status 201 (Created) and a success message.
//   - If any errors occur during parsing, binding, or processing, it responds with an appropriate error message and status.
func (member *MemberController) ChangePassword(ctx *gin.Context) {
	// Get the HTTP request method (GET, POST, etc.) and convert it to lowercase.

	methods := ctx.Request.Method
	method := strings.ToLower(methods)

	// Get the full endpoint URL of the request.
	endpointURL := ctx.FullPath()

	// Get the contextEndpoints map to check if the endpoint exists.
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)

	// Get the endpoint details based on the URL and HTTP method.
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)

	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("ChangePassword failed, validation error, endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Change Password failed, Failed to fetch error values from context")
		return
	}

	// Parse the member_id from the URL parameters.
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	// Check for errors in member_id parsing.
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Password changing failed, invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Bind the JSON request body to the passwordChangeRequest struct.
	var passwordChangeRequest entities.PasswordChangeRequest
	if err := ctx.BindJSON(&passwordChangeRequest); err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Password change failed, invalid JSON data: %s", err.Error())

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid JSON data",
			"errors": err.Error(),
		})
		return
	}

	// Call the ChangePassword use case with the memberID and new password.
	fieldsMap, err := member.useCases.ChangePassword(ctx, memberID, passwordChangeRequest.Key, passwordChangeRequest.NewPassword, passwordChangeRequest.CurrentPassword)

	// Handle errors during password change.
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Password change failed: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Handle validation errors.
	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx).Errorf("Password Changing failed, validation error, err = %s", fields)
		val, hasValue, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		// Check if the ParseFields function returned a value (hasValue).
		if hasValue {
			logger.Log().WithContext(ctx.Request.Context()).Errorf("ChangePassword failed,: %s", val)
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Log password change completion.
	logger.Log().WithContext(ctx.Request.Context()).Info("Password change: Successfull")

	// Respond with success status.
	ctx.JSON(http.StatusCreated, gin.H{
		"message": consts.Success,
	})
}

// ResetPassword is a controller function responsible for initiating the password reset process for a member.
//
// Inputs:
// - ctx: The Gin context containing information about the HTTP request and response.
// Functionality:
// 1. Retrieves the HTTP request method and endpoint URL.,. Validates the endpoint's existence.
// 2. Parses the member ID from the URL parameters and validates its format.
// 3. Binds the incoming JSON request body to the `entities.ResetPassword` struct to fetch the email. Calls the `InitiatePasswordReset` use case to generate a reset key for the member based on the provided email.
// 4. Handles any errors or validation issues that may arise during the process., Logs the outcome of the password reset process, whether successful or with errors.
//
// Outputs:
// - JSON response containing the generated reset key if the process is successful.
// - Appropriate error responses if any validation fails or if there's an internal server error.
//
// HTTP Method: POST (since the function initiates a password reset which is typically achieved via a POST request).
//
// Error Handling:
// The function handles potential errors such as invalid member IDs, invalid JSON data, and any errors returned during the password reset process.
func (member *MemberController) ResetPassword(ctx *gin.Context) {
	// Get the HTTP request method (GET, POST, etc.) and convert it to lowercase.
	methods := ctx.Request.Method
	method := strings.ToLower(methods)

	// Get the full endpoint URL of the request.
	endpointURL := ctx.FullPath()

	// Get the contextEndpoints map to check if the endpoint exists.
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)

	// Get the endpoint details based on the URL and HTTP method.
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	_ = endpoint
	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Reset Password failed, validation error, endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Reset Password failed, Failed to fetch error values from context")
		return
	}

	// Parse the member_id from the URL parameters.
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	// Check for errors in member_id parsing.
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Reset Password  failed, invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Bind the JSON request body to the passwordChangeRequest struct.
	var initiatePasswordReset entities.ResetPassword
	if err := ctx.BindJSON(&initiatePasswordReset); err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Password Reset  failed, invalid JSON data: %s", err.Error())

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid JSON data",
			"errors": err.Error(),
		})
		return
	}

	// Call the ChangePassword use case with the memberID and new password.
	key, fieldsMap, err := member.useCases.InitiatePasswordReset(ctx, memberID, initiatePasswordReset.Email)

	// Handle errors during password change.
	if err != nil {

		logger.Log().WithContext(ctx.Request.Context()).Errorf("Password reset failed: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	// Handle validation errors.
	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx).Errorf("Password Reset Initiation failed, validation error, err = %s", fields)
		val, hasValue, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		// Check if the ParseFields function returned a value (hasValue).
		if hasValue {
			logger.Log().WithContext(ctx.Request.Context()).Errorf("Reset Initiation failed,: %s", val)
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Log password change completion.
	logger.Log().WithContext(ctx.Request.Context()).Info("Password reset: Initiated Succesfully")

	// Respond with success status.
	ctx.JSON(http.StatusCreated, gin.H{
		"Key":     key,
		"message": consts.SuccessfullyInitiated,
	})
}

// GetAllBillingAddresses handles the HTTP request to add a billing address to a member's account.
// It validates the endpoint and method, parses the member_id from the URL parameters,
// and binds the JSON request body to the GetAllBillingAddresses struct.
// It then calls the AddBillingAddress use case with the extracted memberID and billingAddressRequest.
// It handles validation errors, "No record found" errors, and internal server errors.
// Finally, it responds with a success message if the billing address is added successfully.
//
// Parameters:
//
//	@ctx (*gin.Context): The Gin context for handling the HTTP request.
//
// Returns:
//   - If successful, it responds with HTTP status 200 (OK) and a success message.
//   - If any errors occur during validation, "No record found" errors, or processing, it responds with an appropriate error message and status.

func (controller *MemberController) GetAllBillingAddresses(ctx *gin.Context) {
	// Get the HTTP request method (GET, POST, etc.) and convert it to lowercase.
	methods := ctx.Request.Method
	method := strings.ToLower(methods)
	endpointUrl := ctx.FullPath()

	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)

	// Attempt to retrieve the specific endpoint associated with the current request
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)
	_ = endpoint

	// If the endpoint does not exist, log an error and return a 400 Bad Request response
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Viewing Billing Address Failed, validation error, endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Viewing all Billing Address failed:Failed to load Context errors")
		return
	}

	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)
	// If parsing fails, log an error and return an appropriate JSON response

	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("View All Billing Address: Invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	params := entities.Params{}
	if page, exists := ctx.GetQuery("page"); !exists || page == "" {
		params.Page = consts.DefaultPage
	}
	if limit, exists := ctx.GetQuery("limit"); !exists || limit == "" {
		params.Limit = consts.DefaultLimit
	}
	// Call the GetAllBillingAddresses use case with the extracted memberID
	fieldsMap, billingAddresses, _, err := controller.useCases.GetAllBillingAddresses(ctx, memberID, params)

	//Checks the length of fieldMap for checking is there any validation error reported or not.
	if len(fieldsMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Viewing All Billing Addresses Failed")
		fields := utils.FieldMapping(fieldsMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		if hasError {
			logger.Log().WithContext(ctx).Errorf("Viewing Billing Address")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	if err != nil {
		// Check if the error message contains "No record found"
		if strings.Contains(err.Error(), "No record found") {
			// Handle the "No record found" error
			logger.Log().WithContext(ctx.Request.Context()).Errorf("View All Billing Address: Failed to retrieve billing address, err=%s", err.Error())
			val, hasError, errorCode := utils.ParseFields(ctx, consts.NotFoundErr, "", contextError, "", "")
			if hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		} else {
			// Handle other errors, e.g., internal server error
			logger.Log().WithContext(ctx.Request.Context()).Errorf("View All Billing Address: Failed to retrieve billing address, err=%s", err.Error())
			val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
			if hasError {
				ctx.JSON(int(errorCode), val)
				return
			}
		}
	}
	if len(billingAddresses) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No Billing details to show",
			"data": billingAddresses})
		return
	}
	// Send the response as JSON with a success message

	ctx.JSON(http.StatusOK, entities.Response{
		BillingAddress: billingAddresses,
		Message:        consts.SuccessfullyListed,
	})

	// Log the success message
	logger.Log().WithContext(ctx.Request.Context()).Info("View All Billing Address: Billing Address Listed successfully")
}

// RegisterMember is for registraion of member
// This endpoint expects a JSON payload representing the member's registration data.
// Request JSON Body:
//
//	-Member: The JSON object containing member details.
//
// Response:
//   - If successful, it returns a JSON response with status code 201 (Created).
//   - If there is an error in the request (e.g., invalid JSON) or processing, it returns an
//     error JSON response with the appropriate status code and error message.
func (member *MemberController) RegisterMember(ctx *gin.Context) {
	ctxt := ctx.Request.Context()
	methods := ctx.Request.Method
	method := strings.ToLower(methods)
	endpointUrl := ctx.FullPath()

	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists in the context
	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("Member registration failed,Invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Member registration failed:Failed to load Context errors")
		return
	}

	var reqIn entities.Member
	if err := ctx.BindJSON(&reqIn); err != nil {
		logger.Log().WithContext(ctx).Errorf("Member registration failed, Invalid JSON data, err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}
	// Retrieve PartnerID from the Gin context
	partnerID, exists := ctx.Get("partner_id")

	if !exists {
		// Handle the case where PartnerID is not found in the context
		logger.Log().WithContext(ctx).Error(" ")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to get PartnerID from context",
			"errors":    nil,
		})
		return
	}
	// Convert the PartnerID to a string
	partnerIDStr, ok := partnerID.(string)
	if !ok || len(partnerIDStr) == 0 {
		// Log the error
		logger.Log().WithContext(ctx).Error("Failed to Extract PartnerID from Context")

		// Return an error response
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to set PartnerID in context",
			"errors":    nil,
		})
		return
	}

	//Call the RegisterMember
	fieldMap, err := member.useCases.RegisterMember(ctxt, reqIn, contextError, partnerIDStr, endpoint, method)

	//Checks the length of fieldMap for checking is there any validation error reported or not.
	if len(fieldMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Member registration failed: validation error")
		// to build the field format.
		fields := utils.FieldMapping(fieldMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	if err != nil {
		//For logging error message
		logger.Log().WithContext(ctx).Errorf("Member registration failed: database error err=%s", err.Error())
		//Retrieve error based on the api's requirement from the errors which are stored in context
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	ctx.JSON(http.StatusCreated, consts.SuccessfullyRegistered)

	// Log the success message
	logger.Log().WithContext(ctx).Info("Member Registration: Member Registered successfully")
}

// UpdateMember is for registraion of member
// This endpoint expects a JSON payload representing the member's data updation.
// Request JSON Body:
//
//	-Member: The JSON object containing member details.
//
// Response:
//   - If successful, it returns a JSON response with status code 201 (Created).
//   - If there is an error in the request (e.g., invalid JSON) or processing, it returns an
//     error JSON response with the appropriate status code and error message.
func (member *MemberController) UpdateMember(ctx *gin.Context) {
	var memberID uuid.UUID

	methods := ctx.Request.Method
	method := strings.ToLower(methods)
	endpointUrl := ctx.FullPath()
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists in the context
	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("Member profile updation failed,Invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Member profile updation  failed:Failed to load Context errors")
		return
	}

	var args entities.Member
	if err := ctx.BindJSON(&args); err != nil {
		logger.Log().WithContext(ctx).Errorf("Member profile updation failed, Invalid JSON data, err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	//Call the RegisterMember
	fieldMap, err := member.useCases.UpdateMember(ctx, memberID, args)

	//Checks the length of fieldMap for checking is there any validation error reported or not.
	if len(fieldMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Member profile updation failed: validation error")
		fields := utils.FieldMapping(fieldMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		if hasError {
			logger.Log().WithContext(ctx).Errorf("Member profile updation failed")
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	if err != nil {
		//For logging error message
		logger.Log().WithContext(ctx).Errorf("Member profileupdation  failed: database error err=%s", err.Error())
		//Retrieve error based on the api's requirement from the errors which are stored in context
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	ctx.JSON(http.StatusOK, consts.SuccessfullyUpdated)
	logger.Log().WithContext(ctx).Info("Member Profile Updated: Profile successfully Updated ")
}

// ViewMemberProfile retrieves a member's profile based on the provided member_id.
// Parameters:
//
//	- ctx (gin.Context): The Gin context for handling the HTTP request.
//	- memberID: The unique identifier of the member to retrieve.
//
// Returns:
//  - memberProfile: The member's profile information.
//	- error: An error, if any, during the database operation.

func (member *MemberController) ViewMemberProfile(ctx *gin.Context) {
	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()
	ctxt := ctx.Request.Context()

	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists in the context
	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("View Member Profile:Invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("View Member Profile failed:Failed to load Context errors")
		return
	}

	// Extract member_id from the URL parameters (UUID parsing)
	fieldsMap := map[string][]string{}
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("View Member Profile failed-invalid member_id: err=%s", err.Error())
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Invalid)
		fields := utils.FieldMapping(fieldsMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	fieldMap, memberProfile, err := member.useCases.ViewMemberProfile(ctx, ctxt, memberID, contextError, endpoint, method)

	//Checks the length of fieldMap for checking is there any validation error reported or not.
	if len(fieldMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("View Member Profile failed:  validation error")
		// to build the field format.
		fields := utils.FieldMapping(fieldMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	//Checks any database error reported or not.
	if err != nil {
		//For logging error message
		logger.Log().WithContext(ctx).Errorf("View Member Profile failed: err=%s", err.Error())
		//Retrieve error based on the api's requirement from the errors which are stored in context
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Member profile details retreived successfully", "data": memberProfile})
	// Log the success message
	logger.Log().WithContext(ctx).Info("View Member Profile: Member profile details retreived successfully")
}

// ViewMembers retrieves a list of members based on the provided query parameters.
//
// Parameters:
//   - ctx (gin.Context): The Gin context for handling the HTTP request.
//
// Returns:
//   - data ([]entities.Member): A list of member data matching the query parameters.
//   - error: An error, if any, during the database operation.

func (member *MemberController) ViewMembers(ctx *gin.Context) {
	// ctxt := ctx.Request.Context()
	_, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)

	// Check if the endpoint exists in the context
	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("View Members failed:Invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("View Members failed:Failed to load Context errors")
		return
	}

	params := entities.Params{}

	// Check and set default values if any of the headers are empty
	// Check if the header is empty AND there is no corresponding query parameter
	if status, exists := ctx.GetQuery("status"); !exists || status == "" {
		params.Status = consts.DefaultStatus
	} else {
		params.Status = status
	}

	if page, exists := ctx.GetQuery("page"); !exists || page == "" {
		params.Page = consts.DefaultPage
	}
	if limit, exists := ctx.GetQuery("limit"); !exists || limit == "" {
		params.Limit = consts.DefaultLimit
	}
	if country, exists := ctx.GetQuery("country"); !exists || country == "" {
		params.Country = consts.DefaultCountry
	}

	if gender, exists := ctx.GetQuery("gender"); !exists || gender == "" {
		params.Gender = consts.DefaultGender
	} else {
		gender = strings.ToUpper(gender)
		if gender != "M" && gender != "F" {
			ctx.JSON(http.StatusBadRequest, entities.FailureResponse{
				Status:  "failure",
				Code:    http.StatusBadRequest,
				Message: "Invalid gender parameter",
				Errors:  nil,
			})
			return
		}
		params.Gender = gender
	}

	if sortBy, exists := ctx.GetQuery("sort"); !exists || sortBy == "" {
		params.SortBy = consts.DefaultSortBy
		params.Order = consts.DefaultOrder
	} else {
		// Check if the sortby parameter is provided and set the value accordingly
		switch sortBy {
		case "firstname":
			params.SortBy = "firstname"
			params.Order = "ASC"
		case "lastname":
			params.SortBy = "lastname"
			params.Order = "ASC"
		case "created_on":
			params.SortBy = "created_on"
		case "email":
			params.SortBy = "email"
		default:
			// If an invalid value is provided, set the default to "created_on"
			params.SortBy = consts.DefaultSortBy
			params.Order = consts.DefaultOrder
			// If an invalid value is provided, return a Bad Request response
			ctx.JSON(http.StatusBadRequest, entities.FailureResponse{
				Status:  "failure",
				Code:    http.StatusBadRequest,
				Message: "Invalid sort parameter",
				Errors:  nil,
			})
			return
		}
	}

	if order, exists := ctx.GetQuery("order"); exists {
		switch order {
		case "asc":
			params.Order = "ASC"
		case "desc":
			params.Order = "DESC"
		default:
			params.Order = "ASC"
		}
	} else {
		// If "order" is not present, default to sorting by asc of firstname
		params.Order = "ASC"

	}

	if partner, exists := ctx.GetQuery("partner"); !exists || partner == "" {
		params.Partner = consts.DefaultPartner
	}
	if role, exists := ctx.GetQuery("role"); !exists || role == "" {
		params.Role = consts.DefaultRole
	}
	if search, exists := ctx.GetQuery("search"); !exists || search == "" {
		params.Search = consts.DefaultSearch
	}

	// Call ViewMembers to retrieve member data
	membersData, metaData, err := member.useCases.ViewMembers(ctx, params)

	if err != nil {
		// Check if the error message matches the predefined constant error
		if err.Error() == "invalid param value" {
			ctx.JSON(http.StatusBadRequest, entities.FailureResponse{
				Status:  "failure",
				Code:    http.StatusBadRequest,
				Message: "invalid parameter value",
				Errors:  nil,
			})
			return
		}
	}
	if err != nil {
		// Log the error message
		logger.Log().WithContext(ctx).Errorf("View members failed: err=%s", err.Error())

		// Check if the error message matches the predefined constant error
		if err.Error() == "Exceeds maximum record limit" {
			ctx.JSON(http.StatusTooManyRequests, entities.FailureResponse{
				Status:  "failure",
				Code:    http.StatusTooManyRequests,
				Message: "Exceeds maximum record limit",
				Errors:  nil,
			})
			return
		}

		// Retrieve error based on the API's requirement from the errors stored in context
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	if len(membersData) == 0 {
		ctx.JSON(http.StatusOK, entities.MemberResponse{
			Status:  "success",
			Code:    200,
			Message: "successfully retreived all member details",
			DataResp: entities.DataResponse{
				Metadata: metaData,
				// Data:     membersData,
			},
		})
	} else {
		ctx.JSON(http.StatusOK, entities.MemberResponse{
			Status:  "success",
			Code:    200,
			Message: "Successfully Listed All Member Details",
			DataResp: entities.DataResponse{
				Metadata: metaData,
				Data:     membersData,
			},
		})
	}

	logger.Log().WithContext(ctx).Info("View members: Members details retreived successfully")
}

// GetBasicMemberDetailsByEmail retrieves the basic details of a member based on their email.
//
// Parameters:
//   - ctx (gin.Context): The Gin context for handling the HTTP request.
//
// Returns:
//   - None: The function sends a JSON response with the member's basic details.
func (member *MemberController) GetBasicMemberDetailsByEmail(ctx *gin.Context) {
	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()
	//ctxt := ctx.Request.Context()

	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists in the context
	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("View Basic Member Details:Invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("View Basic Member Details failed:Failed to load Context errors")
		return
	}

	// Retrieve PartnerID from the Gin context
	partnerID, exists := ctx.Get("partner_id")

	if !exists {
		// Handle the case where PartnerID is not found in the context
		logger.Log().WithContext(ctx).Error("Failed to get PartnerID from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to get PartnerID from context",
			"errors":    nil,
		})
		return
	}
	// Convert the PartnerID to a string
	partnerIDStr, ok := partnerID.(string)

	if len(partnerIDStr) == 0 {
		logger.Log().WithContext(ctx).Error("Failed to Extract PartnerID from Context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to set PartnerID in context",
			"errors":    nil,
		})
		return

	}
	if !ok {
		// Handle the case where PartnerID is not in the expected format (uuid)
		logger.Log().WithContext(ctx).Error("PartnerID is not in the expected format (uuid)")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "PartnerID is not in the expected format (uuid)",
			"errors":    nil,
		})
		return
	}

	var reqIn entities.MemberPayload
	if err := ctx.BindJSON(&reqIn); err != nil {
		logger.Log().WithContext(ctx).Errorf("View failed, Invalid JSON data, err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	fieldMap, basicMemberData, err := member.useCases.GetBasicMemberDetailsByEmail(ctx, partnerIDStr, reqIn, contextError, endpoint, method)

	if err != nil {
		// For logging error message
		logger.Log().WithContext(ctx).Errorf("View Basic Member Details failed: err=%s", err.Error())
		// Retrieve error based on the API's requirement from the errors stored in context
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	if len(fieldMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("View Basic Member Details failed:  validation error")
		// Build the field format.
		fields := utils.FieldMapping(fieldMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Member basic details retrieved successfully", "data": basicMemberData, "partnerId": partnerIDStr})
	// Log the success message
	logger.Log().WithContext(ctx).Info("View Basic Member Details: Member basic details retrieved successfully")
}

// SubscriptionCheckout handles the process of checking out a subscription for a member.
// This function performs the following steps:
//  1. Extracts the memberID from the URL.
//  2. Deserializes the JSON request body into a CheckoutSubscription object.
//  3. Performs any necessary validation on the checkoutData.
//  4. Calls the HandleSubscriptionCheckout use case function.
//  5. Responds with appropriate JSON results based on success or failure.
//
// Parameters:
//   - ctx (*gin.Context): The Gin context for handling the HTTP request.
func (member *MemberController) SubscriptionCheckout(ctx *gin.Context) {

	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription checkout failed , endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription checkout failed,Failed to load context errors")
		return
	}

	// Extract memberID from the URL
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)
	// Retrieve PartnerID from the Gin context
	partnerID, exists := ctx.Get("partner_id")

	if !exists {
		// Handle the case where PartnerID is not found in the context
		logger.Log().WithContext(ctx).Error(" ")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to get PartnerID from context",
			"errors":    nil,
		})
		return
	}
	// Convert the PartnerID to a string
	partnerIDStr, ok := partnerID.(string)
	if !ok || len(partnerIDStr) == 0 {
		// Log the error
		logger.Log().WithContext(ctx).Error("Failed to Extract PartnerID from Context")

		// Return an error response
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to set PartnerID in context",
			"errors":    nil,
		})
		return
	}
	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription checkout  failed, invalid member_id: %s", err.Error())
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription checkout/purchase  failed")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	// Deserialize the JSON request body into a CheckoutSubscription object
	var checkoutData entities.CheckoutSubscription
	if err := ctx.BindJSON(&checkoutData); err != nil {
		// Respond with an error message if JSON data is invalid
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription checkout  failed, invalid JSON data: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}
	// Call the HandleSubscriptionCheckout use case function
	fieldsMap, err := member.useCases.HandleSubscriptionCheckout(ctx, memberID, checkoutData, partnerIDStr)

	if err != nil {
		// If an error occurs during the use case execution, log the error and return a 500 Internal Server Error response
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Checkout  failed, failed to checkout subscription: %s", err.Error())
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasVal {
			ctx.JSON(int(errorCode), val)
			return
		}
		return
	}
	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Checkout  failed, failed to checkout subscription")
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		if hasVal {
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	// If successful, return a JSON response indicating success
	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.SuccessfullyCheckedout,
	})
}

// SubscriptionRenewal handles the process of checking out a subscription for a member.
// This function performs the following steps:
//  1. Extracts the memberID from the URL.
//  2. Deserializes the JSON request body into a CheckoutSubscription object(same for renewal also).
//  3. Performs any necessary validation on the plan renewal
//  4. Calls the HandleSubscriptionRenewal use case function.
//  5. Responds with appropriate JSON results based on success or failure.
//
// Parameters:
//   - ctx (*gin.Context): The Gin context for handling the HTTP request.
func (member *MemberController) SubscriptionRenewal(ctx *gin.Context) {

	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription renewal failed , endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription renewal failed,Failed to load context errors")
		return
	}
	// Extract memberID from the URL
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)
	// Retrieve PartnerID from the Gin context
	partnerID, exists := ctx.Get("partner_id")

	if !exists {
		// Handle the case where PartnerID is not found in the context
		logger.Log().WithContext(ctx).Error(" ")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to get PartnerID from context",
			"errors":    nil,
		})
		return
	}
	// Convert the PartnerID to a string
	partnerIDStr, ok := partnerID.(string)
	if !ok || len(partnerIDStr) == 0 {
		// Log the error
		logger.Log().WithContext(ctx).Error("Failed to Extract PartnerID from Context")

		// Return an error response
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to set PartnerID in context",
			"errors":    nil,
		})
		return
	}
	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription renewal   failed, invalid member_id: %s", err.Error())
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription renewal  failed")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	// Deserialize the JSON request body into a CheckoutSubscription object
	var checkoutData entities.SubscriptionRenewal
	if err := ctx.BindJSON(&checkoutData); err != nil {
		// Respond with an error message if JSON data is invalid
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription renewal  failed, invalid JSON data: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}
	// Call the HandleSubscriptionCheckout use case function
	fieldsMap, err := member.useCases.HandleSubscriptionRenewal(ctx, memberID, checkoutData, partnerIDStr)

	if err != nil {
		// If an error occurs during the use case execution, log the error and return a 500 Internal Server Error response
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Renewal  failed, failed to renew subscription: %s", err.Error())
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription checkout/purchase  failed")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Renewal  failed, failed to renew subscription")
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription renewal failed")
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// If successful, return a JSON response indicating success
	ctx.JSON(http.StatusCreated, gin.H{
		"message": consts.SuccessfullyRenewed,
	})
}

// SubscriptionCancellation handles the process of canceling a subscription for a member.
// This function performs the following steps:
//  1. Extracts the memberID from the URL.
//  2. Deserializes the JSON request body into a CancelSubscription object.
//  3. Performs any necessary validation on the plan cancellation.
//  4. Calls the HandleSubscriptionCancellation use case function.
//  5. Responds with appropriate JSON results based on success or failure.
//
// Parameters:
//   - ctx (*gin.Context): The Gin context for handling the HTTP request.
func (member *MemberController) SubscriptionCancellation(ctx *gin.Context) {

	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription cancellation failed, endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription cancellation failed,Failed to load context errors")
		return
	}

	// Extract memberID from the URL
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)
	// Retrieve PartnerID from the Gin context
	partnerID, exists := ctx.Get("partner_id")

	if !exists {
		// Handle the case where PartnerID is not found in the context
		logger.Log().WithContext(ctx).Error(" ")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to get PartnerID from context",
			"errors":    nil,
		})
		return
	}
	// Convert the PartnerID to a string
	partnerIDStr, ok := partnerID.(string)
	if !ok || len(partnerIDStr) == 0 {
		// Log the error
		logger.Log().WithContext(ctx).Error("Failed to Extract PartnerID from Context")

		// Return an error response
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": http.StatusInternalServerError,
			"message":   "Failed to set PartnerID in context",
			"errors":    nil,
		})
		return
	}
	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription cancellation failed, invalid member_id: %s", err.Error())
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription cancellation failed")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	// Deserialize the JSON request body into a cancellationData object
	var cancellationData entities.CancelSubscription
	if err := ctx.BindJSON(&cancellationData); err != nil {
		// Respond with an error message if JSON data is invalid
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Subscription cancellation failed, invalid JSON data: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}
	// Call the HandleSubscriptionCancellation use case function
	fieldsMap, err := member.useCases.HandleSubscriptionCancellation(ctx, memberID, cancellationData, partnerIDStr)

	if err != nil {
		// If an error occurs during the use case execution, log the error and return a 500 Internal Server Error response
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Cancellation failed, failed to cancel subscription: %s", err.Error())
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription cancellation failed:Internal Server Error ")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Cancellation failed, failed to cancel subscription")
		val, hasVal, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)
		if hasVal {
			logger.Log().WithContext(ctx.Request.Context()).Error("Subscription cancellation failed,Validations Erros")
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	// If successful, return a JSON response indicating success
	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.SuccessfullyCancelled,
	})

}

// SubscriptionProductSwitch handles the product switch between subscription plans choosed by a member.
func (member *MemberController) SubscriptionProductSwitch(ctx *gin.Context) {

	var data entities.SwitchSubscriptions
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Invalid member id err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest,
			gin.H{
				"message": "inavalid member id",
				"error":   err.Error(),
			})
		return
	}

	ctxt := ctx.Request.Context()
	err = ctx.BindJSON(&data)

	if err != nil {
		logger.Log().WithContext(ctxt).Errorf("Switch products between subscriptions failed err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest,
			gin.H{
				"message": "Binding failed",
				"error":   err.Error(),
			})
		return
	}

	methods := ctx.Request.Method
	method := strings.ToLower(methods)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)

	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("Switch product between subscriptions failed: invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, isErrorExists := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)

	if !isErrorExists {
		logger.Log().WithContext(ctxt).Errorf("Switch products between subscriptions failed, err = %s", consts.ContextErr)
		ctx.JSON(http.StatusBadRequest, gin.H{

			"errorCode": http.StatusBadRequest,
			"message":   consts.ContextErr,
			"errors":    nil,
		})
		return
	}
	// Call the SubscriptionProductSwitch use case with the extracted memberID and data
	validationErrors, err := member.useCases.SubscriptionProductSwitch(ctx, memberID, data)

	if len(validationErrors) != 0 {
		fields := utils.FieldMapping(validationErrors)
		val, errVal, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)

		if errVal {
			ctx.JSON(int(errorCode), val)
		}
		return
	}

	if err != nil {
		logger.Log().WithContext(ctxt).Errorf("Switch products between subscriptions failed err=%s", err.Error())
		val, errVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")
		if errVal {
			ctx.JSON(int(errorCode), val)
		}
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Switch products between subscriptions failed",
			"error":   err.Error(),
		})
		return
	}

	// Product switching between subscriptions success
	ctx.JSON(http.StatusOK, gin.H{
		"message": constant.SuccessSwitched,
	})

	// Log the success message
	logger.Log().WithContext(ctx.Request.Context()).Info("Switch products between subscriptions: Successfully Switched Products between Subscriptions plans")
}

// ViewAllSubscriptions handles the view of all subscription plans choosed by a member.
func (member *MemberController) ViewAllSubscriptions(ctx *gin.Context) {

	//var metaData entities.MetaData
	var (
		reqParam     entities.ReqParams
		responseData = api.Response{
			Message: consts.ValidationErr,
			Code:    http.StatusBadRequest,
			Errors:  struct{}{},
			Data:    struct{}{},
			Status:  "failure",
		}
	)

	if err := ctx.BindQuery(&reqParam); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	methods := ctx.Request.Method
	method := strings.ToLower(methods)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)

	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("Get subscriptions for member failed: invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, isErrorExists := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)

	if !isErrorExists {
		logger.Log().WithContext(ctx).Errorf("Switch products between subscriptions failed, err = %s", consts.ContextErr)
		ctx.JSON(http.StatusBadRequest, gin.H{

			"errorCode": http.StatusBadRequest,
			"message":   consts.ContextErr,
			"errors":    nil,
		})
		return
	}

	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)

	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Get subscriptions for member failed: Invalid member_id: %s", err.Error())
		val, errVal, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")
		if errVal {
			ctx.JSON(int(errorCode), val)
		}
		return
	}

	if status, exists := ctx.GetQuery("status"); !exists || status == "" {
		reqParam.Status = consts.DefaultStatus
	}

	if search, exists := ctx.GetQuery("search"); !exists || search == "" {
		reqParam.Search = consts.DefaultSearch
	}

	if sortBy, exists := ctx.GetQuery("sort"); !exists || sortBy == "" {
		reqParam.Sort = consts.DefaultSortByID
	}

	// Call the ViewAllSubscriptions use case with the extracted memberID
	memberSubscriptions, metadata, validationErrors, err := member.useCases.ViewAllSubscriptions(ctx, memberID, reqParam)

	if err != nil {
		if err.Error() == consts.MaximumRequestError {

			resp := api.Response{
				Status:  "failure",
				Message: "Too many request",
				Code:    http.StatusTooManyRequests,
				Data:    struct{}{},
				Errors:  err.Error(),
			}
			ctx.JSON(http.StatusTooManyRequests, resp)
			return
		}

		if strings.Contains(err.Error(), "No record found") {
			// Handle the "No record found" error
			_, errVal, _ := utils.ParseFields(ctx, consts.NotFound, "", contextError, "", "")
			if errVal {
				logger.Log().WithContext(ctx).Errorf("Get subscriptions for member, Failed to retrieve subscriptions")
				responseData.Message = consts.NotFound
				responseData.Code = http.StatusNotFound
				responseData.Status = "failure"
				ctx.JSON(http.StatusNotFound, responseData)
				return
			}
			return
		}

		_, valStatus, _ := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")

		if valStatus {
			// Handle Internal server error
			logger.Log().WithContext(ctx).Errorf("Get subscriptions for member, error while parsing data")
			responseData.Message = consts.InternalServerErr
			responseData.Code = http.StatusInternalServerError
			responseData.Status = "failure"
			ctx.JSON(http.StatusInternalServerError, responseData)
			return
		}

		return
	}

	if len(validationErrors) != 0 {
		logger.Log().WithContext(ctx).Errorf("Get subscriptions for member failed: validation error")
		fields := utils.FieldMapping(validationErrors)
		val, errVal, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)
		if errVal {
			ctx.JSON(int(errorCode), val)
		}
		return
	}

	successResponse := entities.SuccessResponse{
		Code:     constant.StatusOk,
		Message:  constant.SuccessfullyListedPlans,
		Metadata: metadata,
		Data:     memberSubscriptions,
	}

	ctx.JSON(http.StatusOK, successResponse)

	// Log the success message
	logger.Log().WithContext(ctx.Request.Context()).Info("View All Subscriptions: Subscription plans listed successfully")
}

//DeleteBillingAddress

func (member *MemberController) DeleteBillingAddress(ctx *gin.Context) {
	// Retrieve and preprocess request details
	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()

	// Check if the endpoint exists in the context
	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)
	// Check if the endpoint exists, if not, respond with a validation error.
	if !isEndpointExists {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Delete billing address failed ,endpoint does not exist in the database.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}
	// Get the contextError map to handle error responses.
	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Deleting Billing Address failed, Failed to fetch error values from context")
		return
	}

	// Extract member_id from the URL parameters (UUID parsing)
	memberIDStr := ctx.Param("member_id")
	memberID, err := uuid.Parse(memberIDStr)
	memberBillingStr := ctx.Param("billing_address_id")
	memberBillingID, err := uuid.Parse(memberBillingStr)
	// If parsing fails, log an error and return an appropriate JSON response
	if err != nil {
		logger.Log().WithContext(ctx.Request.Context()).Errorf("Delete BillingAddress failed, invalid member_id: %s", err.Error())
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr, "", contextError, "", "")

		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Call the use case to update the billing address
	fieldsMap, err := member.useCases.DeleteBillingAddress(ctx, memberID, memberBillingID)

	if len(fieldsMap) > 0 {
		fields := utils.FieldMapping(fieldsMap)
		logger.Log().WithContext(ctx).Errorf("Delete billing address failed, err = %s", fields)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr, fields, contextError, endpoint, method)

		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	// Billing address updated successfully
	logger.Log().WithContext(ctx.Request.Context()).Info("Delete BillingAddress: Billing address updated successfully")

	// Data updated successfully
	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.SuccessDelete,
	})
}

//DeleteMember delets a member from database

func (member *MemberController) DeleteMember(ctx *gin.Context) {
	methods := ctx.Request.Method
	method := strings.ToLower(methods)
	endpointUrl := ctx.FullPath()

	contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)

	// Check if the endpoint exists in the context
	if !isEndpointExists {
		logger.Log().WithContext(ctx).Errorf("Member Deletion failed,Invalid endpoint")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.EndpointErr,
			"errors":    nil,
		})
		return
	}

	contextError, errVal := utils.GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	if !errVal {
		logger.Log().WithContext(ctx.Request.Context()).Error("Member Deletion failed:Failed to load Context errors")
		return
	}

	memberID := ctx.Param("member_id")
	//Call the RegisterMember
	fieldMap, err := member.useCases.DeleteMember(ctx, memberID)
	if err != nil {
		//For logging error message
		logger.Log().WithContext(ctx).Errorf("Member registration failed: database error err=%s", err.Error())
		//Retrieve error based on the api's requirement from the errors which are stored in context
		val, hasError, errorCode := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "")
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}
	//Checks the length of fieldMap for checking is there any validation error reported or not.
	if len(fieldMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Member Deletion failed: validation error")
		// to build the field format.
		fields := utils.FieldMapping(fieldMap)
		val, hasError, errorCode := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method)
		if hasError {
			ctx.JSON(int(errorCode), val)
			return
		}
	}

	ctx.JSON(http.StatusCreated, consts.SuccessfullyDeleted)

	// Log the success message
	logger.Log().WithContext(ctx).Info("Member Deletion: Member Deleted successfully")
}
