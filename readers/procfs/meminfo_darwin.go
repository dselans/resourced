// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcMemInfo is ProcMemInfo constructor.
func NewProcMemInfo() *ProcMemInfo {
	p := &ProcMemInfo{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcMemInfo is a reader that scrapes /proc/meminfo data.
type ProcMemInfo struct {
	Data map[string][]interface{}
}

func (p *ProcMemInfo) Run() error {
	return errors.New("/proc/meminfo is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcMemInfo) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
