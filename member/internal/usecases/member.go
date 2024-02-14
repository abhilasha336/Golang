package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"member/internal/consts"
	"member/internal/entities"
	"member/internal/repo"
	"member/utilities"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
	cryptoHash "gitlab.com/tuneverse/toolkit/utils/crypto"
)

// MemberUseCases defines use cases related to member operations.
type MemberUseCases struct {
	repo repo.MemberRepoImply
}

// MemberUseCaseImply interface
type MemberUseCaseImply interface {
	// AddBillingAddress adds a billing address to a member's profile.
	// It takes the context, memberID, and billingAddress as input, and returns a map of validation error messages and an error, if any.
	AddBillingAddress(ctx *gin.Context, memberID uuid.UUID, billingAddress entities.BillingAddress) (map[string][]string, error)

	// UpdateBillingAddress updates an existing billing address in a member's profile.
	// It takes the context, memberID, and billingAddress as input, and returns a map of validation error messages and an error, if any.
	UpdateBillingAddress(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID, billingAddress entities.BillingAddress) (map[string][]string, error)
	//DeleteBillingAddress deletes a members billing address
	DeleteBillingAddress(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID) (map[string][]string, error)
	// UpdateMember updates a member's profile.
	// It takes the context, memberID, and updated member profile details as input, and returns a map of validation error messages and an error, if any.
	UpdateMember(ctx *gin.Context, memberID uuid.UUID, args entities.Member) (map[string][]string, error)
	// GetAllBillingAddresses retrieves all billing addresses associated with a member.
	// It takes the context and memberID as input, and returns a slice of BillingAddress entities and an error, if any.
	GetAllBillingAddresses(ctx *gin.Context, memberID uuid.UUID, params entities.Params) (map[string][]string, []entities.BillingAddress, models.MetaData, error)

	// ChangePassword updates a member's password.
	// It takes the context, memberID, new password, and current password as input, and returns a map of validation error messages and an error, if any.
	ChangePassword(ctx *gin.Context, memberID uuid.UUID, key string, newPassword string, currentPassword string) (map[string][]string, error)

	// InitiatePasswordReset initiates a password reset for a member.
	// It takes the context, memberID, and email as input, and returns a reset token, a map of validation error messages, and an error, if any.
	InitiatePasswordReset(ctx *gin.Context, memberID uuid.UUID, email string) (string, map[string][]string, error)

	// IsMemberExists checks if a member with the specified memberId exists.
	// It takes the memberId and context as input, and returns true if the member exists, false otherwise, and an error if one occurs.
	IsMemberExists(memberId uuid.UUID, ctx context.Context) (bool, error)

	// RegisterMember registers a new member.
	// It takes the context, member details, contextError, partnerID, endpoint, and method as input, and returns headers and an error.
	RegisterMember(ctxt context.Context, args entities.Member, contextError map[string]any, partnerID, endpoint, method string) (map[string][]string, error)

	// ViewMemberProfile retrieves a member's profile.
	// It takes the context, memberId, contextError, endpoint, and method as input, and returns headers, the member profile, and an error.
	ViewMemberProfile(ctx *gin.Context, ctxt context.Context, memberId uuid.UUID, contextError map[string]any, endpoint string, method string) (map[string][]string, entities.MemberProfile, error)

	// ViewMembers retrieves a list of members.
	// It takes the context and parameters as input, and returns a slice of ViewMembers and metadata, along with an error if any.
	ViewMembers(ctx *gin.Context, params entities.Params) ([]entities.ViewMembers, models.MetaData, error)

	// GetBasicMemberDetailsByEmail retrieves basic member details by email.
	// It takes the context, partnerID, member payload, contextError, endpoint, and method as input, and returns headers, basic member data, and an error.
	GetBasicMemberDetailsByEmail(ctx *gin.Context, partnerID string, args entities.MemberPayload, contextError map[string]any, endpoint string,
		method string) (map[string][]string, entities.BasicMemberData, error)

	// HandleSubscriptionCheckout handles the checkout process for a subscription.
	HandleSubscriptionCheckout(ctx *gin.Context, memberID uuid.UUID, checkoutData entities.CheckoutSubscription, partnerIDStr string) (map[string][]string, error)

	// HandleSubscriptionRenewal handles the renewal process for an existing subscription.
	HandleSubscriptionRenewal(ctx *gin.Context, memberID uuid.UUID, checkoutData entities.SubscriptionRenewal, partnerIDStr string) (map[string][]string, error)

	// HandleSubscriptionCancellation handles the cancellation process for an active subscription.
	HandleSubscriptionCancellation(ctx *gin.Context, memberID uuid.UUID, checkoutData entities.CancelSubscription, partnerIDStr string) (map[string][]string, error)
	//SubscriptionProductSwitch switches a product from one active subscription plan to another(based on criterias)
	SubscriptionProductSwitch(context.Context, uuid.UUID, entities.SwitchSubscriptions) (map[string][]string, error)
	//ViewAllSubscriptions list all subscriptions of a member
	ViewAllSubscriptions(*gin.Context, uuid.UUID, entities.ReqParams) ([]entities.ListAllSubscriptions, models.MetaData, map[string][]string, error)
	// GetSubscriptionRecordCount gets the total count of a members subscription.
	GetSubscriptionRecordCount(context.Context, uuid.UUID) (int64, error)
	// IsMemberExists checks if a member with the specified memberId exists.
	IsMemberExist(context.Context, uuid.UUID) (bool, error)
	DeleteMember(ctx *gin.Context, memberID string) (map[string][]string, error)
	AddMemberStores(ctx *gin.Context, memberID uuid.UUID, stores []string) (map[string][]string, error)
}

// GracePeriodError represents an error indicating that the subscription is in the grace period.
type GracePeriodError struct {
	Message string
}

// Error implements the error interface for GracePeriodError.
func (e GracePeriodError) Error() string {
	return e.Message
}

// NewMemberUseCases is a constructor for creating an instance of MemberUseCases.
func NewMemberUseCases(memberRepo repo.MemberRepoImply) MemberUseCaseImply {
	return &MemberUseCases{
		repo: memberRepo,
	}
}

// AddBillingAddress handles adding a billing address for a member.
//
// This function validates the required fields in the billing address, validates the ZIP code,
// and then calls the repository function to add the billing address.
//
// Parameters:
//   - ctxt: The context for the operation.
//   - memberID: The UUID of the member for whom the billing address is being added.
//   - billingAddress: The billing address details to be added.
//   - contextError: A map used to collect error messages related to context validation.
//   - endpoint: A string representing the endpoint used for logging purposes.
//   - method: A string representing the HTTP method used for logging purposes.
//
// Returns:
//   - If successful, it returns nil (no error).
//   - If any required fields are missing or there's an error in adding the billing address,
//     it returns an error with an appropriate error message.
func (member *MemberUseCases) AddBillingAddress(ctx *gin.Context, memberID uuid.UUID, billingAddress entities.BillingAddress) (map[string][]string, error) {
	// Initialize the map to store validation error messages
	fieldsMap := map[string][]string{}
	// Check if the member exists
	exists, err := member.repo.IsMemberExists(memberID, ctx)
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		return fieldsMap, nil
	}
	if err != nil {
		// Log the error
		logger.Log().WithContext(ctx).Errorf("Failed to check if member exists: %s", err.Error())
		return nil, err
	}
	// Step 1: Check if member already has the maximum allowed billing addresses
	billingAddressCount, err := member.repo.GetBillingAddressCountForMember(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check billing address count: %v", err)
		return nil, err
	}
	if billingAddressCount >= 5 {
		utils.AppendValuesToMap(fieldsMap, consts.MaximumAddressCount, consts.Limit)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Maximum limit reached")
		return fieldsMap, nil
	}

	// Validate the required fields in the billing address

	// Check if address is empty
	if billingAddress.Address == "" || len(billingAddress.Address) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.Address, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Missing required fields")
		return fieldsMap, nil
	}
	if len(billingAddress.Address) > 200 {
		utils.AppendValuesToMap(fieldsMap, consts.Address, consts.Maximum)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Address exceeds maximum length of 150 characters")
		return fieldsMap, nil
	}
	// Check if zipcode is empty
	if len(billingAddress.Zipcode) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.Zipcode, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Missing required fields")
		return fieldsMap, nil
	}

	// Validate the ZIP code
	// Check if ZIP code length is not between 4-10 characters
	if len(billingAddress.Zipcode) < 4 || len(billingAddress.Zipcode) > 10 {
		utils.AppendValuesToMap(fieldsMap, consts.Zipcode, consts.ZipFormat)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Invalid ZIP code length")
	} else {
		// Validate ZIP code format using utility function
		if isValidZip := utilities.IsValidZIPCode(billingAddress.Zipcode); !isValidZip {
			utils.AppendValuesToMap(fieldsMap, consts.Zipcode, consts.Invalid)
		}
	}

	// Validate country
	if len(billingAddress.Country) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.Country, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Country cannot be empty")
		return fieldsMap, nil
	}
	// Check if the country exists
	countryExists, err := member.repo.CountryExists(billingAddress.Country)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check country existence: %v", err)
		return nil, err
	}
	if !countryExists {
		utils.AppendValuesToMap(fieldsMap, consts.Country, consts.CountryNotExist)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Country does not exist")
		return fieldsMap, nil
	}

	// Check if the state exists within the provided country
	stateExists, err := member.repo.StateExists(billingAddress.State, billingAddress.Country)
	if len(billingAddress.State) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.State, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Country cannot be empty")
	}
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check state existence: %v", err)
		return nil, err
	}
	if !stateExists {
		utils.AppendValuesToMap(fieldsMap, consts.State, consts.StateNotExist)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: State does not exist")
		return fieldsMap, nil
	}
	// If any validation errors are found, return the error map
	if len(fieldsMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address due to validation errors")
		return fieldsMap, nil
	}

	// Check if the billing address already exists for the member
	exists, err = member.repo.BillingAddressExists(ctx, memberID, billingAddress)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check billing address existence: %v", err)
		return nil, err
	}
	if exists {
		utils.AppendValuesToMap(fieldsMap, consts.BillingAddress, consts.AlreadyExists)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: This billing address already exists for the member")
		return fieldsMap, nil
	}

	// Check if there's already a primary billing address for the member
	hasPrimary, err := member.repo.HasPrimaryBilling(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check billing address existence: %v", err)
		return nil, err
	}
	if hasPrimary && billingAddress.Primary {
		utils.AppendValuesToMap(fieldsMap, consts.Primary, consts.HasPrimary)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Member already has a primary billing address")
		return fieldsMap, nil
	}
	countPrimary, err := member.repo.CountPrimaryBillingAddresses(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address:failed to get primary address count")
		return fieldsMap, nil
	}
	if billingAddressCount == 4 && countPrimary < 1 && !billingAddress.Primary {
		utils.AppendValuesToMap(fieldsMap, consts.Primary, consts.PrimaryMandatory)
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address:One primary billing is mandatory")
		return fieldsMap, nil
	}
	// If all validations pass, call the repository function to add the billing address
	err = member.repo.AddBillingAddress(ctx, memberID, billingAddress)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to add billing address: %s", err.Error())
		return nil, err
	}

	// Return the fields map if needed for further processing or reporting
	return fieldsMap, nil
}

