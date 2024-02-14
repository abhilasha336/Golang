package repo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"member/internal/entities"
	"member/utilities"
	"strings"
	"time"

	"member/internal/consts"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/tuneverse/toolkit/core/logger"
	log "gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/utils"

	"gitlab.com/tuneverse/toolkit/utils/crypto"
)

// MemberRepo defines a repository for member-related operations.
type MemberRepo struct {
	db  *sql.DB
	Cfg *entities.EnvConfig
}

// MemberRepoImply represents the interface for interacting with the Member repository.
type MemberRepoImply interface {
	// Member Registration and Profile Management
	RegisterMember(ctx context.Context, args entities.Member, partnerID string) error
	IsMemberExists(memberID uuid.UUID, ctx context.Context) (bool, error)
	UpdateMember(ctx context.Context, memberID uuid.UUID, args entities.Member) error
	ViewMemberProfile(memberId uuid.UUID, ctx context.Context) (entities.MemberProfile, error)
	GetMemberByID(ctx context.Context, memberID uuid.UUID) (entities.MemberByID, error)
	ViewMembers(ctx context.Context, params entities.Params) ([]entities.ViewMembers, error)
	GetBasicMemberDetailsByEmail(partnerID string, args entities.MemberPayload, ctx context.Context) (entities.BasicMemberData, error)
	DeleteMember(ctx *gin.Context, MemberID uuid.UUID) error
	IsDeleted(ctx *gin.Context, MemberID uuid.UUID) (bool, error)
	IsActive(ctx *gin.Context, MemberID uuid.UUID) (bool, error)
	IsMemberExist(context.Context, uuid.UUID) (bool, error)
	CheckLanguageExist(ctx *gin.Context, Language string) (bool, error)
	GetFilteredRecordCount(ctx context.Context, params entities.Params) (int64, error)
	CheckPartnerIDExists(ctx *gin.Context, partnerID string) (bool, error)

	//Member-Store functionalities

	GetStoreIDsByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]uuid.UUID, error)
	GetPartnerIDByMemberID(ctx context.Context, memberID uuid.UUID) (uuid.UUID, error)
	GetStoreIDByCustomName(ctx context.Context, customName []string) ([]uuid.UUID, error)
	CheckNonExistingMemberStores(ctx *gin.Context, memberID uuid.UUID, storeIDs []uuid.UUID) ([]uuid.UUID, error)
	CheckStoreNameExistsAndReturnIDs(ctx context.Context, storeNames []string) (bool, []uuid.UUID, error)
	CheckPartnerStores(ctx context.Context, partnerID uuid.UUID, storeIDs []uuid.UUID) (bool, error)
	StorePartnerRelation(ctx context.Context, partnerID uuid.UUID, relatedStoreIDs []uuid.UUID) (map[string]bool, error)

	// Password and Security

	UpdatePassword(ctx context.Context, memberID uuid.UUID, key string, newPasswordHash string) error
	GetPasswordHash(ctx context.Context, memberID uuid.UUID) (string, error)
	InitiatePasswordReset(ctx *gin.Context, memberID uuid.UUID, email string) (string, error)
	CheckResetKeyMatch(ctx context.Context, memberID uuid.UUID, key string) (bool, error)

	// Billing Address Management

	AddBillingAddress(ctx *gin.Context, memberID uuid.UUID, billingAddress entities.BillingAddress) error
	UpdateBillingAddress(ctx context.Context, memberID uuid.UUID, memberBillingID uuid.UUID, billingAddress entities.BillingAddress) error
	GetAllBillingAddresses(ctx context.Context, memberID uuid.UUID) ([]entities.BillingAddress, error)
	CheckBillingAddressRelation(ctx context.Context, memberID, billingAddressID uuid.UUID) (bool, error)
	GetBillingAddressCountForMember(ctx *gin.Context, memberID uuid.UUID) (int, error)
	BillingAddressExists(ctx *gin.Context, memberID uuid.UUID, billingAddress entities.BillingAddress) (bool, error)
	CountPrimaryBillingAddresses(ctx *gin.Context, memberID uuid.UUID) (int, error)
	HasPrimaryBilling(ctx *gin.Context, memberID uuid.UUID) (bool, error)
	GetBillingAddressByID(ctx context.Context, memberBillingID uuid.UUID) (*entities.BillingAddress, error)
	AddMemberStoresById(ctx *gin.Context, memberID uuid.UUID, stores []uuid.UUID) error

	// Country and State Checks

	CountryExists(countryName string) (bool, error)
	StateExists(stateName string, countryCode string) (bool, error)

	// Provider and Middleware

	ProviderExists(ctx context.Context, providerName string) (bool, error)
	Middleware(ctx context.Context, token string) (string, error)

	// Email Checks and Record Counts

	CheckEmailExists(ctx context.Context, partnerID, email string) (bool, error)
	GetMemberRecordCount(context.Context) (int64, error)
	GetResetKey(ctx context.Context, memberID uuid.UUID) string
	CheckEmailForMemberID(ctx *gin.Context, memberID uuid.UUID, email string) (bool, error)
	CheckEmailProviderRelation(ctx *gin.Context, email string, provider string) (bool, error)
	PasswordMemberRelation(ctx *gin.Context, memberID uuid.UUID, hashedPassword string) (bool, error)
	CheckMemberPartner(ctx *gin.Context, memberID uuid.UUID, partnerIDStr string) (bool, error)

	// Subscription Handling

	HandleSubscriptionCheckout(ctx context.Context, memberID uuid.UUID, checkoutData entities.CheckoutSubscription) error
	GetSubscriptionStatusName(ctx *gin.Context, subscriptionID string) (string, error)
	GetSubscriptionCountForLastYear(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (int, error)
	GetMaxSubscriptionLimitForID(ctx *gin.Context, subscriptionID string) (int, error)
	IsFreeSubscription(ctx context.Context, subscriptionID string) (bool, error)
	CheckIfPayoutGatewayExists(ctx *gin.Context, paymentGatewayID int) (bool, error)
	CheckIfMemberSubscribedToFreePlan(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (bool, error)
	HasSubscribedToOneTimePlan(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (bool, error)
	IsSubscriptionFree(ctx context.Context, subscriptionID string) (bool, error)
	IsMemberSubscribedToPlan(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (bool, error)
	IsMemberSubscribedToFreePlan(ctx *gin.Context, memberID uuid.UUID, MemberSubscriptionID string) (bool, error)
	HandleSubscriptionRenewal(ctx context.Context, memberID uuid.UUID, checkoutData entities.SubscriptionRenewal) error
	CheckCancellationEnabled(ctx *gin.Context, subscriptionID string) (bool, error)
	HandleSubscriptionCancellation(ctx context.Context, memberID uuid.UUID, checkoutData entities.CancelSubscription) error
	GetPaymentDetailsByPartnerAndGateway(ctx context.Context, partnerID string, paymentGatewayID int) (string, error)
	DecryptPaymentData(ctx *gin.Context, data string) (string, error)

	// Address Updates and Switching

	UpdatePrimaryBillingAddressToFalseAndRandom(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID) error
	CountTotalAddressesForMember(ctx *gin.Context, memberID uuid.UUID) (int, error)
	UpdateRandomBillingAddressToPrimary(ctx *gin.Context, memberID, memberBillingID uuid.UUID) error
	SubscriptionProductSwitch(context.Context, uuid.UUID, entities.SwitchSubscriptions, *map[string][]string) (map[string][]string, error)
	ViewAllSubscriptions(context.Context, uuid.UUID, entities.ReqParams, *map[string][]string) ([]entities.ListAllSubscriptions, error)
	GetSubscriptionRecordCount(context.Context, uuid.UUID) (int64, error)

	// Existence and Relation Checks

	DeleteBillingAddress(*gin.Context, uuid.UUID, uuid.UUID) error
	IsPartnerIdCorrespondsToGateway(ctx context.Context, partnerID string, paymentGatewayID int) (bool, error)
	HasProductsReleaseEndDateGreaterThanToday(ctx *gin.Context, memberSubscriptionID string) (bool, error)
	IsMemberRelatedToSubscription(ctx *gin.Context, memberID uuid.UUID, memberSubscriptionID string) (bool, error)
	GetSubscriptionIDByMemberSubscriptionID(ctx context.Context, memberSubscriptionID string) (uuid.UUID, error)
	IsSubscriptionAboutInWarning(ctx context.Context, memberSubscriptionID string) (bool, string, error)
	IsSubscriptionInGracePeriod(ctx context.Context, memberID uuid.UUID, memberSubscriptionID string) (bool, time.Time, time.Time, time.Duration, bool, error)
	CheckSubscriptionExistenceAndStatusForCheckout(ctx *gin.Context, SubscriptionID string) (exists bool, isActive bool, err error)
	CheckSubscriptionExistenceAndStatusForRenewal(ctx *gin.Context, memberSubscriptionID string) (bool, bool, error)
}

// NewMemberRepo creates a new instance of MemberRepo.
// NewMemberRepo creates a new instance of MemberRepo.
func NewMemberRepo(db *sql.DB, cfg *entities.EnvConfig) *MemberRepo {
	return &MemberRepo{
		db:  db,
		Cfg: cfg,
	}
}

// IsMemberExists checks if a member exists in the database.
// Parameters:
//	- memberID: The UUID of the member to check for existence.
//	- ctx: The context for the database operation.

// Returns:
// - bool: A boolean indicating whether the member exists.
// - error: An error, if any, during the database operation.
func (member *MemberRepo) IsMemberExists(memberId uuid.UUID, ctxt context.Context) (bool, error) {
	var exists int
	//Checking if member with the passed ID exists
	isMemberExistsQ := `select 1 from member where id = $1`
	row := member.db.QueryRowContext(ctxt, isMemberExistsQ, memberId)
	err := row.Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("member with ID %s does not exist", memberId)
		}
		return false, fmt.Errorf("error checking member existence: %v", err)
	}

	return true, nil
}

// AddBillingAddress inserts a new billing address for a member into the database.

// Parameters:
//	- ctx: The context for the database operation.
//	- memberID: The UUID of the member for whom the billing address is being added.
//	- billingAddress: The billing address details to be added.

// Returns:
//
//   - int: The result code (not used in this context).
//   - error: An error, if any, during the database operation.
func (member *MemberRepo) AddBillingAddress(ctx *gin.Context, memberID uuid.UUID, billingAddress entities.BillingAddress) error {
	//Checking if member already has maximum billing address count
	// Checking if member already has maximum billing address count
	billingAddressCount, err := member.GetBillingAddressCountForMember(ctx, memberID)
	if err != nil {
		return fmt.Errorf("failed to check billing address count: %v", err)
	}

	// If the member already has 5 billing addresses, return an error.
	if billingAddressCount >= 5 {
		return errors.New("member already has the maximum allowed number of billing addresses")
	}

	// If no errors so far, proceed with inserting the new billing address.
	_, err = member.db.ExecContext(ctx, `
							INSERT INTO member_billing_address (member_id, address, zip, country_code, state_code, is_primary_billing)
							VALUES ($1, $2, $3, $4, $5, $6)
						`, memberID, billingAddress.Address, billingAddress.Zipcode, billingAddress.Country, billingAddress.State, billingAddress.Primary)

	if err != nil {
		return fmt.Errorf("failed to insert billing address: %v", err)
	}
	// If insertion is successful, return nil to indicate success.
	return nil
}

// UpdateBillingAddress updates an existing billing address for a member in the database.

// Parameters:
//   - ctx: The context for the database operation.
//   - memberID: The UUID of the member for whom the billing address is being updated.
//   - billingAddress: The updated billing address details.
//
// Returns:
//   - error: An error, if any, during the database operation.
func (member *MemberRepo) UpdateBillingAddress(ctx context.Context, memberID uuid.UUID, memberBillingID uuid.UUID, billingAddress entities.BillingAddress) error {
	// Check if such a member exists
	memberExists, err := member.IsMemberExists(memberID, ctx)
	if err != nil {
		return err
	}
	if !memberExists {
		return errors.New("member does not exist")
	}

	// Prepare the dynamic update query and parameters
	updateQry := "UPDATE member_billing_address SET "
	var params []interface{}
	paramCount := 1

	// Helper function to append fields to the update query
	appendField := func(field string, value interface{}) {
		if paramCount > 1 {
			updateQry += ", "
		}
		updateQry += fmt.Sprintf("%s = $%d", field, paramCount)
		params = append(params, value)
		paramCount++
	}

	// Append fields to the update query if they are not empty
	if billingAddress.Address != "" {
		appendField("address", billingAddress.Address)
	}
	if billingAddress.Zipcode != "" {
		appendField("zip", billingAddress.Zipcode)
	}
	if billingAddress.Country != "" {
		appendField("country_code", billingAddress.Country)
	}
	if billingAddress.State != "" {
		appendField("state_code", billingAddress.State)
	}
	// Handle the Primary field, default to false if not provided
	appendField("is_primary_billing", billingAddress.Primary)

	// Add the WHERE clause
	updateQry += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, memberBillingID)

	// Execute the dynamic update query
	_, err = member.db.ExecContext(ctx, updateQry, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMember updates a member's information in the repository.

// Parameters:
//   - ctx: The context for the database operation.
//   - memberID: The UUID of the member whose information is being updated.
//   - args: The updated member information.

// Returns:
//   - error: An error, if any, during the database operation.
func (member *MemberRepo) UpdateMember(ctx context.Context, memberID uuid.UUID, args entities.Member) error {
	// Check if the member with the given memberID exists
	checkMemberQry := `SELECT 1 FROM member WHERE id = $1`
	var memberExists int
	row := member.db.QueryRowContext(ctx, checkMemberQry, memberID)
	if err := row.Scan(&memberExists); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err // Handle other database errors if needed
	}

	// Prepare the dynamic update query and parameters
	updateQry := "UPDATE member"
	var params []interface{}

	// Helper function to append fields to the update query
	appendField := func(field string, value interface{}) {
		if len(params) > 0 {
			updateQry += ", "
		} else {
			updateQry += " SET "
		}
		updateQry += fmt.Sprintf("%s = $%d", field, len(params)+1)
		params = append(params, value)
	}

	// Append fields to the update query if they are not empty
	if args.Title != "" {
		appendField("title", args.Title)
	}
	if args.FirstName != "" {
		appendField("firstname", args.FirstName)
	}
	if args.LastName != "" {
		appendField("lastname", args.LastName)
	}
	if args.Gender != "" {
		appendField("gender", args.Gender)
	}
	if args.Language != "" {
		appendField("language_code", args.Language)
	}
	if args.Country != "" {
		appendField("country_code", args.Country)
	}
	if args.State != "" {
		appendField("state_code", args.State)
	}
	if args.Address1 != "" {
		appendField("address1", args.Address1)
	}
	if args.City != "" {
		appendField("city", args.City)
	}
	if args.Zipcode != "" {
		appendField("zip", args.Zipcode)
	}
	if args.Phone != "" {
		appendField("mobile", args.Phone)
	}
	if args.Address2 != "" {
		appendField("address2", args.Address2)
	}

	// Check if there are any fields to update
	if len(params) == 0 {
		return errors.New("No fields to update")
	}

	// Add the WHERE clause
	updateQry += fmt.Sprintf(" WHERE id = $%d", len(params)+1)
	params = append(params, memberID)

	// Execute the dynamic update query
	_, err := member.db.ExecContext(ctx, updateQry, params...)
	if err != nil {
		return err
	}

	return nil
}

// GetPasswordHash retrieves the password hash for a member.
// This function queries the database to fetch the password hash associated with a member's ID.
// Parameters:
//   - ctx: The context for the operation.
//   - memberID: The UUID of the member whose password hash is being retrieved.
//
// Returns:
//   - If successful, it returns the password hash as a string and nil (no error).
//   - If there's an error in the database operation, it returns an empty string and an error.
func (m *MemberRepo) GetPasswordHash(ctx context.Context, memberID uuid.UUID) (string, error) {
	var passwordHash string
	err := m.db.QueryRowContext(ctx, `
        SELECT password FROM member WHERE id = $1`, memberID).Scan(&passwordHash)
	if err != nil {
		return "", err
	}

	return passwordHash, nil
}

// CheckResetKeyMatch checks if the key in the database matches the entered key and verifies the timestamp.
func (m *MemberRepo) CheckResetKeyMatch(ctx context.Context, memberID uuid.UUID, key string) (bool, error) {
	var storedKey string
	var expirationTimestamp time.Time

	// Fetch stored reset key and its associated expiration timestamp from the database
	err := m.db.QueryRowContext(ctx, `
					SELECT reset_password_key, password_expiry
					FROM member 
					WHERE id = $1;
				`, memberID).Scan(&storedKey, &expirationTimestamp)

	if err != nil {
		return false, err
	}

	// Check if the stored key matches the entered key
	if storedKey != key {
		return false, errors.New("reset key does not match")
	}

	// Verify if the reset key has expired based on its expiration timestamp
	currentTimestamp := time.Now()
	if currentTimestamp.After(expirationTimestamp) {
		return false, fmt.Errorf("timeout: reset key has expired, please try resetting again")
	}

	// If everything is valid, return true indicating the key matches and is not expired
	return true, nil
}

// UpdatePassword updates the password hash for a member.
// This function updates the password hash in the database for the specified member.
// Parameters:
//   - ctx: The context for the operation.
//   - memberID: The UUID of the member whose password hash is being updated.
//   - newPasswordHash: The new password hash to be set for the member.
//
// Returns:
//   - If successful, it returns nil (no error).
//   - If there's an error in the database operation, it returns an error.
func (m *MemberRepo) UpdatePassword(ctx context.Context, memberID uuid.UUID, key string, newPasswordHash string) error {
	// Begin a new transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// First, update the password for the member
	_, err = tx.ExecContext(ctx, `
        UPDATE member SET password = $1 WHERE id = $2`, newPasswordHash, memberID)
	if err != nil {
		return err
	}

	// Then, set the reset_password_key to 'invalidated'
	_, err = tx.ExecContext(ctx, `
        UPDATE member SET reset_password_key = 'invalidated' WHERE id = $1`, memberID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBillingAddresses retrieves all billing addresses associated with a member.
// If no billing addresses are found for the member, it returns an error "No record found."

// Parameters:
//   - ctx: The context for the database operation.
//   - memberID: The UUID of the member for whom billing addresses are being retrieved.
//
// Returns:
//
//   - []entities.BillingAddress: A slice of billing addresses associated with the member.
//   - error: An error, if any, during the retrieval process.
func (member *MemberRepo) GetAllBillingAddresses(ctx context.Context, memberID uuid.UUID) ([]entities.BillingAddress, error) {
	// Check if given memberID exists
	memberExists, err := member.IsMemberExists(memberID, ctx)
	if err != nil {
		return nil, err
	}
	if !memberExists {
		return nil, errors.New("Member not found") // Return custom error message
	}

	// Initialize a slice to store billing addresses
	var billingAddresses []entities.BillingAddress

	// SQL query to retrieve billing addresses for a specific member_id
	query := `
       SELECT COALESCE(address, ''), COALESCE(zip, ''), 
       COALESCE(country_code, ''), COALESCE(state_code, ''), is_primary_billing
       FROM member_billing_address
       WHERE member_id = $1
    `

	rows, err := member.db.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var billingAddress entities.BillingAddress
		err := rows.Scan(
			&billingAddress.Address,
			&billingAddress.Zipcode,
			&billingAddress.Country,
			&billingAddress.State,
			&billingAddress.Primary,
		)
		if err != nil {
			return nil, err
		}

		// Append the retrieved billing address to the slice
		billingAddresses = append(billingAddresses, billingAddress)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return billingAddresses, nil
}

// CheckEmailExists checks if an email address already exists in the database.

// Parameters:

//		@ ctx:   The context for the database operation.
//		@ email: The email string to be checked.

// Returns:

// @ emailVal: The email value found in the database (if it exists).
// @ err:      An error, if any, during the database operation.
func (member *MemberRepo) CheckEmailExists(ctx context.Context, partnerID, email string) (bool, error) {
	// var ctxt *gin.Context
	var exists bool

	// SQL query for checks if an email address already exists.
	checkEmailExistsQ := `SELECT EXISTS (SELECT 1 FROM member WHERE email = $1 AND partner_id = $2)`

	row := member.db.QueryRowContext(ctx, checkEmailExistsQ, email, partnerID)

	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// RegisterMember registers a new member in the repository.

// Parameters:
//
//	@ ctx: The context for the database operation.
//	@ args: An instance of the Member struct containing the member's information
//	  including first name, last name, email, password, terms condition status,
//	  and tax payment status.
//
// Returns:
//
//	@ err: An error, if any, during the database operation.
//
// RegisterMember registers a new member after checking the existence of the provider.
func (member *MemberRepo) RegisterMember(ctx context.Context, args entities.Member, partnerID string) (err error) {

	// Hash the password before storing it in the database.
	var hashedPassword string
	if args.Provider == consts.ProviderInternal {
		hashedPassword, err = crypto.Hash(args.Password)

		if err != nil {
			return
		}
	}

	getOauthProviderID := `SELECT id FROM oauth_provider WHERE name = $1`
	row := member.db.QueryRowContext(ctx, getOauthProviderID, args.Provider)

	var providerId uuid.UUID
	err = row.Scan(&providerId)

	if err != nil {
		return
	}

	// SQL query for inserting a new member record.
	insertQry := fmt.Sprintf(`INSERT INTO member 
				(firstname,lastname,email,password,is_terms_condition_checked,is_paying_tax,partner_id,oauth_provider_id)
				values(%s)`, utils.PreparePlaceholders(8))

	_, err = member.db.Exec(insertQry, args.FirstName,
		args.LastName, args.Email, hashedPassword,
		args.TermsConditionChecked, args.PayingTax,
		partnerID, providerId,
	)

	// Return any error encountered during the database operation.
	if err != nil {
		return
	}

	return
}

// ProviderExists checks if a provider with the given name exists in the database.
func (member *MemberRepo) ProviderExists(ctx context.Context, providerName string) (bool, error) {
	getOauthProviderID := `SELECT EXISTS(SELECT 1 FROM oauth_provider WHERE name = $1)`
	var exists bool
	err := member.db.QueryRowContext(ctx, getOauthProviderID, providerName).Scan(&exists)

	// Return any error encountered during the database operation.
	if err != nil {
		return false, err
	}

	return exists, nil
}

// ViewMemberProfile retrieves a member's profile information, including billing addresses,
// from the database.
//
// Parameters:
//
//	@memberId: The unique identifier of the member whose profile is being retrieved.
//	@ctx: The context for the database operation.
//
// Returns:
//
//	@memberProfile: An instance of entities.MemberProfile containing the member's profile details.
//	@err: An error, if any, encountered during the database operation.

func (member *MemberRepo) ViewMemberProfile(memberId uuid.UUID, ctx context.Context) (entities.MemberProfile, error) {

	// Query to fetch member profile details
	getMemberProfileQ := `
		SELECT
			COALESCE(m.title, ''),
			COALESCE(m.firstname, ''),
			COALESCE(m.lastname, ''),
			COALESCE(m.gender, ''),
			m.email,
			COALESCE(m.mobile, ''),
			COALESCE(m.address1, ''),
			COALESCE(m.address2, ''),
			COALESCE(m.country_code, ''),
			COALESCE(m.state_code, ''),
			COALESCE(m.city, ''),
			COALESCE(m.zip, ''),
			COALESCE(m.language_code, ''),
			m.is_mail_subscribed
		FROM
			member m
		WHERE
			m.id = $1
	`

	// Query to fetch billing address details
	getBillingAddressQ := `
		SELECT
			COALESCE(address, ''),
			COALESCE(zip, ''),
			COALESCE(country_code, ''),
			COALESCE(state_code, ''),
			is_primary_billing
		FROM
			member_billing_address
		WHERE
			member_id = $1
	`

	var memberProfile entities.MemberProfile

	// Fetch member profile details
	err := member.db.QueryRowContext(ctx, getMemberProfileQ, memberId).Scan(
		&memberProfile.MemberDetails.Title,
		&memberProfile.MemberDetails.FirstName,
		&memberProfile.MemberDetails.LastName,
		&memberProfile.MemberDetails.Gender,
		&memberProfile.MemberDetails.Email,
		&memberProfile.MemberDetails.Phone,
		&memberProfile.MemberDetails.Address1,
		&memberProfile.MemberDetails.Address2,
		&memberProfile.MemberDetails.Country,
		&memberProfile.MemberDetails.State,
		&memberProfile.MemberDetails.City,
		&memberProfile.MemberDetails.Zipcode,
		&memberProfile.MemberDetails.Language,
		&memberProfile.EmailSubscribed,
	)

	if err != nil {
		return memberProfile, err
	}

	// Fetch billing address details
	rows, err := member.db.QueryContext(ctx, getBillingAddressQ, memberId)
	if err != nil {
		return memberProfile, err
	}
	defer rows.Close()

	var billingAddresses []entities.BillingAddress
	for rows.Next() {
		var billingAddress entities.BillingAddress
		err := rows.Scan(
			&billingAddress.Address,
			&billingAddress.Zipcode,
			&billingAddress.Country,
			&billingAddress.State,
			&billingAddress.Primary,
		)

		if err != nil {
			return memberProfile, err
		}

		billingAddresses = append(billingAddresses, billingAddress)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		return memberProfile, err
	}

	// Assign billing addresses to member profile if available
	if len(billingAddresses) > 0 {
		memberProfile.MemberBillingAddress = billingAddresses
	} else {
		// If no billing address is found, you can handle it here.
		// For example, you can log a message or return a specific error.
		fmt.Println("No billing address found for member:", memberId)
		// Here, you can either return an error or set some default value or handle as per your application's requirement.
	}

	return memberProfile, nil
}

// ViewMembers retrieves a list of member profiles with additional information,
// such as their roles, partner names, album, track, and artist counts, from the database.
//
// Parameters:
//
//	@ctx: The context for the database operation.
//	@params: An instance of entities.Params containing filtering and pagination parameters.
//
// Returns:
//
//	@members: A slice of entities.ViewMembers, each representing a member's profile.
//	@err: An error, if any, encountered during the database operation.
func (member *MemberRepo) ViewMembers(ctx context.Context, params entities.Params) ([]entities.ViewMembers, error) {
	members := make([]entities.ViewMembers, 0)
	viewMembersQ := `
    SELECT
        m.id,
        COALESCE(CONCAT(m.firstname, ' ', m.lastname), '') AS memberName,
        COALESCE(m.gender, '') AS gender,
        m.member_role_id,
        l.name AS roleName,
        COALESCE(p.name ,'') AS partnerName,
        m.email,
        COALESCE(m.country_code,''),
        COALESCE(c.name,'') AS countryName,
        m.is_active,
        (
            SELECT COUNT(id)
            FROM product
            WHERE member_id = m.id
            AND is_active = $1
            AND is_deleted = $2
        ) AS albumCount,
        (
            SELECT COUNT(id)
            FROM track
            WHERE member_id = m.id
            AND is_deleted = $3
        ) AS trackCount,
        (
            SELECT COUNT(id)
            FROM artist
            WHERE member_id = m.id
            AND is_active = $4
            AND is_deleted = $5
        ) AS artistCount
    FROM
        member m
    INNER JOIN
        lookup l ON l.id = m.member_role_id
    INNER JOIN
        partner p ON p.id = m.partner_id
    LEFT JOIN
        country c ON c.iso = m.country_code
    WHERE
        m.is_deleted = $6
    `

	// Build the query based on the provided parameters
	conditions := []string{"1 = 1"}
	parameters := []interface{}{true, false, false, true, false, false}
	// Check if params.Limit is empty
	if params.Limit == 0 {
		logger.Log().Error("Limit parameter is empty")
		return nil, errors.New("limit parameter is empty")
	}
	// Convert params.Limit to an integer
	if params.Limit > consts.MaximumLimit {
		err := fmt.Errorf("Exceeds maximum record limit")
		fmt.Println(err.Error())
		return nil, err
	}

	// Add conditions based on params
	if params.Status != "" {
		conditions = append(conditions, fmt.Sprintf("m.is_active = %t", params.Status == consts.Active))
	} else {
		// If no Status param is provided, default to showing only active members
		conditions = append(conditions, "m.is_active = true")
	}

	if params.Country != "" {
		conditions = append(conditions, fmt.Sprintf("m.country_code = '%s'", params.Country))
	}

	if params.Partner != "" && params.Partner != "-1" {
		partner, err := uuid.Parse(params.Partner)
		if err != nil {
			return members, err
		}
		conditions = append(conditions, fmt.Sprintf("m.partner_id = '%s'", partner))
	}

	if params.Role != "" && params.Role != "-1" {
		role, err := strconv.Atoi(params.Role)
		if err != nil {
			return members, err
		}
		conditions = append(conditions, fmt.Sprintf("m.member_role_id = %d", role))
	}

	if params.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("m.gender = '%s'", params.Gender))
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(p.name LIKE '%%%s%%' OR CONCAT(m.firstname, ' ', m.lastname) LIKE '%%%s%%' OR m.email LIKE '%%%s%%' OR m.firstname LIKE '%%%s%%')", params.Search, params.Search, params.Search, params.Search)
		conditions = append(conditions, searchCondition)
	}

	// Construct the WHERE clause
	if len(conditions) > 0 {
		viewMembersQ += " AND " + strings.Join(conditions, " AND ")
	}

	// Add sorting logic
	if params.SortBy == "" || (params.SortBy == "name" && params.Order == "") {
		// If SortBy is not provided or if explicitly sorting by name with no order specified, default to sorting by first name in ascending order
		viewMembersQ += " ORDER BY NULLIF(COALESCE(m.firstname, ''), '') ASC, NULLIF(COALESCE(m.lastname, ''), '') ASC"
	} else {
		switch params.SortBy {
		case "firstname":
			viewMembersQ += fmt.Sprintf(" ORDER BY NULLIF(m.firstname, '') %s", params.Order)
		case "lastname":
			viewMembersQ += fmt.Sprintf(" ORDER BY NULLIF(m.lastname, '') %s", params.Order)
		default:
			viewMembersQ += fmt.Sprintf(" ORDER BY m.%s %s", params.SortBy, params.Order)
		}
	}

	// Add pagination
	if params.Page != 0 && params.Limit != 0 {
		page := int(params.Page)
		limit := int(params.Limit)
		offset := (page - 1) * limit

		viewMembersQ += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, limit)
	}

	rows, err := member.db.QueryContext(ctx, viewMembersQ, parameters...)
	if err != nil {
		return members, err
	}

	defer rows.Close()

	for rows.Next() {
		var member entities.ViewMembers
		err := rows.Scan(
			&member.MemberId,
			&member.Name,
			&member.Gender,
			&member.Role.Id,
			&member.Role.Name,
			&member.PartnerName,
			&member.Email,
			&member.Country.Code,
			&member.Country.Name,
			&member.Active,
			&member.AlbumCount,
			&member.TrackCount,
			&member.ArtistCount,
		)

		if err != nil {
			return members, err
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return members, err
	}

	return members, nil
}

//GetFilteredRecordCount

func (member *MemberRepo) GetFilteredRecordCount(ctx context.Context, params entities.Params) (int64, error) {
	viewMembersCountQ := `
	SELECT COUNT(*) AS recordCount
	FROM
		member m
	INNER JOIN
		lookup l ON l.id = m.member_role_id
	INNER JOIN
		partner p ON p.id = m.partner_id
	LEFT JOIN
		country c ON c.iso = m.country_code
	WHERE
		m.is_deleted = $1
`

	// Build the query based on the provided parameters
	conditions := []string{"1 = 1"}
	parameters := []interface{}{false}

	// Add conditions based on params
	if params.Status != "" {
		conditions = append(conditions, fmt.Sprintf("m.is_active = %t", params.Status == consts.Active))
	} else {
		// If no Status param is provided, default to showing only active members
		conditions = append(conditions, "m.is_active = true")
	}

	if params.Country != "" {
		conditions = append(conditions, fmt.Sprintf("m.country_code = '%s'", params.Country))
	}

	if params.Partner != "" && params.Partner != "-1" {
		partner, err := uuid.Parse(params.Partner)
		if err != nil {
			return 0, err
		}
		conditions = append(conditions, fmt.Sprintf("m.partner_id = '%s'", partner))
	}

	if params.Role != "" && params.Role != "-1" {
		role, err := strconv.Atoi(params.Role)
		if err != nil {
			return 0, err
		}
		conditions = append(conditions, fmt.Sprintf("m.member_role_id = %d", role))
	}

	if params.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("m.gender = '%s'", params.Gender))
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(p.name LIKE '%%%s%%' OR CONCAT(m.firstname, ' ', m.lastname) LIKE '%%%s%%' OR m.email LIKE '%%%s%%')", params.Search, params.Search, params.Search)
		conditions = append(conditions, searchCondition)
	}

	// Construct the WHERE clause
	if len(conditions) > 0 {
		viewMembersCountQ += " AND " + strings.Join(conditions, " AND ")
	}

	rows, err := member.db.QueryContext(ctx, viewMembersCountQ, parameters...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var totalCount int64
	if rows.Next() {
		err := rows.Scan(&totalCount)
		if err != nil {
			return 0, err
		}
	}

	return totalCount, nil
}

//GetBasicMemberDetailsByEmail retrieves basic member details from the database.
//
// Parameters:
//   - args (entities.MemberPayload): An instance of entities.MemberPayload containing member information.
//   - ctx (context.Context): The context for the database operation.
//
// Returns:
//   - basicMemberData (entities.BasicMemberData): A struct containing the basic member details.
//   - error: An error, if any, encountered during the database operation.

func (member *MemberRepo) GetBasicMemberDetailsByEmail(partnerID string, args entities.MemberPayload, ctx context.Context) (entities.BasicMemberData, error) {

	var basicMemberData entities.BasicMemberData

	getBasicMemberDataQ := `
		SELECT
			m.id,
			CONCAT(m.firstname, ' ', m.lastname) AS memberName,
			p.name,
			p.id,
			m.email,
			m.oauth_provider_id,
			l.name AS user_type,
			ar.name AS user_roles
		FROM
			member m
		INNER JOIN
			partner p
			ON p.id = m.partner_id
		INNER JOIN
			lookup l
			ON l.id = m.member_role_id
		INNER JOIN
			member_access_role mar
			ON mar.member_id = m.id
		INNER JOIN
			access_role ar
			ON ar.id = mar.role_id
		WHERE
			m.partner_id = $1
		AND 
			m.email = $2
	`

	hashedPassword, err := crypto.Hash(args.Password)
	if err != nil {
		return basicMemberData, err
	}

	if args.Provider == consts.ProviderInternal {
		getBasicMemberDataQ = fmt.Sprintf("%s AND m.password = '%s'", getBasicMemberDataQ, hashedPassword)
	}

	rows, err := member.db.QueryContext(ctx, getBasicMemberDataQ, partnerID, args.Email)

	if err != nil {
		return basicMemberData, err
	}

	defer rows.Close()

	for rows.Next() {
		var roles sql.NullString // Use sql.NullString to handle potential NULL values

		err := rows.Scan(
			&basicMemberData.MemberID,
			&basicMemberData.Name,
			&basicMemberData.PartnerName,
			&basicMemberData.PartnerID,
			&basicMemberData.Email,
			&basicMemberData.ProviderID,
			&basicMemberData.MemberType,
			&roles,
		)
		if err != nil {
			return basicMemberData, err
		}

		// Check for NULL roles and append them if they are not NULL
		if roles.Valid {
			basicMemberData.MemberRoles = append(basicMemberData.MemberRoles, roles.String)
		}
	}

	if err := rows.Err(); err != nil {
		return basicMemberData, err
	}

	return basicMemberData, nil
}

func (member *MemberRepo) Middleware(ctx context.Context, token string) (string, error) {

	var (
		blacklistedToken string
		log              = log.Log().WithContext(ctx)
	)

	tokenBlacklisted := `
	SELECT  active_token
	FROM public.refresh_token
	WHERE is_revoked = false 
	AND active_token=$1;`

	row := member.db.QueryRowContext(
		ctx,
		tokenBlacklisted,
		token,
	)

	err := row.Scan(&blacklistedToken)
	if err != nil {
		log.Printf("scan:%v", err)
		return blacklistedToken, err
	}

	return blacklistedToken, nil
}

// GetSubscriptionRecordCount function is used to calculate and return total count of subscriptions
func (member *MemberRepo) GetMemberRecordCount(ctx context.Context) (int64, error) {
	var totalCount int64

	query := `
	SELECT COUNT(*) AS recordCount
FROM (
    SELECT
        m.id,
        COALESCE(CONCAT(m.firstname, ' ', m.lastname), '') AS memberName,
        COALESCE(m.gender, '') AS gender,
        m.member_role_id,
        l.name AS roleName,
        COALESCE(p.name ,'') AS partnerName,
        m.email,
        COALESCE(m.country_code,''),
        COALESCE(c.name,'') AS countryName,
        m.is_active,
        (
            SELECT COUNT(id)
            FROM product
            WHERE member_id = m.id
            AND is_active = $1
            AND is_deleted = $2
        ) AS albumCount,
        (
            SELECT COUNT(id)
            FROM track
            WHERE member_id = m.id
            AND is_deleted = $3
        ) AS trackCount,
        (
            SELECT COUNT(id)
            FROM artist
            WHERE member_id = m.id
            AND is_active = $4
            AND is_deleted = $5
        ) AS artistCount
    FROM
        member m
    INNER JOIN
        lookup l ON l.id = m.member_role_id
    INNER JOIN
        partner p ON p.id = m.partner_id
    LEFT JOIN
        country c ON c.iso = m.country_code
    WHERE
        m.is_deleted = $6
) AS subquery;

`
	row := member.db.QueryRowContext(ctx, query, true, false, false, true, false, false)

	if err := row.Scan(&totalCount); err != nil {
		logger.Log().WithContext(ctx).Errorf("Getting Member count failed, QueryRowContext failed, err=%s", err.Error())
		return 0, err
	}

	return totalCount, nil
}

// Get Count of Billing Address
// getBillingAddressCountForMember fetches the count of billing addresses for a given member ID.
func (member *MemberRepo) GetBillingAddressCountForMember(ctx *gin.Context, memberID uuid.UUID) (int, error) {
	var billingAddressCount int

	err := member.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM member_billing_address
		WHERE member_id = $1
	`, memberID).Scan(&billingAddressCount)

	if err != nil {
		return 0, fmt.Errorf("failed to check billing address count: %v", err)
	}

	return billingAddressCount, nil
}

// InitiatePasswordReset initiates the password reset process for a member.
func (member *MemberRepo) InitiatePasswordReset(ctx *gin.Context, memberID uuid.UUID, email string) (string, error) {
	var dbEmail string

	// Execute the SQL query to fetch the email for the given memberID
	err := member.db.QueryRowContext(ctx, `
					SELECT email FROM public.member WHERE id = $1;
				`, memberID).Scan(&dbEmail)

	if err != nil {
		return "", err
	}

	// Check if the fetched email from the database matches the incoming email
	if dbEmail != email {
		return "", fmt.Errorf("email mismatch: provided email does not match with current email ")
	}

	// If the emails match, proceed to generate a reset key
	NewResetKey, err := GenerateResetKey(ctx, dbEmail) // Implement the logic to generate a reset key

	if err != nil {
		return "", fmt.Errorf("failed to generate reset key: %v", err)
	}

	// Update the member table to set the reset password key and its expiration timestamp for the row with matching email
	expirationTime := time.Now().Add(50 * time.Minute) // Calculate expiration time: current time + 10 minutes
	_, err = member.db.ExecContext(ctx, `
					UPDATE public.member 
					SET reset_password_key = $1, 
					password_expiry = $2 
					WHERE email = $3;
				`, NewResetKey, expirationTime, dbEmail)

	if err != nil {
		panic(fmt.Errorf("failed to update reset key and expiration timestamp in member table: %v", err))
	}

	// Return the generated reset key
	return NewResetKey, nil
}

// generateResetKey generates a reset key using the user's email and current timestamp.
// the length of reset key is 6 digits

func GenerateResetKey(ctx *gin.Context, email string) (string, error) {
	// Generate a unique string using the email and current timestamp
	uniqueString := fmt.Sprintf("%s-%d", email, time.Now().UnixNano())

	// Hash the unique string using SHA-256 to create a secure reset key
	hash := sha256.New()
	_, err := hash.Write([]byte(uniqueString))
	if err != nil {
		return "", fmt.Errorf("failed to generate hash: %v", err)
	}

	// Convert the hashed bytes to a hexadecimal string (reset key)
	resetKey := hex.EncodeToString(hash.Sum(nil))
	hasher := sha256.New()
	hasher.Write([]byte(resetKey))
	NewKey := hex.EncodeToString(hasher.Sum(nil))
	NewResetKey := NewKey[:16]
	return NewResetKey, nil
}

// Function to check if country exists
func (member *MemberRepo) CountryExists(countryName string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM country WHERE iso = $1)`

	var exists bool
	err := member.db.QueryRow(query, countryName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Function to check if state exists and is related to a country
func (member *MemberRepo) StateExists(stateName string, countryCode string) (bool, error) {
	// Query to check if the state exists and is related to the provided country code
	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM public.country_state 
            WHERE iso = $1 AND country_code = $2
        )
    `

	var exists bool
	err := member.db.QueryRow(query, stateName, countryCode).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Function to check if the billing address already exists for the member.
func (member *MemberRepo) BillingAddressExists(ctx *gin.Context, memberID uuid.UUID, billingAddress entities.BillingAddress) (bool, error) {
	var addressExists int
	err := member.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM member_billing_address
		WHERE member_id = $1 AND address = $2 AND zip = $3 AND country_code = $4 AND state_code = $5
	`, memberID, billingAddress.Address, billingAddress.Zipcode, billingAddress.Country, billingAddress.State).Scan(&addressExists)
	if err != nil {
		return false, fmt.Errorf("failed to check address existence: %v", err)
	}
	return addressExists > 0, nil
}

// Function to ensure there isn't already a primary billing address for the member.
func (member *MemberRepo) HasPrimaryBilling(ctx *gin.Context, memberID uuid.UUID) (bool, error) {
	var hasPrimaryBilling bool
	err := member.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 
			FROM public.member_billing_address 
			WHERE member_id = $1 AND is_primary_billing = TRUE
		)
	`, memberID).Scan(&hasPrimaryBilling)
	if err != nil {
		return false, fmt.Errorf("failed to check primary billing address existence: %v", err)
	}
	return hasPrimaryBilling, nil
}

// CheckBillingAddressRelation checks if a billing address with the specified ID exists and is associated with the given member.
// It takes a context, member ID, and billing address ID. Returns (true, nil) if the relationship is valid, (false, nil) if not,
// and an error if any issues occur during the checks.

func (member *MemberRepo) CheckBillingAddressRelation(ctx context.Context, memberID, billingAddressID uuid.UUID) (bool, error) {
	var billingExists bool

	// Check if the billing address exists by its ID.
	err := member.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1 
            FROM public.member_billing_address 
            WHERE id = $1
        )
    `, billingAddressID).Scan(&billingExists)

	if err != nil {
		return false, fmt.Errorf("failed to check if billing address exists: %v", err)
	}

	if !billingExists {
		return false, nil // Billing address does not exist.
	}

	// Check if the billing address ID is associated with the given member ID.
	var validRelation bool
	err = member.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM public.member_billing_address 
            WHERE id = $1 AND member_id = $2
        )
    `, billingAddressID, memberID).Scan(&validRelation)

	if err != nil {
		return false, fmt.Errorf("failed to verify relationship between billing address and member: %v", err)
	}

	return validRelation, nil
}

// GetBillingAddressByID retrieves a billing address based on its ID for a specific member.
func (member *MemberRepo) GetBillingAddressByID(ctx context.Context, memberBillingID uuid.UUID) (*entities.BillingAddress, error) {
	// Here's where you would typically make a database query or any other data source operation
	// to fetch the billing address by its ID.

	// Example: Fetch from a database
	var billingAddress entities.BillingAddress
	query := "SELECT address, zip, country_code, state_code, is_primary_billing FROM member_billing_address WHERE id = $1"
	err := member.db.QueryRowContext(ctx, query, memberBillingID).Scan(
		&billingAddress.Address,
		&billingAddress.Zipcode,
		&billingAddress.Country, // Assuming you have a corresponding field in your entities.BillingAddress struct
		&billingAddress.State,   // Assuming you have a corresponding field in your entities.BillingAddress struct
		&billingAddress.Primary, // Assuming you have a corresponding field in your entities.BillingAddress struct
	)

	if err != nil {
		// Handle error, e.g., log it or return an error
		return nil, err
	}

	return &billingAddress, nil
}

//GetMember details by ID

// GetMemberByID fetches a member by their ID from the database
func (member *MemberRepo) GetMemberByID(ctx context.Context, memberID uuid.UUID) (entities.MemberByID, error) {
	query := `SELECT title, firstname, lastname, country_code, state_code, zip, mobile, city, address1, address2 
              FROM member 
              WHERE id = $1`

	var memberInfo entities.MemberByID
	err := member.db.QueryRowContext(ctx, query, memberID).Scan(
		&memberInfo.Title,
		&memberInfo.FirstName,
		&memberInfo.LastName,
		&memberInfo.Country,
		&memberInfo.State,
		&memberInfo.Zipcode,
		&memberInfo.Phone,
		&memberInfo.City,
		&memberInfo.Address1,
		&memberInfo.Address2,
	)

	if err != nil {
		return entities.MemberByID{}, err
	}

	return memberInfo, nil
}

// GetResetKey returns the reset_password_key for a specific member ID.
func (member *MemberRepo) GetResetKey(ctx context.Context, memberID uuid.UUID) string {
	var resetKey sql.NullString // Using sql.NullString to handle NULL values

	// Execute the SQL query to fetch the reset_password_key for the given memberID
	err := member.db.QueryRowContext(ctx, `
		SELECT reset_password_key FROM member WHERE id = $1;
	`, memberID).Scan(&resetKey)

	if err != nil {
		return ""
	}

	// Check if the resetKey is valid (not NULL)
	if resetKey.Valid {
		return resetKey.String
	}

	// If resetKey is NULL or invalid, return an appropriate error message
	return ""
}

// Returns true if the email is associated with the memberID, otherwise false.
func (member *MemberRepo) CheckEmailForMemberID(ctx *gin.Context, memberID uuid.UUID, email string) (bool, error) {
	query := `SELECT id FROM member WHERE id = $1 AND email = $2`

	var foundID uuid.UUID
	err := member.db.QueryRowContext(ctx, query, memberID, email).Scan(&foundID)
	if err == sql.ErrNoRows {
		// No matching record found
		return false, nil
	}
	if err != nil {
		// Handle other errors
		return false, err
	}

	// If a matching memberID is found for the email, return true
	return foundID == memberID, nil
}

// CheckEmailProviderRelation checks if the provided email is associated with the given provider.
func (member *MemberRepo) CheckEmailProviderRelation(ctx *gin.Context, email string, provider string) (bool, error) {
	var exists bool

	// SQL query to check if an email is associated with a specific provider.
	checkRelationQ := `
    SELECT EXISTS (
        SELECT 1 
        FROM member 
        JOIN public.oauth_provider ON member.oauth_provider_id = public.oauth_provider.id
        WHERE member.email = $1 AND public.oauth_provider.name = $2
    )
