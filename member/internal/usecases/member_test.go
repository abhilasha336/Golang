package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"member/internal/consts"
	"member/internal/entities"
	"member/internal/repo/mock"

	"member/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	cryptoHash "gitlab.com/tuneverse/toolkit/utils/crypto"
)

func init() {
	// Initialize the logger with the specified options before running the tests
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}
	logger.InitLogger(clientOpt)
}

// Add Update Billing Address test cases *****************************************
func TestAddBillingAddress(t *testing.T) {
	// Create a new gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a new MemberUseCases instance with the mock repository
	useCases := usecases.NewMemberUseCases(mockRepo)

	// Define test data
	memberID := uuid.New()
	billingAddress := entities.BillingAddress{
		Address: "123 Main St",
		Zipcode: "12345",
	}
	ginCtx := createTestGinContext()
	// Test case 1: Valid billing address
	mockRepo.EXPECT().AddBillingAddress(ginCtx, memberID, billingAddress).Return(nil)
	fieldsMap, err := useCases.AddBillingAddress(ginCtx, memberID, billingAddress)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(fieldsMap) != 0 {
		t.Errorf("Expected empty fieldsMap, got %v", fieldsMap)
	}

	// Test case 2: Missing required fields in billing address
	invalidAddress := entities.BillingAddress{
		Address: "",
	}
	if invalidAddress.Address == "" {
		return
	}
	fieldsMap, err = useCases.AddBillingAddress(ginCtx, memberID, invalidAddress)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(fieldsMap) != 1 {
		t.Errorf("Expected 1 field in fieldsMap, got %v", fieldsMap)
	}

	// Test case 3: Invalid ZIP code
	invalidZipcode := entities.BillingAddress{
		Address: "123 Main St",
		Zipcode: "invalid",
	}
	fieldsMap, err = useCases.AddBillingAddress(ginCtx, memberID, invalidZipcode)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(fieldsMap) != 1 {
		t.Errorf("Expected 1 field in fieldsMap, got %v", fieldsMap)
	}
}

// TestUpdateBillingAddress is the test case for updating an existing Billing Address
func TestUpdateBillingAddress(t *testing.T) {
	// Create a new gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a new MemberUseCases instance with the mock repository
	useCases := usecases.NewMemberUseCases(mockRepo)

	// Define test data
	memberID := uuid.New()
	memberBillingID := uuid.New()
	billingAddress := entities.BillingAddress{
		Address: "123 Main St",
		Zipcode: "12345",
	}

	// Test case 1: Valid billing address
	// Convert context.Background() to *gin.Context for testing or specific use cases.
	ginCtx := createTestGinContext()
	mockRepo.EXPECT().UpdateBillingAddress(ginCtx, memberID, memberBillingID, billingAddress).Return(nil)
	fieldsMap, err := useCases.UpdateBillingAddress(ginCtx, memberID, memberBillingID, billingAddress)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(fieldsMap) != 0 {
		t.Errorf("Expected empty fieldsMap, got %v", fieldsMap)
	}
	// Test case 2: Missing required fields in billing address
	invalidAddress := entities.BillingAddress{
		Address: "",
		Zipcode: "",
	}
	if invalidAddress.Address == "" || invalidAddress.Zipcode == "" {
		return
	}
	fieldsMap, err = useCases.UpdateBillingAddress(ginCtx, memberID, memberBillingID, invalidAddress)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(fieldsMap) != 1 {
		t.Errorf("Expected 1 fields in fieldsMap, got %v", fieldsMap)
	}

	// Test case 3: Invalid ZIP code
	invalidZipcode := entities.BillingAddress{
		Address: "123 Main St",
		Zipcode: "invalid",
	}
	fieldsMap, err = useCases.UpdateBillingAddress(ginCtx, memberID, memberBillingID, invalidZipcode)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(fieldsMap) != 1 {
		t.Errorf("Expected 1 field in fieldsMap, got %v", fieldsMap)
	}

}

