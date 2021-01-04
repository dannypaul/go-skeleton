package primitive

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Id string

func (id Id) String() string {
	return string(id)
}

func (id Id) Equals(anotherId Id) bool {
	return string(id) == string(anotherId)
}

func (id Id) MarshalBSONValue() (bsontype.Type, []byte, error) {
	hex := id
	if len(id.String()) == 0 {
		hex = NewObjectId()
	}
	objectId, err := primitive.ObjectIDFromHex(hex.String())
	if err != nil {
		return bsontype.ObjectID, []byte{}, fmt.Errorf("'%s' is not an ObjectID", id.String())
	}
	return bsonx.ObjectID(objectId).MarshalBSONValue()
}

func (id *Id) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	if t != bsontype.ObjectID || len(raw) != 12 {
		return fmt.Errorf("unable to unmarshal ObjectID. bsontype: %v, length: %v", t, len(raw))
	}

	val := bsonx.Null()
	err := val.UnmarshalBSONValue(bsontype.ObjectID, raw)
	if err != nil {
		return fmt.Errorf("invalid ObjectID from source: %v", err)
	}

	*id = Id(val.ObjectID().Hex())
	return nil
}

type IdList []Id

func (i IdList) Add(id Id) []Id {
	for _, x := range i {
		if x == id {
			return i
		}
	}
	return append(i, id)
}

func (i IdList) Remove(id Id) []Id {
	filtered := make([]Id, 0)
	for _, x := range i {
		if x != id {
			filtered = append(filtered, x)
		}
	}
	return filtered
}

func NewObjectId() Id {
	return Id(primitive.NewObjectID().Hex())
}
