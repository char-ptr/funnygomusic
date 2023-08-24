package databaser

import (
	"github.com/meilisearch/meilisearch-go"
	"os"
)

func NewMeili() *meilisearch.Client {
	mker, _ := os.LookupEnv("MEILI_MASTER_KEY")
	return meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://meili:7700",
		APIKey: mker,
	})
}