`

	// Execute the SQL query using QueryRowContext, passing the email and provider from the payload.
	row := member.db.QueryRowContext(ctx, checkRelationQ, email, provider)

	// Scan the result into the exists variable.
	if err := row.Scan(&exists); err != nil {
		return false, err // Return the error if any.
	}

	// Return the exists boolean value and nil error if successful.
	return exists, nil
}

// PasswordMemberRelation checks if the hashedPassword matches the password for a given memberID
func (member *MemberRepo) PasswordMemberRelation(ctx *gin.Context, memberID uuid.UUID, hashedPassword string) (bool, error) {
	var storedPassword string
	query := `SELECT password FROM public.member WHERE id = $1`

	// Query the database to get the stored password for the memberID
	err := member.db.QueryRowContext(ctx, query, memberID).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no rows are returned (memberID not found)
			return false, nil
		}
		// Handle other errors
		return false, fmt.Errorf("error querying database: %v", err)
	}

	// Check if the stored password matches the hashedPassword
	if storedPassword == hashedPassword {
		return true, nil // Passwords match
	}

	return false, nil // Passwords do not match
}

// CountPrimaryBillingAddresses returns the count of records with is_primary_billing set to true for a given memberID
func (member *MemberRepo) CountPrimaryBillingAddresses(ctx *gin.Context, memberID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM public.member_billing_address WHERE member_id = $1 AND is_primary_billing = true`

	// Query the database to get the count of primary billing addresses for the memberID
	err := member.db.QueryRowContext(ctx, query, memberID).Scan(&count)
	if err != nil {
		// Handle errors
		return 0, fmt.Errorf("error querying database: %v", err)
	}

	return count, nil
}

