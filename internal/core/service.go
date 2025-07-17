package core

import (
	"context"
	"errors"
	"fmt"
	customErrors "ips-lacpass-backend/internal/errors"
	"ips-lacpass-backend/internal/repository/fhir"
	"ips-lacpass-backend/internal/repository/vhl"
	"sort"

	authMiddleware "ips-lacpass-backend/internal/middleware"
)

type UserService struct {
	Repository UserRepository
}

type FhirService struct {
	Repository fhir.FhirRepository
}

type VhlService struct {
	Repository vhl.VhlRepository
}

func NewUserService(r UserRepository) UserService {
	return UserService{
		Repository: r,
	}
}

func NewFhirService(r fhir.FhirRepository) FhirService {
	return FhirService{
		Repository: r,
	}
}

func NewVhlService(r vhl.VhlRepository) VhlService {
	return VhlService{
		Repository: r,
	}
}

func (us *UserService) CreateUser(ctx context.Context, ur UserRequest) (*User, error) {

	user := &User{
		Username:     ur.Identifier,
		Email:        ur.Email,
		FirstName:    ur.FirstName,
		LastName:     ur.LastName,
		Locale:       ur.Locale,
		DocumentType: ur.DocumentType,
		Identifier:   ur.Identifier,
	}

	resp, err := us.Repository.CreateUser(ctx, *user, ur.Password)
	if err != nil {
		var uErr *customErrors.HttpError
		if errors.As(err, &uErr) {
			return nil, uErr
		}
		return nil, &customErrors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "auth_service_error", "message": "Failed to connect to authentication service"}}, Err: err,
		}
	}
	return resp, nil
}

func (us *UserService) UpdateUser(ctx context.Context, ur UserUpdateRequest) (*User, error) {

	userUUID, err := authMiddleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 401,
			Body:       []map[string]interface{}{{"error": "user_uuid_not_found", "message": "User UUID not found in request context"}},
			Err:        err,
		}
	}

	resp, err := us.Repository.UpdateUser(ctx, userUUID, ur)
	if err != nil {
		var uErr *customErrors.HttpError
		if errors.As(err, &uErr) {
			return nil, uErr
		}
		return nil, &customErrors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "auth_service_error", "message": "Failed to connect to authentication service"}}, Err: err,
		}
	}
	return resp, nil
}

func (fs *FhirService) GetIps(ctx context.Context) (map[string]interface{}, error) {
	userId, err := authMiddleware.GetUserDocIDFromContext(ctx)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 401,
			Body:       []map[string]interface{}{{"error": "user_identifier_not_found", "message": "User identifier not found in request context"}},
			Err:        err,
		}
	}

	bundle, err := fs.Repository.GetDocumentReference(userId)
	if err != nil {
		fmt.Printf("Error fetching document reference: %v\n", err)
		return nil, err
	}
	entries := bundle.Entry
	if len(entries) == 0 {
		return nil, &customErrors.HttpError{
			StatusCode: 404,
			Body:       []map[string]interface{}{{"error": "not_found", "message": "No IPS found for the user"}},
			Err:        fmt.Errorf("no IPS found for the user"),
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Resource.Meta.LastUpdated > entries[j].Resource.Meta.LastUpdated
	})

	ipsBundle, err := fs.Repository.GetIpsBundle(entries[0].Resource.Content[0].Attachment.URL)
	if err != nil {
		return nil, err
	}

	return ipsBundle, nil

}

func (vs *VhlService) CreateQrCode(ctx context.Context, expiresOn *string, content *string, passCode *string) (*vhl.QrData, error) {
	if content == nil || *content == "" {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "invalid_request", "message": "Content cannot be empty"}},
			Err:        errors.New("content cannot be empty"),
		}
	}

	qrData, err := vs.Repository.CreateQr(ctx, vhl.CreateQrRequest{
		JsonContent: *content,
		ExpiresOn:   *expiresOn,
		PassCode:    *passCode,
	})
	if err != nil {
		return nil, err
	}
	return qrData, nil
}
