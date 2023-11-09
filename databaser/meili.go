package databaser

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeili() *meilisearch.Client {
	mker, _ := os.LookupEnv("MEILI_MASTER_KEY")
	meiliHost, _ := os.LookupEnv("MEILI_HOST")
	return meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meiliHost,
		APIKey: mker,
	})
}
