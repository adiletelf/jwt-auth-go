package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Generate(t *testing.T) {
	req, _ := http.NewRequest("POST", generatePath, nil)
	q := req.URL.Query()
	q.Add("uuid", uuid.New().String())
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
	validTokenDetails(t, response)
}
