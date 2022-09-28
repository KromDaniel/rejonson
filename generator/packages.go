package main

var packages = []GoredisPackage{
	{
		GoRedisImport: "github.com/go-redis/redis",
		Output:        "v6.generated.go",
	},
	{
		GoRedisImport: "github.com/go-redis/redis/v7",
		Output:        "v7/v7.generated.go",
	},
	{
		GoRedisImport: "github.com/go-redis/redis/v8",
		Output:        "v8/v8.generated.go",
		HasContext:    true,
	},
	{
		GoRedisImport: "github.com/go-redis/redis/v9",
		Output:        "v9/v9.generated.go",
		HasContext:    true,
	},
}
