package main

import (
	"testing"
	"os"
	"log"
	"reflect"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ac *kafka.AdminClient

func TestMain(m *testing.M) {
	var err error
	ac, err = kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Admin client: %s\n", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestFetchTopics(t *testing.T) {
	want := []string{"Test_Topic", "Test_Topic2"}

	got := getTopics(ac)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("ERROR: Expected %s, got %s", want, got)
	}
}

func TestTemplate(t *testing.T) {

	test1 := `resource "kafka_topic" "Test_Topic" {
		name               = "Test_Topic"
		replication_factor = 1
		partitions         = 12
		  
		config = {
				"compression.type" = "producer"
				"retention.ms" = "259200000"
				"segment.bytes" = "1073741824"
				"unclean.leader.election.enable" = "false"
		}
	}
	
	`
	
	test2 := `resource "kafka_topic" "Test_Topic2" {
		name               = "Test_Topic2"
		replication_factor = 1
		partitions         = 21
		  
		config = {
				"compression.type" = "snappy"
				"retention.ms" = "604800000"
				"segment.bytes" = "105314687"
				"unclean.leader.election.enable" = "true"
		}
	}
	
	`
	
	test_topic_1 := getTopicMetadata(ac, "Test_Topic")
	if !reflect.DeepEqual(strings.Fields(test_topic_1), strings.Fields(test1)) {
		t.Errorf("ERROR: Mismatch in Test_Topic2 generated function template and test case\nGenerated:\n%s\nExpected:\n%s", test_topic_1, test1)
	}
	
	test_topic_2 := getTopicMetadata(ac, "Test_Topic2")
	if !reflect.DeepEqual(strings.Fields(test_topic_2), strings.Fields(test2)) {
		t.Errorf("ERROR: Mismatch in Test_Topic2 generated function template and test case\nGenerated:\n%s\nExpected:\n%s", test_topic_2, test2)
	}

	
}