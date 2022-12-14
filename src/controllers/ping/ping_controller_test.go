package ping

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/ping",nil)

	Ping(c)

	pingResponse, error := io.ReadAll(response.Body)

	assert.Nil(t, error)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "pong", string(pingResponse))
	
}


