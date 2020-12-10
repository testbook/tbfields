package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFields(t *testing.T) {

	type Robot struct {
		ID       int    `bson:"ID"`
		Name     string `bson:"Namesss"`
		LastName string `bson:"-"`
	}

	robot := Robot{ID: 9999, Name: "testName"}

	assert := assert.New(t)

	filterFields, fieldsNotFound, err := GetFields(&robot, []string{"ID", "Namesss", "LastName"}, BsonTagKey)
	if err != nil {
		assert.Error(err)
	}

	assert.NotNil(fieldsNotFound)
	assert.Len(fieldsNotFound, 1) //lastName will not be found

	assert.NotNil(filterFields)
	assert.Len(filterFields, 2)
	assert.Equal(filterFields["ID"], 9999)
	assert.Equal(filterFields["Namesss"], "testName")
}
