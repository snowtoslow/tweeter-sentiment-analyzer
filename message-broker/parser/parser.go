package parser

import (
	"fmt"
	"log"
	"strings"
)

func stringToCommand(input string) (jsonBytes []byte, err error) {
	splits := strings.Split(strings.TrimSpace(input), " ")
	if len(splits) > 2 {
		return nil, fmt.Errorf("string len after spliting doesn't feets in requirements for parsing: {command: {value: splits[0]}, opic: {value: splits[1]}")
	}

	jsonBytes = []byte(fmt.Sprintf(`{"command": {"value" : "%s"},"topic": {"value": "%s"}}`, splits[0], splits[1]))
	log.Println(string(jsonBytes))
	return
}

type Subscribe struct {
	Topics Topics `json:"topics"`
}

type Unsubscribe struct {
	Topics `json:"topics"`
}

type Topics []string
