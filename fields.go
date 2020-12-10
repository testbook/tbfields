package main

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

const (
	BsonTagKey = "bson"
	JsonTagKey = "json"
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

func getAllFields(updateModel interface{}, tagKey string) bson.M {
	fieldsToConsider := bson.M{}

	dataValue := reflect.ValueOf(updateModel).Elem()
	typeOfCurrentField := dataValue.Type()

	for i := dataValue.NumField() - 1; i >= 0; i-- {

		currentField := dataValue.Field(i)

		if data, ok := typeOfCurrentField.Field(i).Tag.Lookup(tagKey); ok && currentField.CanInterface() {
			fieldsToConsider[data] = currentField.Interface()
		}
	}

	return fieldsToConsider
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

	allFields := getAllFields(updateModel, tagKey)

	filteredFields, fieldsNotFound = filterFields(allFields, fields)

	return filteredFields, fieldsNotFound, nil
}
