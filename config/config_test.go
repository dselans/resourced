package config

import (
	"encoding/json"
	"testing"
)

func TestNewConfigStorage(t *testing.T) {
	config, err := NewConfigStorage("$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader", "$GOPATH/src/github.com/resourced/resourced/tests/data/config-writer")
	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

	if len(config.Readers) <= 0 {
		t.Errorf("Length of reader config should > 0. config.Readers: %v", config.Readers)
	}
	if len(config.Writers) != 4 {
		t.Errorf("Length of reader config should == 3. config.Writers: %v", config.Writers)
	}
}

func TestNewReaderConfig(t *testing.T) {
	config, err := NewConfig("$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader/gostruct-docker-container-memory.toml", "reader")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}

	if config.GoStruct != "DockerContainersMemory" {
		t.Fatalf("Config is initialized incorrectly. config.GoStruct: %v", config.GoStruct)
	}
	if config.Path != "/docker/containers/memory" {
		t.Fatalf("Config is initialized incorrectly. config.Path: %v", config.Path)
	}
	if config.Interval != "3s" {
		t.Fatalf("Config is initialized incorrectly. config.Interval: %v", config.Interval)
	}
	if config.GoStructFields["CgroupBasePath"] != "/sys/fs/cgroup/cpuacct/docker" {
		inJson, _ := json.Marshal(config.GoStructFields)
		t.Fatalf("Config is initialized incorrectly. config.GoStructFields: %v", string(inJson))
	}
}

func TestNewWriterConfigWithJsonProcessor(t *testing.T) {
	config, err := NewConfig("$GOPATH/src/github.com/resourced/resourced/tests/data/config-writer/gostruct-stdout.toml", "writer")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}
	if config.GoStructFields["JsonProcessor"] == "" {
		inJson, _ := json.Marshal(config.GoStructFields)
		t.Fatalf("Config is initialized incorrectly. config.GoStructFields: %v", string(inJson))
	}
}