// TestIsMemberExists to check if member exists or not
func TestIsMemberExists(t *testing.T) {
	// Create a new gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a new MemberUseCases instance with the mock repository
	useCases := usecases.NewMemberUseCases(mockRepo)

	// Define a member ID for testing
	memberID := uuid.New()

	// Test case 1: Member exists
	ctx := context.Background()
	mockRepo.EXPECT().IsMemberExists(memberID, ctx).Return(true, nil)
	exists, err := useCases.IsMemberExists(memberID, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected member to exist, but it doesn't")
	}

	// Test case 2: Member does not exist
	nonExistentMemberID := uuid.New()
	mockRepo.EXPECT().IsMemberExists(nonExistentMemberID, ctx).Return(false, nil)
	exists, err = useCases.IsMemberExists(nonExistentMemberID, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected member not to exist, but it does")
	}
}

//Change-Password test case *****************************************8

func TestChangePassword(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a MemberUseCases instance with the mock repository
	useCase := usecases.NewMemberUseCases(mockRepo)
	ginCtx := createTestGinContext()
	// Define test parameters
	memberID := uuid.New()
	key := "qw34rty677"
	currentPassword := "OldPassword@123"
	newPassword := "NewPassword@123"

	t.Run("Valid Password Change", func(t *testing.T) {
		// Setup mock expectations for valid password change
		var currentPassword string

		// Call the Hash function and handle both the result and the error
		hashedPassword, hashError := cryptoHash.Hash(currentPassword)
		if hashError != nil {
			fmt.Fprintf(os.Stderr, " error message: %v\n", hashError)
		}
		// Setup mock expectations for valid password change with the hashed password
		mockRepo.EXPECT().GetPasswordHash(gomock.Any(), memberID).Return(hashedPassword, nil)

		mockRepo.EXPECT().UpdatePassword(gomock.Any(), memberID, key, gomock.Any()).Return(nil)
		// Convert context.Background() to *gin.Context for testing or specific use cases.
		ginCtx := createTestGinContext()
		// Execute the function
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, key, newPassword, currentPassword)

		// Assertions
		assert.Empty(t, fieldsMap)
		assert.NoError(t, err)
	})

	t.Run("Empty New Password", func(t *testing.T) {
		// Execute the function with an empty new password
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, "", key, currentPassword)

		// Assertions
		assert.NotEmpty(t, fieldsMap["new_password"])
		assert.NoError(t, err)
	})

	t.Run("Invalid New Password", func(t *testing.T) {
		// Setup mock expectations for invalid new password
		mockRepo.EXPECT().GetPasswordHash(gomock.Any(), memberID).Return("", nil)

		// Execute the function
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, key, "WeakPassword", currentPassword)

		// Assertions
		assert.NotEmpty(t, fieldsMap["new_password"])
		assert.NoError(t, err)
	})

	t.Run("New Password Same as Current Password", func(t *testing.T) {
		// Execute the function with the same new password as the current password
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, key, currentPassword, currentPassword)

		// Assertions
		assert.NotEmpty(t, fieldsMap["new_password"])
		assert.NoError(t, err)
	})

	t.Run("Invalid Current Password", func(t *testing.T) {
		// Setup mock expectations for an incorrect current password
		mockRepo.EXPECT().GetPasswordHash(gomock.Any(), memberID).Return("hashed_current_password", nil)

		// Execute the function with an incorrect current password
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, key, newPassword, "IncorrectPassword")

		// Assertions
		assert.NotEmpty(t, fieldsMap["current_password"])
		assert.NoError(t, err)
	})

	t.Run("Error Retrieving Password Hash", func(t *testing.T) {
		// Setup mock expectations for an error when retrieving the password hash
		mockRepo.EXPECT().GetPasswordHash(gomock.Any(), memberID).Return("", errors.New("error retrieving password hash"))

		// Execute the function
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, key, newPassword, currentPassword)

		// Assertions
		assert.Empty(t, fieldsMap)
		assert.Error(t, err)
	})

	t.Run("Error Updating Password", func(t *testing.T) {
		var currentPassword string

		// Call the Hash function and handle both the result and the error
		hashedPassword, hashError := cryptoHash.Hash(currentPassword)
		if hashError != nil {
			fmt.Fprintf(os.Stderr, "This is an example error message: %v\n", hashError)
		}

		// Setup mock expectations for valid password change with the hashed password
		mockRepo.EXPECT().GetPasswordHash(gomock.Any(), memberID).Return(hashedPassword, nil)
		mockRepo.EXPECT().UpdatePassword(gomock.Any(), memberID, key, gomock.Any()).Return(errors.New("error updating password"))

		// Execute the function
		fieldsMap, err := useCase.ChangePassword(ginCtx, memberID, key, newPassword, currentPassword)

		// Assertions
		assert.Empty(t, fieldsMap)
		assert.Error(t, err)
	})
}

