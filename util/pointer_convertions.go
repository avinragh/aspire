package util

import "github.com/go-openapi/strfmt"

func GetStringPointer(valueString string) *string {
	return &valueString
}

func GetInt64Pointer(valueInt int64) *int64 {
	return &valueInt
}

func GetDateTimePointer(valueDateTime strfmt.DateTime) *strfmt.DateTime {
	return &valueDateTime
}

func GetFloat64Pointer(valueFloat float64) *float64 {
	return &valueFloat
}
