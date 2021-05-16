package utils

import (
	"fmt"
	"log"
	"os"
	"tweeter-sentiment-analyzer/message-broker/typedmsg"
)

func GenerateUuid() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal("ERROR OPENING FILE " + "/dev/urandom")
	}
	b := make([]byte, 16)
	if _, err = f.Read(b); err != nil {
		log.Fatal("ERROR READING FROM FILE: " + "/dev/urandom")
	}
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func Missing(a, b []typedmsg.Topic) (diffs []typedmsg.Topic) {
	// create map with length of the 'a' slice
	ma := make(map[string]struct{}, len(a))

	// Convert first slice to map with empty struct (0 bytes)
	for _, ka := range a {
		ma[ka.Value] = struct{}{}
	}
	// find missing values in a
	for _, kb := range b {
		if _, ok := ma[kb.Value]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}

func ConvertToTopic(strings []string) []typedmsg.Topic {
	var topics []typedmsg.Topic
	for _, v := range strings {
		topics = append(topics, typedmsg.Topic{
			Value:     v,
			IsDurable: true,
		})
	}
	return topics
}