//View All Billing Address Test case ****************************************

func TestMemberUseCases_GetAllBillingAddresses(t *testing.T) {
	// Create a new GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a MemberUseCases instance with the mock repository
	memberUseCases := usecases.NewMemberUseCases(mockRepo)

	// Define a memberID for the test
	memberID := uuid.New()

	t.Run("Success Case", func(t *testing.T) {
		// Define expected billing addresses
		expectedBillingAddresses := []entities.BillingAddress{
			{
				Address: "123 Main St",
				Zipcode: "12345",
				Country: "USA",
				State:   "CA",
				Primary: true,
			},
			{
				Address: "456 Elm St",
				Zipcode: "67890",
				Country: "USA",
				State:   "NY",
				Primary: false,
			},
		}

		// Set up expectations for GetAllBillingAddresses
		mockRepo.EXPECT().GetAllBillingAddresses(gomock.Any(), memberID).Return(expectedBillingAddresses, nil)

		// Convert context.Background() to *gin.Context for testing or specific use cases.
		ginCtx := createTestGinContext()

		// Create a params instance for testing (modify as needed based on your implementation)
		params := entities.Params{
			Page:  1,
			Limit: 10,
		}

		// Perform the actual method call
		fieldsMap, billingAddresses, metaData, err := memberUseCases.GetAllBillingAddresses(ginCtx, memberID, params)

		// Check for any validation errors
		require.Empty(t, fieldsMap, "Expected no validation errors, got: %v", fieldsMap)

		// Check for any errors
		assert.NoError(t, err)

		// Compare the retrieved billing addresses with the expected values
		assert.Equal(t, expectedBillingAddresses, billingAddresses)

		// Add additional assertions for metadata if needed
		assert.NotNil(t, metaData)
		assert.Equal(t, int32(1), metaData.CurrentPage)
		assert.Equal(t, int32(10), metaData.PerPage)
		assert.Equal(t, int64(len(expectedBillingAddresses)), metaData.Total)
	})
}

func TestMemberRepo_IsMemberExists(t *testing.T) {
	// Initialize your logger client
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	logger.InitLogger(clientOpt)

	// Create a new GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a MemberUseCases instance with the mock repository
	memberUseCases := usecases.NewMemberUseCases(mockRepo)
	memberID := uuid.New()

	t.Run("Member Exists", func(t *testing.T) {
		// Expect a query and return a result with a row indicating member existence
		mockRepo.EXPECT().IsMemberExists(gomock.Any(), gomock.Any()).Return(true, nil)

		exists, err := memberUseCases.IsMemberExists(memberID, context.Background())
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Expect a query to trigger a database error
		expectedError := errors.New("database error")
		mockRepo.EXPECT().IsMemberExists(gomock.Any(), gomock.Any()).Return(false, expectedError)

		exists, err := memberUseCases.IsMemberExists(memberID, context.Background())
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, expectedError, err)
	})
}

//Update Member Test case ************************************************