// UpdateMember updates a member's profile.
// Parameters:
//
//	@ ctx: The context for the database operation.
//	@ memberID: The UUID of the member whose profile is being updated.
//	@ args: The updated member profile details.
//
// Returns:
//
//	@ error: An error, if any, during the database operation.
//
// UpdateMember updates a member's profile.
// Parameters:
//
//	@ ctx: The context for the database operation.
//	@ memberID: The UUID of the member whose profile is being updated.
//	@ args: The updated member profile details.
//
// Returns:
//
//	@ fieldsMap: A map of validation errors.
//	@ error: An error, if any, during the database operation.
//
// UpdateMember updates a member's profile.
func (member *MemberUseCases) UpdateMember(ctx *gin.Context, memberID uuid.UUID, args entities.Member) (map[string][]string, error) {
	// Initialize a map to store validation errors
	fieldsMap := map[string][]string{}
	// Check if the member exists
	exists, err := member.repo.IsMemberExists(memberID, ctx)

	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		return fieldsMap, err
	}

	// Fetch the current member details before the update operation
	currentMember, err := member.repo.GetMemberByID(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to fetch current member details: %v", err)
		return nil, err
	}

	// Flag to track if Country has been updated
	countryUpdated := args.Country != currentMember.Country.String

	if len(args.Title) != 0 {
		if args.Title != "Ms" && args.Title != "Mr" {
			utils.AppendValuesToMap(fieldsMap, consts.Title, consts.Invalid)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid salutation")
		}
	}
	if len(args.FirstName) != 0 {
		// Validate First Name
		if len(args.FirstName) > 15 || !utilities.ValidateName(args.FirstName) {
			utils.AppendValuesToMap(fieldsMap, consts.FirstName, consts.Valid)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid First name")
		}
	}
	if len(args.Gender) != 0 {
		if len(args.Gender) != 1 || (args.Gender != "M" && args.Gender != "F") {
			utils.AppendValuesToMap(fieldsMap, consts.Gender, consts.Invalid)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid gender")
		}
	}
	if len(args.Language) != 0 {
		if len(args.Language) > 2 {
			utils.AppendValuesToMap(fieldsMap, consts.Language, consts.InvalidFormat)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid language Format")
		}
		langExists, err := member.repo.CheckLanguageExist(ctx, args.Language)
		if err != nil {
			return nil, err
		}
		if !langExists {
			utils.AppendValuesToMap(fieldsMap, consts.Language, consts.Invalid)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid language/Language Not found")
		}

	}
	// Validate Last Name
	if len(args.LastName) != 0 {
		if len(args.LastName) > 15 || !utilities.ValidateName(args.LastName) {
			utils.AppendValuesToMap(fieldsMap, consts.LastName, consts.Valid)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid last name")
		}
	}
	// Validate Country
	if args.Country != "" {
		countryExists, err := member.repo.CountryExists(args.Country)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to check country existence: %v", err)
			return nil, err
		}
		if !countryExists {
			utils.AppendValuesToMap(fieldsMap, consts.Country, consts.CountryNotExist)
			logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Country does not exist")
			return fieldsMap, nil
		}
	}

	// Validate State only if Country has been updated
	if countryUpdated && args.State == currentMember.State.String {
		utils.AppendValuesToMap(fieldsMap, consts.State, consts.UpdateStateWithCountry)
		logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Update State as Country has been changed")
	}

	// Validate State
	if args.State != "" && args.Country != "" {
		stateExists, err := member.repo.StateExists(args.State, args.Country)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to check state existence: %v", err)
			return nil, err
		}
		if !stateExists {
			utils.AppendValuesToMap(fieldsMap, consts.State, consts.StateNotExist)
			logger.Log().WithContext(ctx).Errorf("Failed to add billing address: State does not exist")
			return fieldsMap, nil
		}
	}
	if len(args.Zipcode) != 0 {
		// Validate ZIP code
		if len(args.Zipcode) < 4 || len(args.Zipcode) > 10 || !utilities.IsValidZIPCode(args.Zipcode) {
			utils.AppendValuesToMap(fieldsMap, consts.Zipcode, consts.ZipFormat)
			logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Invalid ZIP code format or length")
		}
	}
	if len(args.Phone) != 0 || len(args.Phone) > 12 {
		// Validate phone number

		if !utilities.IsValidPhoneNumber(args.Phone, args.Country) {
			utils.AppendValuesToMap(fieldsMap, consts.PhoneNumber, consts.Format)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Invalid phone number format")
			return fieldsMap, nil
		}
	}
	// Validate City, Address1, and Address2 lengths
	if len(args.City) != 0 {
		if len(args.City) > 100 {
			utils.AppendValuesToMap(fieldsMap, consts.City, consts.Maximum)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: City exceeds maximum characters")
		}
	}
	if len(args.Address1) != 0 {
		if len(args.Address1) > 200 {
			utils.AppendValuesToMap(fieldsMap, consts.Address1, consts.Maximum)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Address1 exceeds maximum characters")
		}
	}
	if len(args.Address2) != 0 {
		if len(args.Address2) > 200 {
			utils.AppendValuesToMap(fieldsMap, consts.Address2, consts.Maximum)
			logger.Log().WithContext(ctx).Errorf("Profile updation failed, validation error: Address2 exceeds maximum characters")
		}
	}
	// If there are validation errors, return them
	if len(fieldsMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Profile updation failed due to validation errors")
		return fieldsMap, nil
	}

	// Perform the update operation
	err = member.repo.UpdateMember(ctx, memberID, args)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Profile Updation failed, internal server error: %s", err.Error())
		return nil, err
	}

	// If everything is successful, return nil
	return nil, nil
}

