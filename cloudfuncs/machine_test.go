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

func TestGetAllMachinesZeroResults(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()
	client = getTestClient(ctx, t)

	//zero results
	actual, err := getAllMachines(ctx, client)
	is.NoError(err, "getAllMachines() should not error out if no machines found")
	is.Equal(0, len(actual))
	is.Equal([]Machine(nil), actual, "getAllMachines() should not return any machines if there is nothing in the firestore")

}

func TestGetAllMachinesOneGoodResult(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()
	client = getTestClient(ctx, t)
	coll := client.Collection("machines")

	assembler, _, _ := coll.Add(ctx, map[string]interface{}{
		"name":  "Assembler",
		"power": 15,
	})
	defer deleteDocRef(ctx, t, assembler)
	expectedAssembler := Machine{"Assembler", 15}

	actual, err := getAllMachines(ctx, client)
	is.NoError(err, "getAllMachines() should not error out if one machine found")
	is.Equal(1, len(actual))
	is.Contains(actual, expectedAssembler)
}

func TestGetAllMachinesSeveralGoodResults(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()
	client = getTestClient(ctx, t)
	coll := client.Collection("machines")
	constructor, _, _ := coll.Add(ctx, map[string]interface{}{
		"name":  "Constructor",
		"power": 4,
	})
	defer deleteDocRef(ctx, t, constructor)
	foundry, _, _ := coll.Add(ctx, map[string]interface{}{
		"name":  "Foundry",
		"power": 16,
	})
	defer deleteDocRef(ctx, t, foundry)
	expectedConstructor := Machine{"Constructor", 4}
	expectedFoundry := Machine{"Foundry", 16}

	actual, err := getAllMachines(ctx, client)
	is.NoError(err, "getAllMachines() should not error out if multiple machines found")
	is.Equal(2, len(actual))
	is.Contains(actual, expectedConstructor)
	is.Contains(actual, expectedFoundry)
	//don't do equals because we can't assume retrieval order
}

func TestGetAllMachinesBadResult(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()
	client = getTestClient(ctx, t)
	coll := client.Collection("machines")

	foundry, _, _ := coll.Add(ctx, map[string]interface{}{
		"name":  "Foundry",
		"power": "16",
	})
	defer deleteDocRef(ctx, t, foundry)
	//expectedAssembler := Machine{"Assembler", 15}

	_, err := getAllMachines(ctx, client)
	is.Error(err, "getAllMachines() should error out if it can't convert results to Machine type")
}