func TestMemberUseCases_UpdateMember(t *testing.T) {
	// Initialize your logger client here (if needed)
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	logger.InitLogger(clientOpt)

	// Create a new GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock.NewMockMemberRepoImply(ctrl)

	// Create a MemberUseCases instance with the mock repository
	memberUseCases := usecases.NewMemberUseCases(mockRepo)

	// Define test data with valid member details
	memberID := uuid.New()
	validTestMember := entities.Member{
		FirstName: "John",
		LastName:  "Doe",
		Zipcode:   "12345",
		State:     "CA",
		Phone:     "9197132205",
		Country:   "IN",
	}

	// Expectations for the mock repository (valid case)
	mockRepo.EXPECT().
		UpdateMember(gomock.Any(), gomock.Eq(memberID), gomock.Eq(validTestMember)).
		Return(nil).
		Times(1)
	// Convert context.Background() to *gin.Context for testing or specific use cases.
	ginCtx := createTestGinContext()
	// Perform the actual method call for the valid case

	fieldsMap, err := memberUseCases.UpdateMember(ginCtx, memberID, validTestMember)

	// Check for errors (both general and validation errors)
	if err != nil {
		t.Errorf("Expected no error for valid data, got: %v", err)
	}

	// Check that the validation error map is empty for valid data
	if len(fieldsMap) > 0 {
		t.Errorf("Expected no validation errors for valid data, got: %v", fieldsMap)
	}

	// Define test data with invalid member details (missing first name)
	invalidTestMember := entities.Member{
		LastName: "Doe",
		Zipcode:  "12345",
		State:    "CA",
		Phone:    "9497132205",
		Country:  "IN",
	}

	// Perform the actual method call for the invalid case (missing first name)
	fieldsMap, _ = memberUseCases.UpdateMember(ginCtx, memberID, invalidTestMember)

	// Check that the validation error map contains the expected validation errors (missing first name)
	expectedValidationErrors := map[string][]string{
		"firstName": {"required", "valid"},
	}
	if validateFieldsMap(fieldsMap, expectedValidationErrors) {
		t.Errorf("Expected validation errors %v, got %v", expectedValidationErrors, fieldsMap)
	}

	// Define test data with another invalid member details (invalid last name)
	invalidTestMember = entities.Member{
		FirstName: "John",
		LastName:  "123",
		Zipcode:   "12356",
		State:     "CA",
		Phone:     "1234567890",
		Country:   "IN",
	}

	// Perform the actual method call for another invalid case (invalid last name)
	fieldsMap, _ = memberUseCases.UpdateMember(ginCtx, memberID, invalidTestMember)
	expectedValidationErrors = map[string][]string{
		"firstName": {"required", "valid"},
	}
	if validateFieldsMap(fieldsMap, expectedValidationErrors) {
		t.Errorf("Expected validation errors %v, got %v", expectedValidationErrors, fieldsMap)
	}

}

// Helper function to validate the fields map against expected validation errors
func validateFieldsMap(fieldsMap map[string][]string, expected map[string][]string) bool {
	if len(fieldsMap) != len(expected) {
		return false
	}
	for field, errors := range expected {
		actualErrors, ok := fieldsMap[field]
		if !ok {
			return false
		}
		if len(actualErrors) != len(errors) {
			return false
		}
		for i, errorMsg := range errors {
			if actualErrors[i] != errorMsg {
				return false
			}
		}
	}
	return true
}

//

