package api

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// PeerReviewsService handles peer review-related API calls
type PeerReviewsService struct {
	client *Client
}

// NewPeerReviewsService creates a new peer reviews service
func NewPeerReviewsService(client *Client) *PeerReviewsService {
	return &PeerReviewsService{client: client}
}

// PeerReview represents a Canvas peer review
type PeerReview struct {
	ID                 int64               `json:"id"`
	AssessorID         int64               `json:"assessor_id"`
	AssetID            int64               `json:"asset_id"`
	AssetType          string              `json:"asset_type"`
	UserID             int64               `json:"user_id"`
	WorkflowState      string              `json:"workflow_state"`
	User               *User               `json:"user,omitempty"`
	Assessor           *User               `json:"assessor,omitempty"`
	SubmissionComments []SubmissionComment `json:"submission_comments,omitempty"`
}

// ListPeerReviewsOptions holds options for listing peer reviews
type ListPeerReviewsOptions struct {
	Include []string // submission_comments, user
}

// List retrieves peer reviews for an assignment
func (s *PeerReviewsService) List(ctx context.Context, courseID, assignmentID int64, opts *ListPeerReviewsOptions) ([]PeerReview, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/peer_reviews", courseID, assignmentID)

	if opts != nil && len(opts.Include) > 0 {
		query := url.Values{}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var reviews []PeerReview
	if err := s.client.GetAllPages(ctx, path, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// ListForSubmission retrieves peer reviews for a specific submission
func (s *PeerReviewsService) ListForSubmission(ctx context.Context, courseID, assignmentID, submissionID int64, opts *ListPeerReviewsOptions) ([]PeerReview, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d/peer_reviews", courseID, assignmentID, submissionID)

	if opts != nil && len(opts.Include) > 0 {
		query := url.Values{}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var reviews []PeerReview
	if err := s.client.GetAllPages(ctx, path, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// CreatePeerReviewParams holds parameters for creating a peer review
type CreatePeerReviewParams struct {
	UserID int64 // The reviewer's user ID
}

// Create creates a new peer review assignment
func (s *PeerReviewsService) Create(ctx context.Context, courseID, assignmentID, submissionID int64, params *CreatePeerReviewParams) (*PeerReview, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d/peer_reviews", courseID, assignmentID, submissionID)

	body := make(map[string]interface{})
	body["user_id"] = params.UserID

	var review PeerReview
	if err := s.client.PostJSON(ctx, path, body, &review); err != nil {
		return nil, err
	}

	return &review, nil
}

// Delete removes a peer review assignment
func (s *PeerReviewsService) Delete(ctx context.Context, courseID, assignmentID, submissionID, reviewerID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d/peer_reviews", courseID, assignmentID, submissionID)

	query := url.Values{}
	query.Add("user_id", fmt.Sprintf("%d", reviewerID))
	path += "?" + query.Encode()

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// ListSections retrieves peer review sections for an assignment
func (s *PeerReviewsService) ListSections(ctx context.Context, courseID, assignmentID int64, include []string) ([]PeerReview, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/peer_reviews", courseID, assignmentID)

	if len(include) > 0 {
		path += "?include[]=" + strings.Join(include, "&include[]=")
	}

	var reviews []PeerReview
	if err := s.client.GetAllPages(ctx, path, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}
