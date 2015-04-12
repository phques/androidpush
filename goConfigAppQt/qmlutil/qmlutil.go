// utility funcs for go-qml
package qmlutil

import (
	"fmt"
	"gopkg.in/qml.v1"
	"reflect"
)

// FetchObjects populates all the fields of structPtr with Qt entities
// structPtr must be a pointer to a struct with fields of type qml.Object
// The name used for lookup in QML will be the fieldname or
// the qmlName specified by a field tag : Toto qml.Object `QML:"zeToto"`
// Note that only exported fields can be written to.
func FetchObjects(structPtr interface{}, parentQML qml.Common) error {

	// Get valueOf structPtr by reflection,
	// check that it is of type Ptr
	valueOf := reflect.ValueOf(structPtr)
	if valueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("FetchObjects, parameter structPtr is not ptr to struct")
	}

	// Get value pointed to,
	// check that it is of type Struct
	structValue := valueOf.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("FetchObjects, parameter structPtr is not ptr to struct")
	}

	structType := structValue.Type()

	// Since qml.Object is an interface, getting a TypeOf is a bit hard...
	// trick is to use a pointer to qml.Object, then use Elem() on that
	var qmlObj qml.Object
	qmlObjType := reflect.TypeOf(&qmlObj).Elem()

	// check all fields of struct
	for i := 0; i < structValue.NumField(); i++ {
		// Get info on field
		field := structValue.Field(i)
		structField := structType.Field(i)

		fieldName := structField.Name
		fullFieldName := fmt.Sprintf("%s.%s", structType, fieldName)

		fmt.Println("--", fullFieldName)

		// check that field is of type qml.Object
		if field.Type() != qmlObjType {
			return fmt.Errorf("FetchObjects, field %q is not of type qml.Object", fullFieldName)
		}

		// check that we can set/write-to the field
		if !field.CanSet() {
			return fmt.Errorf("FetchObjects, field %q is not settable", fullFieldName)
		}

		// lookup field tag QML:"qmlName"
		qmlName := structField.Tag.Get("QML")
		if len(qmlName) == 0 {
			// if no QML field tag, use fieldName to lookup the QML entity
			qmlName = fieldName
		}

		// Lookup QML entity, error if nil
		qmlObj = parentQML.ObjectByName(qmlName)
		if qmlObj == nil {
			return fmt.Errorf("FetchObjects, could not find %q for field %q", qmlName, fullFieldName)
		}

		// finally, set the field value = found QML entity
		field.Set(reflect.ValueOf(qmlObj))
		fmt.Println("  ->found control ", qmlName)
	}

	return nil
}
