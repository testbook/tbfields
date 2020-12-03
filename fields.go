package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

/*
using tbjson

1. use projection like in response writer
2. ask for projection from api and data model
3. use field object from tbbson and get updated value (string)
4. using updated value (string) create a map of map[string]interface{}
5. return it to model layer


in tbbson there is a support for only int, float64, obj does is not support arrays or boolean? talking about new field struct added by vinod
it must be working else how is entire array in newtb working?
func (f *Fields) Populate(fields map[string]interface{}) error {
	keys := make(map[string]int8, len(fields))
	subFields := make(map[string]Fields)
	for key, val := range fields {
		keys[key] = 1
		switch val.(type) {
		case int, float64:
		case map[string]interface{}:
			subField := new(Fields)
			err := subField.Populate(val.(map[string]interface{}))
			if err != nil {
				return err
			}
			subFields[key] = *subField
		default:
			return fmt.Errorf("expected values in the map are int or map[string]interface{}, found: type of %s to be %T.\n", key, val)
		}
	}
	f.Keys = keys
	f.SubFields = subFields
	return nil
}

using normal code with type check || if somehow we are able to skip type checking then lets go ahead with it
1. Create a map of all tag to value
2. Create updated map for storing values to be updated
3. Iterate over field array to include values from map to updated map (perform data type check? no need to perform data type check? if we can get the value from struct)
4. return updated map
*/

func GetFeildsBson(updateModel interface{}, fields []string) (bson.M, error) {

	return bson.M{}, nil
}

func f(data interface{}) {

	//structObj := reflect.ValueOf(data).Interface()
	//reflectionBitch := reflect.TypeOf(structObj)
	//
	//fmt.Println(reflectionBitch.Name())
	//
	//firstElement := reflectionBitch.Field(0)
	//
	//fmt.Println(firstElement.Tag.Get("bson"))

	//get the original struct type
	dataValue := reflect.ValueOf(data).Elem()

	for i := dataValue.NumField() - 1; i >= 0; i-- {

		currentField := dataValue.Field(i)

		//fmt.Println("fields: ", currentField.String())
		//fmt.Println("fields: ", temp.)

		typeOfCurrentField := dataValue.Type()

		fmt.Println("bsonTag ", typeOfCurrentField.Field(i).Tag.Get("bson"))
		fmt.Println("type ", currentField.Type())
		fmt.Println("data ", currentField.Interface())

		//fmt.Println(typeOfCurrentField)
		//
		//fmt.Println(currentField.Type())
		//fmt.Println(currentField.Interface())
		//
		//t := reflect.TypeOf(currentField.Interface())
		//f1, _ := t.FieldByName("Name")
		//fmt.Println(f1.Tag) // f one
	}

	v := reflect.ValueOf(data).Elem().FieldByName("ID")
	//fmt.Println("fields: ", reflect.ValueOf(data).Elem().NumField())
	ptr := v.Addr().Interface().(*int)
	*ptr = 100

}

type robot struct {
	ID   int    `bson:"ID"`
	Name string `bson:"Namesss"`
}

func main() {

	robot := robot{ID: 69, Name: "testies"}

	// this is unnecessary, and is equivalent to f(&robot)
	//var iface interface{} =
	f(&robot)
	//fmt.Println(robot.ID) //I want to get here 100
}
