package ping

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"tokenalert_user-api/utils/test_utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/ping",nil)

	Ping(c)

	pingResponse, error := ioutil.ReadAll(response.Body)

	assert.Nil(t, error)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "pong", string(pingResponse))
	
}


