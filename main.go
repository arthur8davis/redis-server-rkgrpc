package main

import (
	"embed"
	"github.com/rookie-ninja/rk-entry/v2/entry"
	"rhyme80/redis-server-rkgrpc/bootstrap"
)

//go:embed boot.yaml
var boot []byte

//go:embed infrastructure/servicegrpc/proto
var docsFS embed.FS

//go:embed infrastructure/servicegrpc/proto
var staticFS embed.FS

func init() {
	rkentry.GlobalAppCtx.AddEmbedFS(rkentry.DocsEntryType, "redisserver", &docsFS)
	rkentry.GlobalAppCtx.AddEmbedFS(rkentry.SWEntryType, "redisserver", &docsFS)
	rkentry.GlobalAppCtx.AddEmbedFS(rkentry.StaticFileHandlerEntryType, "redisserver", &staticFS)
}

func main() {
	bootstrap.Run(boot)
}
