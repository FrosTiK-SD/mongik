package mongik

import (
	"fmt"

	mongik "github.com/FrosTiK-SD/mongik/models"
)

func CacheLog(mongikClient *mongik.Mongik, log string) {
	if mongikClient.Config.Debug == true {
		fmt.Println(log)
	}
}
