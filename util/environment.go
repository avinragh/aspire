package util

import (
	"os"
	"strconv"
)

func GetEnvAsIntOrDefault(envVarname string, defaultValue int) int {
	strVal, ok := os.LookupEnv(envVarname)
	if ok {
		intVal, err := strconv.Atoi(strVal)
		if err == nil {
			return intVal
		}
	}
	return defaultValue
}

func GetEnvAsFloat64OrDefault(envVarname string, defaultValue float64) float64 {
	strVal, ok := os.LookupEnv(envVarname)
	if ok {
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func GetEnvOrDefaultString(envVarname, defaultValue string) string {
	strVal, ok := os.LookupEnv(envVarname)
	if ok {
		return strVal
	}
	return defaultValue
}