// GetAllBillingAddresses retrieves all billing addresses associated with a member.
// Parameters:
//
//	@ ctx: The context for the database operation.
//	@ memberID: The UUID of the member for whom billing addresses are being retrieved.
//
// Returns:
//
//	@ []entities.BillingAddress: A slice of billing addresses associated with the member.
//	@ error: An error, if any, during the retrieval process.
func (member *MemberUseCases) GetAllBillingAddresses(ctx *gin.Context, memberID uuid.UUID, params entities.Params) (map[string][]string, []entities.BillingAddress, models.MetaData, error) {
	// Call the repository method to retrieve all billing addresses
	// Initialize a map to store validation errors

	fieldsMap := map[string][]string{}

	// Declare variables for page and limit
	var pageInt, limitInt int64

	// Check if the page parameter has any value, either from the query or default
	if page, exists := ctx.GetQuery("page"); exists && page != "" {
		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			logger.Log().Errorf("Failed to parse page value: %v", err)

		} else {
			params.Page = int16(pageInt)
		}
	}
	if limit, exists := ctx.GetQuery("limit"); exists && limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			// Handle error
			logger.Log().Errorf("Failed to parse limit value: %v", err)
			// Handle the error as needed, maybe set a default value or return an error
		} else {
			params.Limit = int16(limitInt)
		}
	}

	// Check if the member exists
	exists, err := member.repo.IsMemberExists(memberID, ctx)
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		return fieldsMap, nil, models.MetaData{}, err
	}

	billingAddressCount, err := member.repo.GetBillingAddressCountForMember(ctx, memberID)
	if err != nil {
		return nil, nil, models.MetaData{}, fmt.Errorf("failed to check billing address count: %v", err)
	}
	if billingAddressCount < 1 {
		utils.AppendValuesToMap(fieldsMap, consts.BillingCount, consts.NoBillingAddress)
		logger.Log().WithContext(ctx).Errorf("No billing address found for this member")
		return fieldsMap, nil, models.MetaData{}, err
	}

	billingAddresses, err := member.repo.GetAllBillingAddresses(ctx, memberID)

	metadata := &models.MetaData{
		CurrentPage: int32(pageInt) + 1,
		PerPage:     int32(limitInt) + consts.DefaultLimitAddress,
		Total:       int64(billingAddressCount),
	}

	metadata = utils.MetaDataInfo(metadata)
	if err != nil {
		return fieldsMap, nil, models.MetaData{}, err
	}

	return nil, billingAddresses, *metadata, nil
}

// UpdateBillingAddress handles updating a billing address for a member.

// This function checks for missing required fields in the billing address, validates the ZIP code,
// and then calls the repository function to update the billing address.
//
// Parameters:
//
//	@ ctxt: The context for the operation.
//	@ memberID: The UUID of the member for whom the billing address is being updated.
//	@ billingAddress: The updated billing address details.
//	@ contextError: A map used to collect error messages related to context validation.
//	@ endpoint: A string representing the endpoint used for logging purposes.
//	@ method: A string representing the HTTP method used for logging purposes.
//
// Returns:
//   - If successful, it returns nil (no error).
//   - If any required fields are missing or there's an error in updating the billing address,
//     it returns an error with an appropriate error message.
func (member *MemberUseCases) UpdateBillingAddress(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID, billingAddress entities.BillingAddress) (map[string][]string, error) {

	// Validate the required fields in the billing address
	fieldsMap := map[string][]string{}

	// Check if the member exists
	exists, err := member.repo.IsMemberExists(memberID, ctx)
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
	}

	if err != nil {
		// Log the error
		logger.Log().WithContext(ctx).Errorf("Failed to check if member exists: %s", err.Error())
		return nil, err
	}
	IsValid, err := member.repo.CheckBillingAddressRelation(ctx, memberID, memberBillingID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to update billing address %s", err.Error())
		return nil, err
	}
	if !IsValid {
		utils.AppendValuesToMap(fieldsMap, consts.MemberBilling, consts.InvalidRelation)
		logger.Log().WithContext(ctx).Errorf("Failed to update billing address")
		return fieldsMap, nil
	}
	if len(billingAddress.Address) != 0 {
		if len(billingAddress.Address) > 200 {
			utils.AppendValuesToMap(fieldsMap, consts.Address, consts.Maximum)
			logger.Log().WithContext(ctx).Errorf("Failed to update billing address: Address exceeds maximum length of 150 characters")
			return fieldsMap, nil
		}
	}
	currentBillingAddress, err := member.repo.GetBillingAddressByID(ctx, memberBillingID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to fetch current billing address: %s", err.Error())
		return nil, err
	}
	totalAddress, err := member.repo.CountTotalAddressesForMember(ctx, memberID)
	if totalAddress == 1 {
		if currentBillingAddress.Primary && !billingAddress.Primary {
			utils.AppendValuesToMap(fieldsMap, consts.Primary, consts.PrimaryMandatory)
			logger.Log().WithContext(ctx).Errorf("Failed to update billing address")
			return fieldsMap, nil
		}
	}
	// Check if there's already a primary billing address for the member
	hasPrimary, err := member.repo.HasPrimaryBilling(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Error fetching primary billing: %v", err)
		return fieldsMap, err
	}
	if hasPrimary {
		if billingAddress.Primary {
			utils.AppendValuesToMap(fieldsMap, consts.Primary, consts.HasPrimary)
			logger.Log().WithContext(ctx).Errorf("Failed to update billing address: Already has primary billing address")
			return fieldsMap, nil
		}
	}

	// if currentBillingAddress.Primary && !billingAddress.Primary {
	// 	err = member.repo.UpdatePrimaryBillingAddressToFalseAndRandom(ctx, memberID, memberBillingID)
	// 	if err != nil {
	// 		logger.Log().WithContext(ctx).Errorf("Failed to update billing address: %s", err.Error())
	// 	}
	// }

	// Validate the ZIP code , Check if ZIP code length is not between 4-10 characters
	if len(billingAddress.Zipcode) != 0 {
		if len(billingAddress.Zipcode) < 4 || len(billingAddress.Zipcode) > 10 {
			utils.AppendValuesToMap(fieldsMap, consts.Zipcode, consts.ZipFormat)
			logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Invalid ZIP code length")
		} else {
			// Validate ZIP code format using utility function
			if isValidZip := utilities.IsValidZIPCode(billingAddress.Zipcode); !isValidZip {
				utils.AppendValuesToMap(fieldsMap, consts.Zipcode, consts.Invalid)
			}
		}
	}
	if len(billingAddress.Country) != 0 {
		//  Validate country and state
		countryExists, err := member.repo.CountryExists(billingAddress.Country)
		if len(billingAddress.Country) == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.Country, consts.Required)
			logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Country cannot be empty")
			return fieldsMap, nil
		}
		if err != nil || !countryExists {
			utils.AppendValuesToMap(fieldsMap, consts.Country, consts.CountryNotExist)
			logger.Log().WithContext(ctx).Errorf("Failed to update billing address: Country does not exist")
			return fieldsMap, nil
		}
	}
	if len(billingAddress.State) != 0 {
		stateExists, err := member.repo.StateExists(billingAddress.State, billingAddress.Country)
		if len(billingAddress.State) == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.State, consts.Required)
			logger.Log().WithContext(ctx).Errorf("Failed to add billing address: Country cannot be empty")
			return fieldsMap, nil
		}
		if err != nil || !stateExists {
			utils.AppendValuesToMap(fieldsMap, consts.State, consts.StateNotExist)
			logger.Log().WithContext(ctx).Errorf("Failed to update billing address: State does not exist")
			return fieldsMap, nil
		}
	}

	if len(fieldsMap) != 0 {
		logger.Log().WithContext(ctx).Errorf("Failed to update billing address")
	}

	// Call the repository function to add the billing address
	err = member.repo.UpdateBillingAddress(ctx, memberID, memberBillingID, billingAddress)

	// Check if there was an error
	if err != nil {

		logger.Log().WithContext(ctx).Errorf("Failed to update billing address: %s", err.Error())

		return nil, err
	}

	return fieldsMap, nil
}