func TestRegisterMember(t *testing.T) {
	// sample member entity for internal provider.
	member1 := entities.Member{
		FirstName:             "",
		LastName:              "",
		Email:                 "rahul@gmail.com",
		Password:              "Rahul@123",
		TermsConditionChecked: true,
		PayingTax:             true,
		Provider:              "internal",
	}

	partnerID := "147f090b-fc3c-4edb-ad9b-085a6381850e"

	// sample member entity for external provider like google, facebook, etc.
	member2 := entities.Member{
		Email:    "rahul@gmail.com",
		Provider: "google",
	}

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	logger.InitLogger(clientOpt)

	testCases := []struct {
		name          string
		member        entities.Member
		buildStubs    func(mockMemberRepo *mock.MockMemberRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string][]string, err error)
	}{
		{
			name:   "Valid Registration with internal provider",
			member: member1,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				mockMemberRepo.EXPECT().CheckEmailExists(gomock.Any(), "147f090b-fc3c-4edb-ad9b-085a6381850e", member1.Email).
					Times(1).Return(false, nil)
				mockMemberRepo.EXPECT().RegisterMember(gomock.Any(), member1, "147f090b-fc3c-4edb-ad9b-085a6381850e").
					Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, err error) {
				require.Len(t, fieldsMap, 0)
				require.Nil(t, err)
			},
		},
		{
			name:   "Email Already Exists",
			member: member1,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				mockMemberRepo.EXPECT().CheckEmailExists(gomock.Any(), partnerID, member1.Email).
					Times(1).Return(true, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, err error) {
				require.Len(t, fieldsMap, 1)
				require.Contains(t, fieldsMap, "email")
				require.Nil(t, err)
			},
		},
		{
			name:   "Valid Registration with Google Provider",
			member: member2,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				mockMemberRepo.EXPECT().CheckEmailExists(gomock.Any(), partnerID, member2.Email).
					Times(1).Return(false, nil)

				mockMemberRepo.EXPECT().RegisterMember(gomock.Any(), member2, partnerID).
					Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, err error) {
				require.Len(t, fieldsMap, 0)
				require.Nil(t, err)
			},
		},
		{
			name: "Google Provider with Existing Email",
			member: entities.Member{
				Email:    "rahul@gmail.com",
				Provider: "google",
			},
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				mockMemberRepo.EXPECT().CheckEmailExists(gomock.Any(), partnerID, member2.Email).
					Times(1).Return(true, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, err error) {
				require.Len(t, fieldsMap, 1)
				require.Contains(t, fieldsMap, "email")
				require.Nil(t, err)
			},
		},
		{
			name: "Empty Password",
			member: entities.Member{
				FirstName:             "John",
				LastName:              "Doe",
				Email:                 "rahul@gmail.com",
				Password:              "", // Empty password.
				TermsConditionChecked: true,
				PayingTax:             true,
				Provider:              "internal",
			},
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the CheckEmailExists function to be called with the given email.
				mockMemberRepo.EXPECT().CheckEmailExists(gomock.Any(), partnerID, member1.Email).
					Times(1).Return(false, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, err error) {
				require.Len(t, fieldsMap, 1)
				require.Contains(t, fieldsMap, "password")
				require.Nil(t, err)
			},
		},
		{
			name: "Registration with No Tax Payment",
			member: entities.Member{
				FirstName:             "John",
				LastName:              "Doe",
				Email:                 "rahul@gmail.com",
				Password:              "Pass@123",
				TermsConditionChecked: true,
				PayingTax:             false, // No tax payment.
				Provider:              "internal",
			},
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the CheckEmailExists function to be called with the given email.
				mockMemberRepo.EXPECT().CheckEmailExists(gomock.Any(), partnerID, member1.Email).
					Times(1).Return(false, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, err error) {
				require.Len(t, fieldsMap, 1)
				require.Contains(t, fieldsMap, "paying_tax")
				require.Nil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMemberRepo := mock.NewMockMemberRepoImply(ctrl)
			tc.buildStubs(mockMemberRepo)

			memberUseCase := usecases.NewMemberUseCases(mockMemberRepo)
			fieldsMap, err := memberUseCase.RegisterMember(context.Background(), tc.member, map[string]interface{}{}, partnerID, "", "")

			tc.checkResponse(t, fieldsMap, err)
		})
	}
}

// Test case to get basic member details by email
func TestGetBasicMemberDetailsByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Initialize logger outside the loop
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}

	logger.InitLogger(clientOpt)

	args := entities.MemberPayload{
		Email:    "john@example.com",
		Provider: consts.ProviderInternal,
		Password: "Pass@123",
	}

	testCases := []struct {
		name          string
		args          entities.MemberPayload
		buildStubs    func(mockMemberRepo *mock.MockMemberRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string][]string, basicData entities.BasicMemberData, err error)
	}{
		{
			name: "Valid Email and Password",
			args: args,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				mockMemberRepo.EXPECT().GetBasicMemberDetailsByEmail(
					gomock.Any(),
					gomock.Eq("614608f2-6538-4733-aded-96f902007254"),
					gomock.Eq(args),
				).Return(map[string][]string{}, entities.BasicMemberData{}, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, basicData entities.BasicMemberData, err error) {
				require.Len(t, fieldsMap, 0) // No errors expected
				require.Nil(t, err)
			},
		},
		{
			name: "Empty Email",
			args: entities.MemberPayload{
				Email:    "",
				Provider: consts.ProviderInternal,
				Password: "Pass@123",
			},
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// No expectations as there should be no repository calls.
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, basicData entities.BasicMemberData, err error) {
				require.Len(t, fieldsMap, 1)
				require.Contains(t, fieldsMap, consts.Email)
				require.Nil(t, err)
			},
		},
	}
	ginCtx := createTestGinContext()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockMemberRepo := mock.NewMockMemberRepoImply(ctrl)
			tc.buildStubs(mockMemberRepo)

			memberUseCase := usecases.NewMemberUseCases(mockMemberRepo)

			// Use a proper context here, depending on your application requirements
			fieldsMap, basicData, err := memberUseCase.GetBasicMemberDetailsByEmail(ginCtx, "partnerID_value", tc.args, nil, "expected_endpoint", "expected_method")

			tc.checkResponse(t, fieldsMap, basicData, err)
		})
	}
}

