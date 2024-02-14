package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	gt "github.com/bas24/googletranslatefree"

	"localization/internal/consts"
	"localization/internal/entities"
	"localization/internal/entities/db"
	"localization/utilities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrorCodesRepo is a repository for managing error codes and translations.
type ErrorCodesRepo struct {
	db         *mongo.Database
	Collection Collection
	Cfg        *entities.EnvConfig
}
type Collection struct {
	Endpoint *mongo.Collection
	Error    *mongo.Collection
}

// Document represents a document in the database.
type Document struct {
	ID string         `bson:"_id"`
	En map[string]any `bson:"en"`
}

// ErrorCodesRepoImply defines the interface for ErrorCodesRepo.
type ErrorCodesRepoImply interface {
	GetError(string, string, string, string, string) (any, error)
	AppendError(context.Context, string, string, string, db.ErrorData, string) error
	AddTranslation(context.Context, []string) (any, error)
	AddEndpoint(context.Context, entities.RequestData) error
	GetEndpointName(context.Context) (entities.ResponseData, error)
}

// NewErrorCodesRepo creates a new instance of ErrorCodesRepo.
func NewErrorCodesRepo(db *mongo.Database, Cfg *entities.EnvConfig) ErrorCodesRepoImply {
	return &ErrorCodesRepo{db: db,
		Cfg: Cfg,
		Collection: Collection{
			Endpoint: db.Collection(consts.CollectionEndpoint),
			Error:    db.Collection(consts.CollectionStatusCodes),
		},
	}
}

//GetError is used to get errors based on various params like
//params

// @errType 	- Type of error
// @endpoint 	- Endpoint names
// @lang 		- Language
// @field 		- Field is of the form fieldname:value (address:Valid)
// @method 		- Method of the endpoint (get,post..)
func (repo *ErrorCodesRepo) GetError(errType string, endpoint string, lang string, field string, method string) (any, error) {
	collection := repo.Collection.Error

	// Define a filter to match the document you want
	id, err := primitive.ObjectIDFromHex(repo.Cfg.DocID)
	if err != nil {
		log.Println(consts.MongoCollectionErr)
		return nil, err
	}

	filter := bson.M{"_id": id}

	// Query the document
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(consts.NotFoundErr)
			return nil, errors.New(consts.NoMatchErr)
		} else {
			return nil, err
		}
	}

	if result == nil {
		log.Println(consts.QueryErr)
		return nil, errors.New(consts.QueryResultErr)
	}

	enValue, enExists := result[lang].(bson.M)
	if !enExists {
		log.Println(consts.EnLangErr)
		errMsg := fmt.Sprintf("%s language not found or not of the expected language", lang)
		return nil, errors.New(errMsg)
	}

	// If both error type and endpoint are empty, generate a general error response.
	if errType == "" && endpoint == "" {
		responseJSON := utilities.ErrorResponseFormat(consts.AllErrMsg, "", enValue, consts.ListAllErr)
		return responseJSON, nil

	}

	validationErrorValue, validationErrorExists := enValue[errType].(bson.M)
	if !validationErrorExists {
		log.Println(consts.ValidationErrMessage)
		errMsg := fmt.Sprintf("%s not found or not of the expected type", errType)
		return nil, errors.New(errMsg)
	}

	errorsValue, errorsExists := validationErrorValue[consts.Errors].(bson.M)
	if !errorsExists && !validationErrorExists {
		log.Println(consts.ErrorKeyErr)
		return nil, errors.New(consts.KeyErr)

	}
	if !errorsExists && validationErrorExists {
		responseJSON := utilities.ErrorResponseFormat(consts.AllErrMsg, validationErrorValue[consts.ErrorCode], validationErrorValue, errType)
		return responseJSON, nil
	}

	// If the endpoint is empty, return all  validation error response.
	if endpoint == "" {
		responseJSON := utilities.ErrorResponseFormat(consts.ValidataionErr, validationErrorValue[consts.ErrorCode], errorsValue, errType)
		return responseJSON, nil
	}

	endpointFieldValue, endpointFieldExists := errorsValue[endpoint].(bson.M)
	if !endpointFieldExists {
		log.Println(consts.EndpointErr)
		errMsg := fmt.Sprintf("%s not found or not of the expected endpoint", endpoint)
		return nil, errors.New(errMsg)
	}

	// Check if the method is empty and the field is not empty.
	// This typically indicates a key error related to the method.
	if method == "" && field != "" {
		log.Println(consts.MethodKeyErr)
		return nil, errors.New(consts.MethodErr)
	}

	// If the method is empty, generate a validation error response.
	if method == "" {
		responseJSON := utilities.ErrorResponseFormat(consts.ValidataionErr, validationErrorValue[consts.ErrorCode], endpointFieldValue, errType)
		return responseJSON, nil
	}

	endpointMethodValue, endpointMethodExists := endpointFieldValue[method].(bson.M)
	if !endpointMethodExists {
		log.Println(consts.MethodKeyErr)
		errMsg := fmt.Sprintf("%s not found or not of the expected endpoint", method)
		return nil, errors.New(errMsg)
	}

	// Check if the endpoint field exists and the field is empty.
	// This scenario typically indicates a validation error.
	if endpointFieldExists && field == "" {
		responseJSON := utilities.ErrorResponseFormat(consts.ValidataionErr, validationErrorValue[consts.ErrorCode], endpointMethodValue, endpoint)
		return responseJSON, nil
	}

	fieldErrors, fieldErrorsExist := utilities.ExtractFieldErrors(field, errorsValue, endpoint, method)
	if !fieldErrorsExist {
		log.Println(consts.FieldErr)
		return nil, errors.New(consts.FieldKeyErr)
	}

	responseJSON := utilities.ErrorResponseFormat(consts.ValidataionErr, validationErrorValue[consts.ErrorCode], fieldErrors, endpoint)
	return responseJSON, nil
}

