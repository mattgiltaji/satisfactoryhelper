package cloudfuncs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//look into https://cloud.google.com/sdk/gcloud/reference/beta/emulators/firestore/

func TestMachineHttpBadMethods(t *testing.T) {
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
		is.Equal(testCase.expectedStatus, resp.StatusCode, "%s should fail, not implemented yet", testCase.method)
	}
}

func TestMachineHttpGet(t *testing.T) {

	is := assert.New(t)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	MachineHttp(rr, req)

	resp := rr.Result()
	is.Equal(http.StatusOK, resp.StatusCode, "GET should be ok")
}