// CountTotalAddressesForMember fetches the count of total addresses related to a specific memberID
func (member *MemberRepo) CountTotalAddressesForMember(ctx *gin.Context, memberID uuid.UUID) (int, error) {
	var count int

	// SQL query to count the total addresses for the given member_id
	query := `SELECT COUNT(*) FROM public.member_billing_address WHERE member_id = $1`

	// Execute the SQL query and scan the result into the count variable
	err := member.db.QueryRowContext(ctx, query, memberID).Scan(&count)
	if err != nil {
		// Handle the error if any
		return 0, fmt.Errorf("failed to fetch total addresses for member: %v", err)
	}

	// Return the count of total addresses for the member
	return count, nil
}

// HandleSubscriptionCheckout handles the checkout process for a subscription.
// It performs the following steps:
//  1. Checks if the provided SubscriptionID and PaymentGatewayID exist in their respective tables.
//  2. Queries the currency_id associated with the subscription plan.
//  3. Fetches subscription duration details.
//  4. Calculates the expiration date based on the subscription duration.
//  5. Inserts a new member_subscription record.
//  6. Inserts data into the member_payout_gateway table.
//
// Parameters:
//   - ctx (context.Context): The context for the database operations.
//   - memberID (uuid.UUID): The unique identifier of the member.
//   - checkoutData (entities.CheckoutSubscription): The checkout data, including SubscriptionID and PaymentGatewayID.
//
// Returns:
//   - error: An error if any database operation fails.
func (member *MemberRepo) HandleSubscriptionCheckout(ctx context.Context, memberID uuid.UUID, checkoutData entities.CheckoutSubscription) error {
	// Start a transaction.
	tx, err := member.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var currencyID, subscriptionDurationValue int

	// Query to fetch currency_id
	row := member.db.QueryRowContext(ctx, `
				SELECT currency_id
				FROM subscription_plan
				WHERE id = $1
			`, checkoutData.SubscriptionID)

	if err := row.Scan(&currencyID); err != nil {

		return err
	}

	// Fetch the subscription duration value from the database
	query := `
			SELECT sd.value
			FROM public.subscription_duration AS sd
			WHERE sd.id = (
				SELECT sp.subscription_duration_id
				FROM public.subscription_plan AS sp
				WHERE sp.id = $1
			);	
	`

	err = member.db.QueryRowContext(ctx, query, checkoutData.SubscriptionID).Scan(&subscriptionDurationValue)

	if err != nil {
		return err
	}

	// Check if the subscription is free
	isFree, err := member.IsFreeSubscription(ctx, checkoutData.SubscriptionID)
	if err != nil {
		return err
	}

	// Calculate the expiration date
	expirationDate := time.Now().Add(time.Duration(subscriptionDurationValue) * 24 * time.Hour)

	combinedQuery := `
		INSERT INTO member_subscription 
		(member_id, subscription_id, expiration_date, member_subscription_status_id, custom_name)
		VALUES 
		($1, $2, $3, (SELECT id FROM member_subscription_status WHERE name = 'active'), $4)
		RETURNING member_subscription_status_id
	`

	// Execute the combined query within the transaction.
	result, err := tx.ExecContext(ctx, combinedQuery, memberID, checkoutData.SubscriptionID, expirationDate, checkoutData.CustomName)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 1 {
		// Insertion was successful
		// You can perform additional actions or return a success message if needed
		fmt.Println("Insertion successful")
	} else {
		// No rows were affected, indicating the insertion failed
		return errors.New("Insertion failed")
	}
	// If the subscription is not free, insert into member_payout_gateway
	if !isFree {
		_, err = member.db.ExecContext(ctx, `
			INSERT INTO member_payout_gateway (member_id, payment_gateway_id, currency_id, payment_details)
			SELECT $1, $2, sp.currency_id, 
			('{"payment_amount": ' || sp.amount || ', "tax_percentage": ' || sp.tax_percentage || ' }')::jsonb
			FROM subscription_plan sp
			WHERE sp.id = $3
		`, memberID, checkoutData.PaymentGatewayID, checkoutData.SubscriptionID)

		if err != nil {
			return err
		}
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// CheckSubscriptionExistenceAndStatus checks if the given subscription ID exists in the database
// and if the subscription plan associated with it is currently active.
// It returns a boolean indicating existence and activation status, along with an error if any.
func (member *MemberRepo) CheckSubscriptionExistenceAndStatusForCheckout(ctx *gin.Context, SubscriptionID string) (exists bool, isActive bool, err error) {
	// Query to check if the subscription exists.
	existenceQuery := `
	SELECT EXISTS (
		SELECT 1
		FROM public.subscription_plan
		WHERE id = $1
	) AS exists;`

	// Query to retrieve the is_active status.
	statusQuery := `
	SELECT is_active
	FROM public.subscription_plan
	WHERE id = $1;`

	// Check if the subscription exists.
	err = member.db.QueryRowContext(ctx, existenceQuery, SubscriptionID).Scan(&exists)

	if err != nil {
		return false, false, fmt.Errorf("error checking subscription existence: %s", err)
	}

	// If the subscription doesn't exist, return early.
	if !exists {
		return false, false, nil
	}

	// Retrieve the is_active status.
	err = member.db.QueryRowContext(ctx, statusQuery, SubscriptionID).Scan(&isActive)

	if err != nil {
		return false, false, fmt.Errorf("error retrieving subscription status: %s", err)
	}

	return exists, isActive, nil
}

// Function to get the subscription status name by subscription_id
func (member *MemberRepo) GetSubscriptionStatusName(ctx *gin.Context, MemberSubscriptionID string) (string, error) {
	var statusName string

	// Query to fetch the subscription status name based on subscription_id
	err := member.db.QueryRowContext(ctx, `
        SELECT ms.name
        FROM public.member_subscription_status ms
        JOIN public.member_subscription msub ON ms.id = msub.member_subscription_status_id
        WHERE msub.id = $1
    `, MemberSubscriptionID).Scan(&statusName)

	if err != nil {
		return "", fmt.Errorf("error fetching subscription status name: %s", err)
	}

	// If no error and statusName is empty, it means no matching record found.
	if statusName == "" {
		return "", fmt.Errorf("no matching subscription status found for the given subscription ID")
	}

	return statusName, nil
}

// GetSubscriptionCountForLastYear checks how many times a member has subscribed to a specific subscription plan in the last year.
func (member *MemberRepo) GetSubscriptionCountForLastYear(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (int, error) {
	// Calculate the start and end date of the last year
	currentYearStart := time.Now().AddDate(-1, 0, 0).Format("2006-01-02 15:04:05")
	currentYearEnd := time.Now().Format("2006-01-02 15:04:05")

	// SQL query to count subscriptions for the given memberID and subscriptionID within the last year
	query := `
		SELECT COUNT(*) 
		FROM public.member_subscription 
		WHERE member_id = $1 
		AND subscription_id = $2 
		AND created_on BETWEEN $3 AND $4;
	`

	var count int

	// Execute the SQL query and scan the result into the count variable
	err := member.db.QueryRowContext(ctx, query, memberID, subscriptionID, currentYearStart, currentYearEnd).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch subscription count: %v", err)
	}

	// Fetch the subscription limit for the given subscriptionID from the subscription_plan table
	var subscriptionLimit int
	err = member.db.QueryRowContext(ctx, "SELECT subscription_limit_per_year FROM public.subscription_plan WHERE id = $1", subscriptionID).Scan(&subscriptionLimit)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch subscription limit per year: %v", err)
	}

	// Check if the member's subscription count exceeds the subscription limit per year
	if count >= subscriptionLimit {
		return count, fmt.Errorf("subscription limit exceeded for the member for the given subscription plan")
	}

	return count, nil
}

// GetMaxSubscriptionLimitForID fetches the maximum subscription limit per year for a given subscriptionID.
func (member *MemberRepo) GetMaxSubscriptionLimitForID(ctx *gin.Context, subscriptionID string) (int, error) {
	// SQL query to fetch the subscription_limit_per_year for the given subscriptionID
	query := `
		SELECT subscription_limit_per_year 
		FROM public.subscription_plan 
		WHERE id = $1;
	`

	var subscriptionLimit int

	// Execute the SQL query and scan the result into the subscriptionLimit variable
	err := member.db.QueryRowContext(ctx, query, subscriptionID).Scan(&subscriptionLimit)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch maximum subscription limit for the given subscription ID: %v", err)
	}

	return subscriptionLimit, nil
}