//AppendError is used to add errors based on various params like
//params

// @errType 	- Type of error
// @endpoint 	- Endpoint names
// @lang 		- Language
// @code 		- input error json
// @method 		- Method of the endpoint (get,post..)
func (repo *ErrorCodesRepo) AppendError(ctx context.Context, errType string, endpoint string, lang string, code db.ErrorData, method string) error {
	collection := repo.Collection.Error
	updateData := code.Field
	id, err := primitive.ObjectIDFromHex(repo.Cfg.DocID)
	if err != nil {
		log.Println(consts.MongoCollectionErr)
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{}}

	// Check if both errType and endpoint are empty.
	// This condition is typically used to handle a specific case.
	if errType == "" && endpoint == "" {
		// Iterate over the updateData map to construct dynamic keys and update the corresponding values in the MongoDB update query.
		for key, value := range updateData {
			dynamicKey := fmt.Sprintf("%s.%s", lang, key)
			update["$set"].(bson.M)[dynamicKey] = value
		}
	} else {
		// Construct the key based on method, errType, and endpoint
		fieldKey := fmt.Sprintf("%s.%s.%s.%s.%s", lang, errType, consts.Errors, endpoint, method)
		for key, value := range updateData {
			fieldInnerKey := fmt.Sprintf("%s.%s", fieldKey, key)
			update["$set"].(bson.M)[fieldInnerKey] = value
		}
	}

	options := options.Update().SetUpsert(true)
	// Perform the update operation
	_, err = collection.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		log.Println(consts.FailedErr)
		return err
	}

	return nil
}

