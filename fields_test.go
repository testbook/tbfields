package fields

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFields(t *testing.T) {

	type Robot struct {
		ID       int    `bson:"id"`
		Name     string `bson:"name"`
		LastName string `bson:"-"`
	}

	robot := Robot{ID: 9999, Name: "testName"}

	assert := assert.New(t)

	filterFields, fieldsNotFound, err := GetFields(&robot, []string{"id", "name", "LastName"}, BsonTagKey)
	if err != nil {
		assert.Error(err)
	}

	assert.NotNil(fieldsNotFound)
	assert.Len(fieldsNotFound, 1) //lastName will not be found

	assert.NotNil(filterFields)
	assert.Len(filterFields, 2)
	assert.Equal(filterFields["id"], 9999)
	assert.Equal(filterFields["name"], "testName")
}

func TestGetFieldsWithMapAndNestedObjects(t *testing.T) {
	type Robot struct {
		ID         int    `bson:"id"`
		Name       string `bson:"name"`
		LastName   string `bson:"-"`
		MiddleName string //bson tag is missing
	}

	robot := Robot{ID: 9999, Name: "testName", MiddleName: "mymiddleName"}

	newMap := map[string]interface{}{}
	newMap["robotKey"] = robot

	assert := assert.New(t)

	filterFields, fieldsNotFound, err := GetFields(newMap, []string{"robotKey.id", "robotKey.name", "LastName", "robotKey.MiddleName"}, BsonTagKey)
	if err != nil {
		assert.Error(err)
	}

	assert.NotNil(fieldsNotFound)
	assert.Len(fieldsNotFound, 1) //lastName will not be found

	assert.NotNil(filterFields)
	assert.Len(filterFields, 3)
	assert.Equal(filterFields["robotKey.id"], 9999)
	assert.Equal(filterFields["robotKey.name"], "testName")
	assert.Equal(filterFields["robotKey.MiddleName"], "mymiddleName")
}

func TestGetFieldsWithNestedObjects(t *testing.T) {

	type RobotFingers struct {
		MiddleFinger string                      `bson:"middleFinger"`
		FirstFinger  string                      `bson:"firstFinger"`
		Properties   map[interface{}]interface{} `bson:"properties"`
	}

	type RobotHand struct {
		LeftHand  string        `bson:"leftHand"`
		RightHand string        `bson:"rightHand"`
		Fingers   *RobotFingers `bson:"robotFingers"`
	}

	type Robot struct {
		ID       int       `bson:"id"`
		Name     string    `bson:"name"`
		LastName string    `bson:"-"`
		Hands    RobotHand `bson:"hands"`
	}

	robot := Robot{ID: 9999, Name: "testName"}
	robot.Hands = RobotHand{
		LeftHand:  "LH",
		RightHand: "RH",
		Fingers: &RobotFingers{
			MiddleFinger: "MF",
			FirstFinger:  "FF",
			Properties: map[interface{}]interface{}{
				"kachra": "seth",
				"stock":  "market",
				"100":    10000,
			},
		},
	}

	assert := assert.New(t)

	filterFields, fieldsNotFound, err := GetFields(&robot,
		[]string{
			"id", "name", "lastName", "hands.leftHand", "hands.rightHand", "hands.robotFingers.middleFinger", "hands.robotFingers.firstFinger",
			"hands.robotFingers.properties.kachra", "hands.robotFingers.properties.stock",
		}, BsonTagKey)
	if err != nil {
		assert.Error(err)
	}

	assert.NotNil(fieldsNotFound)
	assert.Len(fieldsNotFound, 1) //lastName will not be found

	assert.NotNil(filterFields)
	assert.Len(filterFields, 8)
	assert.Equal(filterFields["id"], 9999)
	assert.Equal(filterFields["name"], "testName")
	assert.Equal(filterFields["hands.leftHand"], "LH")
	assert.Equal(filterFields["hands.rightHand"], "RH")
	assert.Equal(filterFields["hands.robotFingers.middleFinger"], "MF")
	assert.Equal(filterFields["hands.robotFingers.firstFinger"], "FF")
	assert.Equal(filterFields["hands.robotFingers.properties.kachra"], "seth")
	assert.Equal(filterFields["hands.robotFingers.properties.stock"], "market")
	//assert.Equal(filterFields["hands.robotFingers.properties.100"], 10000)
}
