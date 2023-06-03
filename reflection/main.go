package main

import (
	"fmt"
	"reflect"
)

type Employee struct {
	Name  string
	Ecode int
}

func main() {

	emp := []Employee{
		{
			Name:  "abhi",
			Ecode: 111,
		},
		{
			Name:  "lash",
			Ecode: 112},
	}
	returnTypAndVal(emp)
	returnTypAndVal(10)
	returnTypAndVal("abhi")
}

func returnTypAndVal(dataStructure interface{}) {

	typ := reflect.TypeOf(dataStructure)
	val := reflect.ValueOf(dataStructure)
	fmt.Println("val", val)

	switch typ.Kind() {
	case reflect.Slice:
		if val.Len() > 0 && val.Index(0).Kind() == reflect.Struct {
			for i := 0; i < val.Len(); i++ {
				structVal := val.Index(i)
				structTyp := structVal.Type()

				for j := 0; j < structTyp.NumField(); j++ {
					field := structTyp.Field(j)
					fieldVal := structVal.Field(j)
					fmt.Printf("%s = %v\n", field.Name, fieldVal)
				}
			}
		}
	case reflect.Int:
		fmt.Printf("%s = %v\n", typ, val)
	case reflect.Struct:

		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			fieldVal := val.Field(j)
			fmt.Printf("%s = %v\n", field.Name, fieldVal)
		}
	default:
		fmt.Println("hi")
	}
}
