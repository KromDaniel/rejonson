package rejonson

import (
	"testing"
	"github.com/KromDaniel/jonson"
)

var base_json_object = `
{
  "keyA": 56,
  "keyB": "some string",
  "numbersArray": [
    1,
    2,
    3,
    4,
    5
  ],
  "strArray": [
    "a",
    "b",
    "c"
  ]
}
`

func TestRejsonMethods(t *testing.T) {
	originalJS := jonson.ParseUnsafe([]byte(base_json_object))
	t.Log("Getting redis connection and setting random value with JSON.SET")
	client := GetRedisConnection()

	setErr := client.JsonSet(concatKey("initial_key"), ".", originalJS.ToUnsafeJSONString()).Err()
	if setErr != nil {
		t.Fatal("Couldn't initiate redis basic JSON.set command", setErr)
	}

	jsonBytes, getJsonError := client.JsonGet(concatKey("initial_key")).Bytes()

	if getJsonError != nil {
		t.Fatal("Unable to perform JSON.GET", getJsonError)
	}

	parsedJonson, parsedJonsonError := jonson.Parse(jsonBytes)
	if parsedJonsonError != nil {
		t.Fatal("Unable to parse returned json ", parsedJonsonError)
	}

	if !jonson.EqualsDeep(originalJS, parsedJonson){
		t.Fatal("Not equals deep!")
	}
}
