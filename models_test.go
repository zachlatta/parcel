package parcel

import (
	"encoding/json"
	"reflect"
	"testing"
)

var sampleMsgs = []Message{
	Message{
		"method1",
		map[string]interface{}{
			"arg1": "arg1data",
			"arg2": "arg2data",
		},
		"#1",
	},
	Message{
		"method2",
		map[string]interface{}{
			"arg1": "arg1data",
		},
		"#2",
	},
	Message{
		"method3",
		map[string]interface{}{},
		"#3",
	},
}

var sampleMsgJSON = `[["method1",{"arg1":"arg1data","arg2":"arg2data"},"#1"],["method2",{"arg1":"arg1data"},"#2"],["method3",{},"#3"]]`

func TestMessageMarshalJSON(t *testing.T) {
	expected := sampleMsgJSON

	actual, err := json.Marshal(sampleMsgs)
	if err != nil {
		t.Error(err)
	}

	if expected != string(actual) {
		t.Errorf("expected: %v, got: %v", expected, actual)
	}
}

func TestMessageUnmarshalJSON(t *testing.T) {
	expected := sampleMsgs
	actual := []Message{}

	if err := json.Unmarshal([]byte(sampleMsgJSON), &actual); err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error("expected: %v, got: %v", expected, actual)
	}
}
