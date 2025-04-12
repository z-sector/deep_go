package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const (
	keyTag        = "properties"
	emptyTagValue = "omitempty"
)

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	personType := reflect.TypeOf(person)
	personValue := reflect.ValueOf(person)

	sb := strings.Builder{}
	for i := 0; i < personType.NumField(); i++ {
		field := personType.Field(i)
		tagValue := field.Tag.Get(keyTag)
		if tagValue == "" {
			continue
		}

		fieldValue := personValue.Field(i)
		fieldName, omitEmpty := parseTagValue(tagValue)
		if fieldName == "" || fieldValue.IsZero() && omitEmpty {
			continue
		}

		var strValue string
		switch fieldValue.Kind() {
		case reflect.String:
			strValue = fieldValue.String()
		case reflect.Int:
			strValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Bool:
			strValue = strconv.FormatBool(fieldValue.Bool())
		default:
			strValue = fmt.Sprintf("%v", fieldValue)
		}

		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fieldName)
		sb.WriteByte('=')
		sb.WriteString(strValue)
	}

	return sb.String()
}

func parseTagValue(tagValue string) (fieldName string, omitEmpty bool) {
	for _, part := range strings.Split(tagValue, ",") {
		str := strings.TrimSpace(part)

		if str == "" {
			continue
		}

		if str == emptyTagValue {
			omitEmpty = true
			continue
		}

		fieldName = str
	}
	return
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
