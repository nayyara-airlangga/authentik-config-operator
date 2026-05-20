package authentik_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nayyara-airlangga/authentik-config-operator/internal/authentik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v3/core/groups", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req authentik.GroupRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "test-group", req.Name)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(authentik.Group{
			PK:          "some-uuid",
			Name:        "test-group",
			IsSuperuser: false,
		})
	}))

	defer server.Close()

	client, err := authentik.New(server.URL, "test-token", []byte{}, true)

	require.NoError(t, err)

	group, err := client.CreateGroup(t.Context(), authentik.GroupRequest{
		Name: "test-group",
	})

	require.NoError(t, err)

	assert.Equal(t, "some-uuid", group.PK)
	assert.Equal(t, "test-group", group.Name)
}

func TestCreateGroup_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req authentik.GroupRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Empty(t, req.Name)

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"name": []string{"This field is required."},
		})
	}))

	defer server.Close()

	client, err := authentik.New(server.URL, "test-token", []byte{}, true)

	require.NoError(t, err)

	_, err = client.CreateGroup(t.Context(), authentik.GroupRequest{})

	require.Error(t, err)

	var apiErr *authentik.ApiError

	require.ErrorAs(t, err, &apiErr)

	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
}
