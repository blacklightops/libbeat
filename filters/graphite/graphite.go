package graphite

import (
	"errors"
	"github.com/blacklightops/libbeat/common"
	"github.com/blacklightops/libbeat/filters"
	"regexp"
	"fmt"
	"strings"
)

type GraphiteMetricExp struct {
	*regexp.Regexp
}

var metricExp = GraphiteMetricExp{regexp.MustCompile(`^[\s]*(?:put)?[\s]*(?P<datacenter>[\w]+)\.(?P<host>[^_]+)_(?P<domain>[\w]+).(?P<service_type>[A-Za-z\-]+)(?P<service_id>[0-9]{3})?(?P<service_id>[a-z\_\d]+)?.(?P<component>[\w]+).(?P<metric>.*)$`)}

func (r *GraphiteMetricExp) FindStringSubmatchMap(s string) (map[string]string, error) {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures, errors.New("Line did not match regex")
	}

	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
    if name == "metric_tags" {
      match[i] += " datacenter=" + captures["datacenter"]
      match[i] += " host=" + captures["host"]
      match[i] += " service=" + captures["service_type"] + captures["service_id"]
      match[i] += " service_type=" + captures["service_type"]
      match[i] += " service_id=" + captures["service_id"]
      match[i] += " component=" + captures["component"]
    }
		captures[name] = strings.TrimSpace(match[i])

	}
	return captures, nil
}


type Graphite struct {
	name string
}

func (graphite *Graphite) New(name string, config map[string]interface{}) (filters.FilterPlugin, error) {
	return &Graphite{name: name}, nil
}

//TODO: Check for Errors Here
func (graphite *Graphite) Filter(event common.MapStr) (common.MapStr, error) {
	text := event["message"]
	text_string := text.(*string)

	metric_data, err := metricExp.FindStringSubmatchMap(*text_string)
	if err != nil {
		return event, nil
	}

	parsed_tags := strings.Fields(metric_data["metric_tags"])
	tags := make(map[string]string)

	for _, v := range parsed_tags {
		tag := strings.Split(v, "=")
		tags[tag[0]] = tag[1]
	}

	event["metric_name"] = metric_data["metric_name"]
	event["metric_value"] = metric_data["metric_value"]
	event["metric_timestamp"] = metric_data["metric_timestamp"]
	event["metric_tags"] = metric_data["metric_tags"]
	event["metric_tags_map"] = tags
	event["index"] = fmt.Sprintf("%s-%s", tags["datacenter"], tags["host"])
	return event, nil
}

func (graphite *Graphite) String() string {
	return graphite.name
}

func (graphite *Graphite) Type() filters.Filter {
	return filters.GraphiteFilter
}
