package main

import (
	"context"
	"fmt"

	"github.com/cosmos/go-bip39"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
	"golang.org/x/crypto/curve25519"
)

func main() {
	seed := bip39.NewSeed("mnemonic", "")
	priv, err := curve25519.X25519(seed[:32], curve25519.Basepoint)
	fmt.Println(priv)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := ydb.Open(ctx,
		"grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1gma7vam0olbnsoco9c/etn8ir0495vp4r28fdns",
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials("cmd/server/authorized_key.json"),
		// ydb.WithAccessTokenCredentials("t1.9euelZrHlM-emI3Jy46Rm8ial5iXze3rnpWaz8rIjJKXk42Pl82Uz8eQk5bl8_dCMSll-e83EgwZ_t3z9wJgJmX57zcSDBn-.kK4cxBpaIJn0Oz3TPGZJSq3ztfDj8W4GeRI96XUOWsoMIA7te_7b5svpGkruNF4SOvwEiKv826ZTp_T13wyYDg"),

	)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close(ctx)
	}()
}
