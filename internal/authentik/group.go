package authentik

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	createGroupPath string = "/api/v3/core/groups"
)

// GroupRequest is the request body for creating/updating a group.
type GroupRequest struct {
	Name        string         `json:"name"`
	IsSuperuser bool           `json:"is_superuser,omitempty"`
	Parents     []string       `json:"parent,omitempty"`
	Users       []int          `json:"users,omitempty"`
	Attributes  map[string]any `json:"attributes,omitempty"`
}

// Group is the response from the Authentik API.
type Group struct {
	PK          string         `json:"pk"`
	Name        string         `json:"name"`
	IsSuperuser bool           `json:"is_superuser"`
	Parents     []string       `json:"parent"`
	Attributes  map[string]any `json:"attributes"`
	Users       []int          `json:"users"`
}

func (c *Authentik) CreateGroup(ctx context.Context, req GroupRequest) (*Group, error) {
	resp, err := c.do(ctx, http.MethodPost, createGroupPath, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create group: %w", parseApiError(resp))
	}

	var group Group

	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return nil, fmt.Errorf("failed to decode response for create group request: %w", err)
	}

	return &group, nil
}