// IsFreeSubscription checks if a subscription is free.
// Parameters:
//   - ctx (context.Context): The context for the database operation.
//   - subscriptionID (string): The unique identifier of the subscription to check.
//
// Returns:
//   - (bool): True if the subscription is free, false otherwise.
//   - (error): An error if the operation fails.
func (member *MemberRepo) IsFreeSubscription(ctx context.Context, subscriptionID string) (bool, error) {
	var isFreeSubscription bool
	var subscriptionExists bool

	// Query to check if the subscription exists and if it's free.
	err := member.db.QueryRowContext(ctx, `
				SELECT EXISTS (SELECT 1 FROM subscription_plan WHERE id = $1) AS subscription_exists,
				is_free_subscription
				FROM subscription_plan
				WHERE id = $1
			`, subscriptionID).Scan(&subscriptionExists, &isFreeSubscription)

	if !subscriptionExists {
		return false, fmt.Errorf("subscription does not exist")
	}

	if err != nil {
		return false, err
	}

	return isFreeSubscription, nil
}

// CheckIfPayoutGatewayExists checks if a payment gateway with the given ID exists in the partner_payment_gateway table.
// It takes a Gin context and the paymentGatewayID. Returns (true, nil) if the gateway exists, an error if there's an issue
// with the database query, and (false, error) if the gateway does not exist.
func (member *MemberRepo) CheckIfPayoutGatewayExists(ctx *gin.Context, paymentGatewayID int) (bool, error) {
	var paymentGatewayExists bool

	// Check if PaymentGatewayID exists
	err := member.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM partner_payment_gateway WHERE payment_gateway_id = $1)`,
		paymentGatewayID).Scan(&paymentGatewayExists)

	if err != nil {
		return false, fmt.Errorf("error checking payout gateway existence: %s", err)
	}

	if !paymentGatewayExists {
		return false, nil
	}

	return true, nil
}

// CheckIfMemberSubscribedToFreePlan checks if a member has already subscribed to a specific free subscription.
// It takes the context, memberID, and subscriptionID as parameters.
// Returns a boolean indicating whether the member has subscribed to the free plan and an error if any.
func (member *MemberRepo) CheckIfMemberSubscribedToFreePlan(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (bool, error) {
	// Check if the subscription is free

	isFree, err := member.IsSubscriptionFree(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	// If the subscription is not free, return false
	if !isFree {
		return false, nil
	}

	// Check if the member is subscribed to this free subscription
	var hasSubscribed bool
	err = member.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM public.member_subscription
            WHERE member_id = $1
            AND subscription_id = $2
        );
    `, memberID, subscriptionID).Scan(&hasSubscribed)

	if err != nil {
		return false, fmt.Errorf("error checking member subscription: %s", err)
	}

	return hasSubscribed, nil
}

