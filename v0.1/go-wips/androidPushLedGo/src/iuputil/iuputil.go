// utility funcs for go-IUP
package iuputil

import (
	"fmt"
	"github.com/grd/iup"
	"reflect"
)

// FetchWidgets populates all the fields of structPtr with IUP entities (ie. from LED file)
// structPtr must be a pointer to a struct of fields of type *iup.Ihandle
// The name used for lookup in IUP will be the fieldname or
// the iupName specified by a field tag : Toto *Ihandle `IUP:"iupName"`
// Note that only exported fields can be written to.
func FetchWidgets(structPtr interface{}) error {

	// Get valueOf structPtr by reflection,
	// check that it is of type Ptr
	valueOf := reflect.ValueOf(structPtr)
	if valueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("parameter is not ptr to struct of *Ihandles")
	}

	// Get value pointed to,
	// check that it is of type Struct
	elem := valueOf.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("parameter is not ptr to struct of *Ihandles")
	}

	typeOfElem := elem.Type()

	// check all fields of struct
	for i := 0; i < elem.NumField(); i++ {
		var ihandlePtr *iup.Ihandle

		// Get info on field
		field := elem.Field(i)
		typeOfField := typeOfElem.Field(i)

		fieldName := typeOfField.Name
		fullFieldName := fmt.Sprintf("%s.%s", typeOfElem, fieldName)

		fmt.Println("--", fullFieldName)

		// check that field is of type *iup.Ihandle
		if field.Type() != reflect.TypeOf(ihandlePtr) {
			return fmt.Errorf("field %q is not of type *Ihandle", fullFieldName)
		}

		// check that we can set/write-to the field
		if !field.CanSet() {
			return fmt.Errorf("field %q is not settable", fullFieldName)
		}

		// lookup field tag IUP:"iupName"
		iupName := typeOfField.Tag.Get("IUP")
		if len(iupName) == 0 {
			// if no IUP field tag, use fieldName to lookup the IUP entity
			iupName = fieldName
		}

		// Lookup IUP entity, error if nil
		ihandlePtr = iup.GetHandle(iupName)
		if ihandlePtr == nil {
			return fmt.Errorf("could not find %q for field %q", iupName, fullFieldName)
		}

		// finally, set the field value = found IUP entity
		field.Set(reflect.ValueOf(ihandlePtr))
		fmt.Println("  ->found widget ", iupName)
	}

	return nil
}
