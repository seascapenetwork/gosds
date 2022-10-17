package topic

import (
	"fmt"
	"strings"
)

type TopicFilter struct {
	Organization []string
	Project      []string
	NetworkId    []string
	Group        []string
	Name         []string
}

func NewFilterTopic(organization []string, project []string, networkId []string, group []string, name []string) TopicFilter {
	return TopicFilter{
		Organization: organization,
		Project:      project,
		NetworkId:    networkId,
		Group:        group,
		Name:         name,
	}
}

func (t *TopicFilter) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"organization": t.Organization,
		"project:":     t.Project,
		"network_id":   t.NetworkId,
		"group":        t.Group,
		"name":         t.Name,
	}
}

func (t *TopicFilter) Len(property string) int {
	switch property {
	case "organization":
		return len(t.Organization)
	case "project":
		return len(t.Project)
	case "network_id":
		return len(t.NetworkId)
	case "group":
		return len(t.Group)
	case "name":
		return len(t.Name)
	default:
		return t.Len("organization") + t.Len("project") + t.Len("network_id") + t.Len("group") + t.Len("name")
	}
}

func (t *TopicFilter) Key() TopicKey {
	return TopicKey(t.ToString())
}

func list(properties []string) string {
	str := ""
	for _, v := range properties {
		str += "," + v
	}

	return str
}

func (t *TopicFilter) ToString() string {
	str := ""
	if len(t.Organization) > 0 {
		str += "o:" + list(t.Organization) + ";"
	}
	if len(t.Project) > 0 {
		str += "p:" + list(t.Project) + ";"
	}
	if len(t.NetworkId) > 0 {
		str += "n:" + list(t.NetworkId) + ";"
	}
	if len(t.Group) > 0 {
		str += "g:" + list(t.Group) + ";"
	}
	if len(t.Name) > 0 {
		str += "s:" + list(t.Name) + ";"
	}

	return str
}

func ParseJSONToTopicFilter(obj map[string]interface{}) TopicFilter {
	topic := TopicFilter{
		Organization: []string{},
		Project:      []string{},
		NetworkId:    []string{},
		Group:        []string{},
		Name:         []string{},
	}

	if obj["network_id"] != nil {
		network_id := obj["network_id"].([]interface{})
		topic.NetworkId = make([]string, len(network_id))
		for i, o := range network_id {
			topic.NetworkId[i] = o.(string)
		}
	}

	if obj["organization"] != nil {
		organizations := obj["organization"].([]interface{})
		topic.Organization = make([]string, len(organizations))
		for i, o := range organizations {
			topic.Organization[i] = o.(string)
		}
	}

	if obj["project"] != nil {
		project := obj["project"].([]interface{})
		topic.Project = make([]string, len(project))
		for i, o := range project {
			topic.Project[i] = o.(string)
		}
	}

	if obj["group"] != nil {
		group := obj["group"].([]interface{})
		topic.Group = make([]string, len(group))
		for i, o := range group {
			topic.Group[i] = o.(string)
		}
	}

	if obj["name"] != nil {
		name := obj["name"].([]interface{})
		topic.Name = make([]string, len(name))
		for i, o := range name {
			topic.Name[i] = o.(string)
		}
	}

	return topic
}

func ParseStringToTopicFilter(topicString string) (TopicFilter, error) {
	parts := strings.Split(topicString, ";")
	if len(parts) < 2 {
		return TopicFilter{}, fmt.Errorf("atleast organization and project should be provided")
	}

	if len(parts) > 6 {
		return TopicFilter{}, fmt.Errorf("at most topic shuld be 6 level")
	}

	return TopicFilter{}, nil
	// organization := parts[0]
	// project := parts[1]
	// networkId := ""
	// group := ""
	// name := ""
	// method := ""
	// if len(parts) > 2 {
	// 	networkId = parts[2]
	// }
	// if len(parts) > 3 {
	// 	group = parts[3]
	// }
	// if len(parts) > 4 {
	// 	name = parts[4]
	// }
	// if len(parts) > 5 {
	// 	method = parts[5]
	// }

	// return TopicFilter{
	// 	Organization: organization,
	// 	Project:      project,
	// 	NetworkId:    networkId,
	// 	Group:        group,
	// 	Name:         name,
	// 	Method:       method,
	// }, nil
}