// HasSubscribedToOneTimePlan checks if a member has already subscribed to a specific one-time subscription plan.
// It takes the context, memberID, and subscriptionID as parameters.
// Returns a boolean indicating whether the member has subscribed to the one-time plan and an error if any.
func (member *MemberRepo) HasSubscribedToOneTimePlan(ctx *gin.Context, memberID uuid.UUID, subscriptionID string) (bool, error) {
	// Check if it's a one-time subscription
	isOneTime, err := member.IsOneTimeSubscription(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	// If it's not a one-time subscription, return false
	if !isOneTime {
		return false, nil
	}

	// Check if the member is subscribed to this one-time subscription
	var hasSubscribed bool
	err = member.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM public.member_subscription
            WHERE member_id = $1
            AND subscription_id = $2
        );
    `, memberID, subscriptionID).Scan(&hasSubscribed)

	if err != nil {
		return false, fmt.Errorf("error checking member subscription: %s", err)
	}

	return hasSubscribed, nil
}

// IsOneTimeSubscription checks if the provided subscription ID corresponds to a one-time subscription.
// Returns a boolean indicating if it's a one-time subscription and an error if any.
func (member *MemberRepo) IsOneTimeSubscription(ctx context.Context, subscriptionID string) (bool, error) {
	var isOneTime bool
	err := member.db.QueryRowContext(ctx, `
        SELECT is_one_time_subscription
        FROM public.subscription_plan
        WHERE id = $1;
    `, subscriptionID).Scan(&isOneTime)

	if err != nil {
		return false, fmt.Errorf("error checking if it's a one-time subscription: %s", err)
	}

	return isOneTime, nil
}

// IsSubscriptionFree checks if the provided subscription ID corresponds to a free subscription.
// Returns a boolean indicating if it's a free subscription and an error if any.
func (member *MemberRepo) IsSubscriptionFree(ctx context.Context, subscriptionID string) (bool, error) {
	var isFree bool
	err := member.db.QueryRowContext(ctx, `
        SELECT is_free_subscription
        FROM subscription_plan
        WHERE id = $1
    `, subscriptionID).Scan(&isFree)

	if err != nil {
		return false, fmt.Errorf("error checking if subscription is free: %s", err)
	}

	return isFree, nil
}

// Function to check if a member is subscribed to a specified plan and if the plan is free.
func (member *MemberRepo) IsMemberSubscribedToFreePlan(ctx *gin.Context, memberID uuid.UUID, MemberSubscriptionID string) (bool, error) {
	// Query to check if the member is subscribed to the specified plan and get the corresponding subscription plan.
	checkMemberSubscriptionQuery := `
        SELECT EXISTS (
            SELECT 1
            FROM member_subscription ms
            INNER JOIN subscription_plan sp ON ms.subscription_id = sp.id
            WHERE ms.member_id = $1 AND ms.id = $2 AND sp.is_free_subscription = true
        ) AS subscribed_to_free_plan;
    `

	var subscribedToFreePlan bool

	// Execute the query and scan the result into the subscribedToFreePlan variable.
	err := member.db.QueryRowContext(ctx, checkMemberSubscriptionQuery, memberID, MemberSubscriptionID).Scan(&subscribedToFreePlan)

	if err != nil {
		return false, fmt.Errorf("error checking subscription status: %w", err)
	}

	return subscribedToFreePlan, nil
}

// Function to check if a member is subscribed to a specified plan.
func (member *MemberRepo) IsMemberSubscribedToPlan(ctx *gin.Context, memberID uuid.UUID, MemberSubscriptionID string) (bool, error) {
	// Query to check if the member is subscribed to the specified plan.
	checkMemberSubscriptionQuery := `
        SELECT EXISTS (SELECT 1 FROM member_subscription WHERE member_id = $1 AND id = $2) AS subscribed;
    `

	var subscribed bool

	// Execute the query and scan the result into the subscribed variable.
	err := member.db.QueryRowContext(ctx, checkMemberSubscriptionQuery, memberID, MemberSubscriptionID).Scan(&subscribed)

	if err != nil {
		return false, fmt.Errorf("error checking subscription status: %w", err)
	}

	return subscribed, nil
}

// HandleSubscriptionRenewal handles the renewal of a subscription for a member.
// Parameters:
//   - ctx (context.Context): The context for the database operation.
//   - memberID (uuid.UUID): The unique identifier of the member renewing the subscription.
//   - renewal (entities.SubscriptionRenewal): The data related to subscription renewal.
// Returns:
//   - (error): An error if the renewal process fails.

// HandleSubscriptionRenewal handles the renewal of a subscription for a member.
func (member *MemberRepo) HandleSubscriptionRenewal(ctx context.Context, memberID uuid.UUID, checkoutData entities.SubscriptionRenewal) error {
	tx, err := member.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	SubscriptionID, err := member.GetSubscriptionIDByMemberSubscriptionID(ctx, checkoutData.MemberSubscriptionID)
	if err != nil {
		return err
	}

	// Fetch grace period details.
	inGracePeriod, _, graceEnd, _, _, graceErr := member.IsSubscriptionInGracePeriod(ctx, memberID, checkoutData.MemberSubscriptionID)
	if graceErr != nil {
		return graceErr
	}

	// If in grace period, use the grace end as the new expiration date.
	var expirationDate time.Time
	if inGracePeriod {
		expirationDate = graceEnd
	} else {
		// If not in grace period, calculate the new expiration date based on subscription duration.
		var subscriptionDurationValue int
		query := `
			SELECT sd.value
			FROM public.subscription_duration AS sd
			WHERE sd.id = (
				SELECT sp.subscription_duration_id
				FROM public.subscription_plan AS sp
				WHERE sp.id = $1
			);`

		err = tx.QueryRowContext(ctx, query, SubscriptionID).Scan(&subscriptionDurationValue)
		if err != nil {
			return err
		}

		// Calculate the new expiration date.
		expirationDate = time.Now().Add(time.Duration(subscriptionDurationValue) * 7 * 24 * time.Hour)
	}

	// Insert today's date into the 'renewed_on' field.
	renewedOnQuery := `
		UPDATE member_subscription
		SET expiration_date = $1,
		    renewed_on = current_date
		WHERE id = $2;
	`

	// Execute the update query to set the new expiration date and update 'renewed_on'.
	_, err = tx.ExecContext(ctx, renewedOnQuery, expirationDate, checkoutData.MemberSubscriptionID)
	if err != nil {
		return err
	}

	return nil
}

// IsSubscriptionAboutToExpire checks if the subscription is about to expire.
func (member *MemberRepo) IsSubscriptionAboutInWarning(ctx context.Context, memberSubscriptionID string) (bool, string, error) {
	var subscriptionID uuid.UUID
	var subscriptionDurationID int64

	var createdOn, expirationDate time.Time

	// Query to fetch subscription_id, created_on, and expiration_date from member_subscription table.
	fetchSubscriptionInfoQuery := `
        SELECT subscription_id, created_on, expiration_date
        FROM public.member_subscription
        WHERE id = $1;
    `

	err := member.db.QueryRowContext(ctx, fetchSubscriptionInfoQuery, memberSubscriptionID).Scan(&subscriptionID, &createdOn, &expirationDate)

	if err != nil {
		return false, "", err
	}

	// Query to fetch subscription_duration_id from subscription_plan table.
	fetchSubscriptionDurationIDQuery := `
        SELECT subscription_duration_id
        FROM public.subscription_plan
        WHERE id = $1;
    `

	err = member.db.QueryRowContext(ctx, fetchSubscriptionDurationIDQuery, subscriptionID).Scan(&subscriptionDurationID)
	if err != nil {
		return false, "", err
	}

	// Query to fetch subscription duration value from subscription_duration table.
	fetchSubscriptionDurationValueQuery := `
        SELECT value
        FROM public.subscription_duration
        WHERE id = $1;
    `

	var subscriptionDurationValue int
	err = member.db.QueryRowContext(ctx, fetchSubscriptionDurationValueQuery, subscriptionDurationID).Scan(&subscriptionDurationValue)
	if err != nil {
		return false, "", err
	}

	// Calculate expiry dates based on the created_on date and subscription duration.
	subscriptionDuration := time.Duration(subscriptionDurationValue) * 24 * time.Hour // Assuming subscription duration is in days.
	firstExpiryDate := expirationDate.Add(-subscriptionDuration)
	thirdExpiryDate := expirationDate.Add(-3 * subscriptionDuration)

	// Get the current date.
	currentDate := time.Now()

	// Check if the current date is between the first and third expiry dates.
	if currentDate.After(firstExpiryDate) && currentDate.Before(thirdExpiryDate) {
		message := "Your subscription is warning peroid and about to expire. Please renew by making a payment."
		return true, message, nil
	}

	return false, "", nil
}

// IsSubscriptionInGracePeriod checks if the subscription is in the grace period and can be renewed.
// this checks grace duration , if currently plan is in grace peroid and also when does the grace start and when  it ends
func (member *MemberRepo) IsSubscriptionInGracePeriod(ctx context.Context, memberID uuid.UUID, memberSubscriptionID string) (bool, time.Time, time.Time, time.Duration, bool, error) {
	var canRenewableWithin int
	var expirationDate time.Time

	// Query to fetch can_renewable_within and expiration_date from member_subscription and subscription_plan tables.
	checkGracePeriodQuery := `
        SELECT sp.can_renewable_within, ms.expiration_date
        FROM member_subscription ms
        INNER JOIN subscription_plan sp ON ms.subscription_id = sp.id
        WHERE ms.id = $1 AND ms.member_id = $2;
    `

	err := member.db.QueryRowContext(ctx, checkGracePeriodQuery, memberSubscriptionID, memberID).Scan(&canRenewableWithin, &expirationDate)
	if err != nil {
		return false, time.Time{}, time.Time{}, 0, false, err
	}

	graceStart := expirationDate
	graceEnd := expirationDate.Add(time.Duration(canRenewableWithin) * 7 * 24 * time.Hour)

	// Check if the current date is within the grace period.
	if time.Now().After(graceEnd) {

		graceDuration := graceEnd.Sub(graceStart)

		isWithinGraceDuration := time.Now().Before(graceStart.Add(graceDuration))
		return true, graceStart, graceEnd, graceDuration, isWithinGraceDuration, nil
	}

	return false, graceStart, graceEnd, 0, false, nil
}

// CheckCancellationEnabled checks if a subscription plan allows cancellations.
// It returns true if cancellation is enabled, otherwise false.
func (member *MemberRepo) CheckCancellationEnabled(ctx *gin.Context, MemberSubscriptionID string) (bool, error) {
	// SQL query to fetch the SubscriptionID column for the given MemberSubscriptionID.
	query := `SELECT subscription_id FROM public.member_subscription WHERE id = $1`
	var SubscriptionID string
	err := member.db.QueryRowContext(ctx, query, MemberSubscriptionID).Scan(&SubscriptionID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch subscription ID: %v", err)
	}

	// SQL query to fetch the is_cancellation_enabled column for the retrieved SubscriptionID.
	query = `SELECT is_cancellation_enabled FROM public.subscription_plan WHERE id = $1`
	var isCancellationEnabled bool

	// Execute the SQL query and scan the result into the isCancellationEnabled variable.
	err = member.db.QueryRowContext(ctx, query, SubscriptionID).Scan(&isCancellationEnabled)
	if err != nil {
		return false, fmt.Errorf("failed to fetch cancellation status: %v", err)
	}

	return isCancellationEnabled, nil
}

// HandleSubscriptionCancellation handles the cancellation process for a subscription.
// It performs the following steps:
//  1. Checks if the member has subscribed to the plan.
//  2. Checks if the associated subscription plan is active/on hold.
//  3. Verifies if the subscription is cancellation enabled.
//
// Parameters:
//   - ctx (context.Context): The context for the database operations.
//   - memberID (uuid.UUID): The unique identifier of the member.
//   - checkoutData (entities.CancelSubscription): The cancellation data, including SubscriptionID.
//
// Returns:
//   - error: An error if any database operation fails.

func (member *MemberRepo) HandleSubscriptionCancellation(ctx context.Context, memberID uuid.UUID, checkoutData entities.CancelSubscription) error {

	// If all conditions are met, update the member_subscription_status to "cancelled".
	updateStatusQuery := `
   			 UPDATE member_subscription
    		 SET member_subscription_status_id = (SELECT id FROM member_subscription_status WHERE name = 'cancelled') 
			 WHERE id = $1 ;
`

	_, err := member.db.ExecContext(ctx, updateStatusQuery, checkoutData.MemberSubscriptionID)

	if err != nil {
		return err
	}

	return nil
}

func (member *MemberRepo) UpdatePrimaryBillingAddressToFalseAndRandom(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID) error {
	// Start a transaction
	tx, err := member.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Step 1: Update the primary status of the current primary billing address to false
	_, err = tx.ExecContext(ctx, `UPDATE member_billing_address SET is_primary_billing = false WHERE id = $1`, memberBillingID)
	if err != nil {
		return err
	}

	// Step 2: Fetch all billing addresses for the given member except memberBillingID
	rows, err := tx.QueryContext(ctx, `SELECT id FROM member_billing_address WHERE member_id = $1 AND id != $2`, memberID, memberBillingID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var shuffledAddresses []uuid.UUID

	for rows.Next() {
		var addressID uuid.UUID
		if err := rows.Scan(&addressID); err != nil {
			return err
		}
		shuffledAddresses = append(shuffledAddresses, addressID)
	}

	// Check if there are any addresses available
	if len(shuffledAddresses) == 0 {
		return errors.New("no other billing addresses found for the given member")
	}

	// Shuffle the addresses
	rand.Shuffle(len(shuffledAddresses), func(i, j int) {
		shuffledAddresses[i], shuffledAddresses[j] = shuffledAddresses[j], shuffledAddresses[i]
	})

	// Get the first address after shuffling (randomly selected)
	randomAddressID := shuffledAddresses[0]

	// Step 3: Update the primary status of the randomly selected billing address to true
	_, err = tx.ExecContext(ctx, `UPDATE member_billing_address SET is_primary_billing = true WHERE id = $1`, randomAddressID)
	if err != nil {
		return err
	}

	return nil
}

// SubscriptionProductSwitch function to switch products between subscriptions
func (member *MemberRepo) SubscriptionProductSwitch(ctx context.Context, memberID uuid.UUID, data entities.SwitchSubscriptions, validationErrors *map[string][]string) (map[string][]string, error) {

	var (
		productCount, maxProductCount, newArtistCount, newTrackCount, newMaxTracksPerProduct, newMaxArtistsPerProduct,
		currentArtistCount, currentTrackCount, currentMaxTracksPerProduct, currentMaxArtistsPerProduct int
		newSubscriptionStatus        string
		newIsActive, currentIsActive bool
	)

	// Check if given memberID exists
	memberExists, err := member.IsMemberExist(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if !memberExists {
		utils.AppendValuesToMap(*validationErrors, consts.MemberrID, consts.NotFound)
	}

	// Start a transaction.
	tx, err := member.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	var (
		currentSubscriptionID, productSubscriptionID, newSubscriptionID, getID uuid.UUID
	)

	// Check the product is already in new given subscription
	err = tx.QueryRowContext(ctx, `
	SELECT p.member_subscription_id
	FROM product p
	JOIN member_subscription ms ON p.member_subscription_id = ms.id
	WHERE p.id = $1 AND ms.member_id = $2
		`, data.ProductReferenceID, memberID).Scan(&newSubscriptionID)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Invalid new subscription id err=%s", err.Error())
	}

	err = tx.QueryRowContext(ctx, `
	SELECT subscription_id FROM member_subscription WHERE id=$1
	`, newSubscriptionID).Scan(&getID)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("product is already in given new subscription err=%s", err.Error())
	}

	if getID.String() == data.NewSubscriptionID {
		utils.AppendValuesToMap(*validationErrors, consts.NewSubscriptionID, consts.AlreadyExist)
		return *validationErrors, nil
	}

	// Check the product is under the correct given subscription plan
	err = tx.QueryRowContext(ctx, `
	SELECT p.member_subscription_id
	FROM product p
	JOIN member_subscription ms ON p.member_subscription_id = ms.id
	WHERE p.id = $1 AND ms.member_id = $2
		`, data.ProductReferenceID, memberID).Scan(&currentSubscriptionID)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Invalid current subscription id err=%s", err.Error())
	}

	err = tx.QueryRowContext(ctx, `
	SELECT subscription_id FROM member_subscription WHERE id=$1
	`, currentSubscriptionID).Scan(&productSubscriptionID)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf("Invalid product id err=%s", err.Error())
	}

	if productSubscriptionID.String() != data.CurrentSubscriptionID {
		utils.AppendValuesToMap(*validationErrors, consts.ProductReferenceID, consts.Invalid)
	}

	// Step 1: Check if the product_count is full for the new subscription plan.
	err = tx.QueryRowContext(ctx, `
		SELECT max(product_count)as maxProductCount FROM subscription_plan WHERE id = $1
		`, data.NewSubscriptionID).Scan(&maxProductCount)

	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, `
		SELECT count(id)as productCount FROM product WHERE member_subscription_id = $1
		`, data.NewSubscriptionID).Scan(&productCount)

	if err != nil {
		return nil, err
	}
	if productCount == maxProductCount {
		utils.AppendValuesToMap(*validationErrors, consts.NewSubscriptionID, consts.LimitExceeds)
	}

	// Step 2: Check if the new subscription plan is active.
	err = tx.QueryRowContext(ctx, `
		SELECT name
		FROM member_subscription_status
		WHERE id = (
			SELECT member_subscription_status_id
			FROM member_subscription
			WHERE subscription_id = $1
		)
	`, data.NewSubscriptionID).Scan(&newSubscriptionStatus)

	if err != nil {
		return nil, err
	}
	if newSubscriptionStatus != "active" {
		utils.AppendValuesToMap(*validationErrors, consts.NewSubscriptionID, consts.Inactive)
	}

	// Step 3: Compare fields between the new and current subscription plans.
	err = tx.QueryRowContext(ctx, `
		SELECT sp2.artist_count, sp2.track_count, sp2.max_tracks_per_product, sp2.max_artists_per_product, sp2.is_active
		FROM subscription_plan AS sp2
		WHERE sp2.id = $1
		`, data.CurrentSubscriptionID).Scan(&currentArtistCount, &currentTrackCount, &currentMaxTracksPerProduct, &currentMaxArtistsPerProduct, &currentIsActive)

	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, `
        SELECT sp1.artist_count, sp1.track_count, sp1.max_tracks_per_product, sp1.max_artists_per_product, sp1.is_active
        FROM subscription_plan AS sp1
        WHERE sp1.id = $1
		`, data.NewSubscriptionID).Scan(&newArtistCount, &newTrackCount, &newMaxTracksPerProduct, &newMaxArtistsPerProduct, &newIsActive)

	if err != nil {
		return nil, err
	}
	if !newIsActive {
		utils.AppendValuesToMap(*validationErrors, consts.NewSubscriptionID, consts.NotActive)
	}
	if len(*validationErrors) != 0 {
		return nil, nil
	}

	if (newArtistCount >= currentArtistCount) && (newTrackCount >= currentTrackCount) &&
		(newMaxTracksPerProduct >= currentMaxTracksPerProduct) && (newMaxArtistsPerProduct >= currentMaxArtistsPerProduct) {

		// Step 4: Update the member_subscription_id in the product table.
		_, err = tx.ExecContext(ctx, `
		UPDATE product
		SET member_subscription_id = (SELECT id FROM member_subscription WHERE subscription_id = $2)
		WHERE id = $1
		`, data.ProductReferenceID, data.NewSubscriptionID)

		if err != nil {
			return nil, err
		}

		// Commit the transaction if everything succeeds.
		if err = tx.Commit(); err != nil {
			return nil, err
		}

		return nil, nil
	}

	utils.AppendValuesToMap(*validationErrors, consts.NewSubscriptionID, consts.NotAllowed)
	return nil, fmt.Errorf("product criteria do not fit the new plan")
}

// ViewAllSubscriptions returns list of all member subscriptions
func (member *MemberRepo) ViewAllSubscriptions(ctx context.Context, memberID uuid.UUID, reqParam entities.ReqParams, validationErrors *map[string][]string) ([]entities.ListAllSubscriptions, error) {
	// Check if given memberID exists
	memberExists, err := member.IsMemberExist(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if !memberExists {
		utils.AppendValuesToMap(*validationErrors, consts.MemberrID, consts.NotFound)
	}

	if len(*validationErrors) != 0 {
		return nil, nil
	}

	var (
		subscriptions []entities.ListAllSubscriptions
		data          entities.ListAllSubscriptions
	)

	query := `
        SELECT
            ms.id AS member_subscription_id,
            ms.custom_name,
            mss.name AS subscription_status,
            ms.expiration_date,
            COUNT(p.id) AS product_count,
            COUNT(pt.track_id) AS track_count,
            COUNT(pa.artist_id) AS artist_count,
            sp.id AS subscription_plan_id,
            sp.name AS subscription_plan_name,
            sp.sku AS subscription_plan_sku,
            sd.name AS subscription_duration,
            sp.product_count,
            sp.track_count,
            sp.artist_count
        FROM
            subscription_plan AS sp
        INNER JOIN
            member_subscription AS ms ON sp.id = ms.subscription_id
        LEFT JOIN
            member_subscription_status AS mss ON ms.member_subscription_status_id = mss.id
        LEFT JOIN
            product AS p ON ms.id = p.member_subscription_id AND p.member_id = ms.member_id 

        LEFT JOIN LATERAL (
            SELECT
                pt.track_id
            FROM
                product_track AS pt
            WHERE
                pt.product_id = p.id
        ) AS pt ON true

        LEFT JOIN LATERAL (
            SELECT
                pa.artist_id
            FROM
                product_artist AS pa
            WHERE
                pa.product_id = p.id
        ) AS pa ON true

        LEFT JOIN
            subscription_duration AS sd ON sp.subscription_duration_id = sd.id
        WHERE
            ms.member_id = $1
    `

	if reqParam.Status != "" {
		if reqParam.Status == consts.Active {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.Active)

		} else if reqParam.Status == consts.Expired {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.Expired)

		} else if reqParam.Status == consts.Cancelled {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.Cancelled)

		} else if reqParam.Status == consts.OnHold {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.OnHold)

		} else if reqParam.Status == consts.Processing {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.Processing)

		} else if reqParam.Status == consts.PaymentFailed {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.PaymentFailed)

		} else if reqParam.Status == consts.InGrace {
			query = fmt.Sprintf("%s AND mss.name = '%s'", query, consts.InGrace)

		}
	}

	if reqParam.Search != "" {
		query = fmt.Sprintf("%s AND (ms.custom_name LIKE '%s')", query, reqParam.Search)
	}

	query += `GROUP BY
            ms.id,
            ms.custom_name,
            mss.name,
            ms.expiration_date,
            sp.id,
            sp.name,
            sp.sku,
            sd.name,
            sp.product_count,
            sp.track_count,
            sp.artist_count
        `

	if reqParam.Sort != "" {
		query = fmt.Sprintf("%s ORDER BY ms.%s ASC", query, reqParam.Sort)
	}

	// Checking page and limit not equal to zero and calculate offset
	if reqParam.Page != 0 && reqParam.Limit != 0 {
		offset := (reqParam.Page - 1) * reqParam.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", reqParam.Limit, offset)
	}

	rows, err := member.db.QueryContext(ctx, query, memberID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&data.ID,
			&data.CustomName,
			&data.Status,
			&data.ExpirationDate,
			&data.ProductsAdded,
			&data.TracksAdded,
			&data.ArtistsAdded,
			&data.SubscriptionDetails.SubscriptionID,
			&data.SubscriptionDetails.Name,
			&data.SubscriptionDetails.SKU,
			&data.SubscriptionDetails.Duration,
			&data.SubscriptionDetails.MaximumProducts,
			&data.SubscriptionDetails.MaximumTracks,
			&data.SubscriptionDetails.MaximumArtists,
		)

		if err != nil {
			return nil, err
		}

		// Marshal CustomName using the custom function
		customNameJSON, err := utilities.MarshalNullableString(data.CustomName)

		if err != nil {
			return nil, err
		}

		data.CustomNameJSON = customNameJSON
		// append the newly created instance to subscriptions
		subscriptions = append(subscriptions, data)
		fmt.Println("-----------slice", subscriptions)
	}

	return subscriptions, nil
}

// GetSubscriptionRecordCount function is used to calculate and return total count of subscriptions
func (member *MemberRepo) GetSubscriptionRecordCount(ctx context.Context, memberID uuid.UUID) (int64, error) {

	var totalCount int64

	query := `
		SELECT COUNT(subscription_id) AS total_count
		FROM member_subscription 
		WHERE member_subscription.member_id = $1
	`
	row := member.db.QueryRowContext(ctx, query, memberID)

	if err := row.Scan(&totalCount); err != nil {
		logger.Log().WithContext(ctx).Errorf("GetSubscriptionRecordCount failed, QueryRowContext failed, err=%s", err.Error())
		return 0, err
	}

	return totalCount, nil
}

// IsMemberExists checks if a member exists in the database.
// Parameters:
//	- memberID: The UUID of the member to check for existence.
//	- ctx: The context for the database operation.

// Returns:
// - bool: A boolean indicating whether the member exists.
// - error: An error, if any, during the database operation.

func (member *MemberRepo) IsMemberExist(ctx context.Context, memberID uuid.UUID) (bool, error) {

	var exists int
	isMemberExistsQ := `SELECT 1 FROM member WHERE id = $1`
	row := member.db.QueryRowContext(ctx, isMemberExistsQ, memberID)
	err := row.Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (member *MemberRepo) UpdateRandomBillingAddressToPrimary(ctx *gin.Context, memberID, memberBillingID uuid.UUID) error {
	// Begin a transaction
	tx, err := member.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			// If there was a panic, rollback the transaction and re-panic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// If there was an error, rollback the transaction
			tx.Rollback()
		} else {
			// Commit the transaction if everything went well
			err = tx.Commit()
		}
	}()

	// Fetch all billing addresses for the given member except memberBillingID
	rows, err := tx.QueryContext(ctx, `SELECT id FROM member_billing_address WHERE member_id = $1 AND id != $2`, memberID, memberBillingID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var shuffledAddresses []uuid.UUID

	for rows.Next() {
		var addressID uuid.UUID
		if err := rows.Scan(&addressID); err != nil {
			return err
		}
		shuffledAddresses = append(shuffledAddresses, addressID)
	}

	// Check if there are any addresses available
	if len(shuffledAddresses) == 0 {
		return errors.New("no other billing addresses found for the given member")
	}

	// Shuffle the addresses
	rand.Shuffle(len(shuffledAddresses), func(i, j int) {
		shuffledAddresses[i], shuffledAddresses[j] = shuffledAddresses[j], shuffledAddresses[i]
	})

	// Get the first address after shuffling (randomly selected)
	randomAddressID := shuffledAddresses[0]

	// Update the primary status of the randomly selected billing address to true
	_, err = tx.ExecContext(ctx, `UPDATE member_billing_address SET is_primary_billing = true WHERE id = $1`, randomAddressID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBillingAddress deletes a billing address entry based on memberID and memberBillingID.
func (member *MemberRepo) DeleteBillingAddress(ctx *gin.Context, memberID uuid.UUID, memberBillingID uuid.UUID) error {
	// Execute the DELETE query
	query := "DELETE FROM public.member_billing_address WHERE member_id = $1 AND id = $2"
	result, err := member.db.ExecContext(ctx, query, memberID, memberBillingID)
	if err != nil {
		return err
	}

	// Check the affected rows to determine if the entry was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// The entry was not found
		return fmt.Errorf("billing address not found for memberID: %s and billingID: %s", memberID, memberBillingID)
	}

	return nil
}

// GetPaymentDetailsByPartnerAndGateway retrieves payment details based on partner ID and payment gateway ID.
// It returns the payment details and any error encountered during the retrieval.

func (member *MemberRepo) GetPaymentDetailsByPartnerAndGateway(ctx context.Context, partnerID string, paymentGatewayID int) (string, error) {
	var paymentDetails sql.NullString

	// Query to fetch payment details based on partner ID and payment gateway ID.
	err := member.db.QueryRowContext(ctx, `
        SELECT payment_details
        FROM public.partner_payment_gateway
        WHERE partner_id = $1
          AND payment_gateway_id = $2
    `, partnerID, paymentGatewayID).Scan(&paymentDetails)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no rows were returned (partner/payment gateway not found)
			return "", fmt.Errorf("no payment details found for partnerID %s and paymentGatewayID %d", partnerID, paymentGatewayID)
		}
		return "", err
	}

	// Check if paymentDetails is NULL
	if !paymentDetails.Valid {
		return "", fmt.Errorf("payment details are NULL for partnerID %s and paymentGatewayID %d", partnerID, paymentGatewayID)
	}

	return paymentDetails.String, nil
}

// IsPartnerIdCorrespondsToGateway checks if the payment gateway passed corresponds to given partner.
func (member *MemberRepo) IsPartnerIdCorrespondsToGateway(ctx context.Context, partnerID string, paymentGatewayID int) (bool, error) {
	var exists bool

	// Query to check if the partner ID corresponds to the given payment gateway ID in the partner_payin_gateway table.
	err := member.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM public.partner_payment_gateway
			WHERE partner_id = $1
			  AND payment_gateway_id = $2
		) AS exists
	`, partnerID, paymentGatewayID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// DecryptPaymentData decrypts the given encrypted data using a decryption key
// specified in the MemberRepo configuration.
//
// The function returns the decrypted string and any error encountered during the decryption process.
func (member *MemberRepo) DecryptPaymentData(ctx *gin.Context, data string) (string, error) {
	if member.Cfg == nil {
		return "", errors.New("configuration is nil")
	}

	decryptionKey := member.Cfg.DecryptionKey
	if decryptionKey == "" {
		return "", errors.New("decryption key is empty")
	}

	key := []byte(decryptionKey)
	if len(key) == 0 {
		return "", errors.New("decryption key is empty after conversion to byte slice")
	}

	decryptedString, err := crypto.Decrypt(data, key)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %v", err)
	}

	return decryptedString, nil
}