// ChangePassword updates a member's password.
// This function validates the new password, checks if it meets the minimum length requirement,
// and verifies that it doesn't match the current password. It then hashes the new password,
// compares it with the stored hash of the current password, and updates the password in the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - memberID: The UUID identifying the member.
//   - newPassword: The new password to set.
//   - currentPassword: The current password for validation.
//
// Returns:
//   - A map of field names to error messages, if any validation errors occur.
//   - An error if there is a database operation failure or other error.
func (member *MemberUseCases) ChangePassword(ctx *gin.Context, memberID uuid.UUID, key string, newPassword string, currentPassword string) (map[string][]string, error) {
	// Initialize a fields map to hold validation errors
	fieldsMap := map[string][]string{}

	exists, err := member.repo.IsMemberExists(memberID, ctx)
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, Member not found")
		return fieldsMap, nil
	}
	// Validate the key
	if len(key) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.Key, consts.Required)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error:Not long enough")

	}
	if len(key) != 16 {
		utils.AppendValuesToMap(fieldsMap, consts.Key, consts.Match)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error:check key length")
		return fieldsMap, nil

	}
	// Use checkResetKeyMatch to verify if the provided key matches the stored key
	matches, err := member.repo.CheckResetKeyMatch(ctx, memberID, key)
	// Fetch the reset key value for the member ID
	keyValue := member.repo.GetResetKey(ctx, memberID)
	// Check if the keys do not match or if the key is invalidated
	if !matches || keyValue == "invalidated" {
		if keyValue == "invalidated" {
			utils.AppendValuesToMap(fieldsMap, consts.Key, consts.InvalidKey)
			logger.Log().WithContext(ctx).Errorf("ChangePassword failed, The key is invalid/timeout")
			return fieldsMap, nil
		} else {
			utils.AppendValuesToMap(fieldsMap, consts.Key, consts.Match)
			logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error: Wrong key")
			return fieldsMap, nil
		}
	}
	// Retrieve the hashed current password from the repository
	// Fetch the hashed current password from the repository
	hashedCurrentPassword, err := member.repo.GetPasswordHash(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, failed to fetch password hash")
		return nil, fmt.Errorf("failed to fetch hashed password: %w", err)
	}

	// Hash the provided current password using your cryptoHash package
	hashedPassword, err := cryptoHash.Hash(currentPassword)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error:Problem with current password hash")
		return nil, fmt.Errorf("failed to hash provided password: %w", err)
	}

	if len(currentPassword) == 0 || currentPassword == " " {
		utils.AppendValuesToMap(fieldsMap, consts.CurrentPassword, consts.Required)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error:Current password is empty")
		return fieldsMap, nil
	}
	if hashedPassword != hashedCurrentPassword {
		utils.AppendValuesToMap(fieldsMap, consts.CurrentPassword, consts.Incorrect)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error: Current password is incorrect")

		// Return the fields map without logging the error (already logged)
		return fieldsMap, nil
	}
	if err != nil {
		// Log the error if password retrieval fails
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, internal server error: %s", err.Error())
		return nil, err
	}

	// Check if the new password is the same as the current password
	if newPassword == currentPassword {
		utils.AppendValuesToMap(fieldsMap, consts.NewPassword, consts.InvalidPassword)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error: New password is the same as the current password")

		// Return the fields map without logging the error (already logged)
		return fieldsMap, nil
	}
	// Check if the new password is empty or its length is less than 8
	if len(newPassword) < 8 || newPassword == "" {
		utils.AppendValuesToMap(fieldsMap, consts.NewPassword, consts.MinLength)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error:Not long enough")
		return fieldsMap, nil
	}
	// Additional password format validation (you can customize this)
	if err := utilities.ValidatePassword(newPassword); err != nil {
		utils.AppendValuesToMap(fieldsMap, consts.NewPassword, consts.Format)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, validation error:Doesnot satisfy required format")
		return fieldsMap, nil

	}

	if len(fieldsMap) == 0 {
		// Hash the new password
		hashedNewPassword, err := cryptoHash.Hash(newPassword)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("ChangePassword failed, Failed to update password: %s", err.Error())
			return nil, err
		}
		err = member.repo.UpdatePassword(ctx, memberID, key, hashedNewPassword)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("ChangePassword failed, internal server error: Failed to update password: %s", err.Error())
			return nil, err
		}
	}
	// Return the fields map (if there are any validation errors, it will be non-empty)
	return fieldsMap, nil
}

// IsMemberExists checks if a member with the given memberId exists.
// Parameters:
//
//	@ memberId: The UUID of the member to check for existence.
//	@ ctx: The context for the database operation.
//
// Returns:
//
//	@ bool: True if the member exists, false otherwise.
//	@ error: An error, if any, during the check.
func (member *MemberUseCases) IsMemberExists(memberId uuid.UUID, ctx context.Context) (bool, error) {
	_, err := member.repo.IsMemberExists(memberId, ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// RegisterMember registers a new member and validates the input fields.
//
// Parameters:
//   - ctxt (context.Context): The context for the database operation.
//   - args (entities.Member): The member data to be registered.
//   - contextError (map[string]any): A map containing context-specific error responses.
//   - endpoint (string): The endpoint name.
//   - method (string): The HTTP request method (e.g., "POST" or "PUT").
//
// Returns:
//   - fieldsMap (map[string][]string): A map of validation error messages.
//   - error: An error, if any, encountered during the registration process.
func (member *MemberUseCases) RegisterMember(ctxt context.Context, args entities.Member, contextError map[string]any, partnerID, endpoint,
	method string) (map[string][]string, error) {

	fieldsMap := map[string][]string{}

	if args.FirstName != "" && !utilities.ValidateMaximumNameLength(args.FirstName) {
		utils.AppendValuesToMap(fieldsMap, consts.FirstName, consts.Length)
	}

	if args.FirstName != "" && !utilities.ValidateName(args.FirstName) {
		utils.AppendValuesToMap(fieldsMap, consts.FirstName, consts.Valid)
	}

	if args.LastName != "" && !utilities.ValidateMaximumNameLength(args.LastName) {
		utils.AppendValuesToMap(fieldsMap, consts.LastName, consts.Length)
	}

	if args.LastName != "" && !utilities.ValidateName(args.LastName) {
		utils.AppendValuesToMap(fieldsMap, consts.LastName, consts.Valid)
	}

	if args.Email == "" {
		utils.AppendValuesToMap(fieldsMap, consts.Email, consts.Required)
	}
	if len(args.Email) > 50 {
		utils.AppendValuesToMap(fieldsMap, consts.Email, consts.TooLong)

	}
	if args.Provider != consts.ProviderSpotify {
		if !utils.IsValidEmail(args.Email) {
			utils.AppendValuesToMap(fieldsMap, consts.Email, consts.Valid)
		}
	}

	emailExists, err := member.repo.CheckEmailExists(ctxt, partnerID, args.Email)

	if err != nil {
		logger.Log().WithContext(ctxt).Errorf("Failed to check an email exists: err=%s", err.Error())
		return nil, err
	}

	if emailExists {
		utils.AppendValuesToMap(fieldsMap, consts.Email, consts.Exists)

	}

	if args.Provider == consts.ProviderInternal && args.Password == "" {
		utils.AppendValuesToMap(fieldsMap, consts.Password, consts.Required)
	}

	if args.Provider == consts.ProviderInternal && len(args.Password) < 8 {
		utils.AppendValuesToMap(fieldsMap, consts.Password, consts.MinLengthPassword)
	}

	// Validate the password.
	passwordErr := utilities.ValidatePassword(args.Password)
	if args.Provider == consts.ProviderInternal {
		// Check if the provider is internal and there is a password validation error.
		if passwordErr != nil {
			// Append values to the map based on the password error.
			utils.AppendValuesToMap(fieldsMap, consts.Password, consts.Format)

		}
	}

	if args.Provider == consts.ProviderInternal && !args.TermsConditionChecked {
		utils.AppendValuesToMap(fieldsMap, consts.TermsAndConditions, consts.Required)
	}

	if args.Provider == consts.ProviderInternal && !args.PayingTax {
		utils.AppendValuesToMap(fieldsMap, consts.PayingTax, consts.Required)

	}

	existsOrNot, err := member.repo.ProviderExists(ctxt, args.Provider)

	if !existsOrNot {
		utils.AppendValuesToMap(fieldsMap, consts.ProviderName, consts.Exists)

	}
	if err != nil {
		logger.Log().WithContext(ctxt).Errorf("Failed to check if provider exists: %s", err.Error())
		return nil, err
	}
	if len(fieldsMap) == 0 {
		err := member.repo.RegisterMember(ctxt, args, partnerID)
		if err != nil {
			// Convert the repository error to a string for better handling.
			repoErrorStr := err.Error()

			logger.Log().WithContext(ctxt).Errorf("Failed to register a member: err=%s", err.Error())
			return nil, fmt.Errorf("registration failed: %s", repoErrorStr)
		}
	}

	return fieldsMap, nil
}

// ViewMemberProfile retrieves a member's profile and checks if the member exists.
//
// Parameters:
//   - ctxt (context.Context): The context for the database operation.
//   - memberId (uuid.UUID): The UUID of the member for whom to retrieve the profile.
//   - contextError (map[string]any): A map containing context-specific error responses.
//   - endpoint (string): The endpoint name.
//   - method (string): The HTTP request method (e.g., "GET" or "POST").
//
// Returns:
//   - fieldsMap (map[string][]string): A map of validation error messages.
//   - memberProfile (entities.MemberProfile): A struct containing the member's profile.
//   - error: An error, if any, encountered during the database operation.
func (member *MemberUseCases) ViewMemberProfile(ctx *gin.Context, ctxt context.Context, memberId uuid.UUID, contextError map[string]any, endpoint string,
	method string) (map[string][]string, entities.MemberProfile, error) {

	fieldsMap := map[string][]string{}
	var memberProfile entities.MemberProfile

	exists, err := member.repo.IsMemberExists(memberId, ctx)

	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		return fieldsMap, memberProfile, nil
	}

	if err != nil {
		logger.Log().WithContext(ctxt).Errorf("Checking member exists failed: err=%s", err.Error())
		return nil, memberProfile, err
	}

	if len(fieldsMap) == 0 {
		memberProfile, err = member.repo.ViewMemberProfile(memberId, ctxt)

		if err != nil {

			logger.Log().WithContext(ctxt).Errorf("View member profile failed: %s", err.Error())
			return nil, memberProfile, err
		}

	}

	return fieldsMap, memberProfile, nil
}

// ViewMembers retrieves a list of member profiles with additional information,
// such as their roles, partner names, album, track, and artist counts, from the database.
//
// Parameters:
//   - ctxt (context.Context): The context for the database operation.
//   - params (entities.Params): An instance of entities.Params containing filtering and pagination parameters.
//
// Returns:
//   - memberData ([]entities.ViewMembers): A slice of entities.ViewMembers, each representing a member's profile.
//   - error: An error, if any, encountered during the database operation.
//

func (member *MemberUseCases) ViewMembers(ctx *gin.Context, params entities.Params) ([]entities.ViewMembers, models.MetaData, error) {
	var memberData []entities.ViewMembers

	// Get the total record count before applying any filters
	recordCount, err := member.repo.GetMemberRecordCount(ctx)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to get member record count: %v", err)
		return memberData, models.MetaData{}, err
	}

	// Check if the status parameter has any value, either from the query or default
	if status, exists := ctx.GetQuery("status"); exists && status != "" {
		params.Status = status
	}
	if country, exists := ctx.GetQuery("country"); exists && country != "" {
		params.Country = country

		countryExists, err := member.repo.CountryExists(params.Country)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to check country existence: %v", err)
			return nil, models.MetaData{}, err
		}

		if !countryExists {
			return memberData, models.MetaData{}, errors.New("invalid param value")
		}
	}

	if gender, exists := ctx.GetQuery("gender"); exists && gender != "" {
		gender = strings.ToUpper(gender)
		if gender == "M" || gender == "F" {
			params.Gender = gender
		} else {
			return memberData, models.MetaData{}, errors.New("invalid param value")
		}
	}

	// Check if the page parameter has any value, either from the query or default
	if page, exists := ctx.GetQuery("page"); exists && page != "" {
		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			logger.Log().Errorf("Failed to parse page value: %v", err)
		} else {
			params.Page = int16(pageInt)
		}
	}
	if limit, exists := ctx.GetQuery("limit"); exists && limit != "" {
		validLimitRegex := regexp.MustCompile(`^[1-9]\d*$|^0$`)

		if !validLimitRegex.MatchString(limit) {
			logger.Log().WithContext(ctx).Errorf("Failed get details, invalid limit: %v", err)
			return memberData, models.MetaData{}, errors.New("invalid param value")
		}

		limitInt, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			logger.Log().Errorf("Failed to parse limit value: %v", err)
			return nil, models.MetaData{}, err
		}

		if limitInt < 0 {
			logger.Log().WithContext(ctx).Errorf("Failed get details, invalid limit: %v", err)
			return memberData, models.MetaData{}, errors.New("invalid param value")
		}

		params.Limit = int16(limitInt)
	}

	if sortBy, exists := ctx.GetQuery("sort"); exists && sortBy != "" {
		params.SortBy = sortBy
	}

	if partner, exists := ctx.GetQuery("partner"); exists && partner != "" {
		// Set the Partner parameter
		params.Partner = partner
	}

	if role, exists := ctx.GetQuery("role"); exists && role != "" {
		params.Role = role
	}

	if search, exists := ctx.GetQuery("search"); exists && search != "" {
		params.Search = search
	}
	recordCount, err = member.repo.GetFilteredRecordCount(ctx, params)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to get filtered member record count: %v", err)
		return memberData, models.MetaData{}, err
	}

	if params.Status != "" || params.Page != 0 || params.Limit != 0 || params.SortBy != "" ||
		params.Partner != "" || params.Role != "" || params.Search != "" {
		// Get members based on the applied parameters
		memberData, err = member.repo.ViewMembers(ctx, params)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to fetch members data, err=%s", err.Error())
			return memberData, models.MetaData{}, err
		}
	} else {
		// No parameters provided, you might want to handle this case accordingly
		logger.Log().WithContext(ctx).Info("No parameters provided")
		return memberData, models.MetaData{}, nil
	}

	// Convert int16 values to string first
	pageStr := strconv.FormatInt(int64(params.Page), 10)
	limitStr := strconv.FormatInt(int64(params.Limit), 10)

	// Now, convert the string values to int64
	pageInt, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		logger.Log().Errorf("Failed to parse page value: %v", err)
	}

	limitInt, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		// Handle error
		logger.Log().Errorf("Failed to parse limit value: %v", err)
		// Handle the error as needed, maybe set a default value or return an error
	}
	if len(memberData) == 0 {
		metadata := &models.MetaData{
			CurrentPage: int32(pageInt),
			PerPage:     0,
			Total:       recordCount,
		}
		return memberData, *metadata, nil
	}

	// Update metadata based on the actual data retrieved
	metadata := &models.MetaData{
		CurrentPage: int32(pageInt),
		PerPage:     int32(limitInt),
		Total:       recordCount,
	}

	metadata = utils.MetaDataInfo(metadata)
	return memberData, *metadata, nil
}

