package mongik

import (
	"fmt"

	"github.com/FrosTiK-SD/mongik/constants"
)

func getKey(collectionName string, operation string, query interface{}, option interface{}) string {
	return fmt.Sprintf("%s | %s | %v | %v", collectionName, constants.DB_FINDONE, query, option)
}
