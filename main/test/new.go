package test

import (
	"Scanner/main/format"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func isUtilizeSuccess(pocStruct format.PocStruct, resp *http.Response) bool {
	var elements []map[string]bool

	for range pocStruct.ResponsePackage.Group {
		elements = append(elements, map[string]bool{
			"User-Agent":      false,
			"Accept-Encoding": false,
			"Accept-Language": false,
			"Cookie":          false,
			"Accept":          false,
			"connection":      false,
			"Content-Type":    false,
			"Body":            false,
			"Status":          false,
			"flag":            false,
		})
	}

	for currentTimes, currentGroup := range pocStruct.ResponsePackage.Group {
		groupReflect := reflect.ValueOf(currentGroup)
		typ := groupReflect.Type()

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldName := field.Name
			fieldValue := groupReflect.FieldByName(fieldName)
			fieldValueInterface := fieldValue.Interface()

			if fieldValue.String() != "" {
				if containValue, ok := fieldValueInterface.(format.Contain); ok {
					v := reflect.ValueOf(containValue)
					t := v.Type()

					for tmp := 0; tmp < t.NumField(); tmp++ {
						currentKey := t.Field(tmp).Name
						expectedValue := reflect.ValueOf(pocStruct.ResponsePackage.Group[currentTimes].Contain.Header).FieldByName(currentKey).String()
						actualValue := resp.Header.Get(strings.ToLower(currentKey))
						elements[currentTimes][currentKey] = expectedValue == actualValue
						elements[currentTimes]["flag"] = elements[currentTimes]["flag"] || elements[currentTimes][currentKey]
					}

				} else if fieldName == "Status" && strings.Contains(resp.Status, fieldValue.String()) {
					elements[currentTimes]["Status"] = true
					elements[currentTimes]["flag"] = true
				} else {
					body, err := ioutil.ReadAll(resp.Body)
					if len(body) != 0 {
						requestBody = body
					}
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					elements[currentTimes]["Body"] = ContainsAny(string(requestBody), pocStruct.ResponsePackage.Group[currentTimes].Contain.Body)
					elements[currentTimes]["flag"] = elements[currentTimes]["flag"] || elements[currentTimes]["Body"]
				}
			}
		}
	}

	isConform := pocStruct.ResponsePackage.Operation == "AND"
	for _, el := range elements {
		if !el["flag"] {
			isConform = !isConform
			break
		}
	}

	if isConform {
		fmt.Println("Success!")
	} else {
		fmt.Println("Error!")
	}

	return false
}