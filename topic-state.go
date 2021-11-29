// This binary is to generate the provided kafka cluster's topic state and translate all topics
// into a terraform kafka_topic resource for use with the terraform-provider-kafka module
//
// Author: Yared Mekuria<Yared.Mekuria@bettercloud.com>

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var GitCommit string

//Struct to provide to golang template
type tfTopic struct {
	Name       string
	ReplFactor int
	Partitions int
	DynConfig  map[string]string
}

//Function to sort all topics in alphabetical order
func sortTopics(s []string) []string {
	var r []string

	sort.Strings(s)

	for _, str := range s {
		//Check if string is empty or starts with "__".
		if str != "" && !(len(str) >= len("__") && str[0:len("__")] == "__") {
			r = append(r, str)
		}
	}
	return r
}

func ReplFactor(c *kafka.TopicMetadata) int {
	var r int

	//Find replication factor of topic by counting number of partitions and assigning greatest.
	//This is literally how Kafka provided tools do it internally
	for _, t := range c.Partitions {
		if r < len(t.Replicas) {
			r = len(t.Replicas)
		}
	}
	return r
}

//Golang template and rendering, provided a pointer to the populated struct. Outputs to stdout.
func tmpl(data *tfTopic) string {
	var tpl bytes.Buffer

	t, _ := template.New("").Parse(`resource "kafka_topic" "{{.Name}}" {
	name               = "{{.Name}}"
	replication_factor = {{.ReplFactor}}
	partitions         = {{.Partitions}}
	  
	config = {
		{{- range $k, $v := .DynConfig }}
		"{{ $k }}" = "{{ $v }}"
		{{- end }}
	}
}

`)

	t.Execute(&tpl, data)
	return tpl.String()
}

func getTopics(ac *kafka.AdminClient) []string {
	//Get list of all topics and number of partitions
	topicList, err := ac.GetMetadata(nil, true, 20000)
	if err != nil {
		fmt.Printf("Failed to GetMetadata: %s\n", err)
		os.Exit(1)
	}

	t := make([]string, len(topicList.Topics))

	for _, res := range topicList.Topics {
		t = append(t, res.Topic)
	}
	topics := sortTopics(t)

	return topics
}

func getTopicMetadata(ac *kafka.AdminClient, t string) string {
	var topicData *tfTopic

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dur, _ := time.ParseDuration("20s")

	resourceType, _ := kafka.ResourceTypeFromString("Topic")

	results, err := ac.DescribeConfigs(ctx,
		[]kafka.ConfigResource{{Type: resourceType, Name: t}},
		kafka.SetAdminRequestTimeout(dur))

	if err != nil {
		fmt.Printf("Failed to DescribeConfigs(%s, %s): %s\n",
			resourceType, t, err)
		os.Exit(1)
	}

	for _, result := range results {

		config := make(map[string]string)

		topicMeta, _ := ac.GetMetadata(&result.Name, false, 20000)
		topicInfo := topicMeta.Topics[result.Name]

		for _, entry := range result.Config {
			if int(entry.Source) == 1 {
				config[entry.Name] = entry.Value
			}
		}

		topicData = &tfTopic{
			Name:       result.Name,
			ReplFactor: ReplFactor(&topicInfo),
			Partitions: len(topicInfo.Partitions),
			DynConfig:  config,
		}
	}
	return tmpl(topicData)
}

func main() {

	//Error handling in case an argument isnt provided.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr,
			"Kafka Topics Terraform Sync\n"+
			"Version: %s"+
			"\n"+
			"Usage: %s <kafka-server:port>",
			GitCommit,
			os.Args[0])
		os.Exit(2)
	}

	cluster := os.Args[1]

	//Create AdminClient
	ac, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": cluster})
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	topics := getTopics(ac)
	for _, t := range topics {
		fmt.Print(getTopicMetadata(ac, t))
	}
}