// HasProductsReleaseEndDateGreaterThanToday checks if there are products associated with the given member subscription
// (identified by memberSubscriptionID) having a release_end_date greater than today.
// It returns true if such products exist, false if none are found, and an error for any database-related issues.
func (member *MemberRepo) HasProductsReleaseEndDateGreaterThanToday(ctx *gin.Context, memberSubscriptionID string) (bool, error) {
	var exists int
	// Constructing SQL query to check for products with release_end_date greater than today
	query := `SELECT 1 FROM public.product WHERE member_subscription_id = $1 AND release_end_date > CURRENT_DATE LIMIT 1`
	// Executing the query and scanning the result into the "exists" variable
	row := member.db.QueryRowContext(ctx, query, memberSubscriptionID)
	err := row.Scan(&exists)

	// Handling the case where no products with release_end_date greater than today are found
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No products with release_end_date greater than today found
		}
		// Returning an error if there is any issue other than no rows found
		return false, fmt.Errorf("error checking product release_end_date: %v", err)
	}

	return true, nil // There are products with release_end_date greater than today
}

// IsMemberRelatedToSubscription checks if a given member is related to the specified MemberSubscriptionID.
func (member *MemberRepo) IsMemberRelatedToSubscription(ctx *gin.Context, memberID uuid.UUID, memberSubscriptionID string) (bool, error) {
	var isRelated bool

	// Query the database to check if the member is related to the specified MemberSubscriptionID.
	err := member.db.QueryRowContext(ctx, ` SELECT 1  FROM member_subscription
            WHERE member_id = $1 AND id = $2
        
    `, memberID, memberSubscriptionID).Scan(&isRelated)

	if err != nil {
		// If there was an error querying the database, return an error.
		return false, fmt.Errorf("error checking if member is related to subscription: %s", err)
	}

	// Return the result indicating whether the member is related to the subscription or not.
	return isRelated, nil
}

