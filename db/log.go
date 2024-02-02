package mongik

import (
	"fmt"
	"os"

	"github.com/FrosTiK-SD/mongik/constants"
)

func CacheLog(log string) {
	if os.Getenv(constants.MONGIK_DEBUG) == "1" {
		fmt.Println(log)
	}
}