// GetBasicMemberDetailsByEmail retrieves basic details of a member based on their email.
//
// This function performs validation on the input arguments and returns any validation errors
// in a map, along with the basic member details and an error if the database operation fails.
//
// Parameters:
//   - ctxt (context.Context): The context for the database operation.
//   - args (entities.MemberPayload): An instance of entities.MemberPayload containing member information.
//   - contextError (map[string]any): A map to collect validation errors.
//   - endpoint (string): The endpoint for the request.
//   - method (string): The HTTP method used for the request.
//
// Returns:
//   - fieldsMap (map[string][]string): A map containing validation errors, if any, for different fields.
//   - memberBasic (entities.BasicMemberData): An instance of entities.BasicMemberData representing basic member details.
//   - error: An error, if any, encountered during the database operation.
func (member *MemberUseCases) GetBasicMemberDetailsByEmail(ctx *gin.Context, partnerID string, args entities.MemberPayload, contextError map[string]any, endpoint string,
	method string) (map[string][]string, entities.BasicMemberData, error) {
	fieldsMap := map[string][]string{}
	var (
		memberBasic entities.BasicMemberData
		err         error
	)
	validEmailProvider, err := member.repo.CheckEmailProviderRelation(ctx, args.Email, args.Provider)

	if !validEmailProvider {
		utils.AppendValuesToMap(fieldsMap, consts.Email, consts.InValidProvider)
	}
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("View basic member details failed: %s", err.Error())
	}
	if args.Email == "" {
		utils.AppendValuesToMap(fieldsMap, consts.Email, consts.Required)
	}
	if args.Provider != consts.ProviderSpotify {
		if !utilities.ValidateEmail(args.Email) {
			utils.AppendValuesToMap(fieldsMap, consts.Email, consts.Valid)
		}
	}
	if args.Provider != consts.ProviderSpotify {
		if args.Provider == consts.ProviderInternal && args.Password == "" {
			utils.AppendValuesToMap(fieldsMap, consts.Password, consts.Required)
		}

	}
	if args.Provider == consts.ProviderInternal && len(args.Password) < 8 {
		utils.AppendValuesToMap(fieldsMap, consts.Password, consts.MinLengthPassword)
	}

	if passwordErr := utilities.ValidatePassword(args.Password); args.Provider == consts.ProviderInternal && passwordErr != nil {
		utils.AppendValuesToMap(fieldsMap, consts.Password, consts.Format)
	}

	if len(fieldsMap) == 0 {
		memberBasic, err = member.repo.GetBasicMemberDetailsByEmail(partnerID, args, ctx)

		if err != nil {
			logger.Log().WithContext(ctx).Errorf("View basic member details failed: %s", err.Error())
			return nil, memberBasic, err
		}

	}

	return fieldsMap, memberBasic, nil
}

//Reset Password initiation

// InitiatePasswordReset initiates the password reset process for a member.
func (member *MemberUseCases) InitiatePasswordReset(ctx *gin.Context, memberID uuid.UUID, email string) (string, map[string][]string, error) {
	fieldsMap := make(map[string][]string)
	var key string
	var err error // Declare the err variable here
	exists, err := member.repo.IsMemberExists(memberID, ctx)
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		logger.Log().WithContext(ctx).Errorf("ChangePassword failed, Member not found")
		return "", fieldsMap, nil
	}
	validEmail, err := member.repo.CheckEmailForMemberID(ctx, memberID, email)
	if !validEmail {
		utils.AppendValuesToMap(fieldsMap, consts.ResetEmail, consts.InvalidEmail)
		logger.Log().WithContext(ctx).Errorf("Reset Password failed, validation error: No email found")
		return "", fieldsMap, nil
	}
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Reset Password failed due to database error: %s", err.Error())
		return "", nil, err
	}
	// Validate email presence
	if len(email) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.ResetEmail, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Reset Password failed, validation error: No email found")
		return "", fieldsMap, nil // Return early if email is empty
	}

	// Validate email format
	if !utilities.ValidateEmail(email) {
		utils.AppendValuesToMap(fieldsMap, consts.ResetEmail, consts.Format)
		logger.Log().WithContext(ctx).Errorf("Reset Password failed, validation error: Invalid email format")
		return "", fieldsMap, nil // Return early if email format is invalid
	}

	// If no validation errors, proceed to reset password
	key, err = member.repo.InitiatePasswordReset(ctx, memberID, email)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Resetting password failed: %s", err.Error())
		return key, fieldsMap, err // Ensure you return the error here
	}

	// Return the generated key and no validation errors
	return key, nil, nil
}

