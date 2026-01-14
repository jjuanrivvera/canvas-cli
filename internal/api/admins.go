package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AdminsService handles admin-related API calls
type AdminsService struct {
	client *Client
}

// NewAdminsService creates a new admins service
func NewAdminsService(client *Client) *AdminsService {
	return &AdminsService{client: client}
}

// Admin represents a Canvas account administrator
type Admin struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id,omitempty"`
	User   *User  `json:"user,omitempty"`
	Role   string `json:"role,omitempty"`
	RoleID int64  `json:"role_id,omitempty"`
	Status string `json:"status,omitempty"`
}

// ListAdminsOptions holds options for listing admins
type ListAdminsOptions struct {
	UserID  []int64
	Page    int
	PerPage int
}

// List retrieves admins for an account
func (s *AdminsService) List(ctx context.Context, accountID int64, opts *ListAdminsOptions) ([]Admin, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/admins", accountID)

	if opts != nil {
		query := url.Values{}

		for _, uid := range opts.UserID {
			query.Add("user_id[]", strconv.FormatInt(uid, 10))
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var admins []Admin
	if err := s.client.GetAllPages(ctx, path, &admins); err != nil {
		return nil, err
	}

	return admins, nil
}

// CreateAdminParams holds parameters for creating an admin
type CreateAdminParams struct {
	UserID           int64
	Role             string // AccountAdmin, etc.
	RoleID           int64
	SendConfirmation *bool
}

// Create adds an admin to an account
func (s *AdminsService) Create(ctx context.Context, accountID int64, params *CreateAdminParams) (*Admin, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/admins", accountID)

	body := make(map[string]interface{})
	body["user_id"] = params.UserID

	if params.Role != "" {
		body["role"] = params.Role
	}

	if params.RoleID > 0 {
		body["role_id"] = params.RoleID
	}

	if params.SendConfirmation != nil {
		body["send_confirmation"] = *params.SendConfirmation
	}

	var admin Admin
	if err := s.client.PostJSON(ctx, path, body, &admin); err != nil {
		return nil, err
	}

	return &admin, nil
}

// Delete removes an admin from an account
func (s *AdminsService) Delete(ctx context.Context, accountID, userID int64, roleID *int64) (*Admin, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/admins/%d", accountID, userID)

	if roleID != nil {
		query := url.Values{}
		query.Add("role_id", strconv.FormatInt(*roleID, 10))
		path += "?" + query.Encode()
	}

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &Admin{UserID: userID, Status: "removed"}, nil
}
