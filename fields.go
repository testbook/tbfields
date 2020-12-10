package fields

import (
	"errors"
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

func getAllFields(updateModel interface{}, tagKey string, prefixKey string, fieldsToConsider bson.M) {

	dataValue := reflect.ValueOf(updateModel)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	typeOfCurrentField := dataValue.Type()

	for i := 0; i < dataValue.NumField(); i++ {

		currentField := dataValue.Field(i)
		currentPrefixKey := prefixKey + KeysSeparatorDot

		if prefixKey == EmptyString {
			currentPrefixKey = EmptyString
		}

		//update currentField with respective data type pointed by pointer
		if currentField.Kind() == reflect.Ptr {
			currentField = reflect.Indirect(currentField)
		}

		if currentField.Kind() == reflect.Struct {

			getAllFields(currentField.Interface(), tagKey, currentPrefixKey+typeOfCurrentField.Field(i).Tag.Get(tagKey), fieldsToConsider)

		} else if currentField.Kind() == reflect.Map {

			//todo: implement code for map

		} else {

			key := typeOfCurrentField.Field(i).Tag.Get(tagKey)
			if key != EmptyString && key != HyphenString && currentField.CanInterface() {
				fieldsToConsider[currentPrefixKey+key] = currentField.Interface()
			}
		}
	}
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