// GetSubscriptionIDByMemberSubscriptionID retrieves the subscription ID for a given member subscription ID.
func (member *MemberRepo) GetSubscriptionIDByMemberSubscriptionID(ctx context.Context, memberSubscriptionID string) (uuid.UUID, error) {
	query := `
		SELECT subscription_id
		FROM public.member_subscription
		WHERE id = $1
	`

	var subscriptionID uuid.UUID
	err := member.db.QueryRowContext(ctx, query, memberSubscriptionID).Scan(&subscriptionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("member subscription with ID %s not found", memberSubscriptionID)
		}
		return uuid.Nil, err
	}

	return subscriptionID, nil
}

// CheckMemberPartner checks if the provided partnerID matches the partner_id associated with the given memberID
func (member *MemberRepo) CheckMemberPartner(ctx *gin.Context, memberID uuid.UUID, partnerIDStr string) (bool, error) {
	var storedPartnerID string

	// Query the database to get the partner_id associated with the memberID
	query := `SELECT partner_id FROM public.member WHERE id = $1`

	err := member.db.QueryRowContext(ctx, query, memberID).Scan(&storedPartnerID)
	if err != nil {
		// Handle errors
		if err == sql.ErrNoRows {
			// Handle the case where the memberID does not exist
			return false, fmt.Errorf("member not found with ID: %v", memberID)
		}
		return false, fmt.Errorf("error querying database: %v", err)
	}

	// Compare the stored partnerID with the provided partnerIDStr
	if storedPartnerID == partnerIDStr {
		// Partner authentication succeeded
		return true, nil
	}

	// Partner authentication failed
	return false, nil
}

