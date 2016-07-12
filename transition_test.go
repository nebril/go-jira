package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTransitionGetList(t *testing.T) {
	setup()
	defer teardown()

	testAPIEndpoint := "/rest/api/2/issue/123/transitions"

	raw, err := ioutil.ReadFile("./mocks/transitions.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	transitions, _, err := testClient.Transition.GetList("123")

	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if transitions == nil {
		t.Error("Expected transition list. Got nil.")
	}

	if len(transitions) != 2 {
		t.Errorf("Expected 2 transitions. Got %d", len(transitions))
	}

	if transitions[0].Fields["summary"].Required != false {
		t.Errorf("First transition summary field should not be required")
	}
}

func TestTransitionCreate(t *testing.T) {
	setup()
	defer teardown()

	testAPIEndpoint := "/rest/api/2/issue/123/transitions"

	transitionID := "22"

	fields := map[string]TransitionField{
		"field": TransitionField{
			Name: "fieldname",
		}}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, testAPIEndpoint)

		decoder := json.NewDecoder(r.Body)
		var payload CreateTransitionPayload
		err := decoder.Decode(&payload)
		if err != nil {
			t.Error("Got error: %v", err)
		}

		if payload.Transition.ID != transitionID {
			t.Errorf("Expected %s to be in payload, got %s instead", transitionID, payload.Transition.ID)
		}
		if payload.Fields["field"] != fields["field"] {
			t.Errorf("Expected fields to be in payload")
		}
	})
	_, err := testClient.Transition.Create("123", transitionID, fields)

	if err != nil {
		t.Error("Got error: %v", err)
	}
}
