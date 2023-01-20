package bootstrap

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
)

const _pingResponseDefaultSuccessful = "ping: PONG"

func newRedis() *redis.Client {
	defaultDBString := os.Getenv("CIA_REDIS_DBDEFAULT")
	defaultDB, err := strconv.Atoi(defaultDBString)
	if err != nil {
		log.Fatalln("environment CIA_REDIS_DBDEFAULT must be int")
	}

	isSslString := os.Getenv("CIA_REDIS_SSL")
	isSsl, err := strconv.ParseBool(isSslString)
	if err != nil {
		log.Fatalln("environment CIA_REDIS_SSL must be bool")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", os.Getenv("CIA_REDIS_HOST"), os.Getenv("CIA_REDIS_PORT")),
		Password:  os.Getenv("CIA_REDIS_PASSWORD"),
		DB:        defaultDB,
		TLSConfig: &tls.Config{InsecureSkipVerify: isSsl},
	})

	if stringResponse := rdb.Ping(context.Background()); stringResponse.String() != _pingResponseDefaultSuccessful {
		log.Fatal(stringResponse)
	}

	return rdb
}
