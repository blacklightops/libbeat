package jsonexpander

import (
  "encoding/json"
	"github.com/blacklightops/libbeat/common"
	"github.com/blacklightops/libbeat/filters"
)

func isJSONString(s string) bool {
    var js interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}

func isJSON(s string) bool {
    var js map[string]interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}

type JSONExpander struct {
	name string
}

func (jsonexpander *JSONExpander) New(name string, config map[string]interface{}) (filters.FilterPlugin, error) {
	return &JSONExpander{name: name}, nil
}

//TODO: Check for Errors Here
func (jsonexpander *JSONExpander) Filter(event common.MapStr) (common.MapStr, error) {
	text := event["message"]
	text_string := text.(*string)

  if isJSONString(*text_string) {
    data := []byte(*text_string)
    err := json.Unmarshal(data, &event)
    if err != nil {
      return event, nil
    }
  }
	return event, nil
}

func (jsonexpander *JSONExpander) String() string {
	return jsonexpander.name
}

func (jsonexpander *JSONExpander) Type() filters.Filter {
	return filters.JSONExpanderFilter
}
