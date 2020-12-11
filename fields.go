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

	input := reflect.ValueOf(updateModel)

	//update currentField with respective data type pointed by pointer
	//since pointer itself can have multiple layer of pointers to be pointed
	//also any type implements interface hence touch down from deepest interface{} to concrete type
	for isPointer(input.Kind()) || isInterface(input.Kind()) {
		input = input.Elem()
	}

	if isStruct(input.Kind()) {

		newMap := translateStructToMap(input, tagKey)
		input = reflect.ValueOf(newMap)

	} else if !isMap(input.Kind()) {
		//invalid type for an object. this block should never execute
		return fields
	}

	//current field will always be a type of map
	for _, key := range input.MapKeys() {

		currentValue := input.MapIndex(key)

		//inn case of map input itself keys will be of type interface instead if string
		for isPointer(key.Kind()) || isInterface(key.Kind()) {
			key = key.Elem()
		}

		//invalid to have other ds in keys except string
		if key.Kind() != reflect.String {
			fmt.Println(key.Kind())
			continue
		}

		if !currentValue.CanInterface() {
			continue
		}

		for isPointer(currentValue.Kind()) || isInterface(currentValue.Kind()) {
			currentValue = currentValue.Elem()
		}

		if isPrimitive(currentValue.Kind()) {

			//add primitive data directly in
			newKeyPrefix := key.Interface().(string)
			fields[newKeyPrefix] = currentValue.Interface()

		} else if isMap(currentValue.Kind()) || isStruct(currentValue.Kind()) {

			newMap := getAllFields(currentValue.Interface(), tagKey)

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

	if reflect.Bool == kind ||
		reflect.Int == kind ||
		reflect.Int8 == kind ||
		reflect.Int16 == kind ||
		reflect.Int32 == kind ||
		reflect.Int64 == kind ||
		reflect.Uint == kind ||
		reflect.Uint8 == kind ||
		reflect.Uint16 == kind ||
		reflect.Uint32 == kind ||
		reflect.Uint64 == kind ||
		reflect.Float32 == kind ||
		reflect.Float64 == kind ||
		reflect.String == kind {

		return true
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

func isInterface(kind reflect.Kind) bool {
	return kind == reflect.Interface
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

			//if bson or json tag is not present then key is empty hence roll back to variable name
			if key == EmptyString {
				key = typeOfCurrentField.Field(i).Name
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

	allFields := getAllFields(updateModel, tagKey)

	filteredFields, fieldsNotFound = filterFields(allFields, fields)

	return filteredFields, fieldsNotFound, nil
}
