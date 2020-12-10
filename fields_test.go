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
		}, BsonTagKey)
	if err != nil {
		assert.Error(err)
	}

	assert.NotNil(fieldsNotFound)
	assert.Len(fieldsNotFound, 1) //lastName will not be found

	assert.NotNil(filterFields)
	assert.Len(filterFields, 6)
	assert.Equal(filterFields["id"], 9999)
	assert.Equal(filterFields["name"], "testName")
	assert.Equal(filterFields["hands.leftHand"], "LH")
	assert.Equal(filterFields["hands.rightHand"], "RH")
	assert.Equal(filterFields["hands.robotFingers.middleFinger"], "MF")
	assert.Equal(filterFields["hands.robotFingers.firstFinger"], "FF")
}
