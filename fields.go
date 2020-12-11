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
func getAllFields(updateModel interface{}, tagKey string, prefixKey string) bson.M {

	fields := bson.M{}

	currentField := reflect.ValueOf(updateModel)

	//update currentField with respective data type pointed by pointer
	//since pointer itself can have multiple layer of pointers to be pointed
	for isPointer(currentField.Kind()) {
		currentField = currentField.Elem()
	}

	if isStruct(currentField.Kind()) {

		newMap := translateStructToMap(currentField, tagKey)
		currentField = reflect.ValueOf(newMap)

	} else if !isMap(currentField.Kind()) {
		//invalid type for an object. this block should never execute
		return fields
	}

	//current field will always be a type of map
	for _, key := range currentField.MapKeys() {

		currentValue := currentField.MapIndex(key)

		//invalid to have other ds in keys except string
		if key.Kind() != reflect.String {
			continue
		}

		if !currentValue.CanInterface() {
			continue
		}

		//currentValue is giving interface{}  here instead of primitive type
		fmt.Println(currentValue.Kind())

		if isPrimitive(currentValue.Kind()) || currentValue.Kind() == reflect.Interface {

			//add primitive data directly in
			newKeyPrefix := key.Interface().(string)
			fields[newKeyPrefix] = currentValue.Interface()

		} else if isMap(currentValue.Kind()) || isStruct(currentValue.Kind()) || isPointer(currentValue.Kind()) {
			newMap := getAllFields(currentValue.Interface(), tagKey, "")

			newKeyPrefix := key.Interface().(string)

			for k, v := range newMap {
				fields[fmt.Sprintf(KeysFormatter, newKeyPrefix, k)] = v
			}
		} else {
			//ignore current type like func or chan
		}
	}

	return fields
}

func isPrimitive(kind reflect.Kind) bool {

	switch kind {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.String:
		return true
	default:
		return false
	}

	return false
}

func isMap(kind reflect.Kind) bool {
	return kind == reflect.Map
}

func isStruct(kind reflect.Kind) bool {
	return kind == reflect.Struct
}

func isPointer(kind reflect.Kind) bool {
	return kind == reflect.Ptr
}

func translateStructToMap(structData reflect.Value, tagKey string) interface{} {

	newMap := bson.M{}

	typeOfCurrentField := structData.Type()

	for i := 0; i < structData.NumField(); i++ {

		currentField := structData.Field(i)

		//pointer to pointer
		for currentField.Kind() == reflect.Ptr {
			currentField = reflect.Indirect(currentField)
		}

		key := typeOfCurrentField.Field(i).Tag.Get(tagKey)

		//truncate the object tree if key is hyphen
		if key != HyphenString && currentField.CanInterface() {

			//if key is empty then roll back to variable name
			//todo: verify below method works properly
			if key == EmptyString {
				key = typeOfCurrentField.Name()
			}

			newMap[key] = currentField.Interface()
		}

	}
	return newMap
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

	allFields := getAllFields(updateModel, tagKey, EmptyString)

	filteredFields, fieldsNotFound = filterFields(allFields, fields)

	return filteredFields, fieldsNotFound, nil
}