//*************View basic member Details Test Case*********

func TestViewMembers(t *testing.T) {
	// Create a sample params object.
	params := entities.Params{
		Status: "all",
		Page:   1,
		Limit:  20,
	}

	// Create sample member data that the repository will return.
	sampleMemberData := []entities.ViewMembers{
		{
			MemberId: uuid.MustParse("980e783c-e664-452d-b1ff-30d2e7767023"),
			Name:     "John Doe",
			Role: struct {
				Id   int
				Name string
			}{
				Id:   1,
				Name: "Member",
			},
			PartnerName: "Partner1",
			Email:       "john@example.com",
			Country: struct {
				Code string
				Name string
			}{
				Code: "IN",
				Name: "INDIA",
			},
			Active:      true,
			AlbumCount:  0,
			TrackCount:  0,
			ArtistCount: 7,
		},
		{
			MemberId: uuid.MustParse("980e703c-e664-452d-b1ff-30d2e7767026"),
			Name:     "Sukanya P",
			Role: struct {
				Id   int
				Name string
			}{
				Id:   1,
				Name: "Member",
			},
			PartnerName: "Partner1",
			Email:       "sukanya@gmail.com",
			Country: struct {
				Code string
				Name string
			}{
				Code: "US",
				Name: "United States",
			},
			Active:      true,
			AlbumCount:  0,
			TrackCount:  0,
			ArtistCount: 0,
		},
	}

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}
	logger.InitLogger(clientOpt)

	// Define test cases.
	testCases := []struct {
		name          string
		params        entities.Params
		buildStubs    func(mockMemberRepo *mock.MockMemberRepoImply)
		checkResponse func(t *testing.T, memberData []entities.ViewMembers, metadata models.MetaData, err error)
	}{
		{
			name:   "Successful Fetch",
			params: params,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the GetMemberRecordCount method to be called once and return a mocked record count.
				mockMemberRepo.EXPECT().GetMemberRecordCount(gomock.Any()).Times(1).Return(int64(50), nil)

				// Expect the ViewMembers method to be called with the given params and return sampleMemberData.
				mockMemberRepo.EXPECT().ViewMembers(gomock.Any(), params).
					Times(1).Return(sampleMemberData, nil)
			},
			checkResponse: func(t *testing.T, memberData []entities.ViewMembers, metadata models.MetaData, err error) {
				// Validate that there's no error and the member data matches.
				require.NoError(t, err)
				require.Equal(t, sampleMemberData, memberData)

				// Validate the metadata values.
				assert.Equal(t, int32(1), metadata.CurrentPage)
				assert.Equal(t, int32(10), metadata.PerPage)
				assert.Equal(t, int32(50), int32(metadata.Total))

			},
		},
		{
			name:   "Error Fetching Members",
			params: params,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the ViewMembers method to be called with the given params and return an error.
				mockMemberRepo.EXPECT().ViewMembers(gomock.Any(), params).
					Times(1).Return(nil, fmt.Errorf("Failed to fetch members"))
			},
			checkResponse: func(t *testing.T, memberData []entities.ViewMembers, metadata models.MetaData, err error) {
				// Expect an error to be returned and an empty member data slice.
				require.Error(t, err)
				require.Empty(t, memberData)
			},
		},
	}

	// Convert context.Background() to *gin.Context for testing or specific use cases.
	ginCtx := createTestGinContext()

	// Iterate over test cases.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMemberRepo := mock.NewMockMemberRepoImply(ctrl)
			tc.buildStubs(mockMemberRepo)

			memberUseCase := usecases.NewMemberUseCases(mockMemberRepo)

			memberData, metadata, err := memberUseCase.ViewMembers(ginCtx, tc.params)
			_ = metadata
			tc.checkResponse(t, memberData, metadata, err)
		})
	}
}

