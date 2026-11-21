package doc_test

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/KromDaniel/rejonson/v9/doc"
)

func Example() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr()})
	defer rdb.Close()

	type User struct {
		Name   string   `json:"name"`
		Visits int      `json:"visits"`
		Tags   []string `json:"tags"`
	}
	users := doc.NewStore[User](rdb, doc.Prefix("example:user:"))
	defer users.Delete(ctx, "42")

	_, _ = users.Set(ctx, "42", User{Name: "Ada", Tags: []string{"admin"}}, doc.TTL(time.Hour))
	_, _ = users.Incr(ctx, "42", "$.visits", 1)
	_, _ = users.Append(ctx, "42", "$.tags", "ops")

	u, _ := users.Update(ctx, "42", func(u *User) error {
		u.Visits += 10
		return nil
	})
	name, _ := doc.GetPath[string](ctx, users, "42", "$.name")

	fmt.Println(name, u.Visits, u.Tags)
	// Output: Ada 11 [admin ops]
}
