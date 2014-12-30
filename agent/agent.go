package agent

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/libtime"
	"os"
	"time"
)

func NewAgent() (*Agent, error) {
	agent := &Agent{}

	err := agent.setConfigStorage()
	if err != nil {
		return nil, err
	}

	err = agent.setDb()
	if err != nil {
		return nil, err
	}

	return agent, err
}

type Agent struct {
	ConfigStorage *resourced_config.ConfigStorage
	DbPath        string
	Db            *bolt.DB
}

func (a *Agent) setConfigStorage() error {
	readerDir := os.Getenv("RESOURCED_CONFIG_READER_DIR")
	writerDir := os.Getenv("RESOURCED_CONFIG_WRITER_DIR")

	configStorage, err := resourced_config.NewConfigStorage(readerDir, writerDir)
	if err == nil {
		a.ConfigStorage = configStorage
	}
	return err
}

func (a *Agent) setDb() error {
	var err error

	dbPath := os.Getenv("RESOURCED_DB")
	if dbPath == "" {
		dbPath = "~/resourced/db"

		err = os.MkdirAll(libstring.ExpandTilde("~/resourced"), 0755)
		if err != nil {
			return err
		}
	}

	a.DbPath = libstring.ExpandTilde(dbPath)

	a.Db, err = bolt.Open(a.DbPath, 0644, nil)
	if err != nil {
		return err
	}

	// Create "resources" bucket
	a.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("resources"))
		return nil
	})

	return err
}

func (a *Agent) DbBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("resources"))
}

func (a *Agent) Run(config resourced_config.Config) ([]byte, error) {
	cmd := libprocess.NewCmd(config.Command)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	err = a.SaveRun(config, output)

	return output, err
}

func (a *Agent) SaveRun(config resourced_config.Config, output []byte) error {
	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Command"] = config.Command
	record["Path"] = config.Path
	record["Interval"] = config.Interval

	runData := make(map[string]interface{})
	err := json.Unmarshal(output, &runData)
	if err != nil {
		return err
	}

	record["Data"] = runData

	recordInJson, err := json.Marshal(record)
	if err != nil {
		return err
	}

	err = a.Db.Update(func(tx *bolt.Tx) error {
		return a.DbBucket(tx).Put([]byte(config.Path), recordInJson)
	})

	return err
}

func (a *Agent) GetRun(config resourced_config.Config) ([]byte, error) {
	return a.GetRunByPath(config.Path)
}

func (a *Agent) GetRunByPath(path string) ([]byte, error) {
	var data []byte

	a.Db.View(func(tx *bolt.Tx) error {
		data = a.DbBucket(tx).Get([]byte(path))
		return nil
	})

	return data, nil
}

func (a *Agent) RunForever(config resourced_config.Config) {
	go func(a *Agent, config resourced_config.Config) {
		for {
			a.Run(config)
			libtime.SleepString(config.Interval)
		}
	}(a, config)
}

func (a *Agent) RunAllForever() {
	for _, config := range a.ConfigStorage.Readers {
		a.RunForever(config)
	}
	for _, config := range a.ConfigStorage.Writers {
		a.RunForever(config)
	}
}