// HandleSubscriptionCheckout handles the checkout of a subscription for a member.
// It takes the
//   - context
//   - memberID
//   - checkoutData
//     as parameters.
//
// It returns a map with subscription details and any encountered errors.
func (member *MemberUseCases) HandleSubscriptionCheckout(ctx *gin.Context, memberID uuid.UUID, checkoutData entities.CheckoutSubscription, partnerIDStr string) (map[string][]string, error) {
	fieldsMap := map[string][]string{}
	// Create a map to store the result
	resultMap := make(map[string][]string)

	subExists, isActive, err := member.repo.CheckSubscriptionExistenceAndStatusForCheckout(ctx, checkoutData.SubscriptionID)
	// Ensure the provided partnerID matches the partnerID associated with the member
	isPartnerValid, err := member.repo.CheckMemberPartner(ctx, memberID, partnerIDStr)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check member partner: %s", err.Error())
		return nil, err
	}

	if !isPartnerValid {
		// Partner authentication failed
		utils.AppendValuesToMap(fieldsMap, consts.PartnerID, consts.AuthenticationFailed)
		logger.Log().WithContext(ctx).Errorf("Partner authentication failed for member: %s", memberID)
		return fieldsMap, nil
	}
	if !subExists {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Invalid subscription plan")
		return fieldsMap, nil
	}
	if !isActive {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Inactive)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Inactive subscription plan")
		return fieldsMap, nil
	}
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: %s", err.Error())
	}
	// Check if SubscriptionID is empty
	if checkoutData.SubscriptionID == "" || len(checkoutData.SubscriptionID) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Missing required fields")
		return fieldsMap, nil
	}
	if len(checkoutData.SubscriptionID) > 36 {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.TooLong)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Subscription Id too long/exceeds limit")
	}
	if len(checkoutData.CustomName) > 60 {
		utils.AppendValuesToMap(fieldsMap, consts.CustonName, consts.TooLong)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Custom name too long")
		return fieldsMap, nil
	}
	// Fetch the count of subscriptions for the member within the last year
	count, err := member.repo.GetSubscriptionCountForLastYear(ctx, memberID, checkoutData.SubscriptionID)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Failed to get members subscription count in last year")
	}

	// Fetch the maximum subscription limit for the given subscription ID
	maxCount, err := member.repo.GetMaxSubscriptionLimitForID(ctx, checkoutData.SubscriptionID)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Failed to get maximum allowed count per yearr")

	}
	if count == maxCount {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.LimitReached)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Reached maximum allowed subscription in an year")
		return fieldsMap, nil
	}
	// If the subscription is free, skip the PaymentGatewayID check
	isFree, err := member.repo.IsFreeSubscription(ctx, checkoutData.SubscriptionID)

	if err != nil {
		if err.Error() == "subscription does not exist" {
			logger.Log().WithContext(ctx).Errorf("Invalid Subscription Id : %s", err.Error())
			return nil, err
		} else {
			logger.Log().WithContext(ctx).Errorf("Failed to check if the subscription is free: %s", err.Error())
			return nil, err
		}
	}
	// Continue with the rest of the code, including PaymentGatewayID check, if it's not free
	if !isFree {
		if checkoutData.PaymentGatewayID == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.Required)
			logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Missing required fields")
			return fieldsMap, nil
		}

		gatewayExists, _ := member.repo.CheckIfPayoutGatewayExists(ctx, checkoutData.PaymentGatewayID)

		if !gatewayExists {
			utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.Invalid)
			logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Invalid gateway")
			return fieldsMap, nil
		}
		PartnerGatewayExists, err := member.repo.IsPartnerIdCorrespondsToGateway(ctx, partnerIDStr, checkoutData.PaymentGatewayID)

		if !PartnerGatewayExists {
			utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.NoRelation)
			logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: This payment gateway is not endorsed by partner")
			return fieldsMap, nil
		}

		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan:error checking if member subscribed free plan")
		}
		if PartnerGatewayExists {
			// If the partner ID corresponds to the payment gateway, retrieve payment details.
			paymentDetails, err := member.repo.GetPaymentDetailsByPartnerAndGateway(ctx, partnerIDStr, checkoutData.PaymentGatewayID)
			if err != nil {
				logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan:payment details are NULL for partnerID")
				return nil, err
			}
			detailsString, err := member.repo.DecryptPaymentData(ctx, paymentDetails)
			if err != nil {
				return nil, err
			}

			var paymentInfo entities.PaymentGatewayDetails

			// Unmarshal JSON string into the PaymentGatewayDetails struct
			err = json.Unmarshal([]byte(detailsString), &paymentInfo)
			if err != nil {
				return nil, err
			}
			if !paymentInfo.Payin {
				utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.NoPayin)
				logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: pay-in transaction is not supported by this payment gateway.")
				return fieldsMap, nil
			}

		}
		if checkoutData.PaymentGatewayID > consts.MaxInt {
			utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.TooLong)
			logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Invalid gateway id")
			return fieldsMap, nil
		}
	}

	subscribedOnetime, err := member.repo.HasSubscribedToOneTimePlan(ctx, memberID, checkoutData.SubscriptionID)

	if subscribedOnetime {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.SubscribedOnce)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: This is is one time plan")
		return fieldsMap, nil
	}
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan:Error checking This is is one time plan")
	}
	// Return the 'fieldsMap' with the processed data and a nil error, indicating success.
	if len(fieldsMap) > 0 {
		return fieldsMap, nil
	}

	// Handle the subscription checkout by calling the HandleSubscriptionCheckout method from the repository.
	err = member.repo.HandleSubscriptionCheckout(ctx, memberID, checkoutData)

	// Check if there was an error during the checkout process.
	if err != nil {
		resultMap["error"] = []string{err.Error()}
		logger.Log().WithContext(ctx).Errorf("Failed to checkout Subscription: %s", err.Error())
		return nil, err
	}

	return fieldsMap, nil
}

// HandleSubscriptionRenewal handles the checkout of a subscription for a member.
// It takes the
//   - context
//   - memberID
//   - checkoutData
//     as parameters.
//
// It returns a map with subscription details and any encountered errors.
func (member *MemberUseCases) HandleSubscriptionRenewal(ctx *gin.Context, memberID uuid.UUID, checkoutData entities.SubscriptionRenewal, partnerIDStr string) (map[string][]string, error) {
	fieldsMap := map[string][]string{}

	// Ensure the provided partnerID matches the partnerID associated with the member.
	isPartnerValid, err := member.repo.CheckMemberPartner(ctx, memberID, partnerIDStr)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check member partner: %s", err.Error())
		return nil, err
	}
	if !isPartnerValid {
		utils.AppendValuesToMap(fieldsMap, consts.PartnerID, consts.AuthenticationFailed)
		logger.Log().WithContext(ctx).Errorf("Partner authentication failed for member: %s", memberID)
		return fieldsMap, nil
	}

	//CheckSubscriptionExistenceAndStatusForRenewal fetches the subscription plan from member subscription  id and check
	//if such a plan exists and is currently active
	subExists, isActive, err := member.repo.CheckSubscriptionExistenceAndStatusForRenewal(ctx, checkoutData.MemberSubscriptionID)
	if !subExists {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Invalid subscription plan")
		return fieldsMap, nil
	}
	if !isActive {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Inactive)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Inactive subscription plan")
		return fieldsMap, nil
	}

	//IsMemberSubscribedToPlan checks if member subscribed to the plan, before proceeding with renewal
	yesSubscribed, err := member.repo.IsMemberSubscribedToPlan(ctx, memberID, checkoutData.MemberSubscriptionID)
	if !yesSubscribed {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.NotSubscribedYet)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Not yet subscribed this plan")
		return fieldsMap, nil
	}

	//IsSubscriptionAboutInWarning checks if the subscription is in warning period ,
	// also notifies member about payment to proceed with renewal
	aboutToExpire, _, err := member.repo.IsSubscriptionAboutInWarning(ctx, checkoutData.MemberSubscriptionID)
	if aboutToExpire {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.AboutToExpire)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Plan is about to expire,Make Payment to renew")
		return fieldsMap, nil
	}

	//IsSubscriptionInGracePeriod checks if the subscription is in grace period ,
	// also notifies member about payment to proceed with renewal before expiration.
	inGrace, _, graceEnd, _, _, err := member.repo.IsSubscriptionInGracePeriod(ctx, memberID, checkoutData.MemberSubscriptionID)
	if inGrace {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.InGrace)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Plan is in grace period,Make Payment to renew")
		return fieldsMap, nil
	}
	//checks if subscription plan expired
	isExpired := time.Now().After(graceEnd)
	if isExpired {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Expired)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Plan expired")
		return fieldsMap, nil
	}

	//GetSubscriptionStatusName checks the status of current subscription
	currentStatus, err := member.repo.GetSubscriptionStatusName(ctx, checkoutData.MemberSubscriptionID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to renew current subscription plan, error in checking if the subscription was free: %s", err.Error())
		return nil, err
	}
	if currentStatus != "active" && currentStatus != "on_hold" && currentStatus != "payment_failed" {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.CheckStatus)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Check the current status of the plan")
		return fieldsMap, nil
	}

	//IsMemberSubscribedToFreePlan checks if the plan is free, if so it cannot be renewed
	subscribedFree, err := member.repo.IsMemberSubscribedToFreePlan(ctx, memberID, checkoutData.MemberSubscriptionID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to renew current subscription plan, error in checking if the subscription was free: %s", err.Error())
		return nil, err
	}
	if subscribedFree {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Subscribed)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: cannot renew free plan")
		return fieldsMap, nil

	}

	//Check if PaymentGatewayID is empty
	if checkoutData.PaymentGatewayID == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Missing required fields")
	}
	if len(checkoutData.MemberSubscriptionID) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: Member Subscription Id is mandatory")
	}
	if len(checkoutData.MemberSubscriptionID) == 0 {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Missing required fields")
		return fieldsMap, nil
	}
	if len(checkoutData.MemberSubscriptionID) > 36 {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.TooLong)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Subscription Id too long/exceeds limit")
		return fieldsMap, nil
	}

	//CheckIfPayoutGatewayExists checks if the payment gateway entered is valid or not
	gatewayExists, err := member.repo.CheckIfPayoutGatewayExists(ctx, checkoutData.PaymentGatewayID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to renew Subscription:Error in checking gateway existance %s", err.Error())
		return nil, err
	}

	if !gatewayExists {
		utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan:Invalid gateway")
		return fieldsMap, nil
	}

	//IsPartnerIdCorrespondsToGateway checks if the entred payment gateway is approved/endorsed by the corresponding partner
	PartnerGatewayExists, err := member.repo.IsPartnerIdCorrespondsToGateway(ctx, partnerIDStr, checkoutData.PaymentGatewayID)
	if !PartnerGatewayExists {
		utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.NoRelation)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: This payment gateway is not endorsed by partner")
	}
	if PartnerGatewayExists {
		// If the partner ID corresponds to the payment gateway, retrieve payment details.
		paymentDetails, err := member.repo.GetPaymentDetailsByPartnerAndGateway(ctx, partnerIDStr, checkoutData.PaymentGatewayID)

		if len(paymentDetails) == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.NoDetails)
			logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: No payment details found for the partner and payment gateway")
			return fieldsMap, nil
		}

		detailsString, err := member.repo.DecryptPaymentData(ctx, paymentDetails)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to renew Subscription:Error in decrypting key %s", err.Error())
			return nil, err
		}

		var paymentInfo entities.PaymentGatewayDetails

		// Unmarshal JSON string into the PaymentGatewayDetails struct
		err = json.Unmarshal([]byte(detailsString), &paymentInfo)
		if err != nil {
			return nil, err
		}
		if !paymentInfo.Payin {
			utils.AppendValuesToMap(fieldsMap, consts.PaymentGatewayID, consts.NoPayin)
			logger.Log().WithContext(ctx).Errorf("Failed to checkout this plan: pay-in transaction is not supported by this payment gateway.")
			return fieldsMap, nil
		}

	}
	// Return the 'fieldsMap' with the processed data and a nil error, indicating success.
	if len(fieldsMap) > 0 {
		return fieldsMap, nil
	}

	// Handle the subscription checkout by calling the HandleSubscriptionCheckout method from the repository.
	//HandleSubscriptionRenewal  sets a new expiration time and updates the renewed on date .
	err = member.repo.HandleSubscriptionRenewal(ctx, memberID, checkoutData)

	// Check if there was an error during the checkout process.
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to renew Subscription: %s", err.Error())
		return nil, err
	}

	return fieldsMap, nil
}

