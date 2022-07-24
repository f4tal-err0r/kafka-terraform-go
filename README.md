## kafka-terraform-go

This is a tool to convert unmanaged kafka cluster topic state directly into terraform.

Initial sync of kafka topics is currently handled by the topic-state binary. To build this bin, you just need to run a `make build` within the bin/ directory, and the compiled binary will be stored within the bin/build folder. Then from there you would run a `topic-state <kafka-cluster:port>` and the binary will fetch all topic info from the cluster. 

The topics will then be output in the form of a kafka_topics resource to stdout for use with the terraform-provider-kafka module. From there just grab a list of all topics either using kafka-topics or something and run them all against `terraform import kafka_topics.$TOPIC $TOPIC` to import state. Success criteria here should be the ability to run a `terraform plan` command and having nothing change.
