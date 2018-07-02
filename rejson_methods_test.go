package rejonson

import (
	"testing"
	"github.com/KromDaniel/jonson"
	"time"
	"math/rand"
	"github.com/stretchr/testify/assert"
	"github.com/go-redis/redis"
	"encoding/json"
	"io/ioutil"
	"sort"
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
	redisTestsPrefix string
)

func concatKey(key string)string {
	return redisTestsPrefix + key
}
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
	var redisOptions redis.Options
	config, _ := ioutil.ReadFile("test_config.json")
	allOptions, err := jonson.Parse(config)

	if err != nil {
		panic("error with reading config file " +  err.Error())
	}

	err = json.Unmarshal(allOptions.At("redisConnection").ToUnsafeJSON(), &redisOptions)
	if err != nil {
		panic("error with reading redis config " +  err.Error())
	}
	client = ExtendClient(redis.NewClient(&redisOptions))
	redisTestsPrefix = allOptions.At("redisKeyPrefix").GetUnsafeString()
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
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	typeRes, err := client.JsonType(key,"").Result()
	assert.Nil(t, err)
	assert.Equal(t, "object", typeRes)

	typeRes, err = client.JsonType(key,"keyB").Result()
	assert.Nil(t, err)
	assert.Equal(t, "string", typeRes)
}

func TestRedisProcessor_JsonNumIncrBy(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	incRes, err := client.JsonNumIncrBy(key, "keyA", 4).Result()
	assert.Nil(t, err)
	assert.Equal(t, "60", incRes)
}

func TestRedisProcessor_JsonNumMultBy(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	multRes, err := client.JsonNumMultBy(key, "numbersArray[1]", 4).Result()
	assert.Nil(t, err)
	assert.Equal(t, "8", multRes)
}

func TestRedisProcessor_JsonStrAppend(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	strAppRes, err := client.JsonStrAppend(key, "keyB", " \"hello\"").Result()
	assert.Nil(t, err)

	assert.Equal(t,16, int(strAppRes))
}

func TestRedisProcessor_JsonStrLen(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	strLenRes, err := client.JsonStrLen(key, "keyB").Result()
	assert.Nil(t, err)
	assert.Equal(t, len("some string"), int(strLenRes))
}

func TestRedisProcessor_JsonArrAppend(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	arrAppendRes, err := client.JsonArrAppend(key, "numbersArray", 12).Result()
	assert.Nil(t, err)
	assert.Equal(t, 6, int(arrAppendRes))
}

func TestRedisProcessor_JsonArrIndex(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	jsn := jonson.NewEmptyJSONArray()

	for i:=0; i < 100; i++ {
		jsn.SliceAppend(i)
	}

	setErr := client.JsonSet(key, ".", jsn.ToUnsafeJSONString()).Err()
	assert.Nil(t, setErr)

	arrIndexRes, err := client.JsonArrIndex(key, ".", 5).Result()
	assert.Nil(t, err)
	assert.Equal(t, 5, int(arrIndexRes))


	arrIndexRes, err = client.JsonArrIndex(key, ".", 5, 10, 90).Result()
	assert.Nil(t, err)
	assert.Equal(t, -1, int(arrIndexRes))
}

func TestRedisProcessor_JsonArrInsert(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	arrInsertRes, err := client.JsonArrInsert(key, "numbersArray", 1, "2").Result()
	assert.Nil(t, err)
	assert.Equal(t, 6, int(arrInsertRes))
}

func TestRedisProcessor_JsonArrLen(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	arrLenRes, err := client.JsonArrLen(key, "numbersArray").Result()
	assert.Nil(t, err)
	assert.Equal(t, 5, int(arrLenRes))
}

func TestRedisProcessor_JsonArrPop(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	arrPopRes, err := client.JsonArrPop(key, "numbersArray", 1).Result()
	assert.Nil(t, err)
	assert.Equal(t, "2", arrPopRes)
}

func TestRedisProcessor_JsonArrTrim(t *testing.T) {
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	arrTrimRes, err := client.JsonArrTrim(key, "numbersArray", 1,3).Result()
	assert.Nil(t, err)
	assert.Equal(t, 3, int(arrTrimRes))

	trimArr, err := client.JsonGet(key, "numbersArray").Result()
	assert.Nil(t, err)
	assertJonsonEqual(jonson.New([]float64{2,3,4}), jonson.ParseUnsafe([]byte(trimArr)) ,t, true)
}

func TestRedisProcessor_JsonObjKeys(t *testing.T) {
	originalJS := jonson.ParseUnsafe([]byte(baseJsonTestObject))
	originalKeys := originalJS.GetObjectKeys()
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	objKeysRes, err := client.JsonObjKeys(key, ".").Result()
	assert.Nil(t, err)
	sort.Strings(originalKeys)
	sort.Strings(objKeysRes)
	assert.Equal(t, originalKeys, objKeysRes)
}

func TestRedisProcessor_JsonObjLen(t *testing.T) {
	originalJS := jonson.ParseUnsafe([]byte(baseJsonTestObject))
	key := concatKey(randStringRunes(32))
	defer client.Del(key)

	if !insertBaseJsonToRedis(key, t) {
		t.FailNow()
	}

	objLenRes, err := client.JsonObjLen(key, ".").Result()
	assert.Nil(t, err)
	assert.Equal(t, len(originalJS.GetObjectKeys()), int(objLenRes))
}

func TestClient_Pipeline(t *testing.T) {
	allKeys := make([]string, 0)
	pipeline := client.Pipeline()

	for i:=0; i < 10; i++ {
		key := concatKey(randStringRunes(32))
		pipeline.JsonSet(key, ".", baseJsonTestObject)
		allKeys = append(allKeys, key)
	}

	// here we expected that delete counter will be 0
	delRes, err := client.Del(allKeys...).Result()
	assert.Nil(t, err)
	assert.Equal(t, 0, int(delRes))

	_, err = pipeline.Exec()
	assert.Nil(t, err)
	// now we expect deleted count to be same as allKeysLength
	delRes, err = client.Del(allKeys...).Result()
	assert.Nil(t, err)
	assert.Equal(t, len(allKeys), int(delRes))
}