// HandleSubscriptionCancellation handles the checkout of a subscription for a member.
// It takes the
//   - context
//   - memberID
//     as parameters.
//
// It returns a map with subscription details and any encountered errors.
func (member *MemberUseCases) HandleSubscriptionCancellation(ctx *gin.Context, memberID uuid.UUID, checkoutData entities.CancelSubscription, partnerIDStr string) (map[string][]string, error) {

	fieldsMap := map[string][]string{}

	// Check if SubscriptionID is empty
	// Ensure the provided partnerID matches the partnerID associated with the member.
	isPartnerValid, err := member.repo.CheckMemberPartner(ctx, memberID, partnerIDStr)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check member partner: %s", err.Error())
		return nil, err
	}
	if !isPartnerValid {
		utils.AppendValuesToMap(fieldsMap, consts.PartnerID, consts.AuthenticationFailed)
		logger.Log().WithContext(ctx).Errorf("Partner authentication failed for member: %s", memberID)
		return fieldsMap, nil
	}

	//CheckSubscriptionExistenceAndStatusForRenewal fetches the subscription plan from member subscription  id and check
	//if such a plan exists and is currently active(CheckSubscriptionExistenceAndStatusForRenewal , actually for cancellation,repo function reused)
	subExists, isActive, err := member.repo.CheckSubscriptionExistenceAndStatusForRenewal(ctx, checkoutData.MemberSubscriptionID)
	if !subExists {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Invalid subscription plan")
		return fieldsMap, nil
	}
	if !isActive {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Inactive)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Inactive subscription plan")
		return fieldsMap, nil
	}

	//IsMemberSubscribedToPlan checks if member subscribed to the plan, before proceeding with renewal
	yesSubscribed, err := member.repo.IsMemberSubscribedToPlan(ctx, memberID, checkoutData.MemberSubscriptionID)
	if !yesSubscribed {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.NotSubscribedYet)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Not yet subscribed this plan")
		return fieldsMap, nil
	}

	//HasProductsReleaseEndDateGreaterThanToday checks if  there are released products that havenot yet reached end date.
	yesProduct, err := member.repo.HasProductsReleaseEndDateGreaterThanToday(ctx, checkoutData.MemberSubscriptionID)
	if yesProduct {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.CannotCancel)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: There are released products under this subscription")
		return fieldsMap, nil
	}
	isRelated, err := member.repo.IsMemberRelatedToSubscription(ctx, memberID, checkoutData.MemberSubscriptionID)
	if !isRelated {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to renew this plan: Invalid member subscription id")
	}

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check subscription status: %v", err)
		return fieldsMap, err
	}

	//CheckCancellationEnabled checks if this subscription is cancellation enabled.
	enabled, err := member.repo.CheckCancellationEnabled(ctx, checkoutData.MemberSubscriptionID)
	if !enabled {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.NotEnabled)
		logger.Log().WithContext(ctx).Errorf("Failed to cancel this subscription: This is not cancellation enabled")
		return fieldsMap, nil
	}
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to check subscription status: %v", err)
		return fieldsMap, err
	}

	//GetSubscriptionStatusName gets the current status of the subscription to be cancelled.
	currentStatus, err := member.repo.GetSubscriptionStatusName(ctx, checkoutData.MemberSubscriptionID)

	if currentStatus != "active" && currentStatus != "on_hold" {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.CheckStatus)
		logger.Log().WithContext(ctx).Errorf("Failed to cancel this plan: Check the current status of the plan")
		return fieldsMap, nil
	}
	if err != nil {
		return nil, err

	}
	if checkoutData.MemberSubscriptionID == "" {
		utils.AppendValuesToMap(fieldsMap, consts.SubscriptionID, consts.Required)
		logger.Log().WithContext(ctx).Errorf("Failed to cancel this subscription: Missing required fields")
		return fieldsMap, nil
	}
	// Return the 'fieldsMap' with the processed data and a nil error, indicating success.
	if len(fieldsMap) > 0 {
		return fieldsMap, nil
	}

	// Handle the subscription cancellation by calling the HandleSubscriptionCancellation method from the repository.
	err = member.repo.HandleSubscriptionCancellation(ctx, memberID, checkoutData)

	// Check if there was an error during the checkout process.
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to cancel Subscription: %s", err.Error())
		return nil, err
	}

	return fieldsMap, nil
}

// SubscriptionProductSwitch function
func (member *MemberUseCases) SubscriptionProductSwitch(ctx context.Context, memberID uuid.UUID, data entities.SwitchSubscriptions) (map[string][]string, error) {

	validationErrors := make(map[string][]string)

	if utilities.IsEmpty(data.CurrentSubscriptionID) {
		utils.AppendValuesToMap(validationErrors, consts.CurrentSubscriptionID, consts.Required)
	}

	if utilities.IsEmpty(data.NewSubscriptionID) {
		utils.AppendValuesToMap(validationErrors, consts.NewSubscriptionID, consts.Required)
	}

	if utilities.IsEmpty(data.ProductReferenceID) {
		utils.AppendValuesToMap(validationErrors, consts.ProductReferenceID, consts.Required)
	}

	_, err := member.repo.SubscriptionProductSwitch(ctx, memberID, data, &validationErrors)

	if len(validationErrors) != 0 {
		return validationErrors, nil
	}

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Switch products between subscriptions failed err=%s", err.Error())
		return nil, err
	}

	return nil, nil
}

