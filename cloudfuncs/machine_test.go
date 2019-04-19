package cloudfuncs

import (
	"context"
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

func TestGetAllMachines(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()
	client = getTestClient(ctx, t)
	//TODO: use a test firestore that we can try querying with zero, one, and many machines

	actual, err := getAllMachines(ctx, client)
	is.NoError(err, "getAllMachines() should not error out")
	is.Equal(5, len(actual))
	is.Contains(actual, Machine{"Assembler", 15})
	is.Contains(actual, Machine{"Constructor", 4})
	is.Contains(actual, Machine{"Foundry", 16})
	is.Contains(actual, Machine{"Smelter", 4})
	is.Contains(actual, Machine{"Manufacturer", 55})
	//don't do equals because we can't assume retrieval order

}
