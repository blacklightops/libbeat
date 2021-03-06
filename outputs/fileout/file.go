package fileout

import (
	"encoding/json"
	"time"

	"github.com/blacklightops/libbeat/common"
	"github.com/blacklightops/libbeat/logp"
	"github.com/blacklightops/libbeat/outputs"
)

type FileOutput struct {
	rotator logp.FileRotator
}

func (out *FileOutput) Init(config outputs.MothershipConfig, topology_expire int) error {
	out.rotator.Path = config.Path
	out.rotator.Name = config.Filename
	if out.rotator.Name == "" {
		out.rotator.Name = "packetbeat"
	}

	rotateeverybytes := uint64(config.Rotate_every_kb) * 1024
	if rotateeverybytes == 0 {
		rotateeverybytes = 10 * 1024 * 1024
	}
	out.rotator.RotateEveryBytes = &rotateeverybytes

	keepfiles := config.Number_of_files
	if keepfiles == 0 {
		keepfiles = 7
	}
	out.rotator.KeepFiles = &keepfiles

	err := out.rotator.CreateDirectory()
	if err != nil {
		return err
	}

	err = out.rotator.CheckIfConfigSane()
	if err != nil {
		return err
	}

	return nil
}

func (out *FileOutput) PublishIPs(name string, localAddrs []string) error {
	// not supported by this output type
	return nil
}

func (out *FileOutput) GetNameByIP(ip string) string {
	// not supported by this output type
	return ""
}

func (out *FileOutput) PublishEvent(ts time.Time, event common.MapStr) error {

	json_event, err := json.Marshal(event)
	if err != nil {
		logp.Err("Fail to convert the event to JSON: %s", err)
		return err
	}

	err = out.rotator.WriteLine(json_event)
	if err != nil {
		return err
	}

	return nil
}