// AddTranslation is used to translate the errors in different languages specified in database
// params
// @language	- Language list from database
func (repo *ErrorCodesRepo) AddTranslation(ctx context.Context, language []string) (any, error) {
	collection := repo.Collection.Error

	// Original JSON data
	resultJson := map[string]any{}

	// Define a filter to match the document you want
	id, err := primitive.ObjectIDFromHex(repo.Cfg.DocID)
	if err != nil {
		log.Println(consts.MongoCollectionErr)
		return nil, err
	}

	filter := bson.M{"_id": id}

	var resultData Document
	err = collection.FindOne(context.Background(), filter).Decode(&resultData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(consts.NotFoundErr)
			return nil, errors.New(consts.NoDocumentErr)
		} else {
			return nil, err
		}
	}

	// Create a new variable to store the translated data
	translatedData := make(map[string]any)

	// Translate values recursively and store the translated data in the new variable
	for _, lang := range language {
		translatedData[lang] = translateValues(resultData.En, lang)
	}

	// Convert the translated data to a formatted JSON string and print it
	jsonString, err := json.MarshalIndent(translatedData, "", "  ")
	if err != nil {
		return "", err
	}

	// Update the document with the translated content
	update := bson.M{
		"$set": translatedData,
	}
	// Perform the update
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	fmt.Printf("Matched %v document(s) and modified %v document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	json.Unmarshal(jsonString, &resultJson)
	return resultJson, nil
}

// translateValues is used to translate in given language
func translateValues(data any, lang string) any {
	switch v := data.(type) {
	case string:
		// Translate the string into the target language
		result, err := gt.Translate(v, "en", lang)
		if err != nil {
			fmt.Printf("Error translating to %s: %v\n", lang, err)
			return v // Return the original string in case of an error
		}
		return result
	case map[string]any:
		// Recursively translate values in a map and return a new map
		translatedMap := make(map[string]any)
		for key, value := range v {
			translatedMap[key] = translateValues(value, lang)
		}
		return translatedMap
	default:
		return v
	}
}

// AddEndpoint is used to add endpoints for error mapping with endpoint name
//params

// @endpt 	- Json data includes (url,method and endpoint name)
func (repo *ErrorCodesRepo) AddEndpoint(ctx context.Context, endpt entities.RequestData) error {
	collection := repo.Collection.Endpoint

	if endpt.Field.URL == "" || endpt.Field.Method == "" || endpt.Field.Endpoint == "" {
		return fmt.Errorf(consts.UrlErr)
	}

	id, err := primitive.ObjectIDFromHex(repo.Cfg.EndpointDocID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing document: %v", err)
	}
	filter := bson.M{"_id": id}

	var existingDocument entities.MyDocument
	err = collection.FindOne(ctx, filter).Decode(&existingDocument)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing document: %v", err)
	}

	newData := entities.EndpointData{
		URL:      endpt.Field.URL,
		Method:   endpt.Field.Method,
		Endpoint: endpt.Field.Endpoint,
	}

	existingDocument.Data = append(existingDocument.Data, newData)

	// Replace the existing document with the modified one
	_, err = collection.ReplaceOne(ctx, filter, existingDocument)
	if err != nil {
		return fmt.Errorf("failed to update existing document: %v", err)
	}

	return nil
}

// GetEndpointName is used to get endpoint names from database based on url and method
func (repo *ErrorCodesRepo) GetEndpointName(ctx context.Context) (entities.ResponseData, error) {
	collection := repo.Collection.Endpoint

	// Convert the string representation of ObjectID to ObjectID type
	objectID, err := primitive.ObjectIDFromHex(repo.Cfg.EndpointDocID)
	if err != nil {
		return entities.ResponseData{}, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Create a filter based on the provided criteria
	filter := bson.M{
		"_id": objectID,
	}

	// Find one document that matches the filter
	var result bson.M
	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return entities.ResponseData{}, err
		}
		return entities.ResponseData{}, err
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return entities.ResponseData{}, err
	}
	jsonString := string(jsonBytes)

	var response entities.ResponseData
	if err := json.Unmarshal([]byte(jsonString), &response); err != nil {
		// Handle error
		return entities.ResponseData{}, err
	}

	return response, err

}
