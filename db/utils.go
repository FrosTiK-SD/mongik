package mongik

import (
	"fmt"

	"github.com/FrosTiK-SD/mongik/constants"
)


func getKey(collectionName string, operation string, query interface{}) string {
	return fmt.Sprintf("%s | %s | %v", collectionName, constants.DB_FINDONE, query)
}