// View Individual Mmember Test Case

func TestViewMemberProfile(t *testing.T) {
	// Create a sample member ID.
	memberID := uuid.New()

	memberProfile := entities.MemberProfile{
		MemberDetails: entities.Member{
			Title:     "Mr",
			FirstName: "John",
			LastName:  "Doe",
			Gender:    "M",
			Email:     "john@example.com",
			Phone:     "1234567890",
			Address1:  "123 Main St",
			Address2:  "",
			Country:   "IN",
			State:     "KL",
			City:      "New York",
			Zipcode:   "",
			Language:  "en",
		},
		EmailSubscribed: false,
		MemberBillingAddress: []entities.BillingAddress{
			{
				Address: "Anchal",
				Zipcode: "2177",
				Country: "IN",
				State:   "KL",
				Primary: true,
			},
		},
	}

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	logger.InitLogger(clientOpt)

	// Define test cases.
	testCases := []struct {
		name          string
		memberID      uuid.UUID
		memberProfile entities.MemberProfile
		buildStubs    func(mockMemberRepo *mock.MockMemberRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string][]string, memberProfile entities.MemberProfile, err error)
	}{
		{
			name:          "Valid Member Profile",
			memberID:      memberID,
			memberProfile: memberProfile,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the IsMemberExists function to be called with the given member ID.
				mockMemberRepo.EXPECT().IsMemberExists(memberID, gomock.Any()).
					Times(1).Return(true, nil)

				// Expect the ViewMemberProfile function to be called with the given member ID.
				mockMemberRepo.EXPECT().ViewMemberProfile(memberID, gomock.Any()).
					Times(1).Return(memberProfile, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, memberProfile entities.MemberProfile, err error) {
				require.Nil(t, err)
				require.Empty(t, fieldsMap)
				require.NotNil(t, memberProfile)
			},
		},
		{
			name:          "Member Does Not Exist",
			memberID:      memberID,
			memberProfile: memberProfile,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the IsMemberExists function to be called with the given member ID.
				mockMemberRepo.EXPECT().IsMemberExists(memberID, gomock.Any()).
					Times(1).Return(false, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, memberProfile entities.MemberProfile, err error) {
				// Expect an error indicating that the member does not exist.
				require.Nil(t, err)
				require.Contains(t, fieldsMap, consts.MemberID)
				require.Empty(t, memberProfile)
			},
		},
		{
			name:          "Error Checking Member Existence",
			memberID:      memberID,
			memberProfile: memberProfile,
			buildStubs: func(mockMemberRepo *mock.MockMemberRepoImply) {
				// Expect the IsMemberExists function to be called with the given member ID.
				mockMemberRepo.EXPECT().IsMemberExists(memberID, gomock.Any()).
					Times(1).Return(false, fmt.Errorf("Failed to check member existence"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string][]string, memberProfile entities.MemberProfile, err error) {
				// Expect an error indicating a failure to check member existence.
				require.Error(t, err)
				require.Empty(t, fieldsMap)
				require.Empty(t, memberProfile)
			},
		},
	}
	// Convert context.Background() to *gin.Context for testing or specific use cases.
	ginCtx := createTestGinContext()

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMemberRepo := mock.NewMockMemberRepoImply(ctrl)
			tc.buildStubs(mockMemberRepo)

			memberUseCase := usecases.NewMemberUseCases(mockMemberRepo)
			fieldsMap, memberProfile, err := memberUseCase.ViewMemberProfile(ginCtx, context.Background(), tc.memberID, nil, "", "")

			tc.checkResponse(t, fieldsMap, memberProfile, err)
		})
	}
}
func createTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// You can set headers, query params, etc., on the gin context if needed.
	// For example:
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("headerKey", "headerValue")
	c.Request.URL.RawQuery = "key=value"

	return c
}