// CheckSubscriptionExistenceAndStatusForRenewal checks if the provided memberSubscriptionID exists,
// if the associated subscription plan exists, and if the plan is currently active.
func (member *MemberRepo) CheckSubscriptionExistenceAndStatusForRenewal(ctx *gin.Context, memberSubscriptionID string) (bool, bool, error) {
	var subscriptionExists, isPlanActive bool

	// Check if the member subscription exists
	query := `SELECT EXISTS(SELECT 1 FROM public.member_subscription WHERE id = $1)`

	err := member.db.QueryRowContext(ctx, query, memberSubscriptionID).Scan(&subscriptionExists)
	if err != nil {
		return false, false, fmt.Errorf("error checking member subscription existence: %v", err)
	}

	// If the member subscription does not exist, no need to proceed with further checks
	if !subscriptionExists {
		return false, false, nil
	}

	// Check if the associated subscription plan exists and if it is currently active
	planQuery := `
		SELECT EXISTS(
			SELECT 1
			FROM public.subscription_plan sp
			JOIN public.member_subscription ms ON sp.id = ms.subscription_id
			WHERE ms.id = $1
			AND sp.is_active = true
		)`
	err = member.db.QueryRowContext(ctx, planQuery, memberSubscriptionID).Scan(&isPlanActive)
	if err != nil {
		return false, false, fmt.Errorf("error checking subscription plan existence and status: %v", err)
	}

	return true, isPlanActive, nil
}

// DeleteMember deletes a member
func (member *MemberRepo) DeleteMember(ctx *gin.Context, MemberID uuid.UUID) error {
	// Construct the SQL query to update the is_deleted field
	// Get the current date
	currentDate := time.Now().UTC().Format("2006-01-02")

	query := `
        UPDATE public.member
        SET is_deleted = true,is_active = false, deleted_on = $2
        WHERE id = $1 AND is_deleted = false
    `

	// Execute the SQL query
	_, err := member.db.Exec(query, MemberID, currentDate)
	if err != nil {

		return err
	}
	return nil
}

// IsActive Checks if the member is currently active or not.
func (member *MemberRepo) IsActive(ctx *gin.Context, MemberID uuid.UUID) (bool, error) {
	// Construct the SQL query to update the is_deleted field

	query := `
    SELECT 1 FROM public.member WHERE is_active = true AND id = $1
`

	// Execute the SQL query
	_, err := member.db.Exec(query, MemberID)
	if err != nil {

		return false, err
	}
	return true, nil
}

// IsDeleted Checks if the member is currently active or not.
func (member *MemberRepo) IsDeleted(ctx *gin.Context, MemberID uuid.UUID) (bool, error) {
	// Construct the SQL query to update the is_deleted field

	query := `
    SELECT 1 FROM public.member WHERE is_deleted = true AND id = $1
`

	// Execute the SQL query
	_, err := member.db.Exec(query, MemberID)
	if err != nil {

		return false, err
	}
	return true, nil
}

// AddMemberStores adds stores related to a member
func (member *MemberRepo) AddMemberStoresById(ctx *gin.Context, memberID uuid.UUID, stores []uuid.UUID) error {

	// Insert the member_id and store_id into the member_store table
	for _, storeID := range stores {
		// First, query the partner_store table to get values
		fetchQuery := `
			SELECT is_store, is_active
			FROM public.partner_store
			WHERE store_id = $1
		`

		var isStore bool
		var isActive bool
		var customName string

		err := member.db.QueryRow(fetchQuery, storeID).Scan(&isStore, &isActive)
		if err != nil {
			fmt.Println("Error fetching values from partner_store:", err)
			return err
		}

		// Now, you can use these fetched values to insert into the member_store table
		insertQuery := `
			INSERT INTO public.member_store (member_id, store_id, is_store, is_active, custom_store_name)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = member.db.Exec(insertQuery, memberID, storeID, isStore, isActive, customName)
		if err != nil {
			fmt.Println("Error inserting into member_store:", err)
			return err
		}
	}

	return nil
}

// GetStoreIDsByPartnerID retrieves all store IDs related to the provided partner ID
func (member *MemberRepo) GetStoreIDsByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]uuid.UUID, error) {
	query := `
        SELECT store_id
        FROM public.partner_store
        WHERE partner_id = $1
    `

	rows, err := member.db.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storeIDs []uuid.UUID
	for rows.Next() {
		var storeID uuid.UUID
		if err := rows.Scan(&storeID); err != nil {
			return nil, err
		}
		storeIDs = append(storeIDs, storeID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return storeIDs, nil
}

// GetPartnerIDByMemberID retrieves the partner ID from the member table based on the provided memberID
func (member *MemberRepo) GetPartnerIDByMemberID(ctx context.Context, memberID uuid.UUID) (uuid.UUID, error) {
	query := `
        SELECT partner_id
        FROM public.member
        WHERE id = $1
        LIMIT 1
    `

	var partnerID uuid.UUID
	err := member.db.QueryRowContext(ctx, query, memberID).Scan(&partnerID)
	if err != nil {
		return uuid.Nil, err
	}

	return partnerID, nil
}

// GetStoreIDByCustomName retrieves the store ID from the partner_store table based on the provided custom name
func (member *MemberRepo) GetStoreIDByCustomName(ctx context.Context, customName []string) ([]uuid.UUID, error) {
	query := `
		SELECT store_id
		FROM public.partner_store
		WHERE custom_name = $1
		LIMIT 1
	`

	rows, err := member.db.QueryContext(ctx, query, customName[0])
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer rows.Close()

	var storeIDs []uuid.UUID
	for rows.Next() {
		var storeID uuid.UUID
		if err := rows.Scan(&storeID); err != nil {
			return nil, err
		}
		storeIDs = append(storeIDs, storeID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return storeIDs, nil
}

// CheckCustomNameExists checks if any custom name exists in the public.partner_store table based on a list of store names
// CheckStoreNameExistsAndReturnIDs checks if store names exist and returns the corresponding IDs.
func (member *MemberRepo) CheckStoreNameExistsAndReturnIDs(ctx context.Context, storeNames []string) (bool, []uuid.UUID, error) {
	// SQL query template
	queryTemplate := `
        SELECT id
        FROM public.store
        WHERE name = $1`

	var storeIDs []uuid.UUID
	allExist := true

	// Iterate through each store name
	for _, storeName := range storeNames {
		var storeID uuid.UUID

		// Execute the query with the parameterized store name
		err := member.db.QueryRowContext(ctx, queryTemplate, storeName).Scan(&storeID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Store name does not exist, set allExist to false and continue
				allExist = false
				continue
			}
			return false, nil, err
		}

		// Append the store ID to the result slice
		storeIDs = append(storeIDs, storeID)
	}

	return allExist, storeIDs, nil
}

// StorePartnerRelation checks if given store is related to members designated partner
func (member *MemberRepo) StorePartnerRelation(ctx context.Context, partnerID uuid.UUID, relatedStoreIDs []uuid.UUID) (map[string]bool, error) {
	existingStoreRelations := make(map[string]bool)

	for _, storeId := range relatedStoreIDs {
		query := `
            SELECT EXISTS (
                SELECT 1
                FROM public.partner_store
                WHERE partner_id = $1 AND store_id = $2
            )
        `
		var exists bool
		err := member.db.QueryRowContext(ctx, query, partnerID, storeId).Scan(&exists)
		if err != nil {
			return nil, err
		}

		existingStoreRelations[storeId.String()] = exists
	}

	return existingStoreRelations, nil
}

// CheckMemberStoreExists checks if the specified stores exist for a member
func (member *MemberRepo) CheckNonExistingMemberStores(ctx *gin.Context, memberID uuid.UUID, storeIDs []uuid.UUID) ([]uuid.UUID, error) {
	// Create a slice to store non-existing store IDs
	nonExistingStoreIDs := []uuid.UUID{}

	// Iterate through each store ID
	for _, storeID := range storeIDs {
		// Query the member_store table to check if the store ID already exists for the member
		fetchQuery := `
			SELECT COUNT(*) 
			FROM public.member_store
			WHERE member_id = $1 AND store_id = $2
		`

		var count int
		err := member.db.QueryRow(fetchQuery, memberID, storeID).Scan(&count)
		if err != nil {
			fmt.Println("Error fetching values from member_store:", err)
			return nil, err
		}

		// If the store ID does not exist for the member, add it to the nonExistingStoreIDs slice
		if count == 0 {
			nonExistingStoreIDs = append(nonExistingStoreIDs, storeID)
		}
	}

	// Return the list of store IDs that do not exist for the member
	return nonExistingStoreIDs, nil
}

// CheckPartnerStores checks if the given store IDs are present in the partner_store table for a specific partner.
func (member *MemberRepo) CheckPartnerStores(ctx context.Context, partnerID uuid.UUID, storeIDs []uuid.UUID) (bool, error) {
	// SQL query template
	queryTemplate := `
		SELECT EXISTS (
			SELECT 1
			FROM public.partner_store
			WHERE partner_id = $1 AND store_id = $2
		)`

	// Iterate through each store ID
	for _, storeID := range storeIDs {
		var partnerStoreExists bool

		// Execute the query with the parameterized partnerID and storeID
		err := member.db.QueryRowContext(ctx, queryTemplate, partnerID, storeID).Scan(&partnerStoreExists)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No matching partner store found for the current storeID
				return false, nil
			}
			return false, err
		}

		// If the partner store does not exist for any storeID, return false
		if !partnerStoreExists {
			return false, nil
		}
	}

	// If partner stores exist for all storeIDs, return true
	return true, nil
}

// CheckLanguageExist checks if given language exists in language list in database.
func (member *MemberRepo) CheckLanguageExist(ctx *gin.Context, Language string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM public.language WHERE code = $1)"

	var exists bool
	err := member.db.QueryRow(query, Language).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking language existence: %v", err)
	}

	return exists, nil
}

// CheckPartnerIDExists checks if a partner with the specified ID exists.
func (member *MemberRepo) CheckPartnerIDExists(ctx *gin.Context, partnerID string) (bool, error) {
	var exists bool

	// SQL query for checking if a partner_id exists in the public.partner table.
	checkPartnerIDExistsQ := `SELECT EXISTS (SELECT 1 FROM public.partner WHERE id = $1)`

	row := member.db.QueryRowContext(ctx, checkPartnerIDExistsQ, partnerID)

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
