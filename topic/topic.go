package topic

import (
	"fmt"
	"strings"
)

type (
	TopicKey string
	Topic    struct {
		Organization string
		Project      string
		NetworkId    string
		Group        string
		Name         string
		Method       string
	}
)

func New(organization string, project string, networkId string, group string, name string, method string) Topic {
	return Topic{
		Organization: organization,
		Project:      project,
		NetworkId:    networkId,
		Group:        group,
		Name:         name,
		Method:       method,
	}
}

func (t *Topic) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"organization": t.Organization,
		"project":      t.Project,
		"network_id":   t.NetworkId,
		"group":        t.Group,
		"name":         t.Name,
		"method":       t.Method,
	}
}

func (t *Topic) ToString(level uint8) string {
	if level < 1 || level > 6 {
		return ""
	}
	switch level {
	case 1:
		return t.Organization
	case 2:
		return fmt.Sprintf("%s.%s", t.Organization, t.Project)
	case 3:
		return fmt.Sprintf("%s.%s.%s", t.Organization, t.Project, t.NetworkId)
	case 4:
		return fmt.Sprintf("%s.%s.%s.%s", t.Organization, t.Project, t.NetworkId, t.Group)
	case 5:
		return fmt.Sprintf("%s.%s.%s.%s.%s", t.Organization, t.Project, t.NetworkId, t.Group, t.Name)
	default:
		// full level
		return fmt.Sprintf("%s.%s.%s.%s.%s.%s", t.Organization, t.Project, t.NetworkId, t.Group, t.Name, t.Method)
	}
}

func (t *Topic) Level() uint8 {
	var level uint8 = 0
	if len(t.Organization) > 0 {
		level++
	}
	if len(t.Project) > 0 {
		level++
	}
	if len(t.NetworkId) > 0 {
		level++
	}
	if len(t.Group) > 0 {
		level++
	}
	if len(t.Name) > 0 {
		level++
	}
	if len(t.Method) > 0 {
		level++
	}
	return level
}

func ParseJSON(obj map[string]interface{}) Topic {
	organization := obj["organization"].(string)
	project := obj["project"].(string)
	topic := Topic{
		Organization: organization,
		Project:      project,
		NetworkId:    "",
		Group:        "",
		Name:         "",
		Method:       "",
	}

	networkId := obj["network_id"]
	if networkId != nil {
		topic.NetworkId = obj["network_id"].(string)
	}

	group := obj["group"]
	if group != nil {
		topic.Group = obj["group"].(string)
	}

	name := obj["name"]
	if name != nil {
		topic.Name = obj["name"].(string)
	}

	method := obj["method"]
	if method != nil {
		topic.Method = obj["method"].(string)
	}

	return topic
}

func ParseString(topicString string) (Topic, error) {
	parts := strings.Split(topicString, ".")
	if len(parts) < 2 {
		return Topic{}, fmt.Errorf("atleast organization and project should be provided")
	}

	if len(parts) > 6 {
		return Topic{}, fmt.Errorf("at most topic shuld be 6 level")
	}

	organization := parts[0]
	project := parts[1]
	networkId := ""
	group := ""
	name := ""
	method := ""
	if len(parts) > 2 {
		networkId = parts[2]
	}
	if len(parts) > 3 {
		group = parts[3]
	}
	if len(parts) > 4 {
		name = parts[4]
	}
	if len(parts) > 5 {
		method = parts[5]
	}

	return Topic{
		Organization: organization,
		Project:      project,
		NetworkId:    networkId,
		Group:        group,
		Name:         name,
		Method:       method,
	}, nil
}

const LEVEL_FULL uint8 = 6 // full topic path, till the method name
const LEVEL_NAME uint8 = 5 // smartcontract level path, till the name of the smartcontract
const LEVEL_2 uint8 = 2    // only organization and project.
