package core

import (
	"context"
	"encoding/json"
	"errors"
	"ips-lacpass-backend/internal/users/client"
	customErrors "ips-lacpass-backend/pkg/errors"
	authMiddleware "ips-lacpass-backend/pkg/middleware"
)

type ServiceInterface interface {
	CreateUser(ctx context.Context, ur UserRequest) (*User, error)
	UpdateUser(ctx context.Context, ur UserUpdateRequest) (*User, error)
}

type UserService struct {
	Client client.ClientInterface
}

func NewService(r client.ClientInterface) UserService {
	return UserService{
		Client: r,
	}
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
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

	// TODO fix this workaround for cyclic dependency issue
	userMap, err := structToMap(user)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to convert user to map"}},
			Err:        err,
		}
	}

	resp, err := us.Client.CreateUser(ctx, userMap, ur.Password)
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
	user.ID = resp.ID
	return user, nil
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

	updateMap, err := structToMap(ur)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to convert update request to map"}},
			Err:        err,
		}
	}

	resp, err := us.Client.UpdateUser(ctx, userUUID, updateMap)
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
	return &User{
		Username:     resp.ID,
		Email:        resp.Email,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		Locale:       resp.Attributes["locale"][0],
		DocumentType: AllowedDocumenTypes[resp.Attributes["document_type"][0]],
		Identifier:   resp.ID,
	}, nil
}
