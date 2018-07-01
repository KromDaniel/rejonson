package rejonson

import (
	"testing"
	"github.com/KromDaniel/jonson"
	"time"
	"math/rand"
	"fmt"
)

var (
	letterRunes        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	baseJsonTestObject = `
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
	client *Client
)

func testEqual(left *jonson.JSON, right *jonson.JSON, t *testing.T, shouldEqual bool) {
	var errorMessage string
	if shouldEqual {
		errorMessage = "should equal to"
	} else {
		errorMessage = "shouldn't equal to"
	}

	isEqual := jonson.EqualsDeep(left, right)
	if isEqual != shouldEqual {
		t.Errorf("%s %s %s", left.ToUnsafeJSONString(), errorMessage, right.ToUnsafeJSONString())
	}
}

func assertErrorNotNil(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Failed test %s with error %S", t.Name(), err.Error())
		t.FailNow()
	}
}

func insertBaseJsonToRedis(key string, t *testing.T) (success bool) {
	if err := client.JsonSet(key,"." ,baseJsonTestObject).Err(); err != nil {
		t.Errorf("Couldnt set initial value to redis with error %s", err.Error())
		return false
	}
	return true
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	client = GetRedisConnection()
	defer client.Close()
	m.Run()

}

func TestRedisProcessor_JsonDel(t *testing.T) {
	// originalJS := jonson.ParseUnsafe([]byte(baseJsonTestObject))
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	delRes, err := client.JsonDel(key, "keyA").Result()
	assertErrorNotNil(err, t)


	fmt.Print("DEL RES", delRes)
}
