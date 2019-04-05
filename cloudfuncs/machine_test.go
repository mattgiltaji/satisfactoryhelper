package cloudfuncs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineHttpMethodCheck(t *testing.T) {
	testCases := []struct {
		method         string
		expectedStatus int
	}{
		{method: "DELETE", expectedStatus: http.StatusMethodNotAllowed},
		{method: "POST", expectedStatus: http.StatusMethodNotAllowed},
		{method: "PUT", expectedStatus: http.StatusMethodNotAllowed},
	}

	is := assert.New(t)
	for _, testCase := range testCases {
		req := httptest.NewRequest(testCase.method, "/", nil)
		rr := httptest.NewRecorder()
		MachineHttp(rr, req)

		resp := rr.Result()
		is.Equal(resp.StatusCode, testCase.expectedStatus, fmt.Sprintf("%s should fail, not implemented yet", testCase.method))
	}
}
