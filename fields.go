package fields

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

func getTagKey(tagKeys []string) (string, error) {

	tagKey := BsonTagKey

	if len(tagKeys) > 1 {
		return tagKey, errors.New("only 1 tag key is allowed")
	}

	if len(tagKeys) > 0 {
		tagKey = strings.ToLower(strings.TrimSpace(tagKeys[0]))
		if tagKey != BsonTagKey && tagKey != JsonTagKey {
			return tagKey, errors.New("invalid tag key")
		}
	}

	return tagKey, nil
}

//since golang has 2 objects type
//struct and map hence both the recursive calls can be interlinked or one type has to be changed to another
//here we are changing struct to map because other way round is not possible and to avoid interlinked recursive calls
//struct can have map inside it
//map can have struct inside it
func getAllFields(updateModel interface{}, tagKey string) bson.M {

	fields := bson.M{}

	currentField := reflect.ValueOf(updateModel)

	//update currentField with respective data type pointed by pointer
	//since pointer itself can have multiple layer of pointers to be pointed
	for currentField.Kind() == reflect.Ptr {
		currentField = currentField.Elem()
	}

	if currentField.Kind() == reflect.Map {

		for _, key := range currentField.MapKeys() {

			currentValue := currentField.MapIndex(key)

			//invalid to have other ds in keys except string
			if key.Kind() != reflect.String {
				continue
			}

			if currentValue.CanInterface() {

				newMap := getAllFields(currentValue.Interface(), tagKey)

				newKeyPrefix := key.Interface().(string)

				for k, v := range newMap {
					fields[fmt.Sprintf(KeysFormatter, newKeyPrefix, k)] = v
				}

			}
		}

	} else if currentField.Kind() == reflect.Struct {

		newMap, _ := translateStructToMap(currentField, tagKey)

		newKeyPrefix := key.Interface().(string)

		for k, v := range newMap {
			fields[fmt.Sprintf(KeysFormatter, newKeyPrefix, k)] = v
		}
	}

	return fields

	//if dataValue.Kind() == reflect.Struct {
	//
	//	for i := 0; i < dataValue.NumField(); i++ {
	//
	//		currentField := dataValue.Field(i)
	//		currentPrefixKey := prefixKey + KeysSeparatorDot
	//
	//		if prefixKey == EmptyString {
	//			currentPrefixKey = EmptyString
	//		}
	//
	//		//update currentField with respective data type pointed by pointer
	//		//since pointer itself can have multiple layer of pointers to be pointed
	//		for currentField.Kind() == reflect.Ptr {
	//			currentField = reflect.Indirect(currentField)
	//		}
	//
	//		if currentField.Kind() == reflect.Struct {
	//
	//			getAllFields(currentField.Interface(), tagKey, currentPrefixKey+typeOfCurrentField.Field(i).Tag.Get(tagKey), fieldsToConsider)
	//
	//		} else if currentField.Kind() == reflect.Map {
	//
	//		} else {
	//
	//			key := typeOfCurrentField.Field(i).Tag.Get(tagKey)
	//			if key != EmptyString && key != HyphenString && currentField.CanInterface() && !isValueNonAcceptable(currentField.Kind()) {
	//				fieldsToConsider[currentPrefixKey+key] = currentField.Interface()
	//			}
	//		}
	//	}
	//
	//} else {
	//
	//}
}

//allow only primitive type for values
func isValueNonAcceptable(kind reflect.Kind) bool {
	if kind == reflect.Map || kind == reflect.Interface || kind == reflect.Uintptr || kind == reflect.Ptr || kind == reflect.Chan ||
		kind == reflect.Func || kind == reflect.Struct {
		return true
	}

	return false
}

func translateStructToMap(structData reflect.Value, tagKey string) (bson.M, error) {

	if structData.Kind() != reflect.Struct {
		return nil, errors.New("non struct is passed for conversion")
	}

	newMap := bson.M{}

	typeOfCurrentField := structData.Type()

	for i := 0; i < structData.NumField(); i++ {

		currentField := structData.Field(i)

		//pointer to pointer
		for currentField.Kind() == reflect.Ptr {
			currentField = reflect.Indirect(currentField)
		}

		key := typeOfCurrentField.Field(i).Tag.Get(tagKey)
		if currentField.Kind() == reflect.Struct {

			nestedMap, err := translateStructToMap(currentField, tagKey)
			if err != nil {
				return nil, err
			}

			for k, v := range nestedMap {
				newKey := fmt.Sprintf(KeysFormatter, key, k)
				newMap[newKey] = v
			}

		}

		if key != EmptyString && key != HyphenString && currentField.CanInterface() && !isValueNonAcceptable(currentField.Kind()) {
			newMap[key] = currentField.Interface()
		}

	}

	return newMap, nil
}

func filterFields(allFields bson.M, fields []string) (bson.M, []string) {

	filteredFields := bson.M{}
	var fieldsNotFound []string
	for _, key := range fields {

		if value, ok := allFields[key]; ok {
			filteredFields[key] = value
		} else {
			fieldsNotFound = append(fieldsNotFound, key)
		}
	}

	return filteredFields, fieldsNotFound
}

func GetFields(updateModel interface{}, fields []string, tagKeys ...string) (filteredFields bson.M, fieldsNotFound []string, err error) {

	tagKey, err := getTagKey(tagKeys)
	if err != nil {
		return nil, nil, err
	}

	allFields := bson.M{}
	getAllFields(updateModel, tagKey, EmptyString, allFields)

	filteredFields, fieldsNotFound = filterFields(allFields, fields)

	return filteredFields, fieldsNotFound, nil
}
