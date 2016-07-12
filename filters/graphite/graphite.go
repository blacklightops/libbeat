package graphite

import (
	"errors"
	"github.com/blacklightops/libbeat/common"
	"github.com/blacklightops/libbeat/filters"
  "github.com/blacklightops/libbeat/logp"
	"regexp"
	"fmt"
	"strings"
)

type GraphiteMetricExp struct {
	*regexp.Regexp
}

var metricExp = GraphiteMetricExp{regexp.MustCompile(`^(?P<datacenter>[\w]+)[\.\-](?P<host>[A-Za-z0-9\-]+)(?:([_])(?:[_A-Za-z0-9]+\.(?P<service_type>[A-Za-z\-]+)(?P<service_id>[0-9]{1,3})?(?P<service_name>[a-z\_\d]+)?.(?P<component>[\w]+)\.)|\.)(?P<metric_name>[a-zA-Z0-9\.\-_]+(?: [a-zA-Z\ ]+)?)(?:[\s](?P<metric_value>[0-9.E\-]+))?[\s](?P<metric_timestamp>[0-9.]{10,14})[\s]?(?P<metric_tags>[\S]+)?$`)}

func (r *GraphiteMetricExp) FindStringSubmatchMap(s string) (map[string]string, error) {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)

	if match == nil {
		return captures, errors.New("Line did not match regex")
	}

	logp.Debug("filter_graphite", "Regex Matches: %v", match)

	for i, name := range r.SubexpNames() {
		//the first is the original string, skip it
		if i == 0 {
			continue
		}
		captures[name] = strings.TrimSpace(match[i])
	}
	logp.Debug("filter_graphite", "Completed Captures Array: %v", captures)

	if _, present := captures["metric_tags"]; present == false {
		captures["metric_tags"] = ""
	}

	if _, present := captures["metric_value"]; present == false {
		captures["metric_value"] = "+1"
		captures["metric_verb"] = "increment"
	} else if captures["metric_value"] == "" {
		captures["metric_value"] = "+1"
		captures["metric_verb"] = "increment"
	} else {
		captures["metric_verb"] = "put"
	}

	var tags string = ""
	tags = buildTagString(tags, "datacenter", captures)
	tags = buildTagString(tags, "host", captures)
	tags = buildTagString(tags, "service_type", captures)
	tags = buildTagString(tags, "service_id", captures)
	tags = buildTagString(tags, "component", captures)

	captures["metric_tags"] += fmt.Sprintf(" %s", tags)
	return captures, nil
}

func buildTagString(tags string, tag string, captures map[string]string) string {
	if value, present := captures[tag]; present {
		if value == "" {
			if tag == "service_id" {
				value = captures["host"]
			} else {
				return tags
			}
		} 
		return fmt.Sprintf("%s %s=%s", tags, tag, value)
	} else {
		return tags
	}
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
	event["metric_verb"] = metric_data["metric_verb"]
	event["metric_timestamp"] = metric_data["metric_timestamp"]
	event["metric_tags"] = metric_data["metric_tags"]
	event["metric_tags_map"] = tags
	event["index"] = fmt.Sprintf("%s-%s-%s", tags["datacenter"], tags["service_type"], tags["service_id"])
	return event, nil
}

func (graphite *Graphite) String() string {
	return graphite.name
}

func (graphite *Graphite) Type() filters.Filter {
	return filters.GraphiteFilter
}
