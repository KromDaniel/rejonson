package rejonson

import (
	"testing"
	"github.com/KromDaniel/jonson"
	"time"
	"math/rand"
	"github.com/stretchr/testify/assert"
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

func assertJonsonEqual(left *jonson.JSON, right *jonson.JSON, t *testing.T, shouldEqual bool) {
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

func insertBaseJsonToRedis(key string, t *testing.T) (success bool) {
	if err := client.JsonSet(key,"." ,baseJsonTestObject).Err(); err != nil {
		t.Errorf("Couldnt set initial value to redis with error %s", err.Error())
		return false
	}
	return true
}

func getBaseJsonFromRedis(key string, t *testing.T)*jonson.JSON {
	b, err := client.JsonGet(key).Bytes()
	if err != nil {
		t.Errorf("Unable to get JSON with key %s with error %s", key, err.Error())
		t.FailNow()
		return nil
	}

	return jonson.ParseUnsafe(b)
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
	originalJS := jonson.ParseUnsafe([]byte(baseJsonTestObject))
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	delRes, err := client.JsonDel(key, "keyA").Result()
	assert.Nil(t, err)
	assert.NotEqual(t, 1, delRes)

	originalJS.DeleteMapKey("keyA")
	assertJonsonEqual(originalJS, getBaseJsonFromRedis(key, t), t, true)
}

func TestRedisProcessor_JsonGet(t *testing.T) {
	originalJS := jonson.ParseUnsafe([]byte(baseJsonTestObject))
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	// check first that the entire object returns
	getRes, err := client.JsonGet(key).Bytes()
	assert.Nil(t, err)
	assertJonsonEqual(originalJS, jonson.ParseUnsafe(getRes), t, true)

	// check that nested object returned
	getRes, err = client.JsonGet(key, "numbersArray").Bytes()
	assert.Nil(t, err)
	assertJonsonEqual(originalJS.At("numbersArray"), jonson.ParseUnsafe(getRes), t, true)
}

func TestRedisProcessor_JsonSet(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	setRes, err := client.JsonSet(key, ".", baseJsonTestObject).Result()

	assert.Nil(t, err)
	assert.Equal(t, "OK", setRes)
}

func TestRedisProcessor_JsonMGet(t *testing.T) {
	originalJS := jonson.ParseUnsafe([]byte(baseJsonTestObject))
	keyA := concatKey(randStringRunes(32))
	keyB := concatKey(randStringRunes(32))
	defer client.Del(keyA, keyB)

	if !insertBaseJsonToRedis(keyA, t) {
		t.FailNow()
	}

	if !insertBaseJsonToRedis(keyB, t) {
		t.FailNow()
	}

	mGetRes, err := client.JsonMGet(keyA, keyB, "strArray").Result()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(mGetRes))
	assertJonsonEqual(originalJS.At("strArray"), jonson.ParseUnsafe([]byte(mGetRes[0])),t,true)
	assertJonsonEqual(originalJS.At("strArray"), jonson.ParseUnsafe([]byte(mGetRes[1])),t,true)
}

func TestRedisProcessor_JsonType(t *testing.T) {
	
}

