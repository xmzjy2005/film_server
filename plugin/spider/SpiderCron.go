package spider

import (
	"film_server/config"
	"film_server/model/system"
)

func ClearCache() {
	system.RemoveCache(config.IndexCacheKey)
}
