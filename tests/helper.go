package tests

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func isRabbitmqAdminLocalAvailable() bool {
	cmd := exec.Command("rabbitmqadmin", "--help")
	return cmd.Run() == nil
}

func getOSEnvOrDefaultStr(varName string, defaultValue string) string {
	v := os.Getenv(varName)
	if v == "" {
		return defaultValue
	} else {
		return v
	}
}

func getOSEnvOrDefaultInt(varName string, defaultValue int) int {
	v := os.Getenv(varName)
	if v == "" {
		return defaultValue
	} else {
		intVal, err := strconv.Atoi(v)
		if err != nil {
			fmt.Printf("env var doesn't contain a number value: var=%s, value=%s\n", varName, v)
			return defaultValue
		} else {
			return intVal
		}
	}
}

func getConnParamsToUse() (string, string, string, int) {
	user := getOSEnvOrDefaultStr("RABBIT_USER", "guest")
	password := getOSEnvOrDefaultStr("RABBIT_PASSWORD", "guest")
	host := getOSEnvOrDefaultStr("RABBIT_SERVER", "localhost")
	port := getOSEnvOrDefaultInt("RABBIT_PORT", 5672)
	return user, password, host, port
}
