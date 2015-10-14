package stdout

import (
	"fmt"
	"github.com/blacklightops/libbeat/common"
	"time"
	//  "github.com/blacklightops/libbeat/logp"
	"github.com/blacklightops/libbeat/outputs"
)

type StdOutput struct {
	enabled string
}

func (out *StdOutput) Init(config outputs.MothershipConfig, topology_expire int) error {
	// not supported by this output type
	return nil
}

func (out *StdOutput) PublishIPs(name string, localAddrs []string) error {
	// not supported by this output type
	return nil
}

func (out *StdOutput) GetNameByIP(ip string) string {
	// not supported by this output type
	return ""
}

func (out *StdOutput) PublishEvent(ts time.Time, event common.MapStr) error {
	//json_event, err := json.Marshal(event)
	//if err != nil {
	//  logp.Err("Fail to convert the event to JSON: %s", err)
	//  return err
	//}

	out.Print(event)
	return nil
}

func (out *StdOutput) Print(event common.MapStr) {
	str := event.String()
	fmt.Println(str)
}
