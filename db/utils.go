package mongik

import (
	"fmt"
	"reflect"
)

func getKey[Option any](collectionName string, operation string, query interface{}, option []*Option) string {
	optionKey := ""
	for _, opt := range option {
		optionKey += iterateStructFields(*opt)
	}
	return fmt.Sprintf("%s | %s | %v | %v", collectionName, operation, query, optionKey)
}

func iterateStructFields(input interface{}) string {
	structKey := "{ "
	value := reflect.ValueOf(input)
	numFields := value.NumField()
	structType := value.Type()
	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		fieldValue := reflect.Indirect(value.Field(i))
		if fieldValue.IsValid() == true && fieldValue.IsZero() == false {
			structKey += fmt.Sprintf("%s: %v, ", field.Name, fieldValue)
		}
	}
	return structKey + " }"
}