// ViewAllSubscriptions function
func (member *MemberUseCases) ViewAllSubscriptions(ctx *gin.Context, memberID uuid.UUID, reqParam entities.ReqParams) ([]entities.ListAllSubscriptions, models.MetaData, map[string][]string, error) {

	validationErrors := make(map[string][]string)

	if status, exists := ctx.GetQuery("status"); exists && status != "" {
		reqParam.Status = status
	}

	if search, exists := ctx.GetQuery("search"); exists && search != "" {
		reqParam.Search = search
	}

	// Fetching sortBy from query parameter
	if sortBy, exists := ctx.GetQuery("sort"); exists && sortBy != "" {
		reqParam.Sort = sortBy
	}

	// Call GetSubscriptionRecordCount function in repo
	recordCount, err := member.repo.GetSubscriptionRecordCount(ctx, memberID)

	if err != nil {
		return []entities.ListAllSubscriptions{}, models.MetaData{}, nil, err
	}

	if reqParam.Limit > consts.MaximumLimit {
		err := fmt.Errorf("Cannot exceed maximum page limit")
		return []entities.ListAllSubscriptions{}, models.MetaData{}, nil, err
	}

	if int64(reqParam.Page)*int64(reqParam.Limit)-recordCount >= int64(reqParam.Limit) { //check whether page parameters are valid
		validationErrors := make(map[string][]string)
		utils.AppendValuesToMap(validationErrors, consts.Page, consts.Invalid)
		//return []entities.ListAllSubscriptions{}, models.MetaData{}, validationErrors, nil
	}

	//call Paginate function
	reqParam.Page, reqParam.Limit = utils.Paginate(reqParam.Page, reqParam.Limit, consts.LimitDefault)

	//call ViewAllSubscriptions function in repo
	memberSubscriptions, err := member.repo.ViewAllSubscriptions(ctx, memberID, reqParam, &validationErrors)

	if len(validationErrors) != 0 {
		return []entities.ListAllSubscriptions{}, models.MetaData{}, validationErrors, nil
	}

	if err != nil {
		return []entities.ListAllSubscriptions{}, models.MetaData{}, nil, err
	}

	metadata := &models.MetaData{
		CurrentPage: reqParam.Page,
		PerPage:     reqParam.Limit,
		Total:       recordCount,
	}

	metadata = utils.MetaDataInfo(metadata)

	return memberSubscriptions, *metadata, nil, nil
}

// GetSubscriptionRecordCount function
func (member *MemberUseCases) GetSubscriptionRecordCount(ctx context.Context, memberID uuid.UUID) (int64, error) {

	//GetSubscriptionRecordCount function call to repo
	totalCount, err := member.repo.GetSubscriptionRecordCount(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("GetSubscriptionRecordCount failed, err=%s", err.Error())
		return 0, err
	}
	return totalCount, nil
}

// IsMemberExists checks if a member with the given memberId exists.
// Parameters:
//
//	@ memberId: The UUID of the member to check for existence.
//	@ ctx: The context for the database operation.
//
// Returns:
//
//	@ bool: True if the member exists, false otherwise.
//	@ error: An error, if any, during the check.
func (member *MemberUseCases) IsMemberExist(ctx context.Context, memberID uuid.UUID) (bool, error) {

	_, err := member.repo.IsMemberExist(ctx, memberID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteBillingAddress to delete a members billing address
func (member *MemberUseCases) DeleteBillingAddress(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID) (map[string][]string, error) {
	fieldsMap := map[string][]string{}
	exists, err := member.repo.IsMemberExist(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to delete this billing address: Invalid member")
		return fieldsMap, nil
	}
	IsValid, err := member.repo.CheckBillingAddressRelation(ctx, memberID, memberBillingID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to delete billing address %s", err.Error())
		return nil, err
	}
	if !IsValid {
		utils.AppendValuesToMap(fieldsMap, consts.MemberBilling, consts.InvalidAddress)
		logger.Log().WithContext(ctx).Errorf("Failed to delete billing address")
		return fieldsMap, nil
	}

	// Return the 'fieldsMap' with the processed data and a nil error, indicating success.
	if len(fieldsMap) > 0 {
		return fieldsMap, nil
	}
	err = member.repo.DeleteBillingAddress(ctx, memberID, memberBillingID)
	return nil, err
}

// DeleteMember deletes a member by their ID. It first validates the member's existence,
// checks if the member has already been deleted, and then performs the deletion operation.
func (member *MemberUseCases) DeleteMember(ctx *gin.Context, memberID string) (map[string][]string, error) {
	// Initialize a map to store validation error messages and proceed with the deletion process.
	fieldsMap := map[string][]string{}
	// Parse the memberID string to a UUID.
	MemberID, err := uuid.Parse(memberID)
	// Check if there was an error parsing the memberID.
	if err != nil {
		return nil, err
	}
	// Check if the member with the given ID exists.
	exists, err := member.repo.IsMemberExist(ctx, MemberID)
	if err != nil {
		return nil, err
	}
	// If the member does not exist, add an error message to the map.
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Invalid)
		logger.Log().WithContext(ctx).Errorf("Failed to delete this member: Invalid member")
		return fieldsMap, nil
	}
	// Check if the member has already been deleted.
	deleted, err := member.repo.IsDeleted(ctx, MemberID)
	if deleted {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Deleted)
		logger.Log().WithContext(ctx).Errorf("Member already deleted")
		return fieldsMap, nil
	}
	// Delete the member.
	err = member.repo.DeleteMember(ctx, MemberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to delete member %s", err.Error())
	}
	// Return nil for success (no validation errors).
	return nil, nil
}

// AddMemberStores adds stores to a member, based on their partner
func (member *MemberUseCases) AddMemberStores(ctx *gin.Context, memberID uuid.UUID, stores []string) (map[string][]string, error) {
	fieldsMap := map[string][]string{}
	exists, err := member.repo.IsMemberExist(ctx, memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Failed to add stores: %s", err.Error())
		return nil, err
	}
	if !exists {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.Exists)
		logger.Log().WithContext(ctx).Errorf("Failed to add member store: Member not found")
		return fieldsMap, nil
	}
	memberStatus, err := member.repo.IsActive(ctx, memberID)
	if !memberStatus {
		utils.AppendValuesToMap(fieldsMap, consts.MemberID, consts.NotActive)
		logger.Log().WithContext(ctx).Errorf("Failed to add member store: Member not Active")
		return fieldsMap, nil
	}

	// Check if payload(store names) is present
	if ctx.Request.ContentLength == 0 {
		partnerID, err := member.repo.GetPartnerIDByMemberID(ctx, memberID)

		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add stores: %s", err.Error())
			return nil, err

		}
		storeID, err := member.repo.GetStoreIDsByPartnerID(ctx, partnerID)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add stores, failed in getting store ids!: %s", err.Error())
			return nil, err
		}
		newIdList, err := member.repo.CheckNonExistingMemberStores(ctx, memberID, storeID)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add stores, failed in getting non redunatant store ids!: %s", err.Error())
			return nil, err
		}
		if len(newIdList) == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.Name, consts.StoresExists)
			logger.Log().WithContext(ctx).Errorf("Failed to add member store: Stores already exists for this member")
			return fieldsMap, nil
		}
		err = member.repo.AddMemberStoresById(ctx, memberID, newIdList)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add stores!!: %s", err.Error())
			return nil, err
		}
	} else {
		if len(stores) == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.Name, consts.Empty)
			logger.Log().WithContext(ctx).Errorf("Failed to add member store: Empty store name list")
			return fieldsMap, nil
		}
		storeExist, relatedStoreIDs, err := member.repo.CheckStoreNameExistsAndReturnIDs(ctx, stores)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add member stores : %s", err.Error())
			return nil, err
		}
		if !storeExist {
			utils.AppendValuesToMap(fieldsMap, consts.Name, consts.Invalid)
			logger.Log().WithContext(ctx).Errorf("Failed to add member store: Invalid store name")
			return fieldsMap, nil
		}
		partnerID, err := member.repo.GetPartnerIDByMemberID(ctx, memberID)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add member stores: %s", err.Error())
			return nil, err
		}

		related, err := member.repo.StorePartnerRelation(ctx, partnerID, relatedStoreIDs)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to  add member stores: %s", err.Error())
			return nil, err
		}
		hasFalse := false
		for _, value := range related {
			if value == false {
				hasFalse = true
				break
			}
		}
		if hasFalse {
			utils.AppendValuesToMap(fieldsMap, consts.Name, consts.NoRelation)
			logger.Log().WithContext(ctx).Errorf("Failed to add member store: Store not related to partner")
			return fieldsMap, nil
		}
		newIdList, err := member.repo.CheckNonExistingMemberStores(ctx, memberID, relatedStoreIDs)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add stores, failed in getting non redunatant store ids!: %s", err.Error())
			return nil, err
		}
		if len(newIdList) == 0 {
			utils.AppendValuesToMap(fieldsMap, consts.Name, consts.StoresExists)
			logger.Log().WithContext(ctx).Errorf("Failed to add member store: Stores already exists for this member")
			return fieldsMap, nil
		}
		//AddMemberStoresById addes stores under given member
		err = member.repo.AddMemberStoresById(ctx, memberID, newIdList)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Failed to add member stores: %s", err.Error())
			return nil, err
		}
	}

	return nil, nil
}
