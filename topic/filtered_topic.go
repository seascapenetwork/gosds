package topic

type TopicFilter struct {
	Organizations  []string
	Projects       []string
	NetworkIds     []string
	Groups         []string
	Smartcontracts []string
	Methods        []string
	Events         []string
}

func NewFilterTopic(o []string, p []string, n []string, g []string, s []string, m []string, e []string) TopicFilter {
	return TopicFilter{
		Organizations:  o,
		Projects:       p,
		NetworkIds:     n,
		Groups:         g,
		Smartcontracts: s,
		Methods:        m,
		Events:         e,
	}
}

func (t *TopicFilter) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"o":  t.Organizations,
		"p:": t.Projects,
		"n":  t.NetworkIds,
		"g":  t.Groups,
		"s":  t.Smartcontracts,
		"m":  t.Methods,
		"e":  t.Events,
	}
}

func (t *TopicFilter) Len(level uint8) int {
	switch level {
	case ORGANIZATION_LEVEL:
		return len(t.Organizations)
	case PROJECT_LEVEL:
		return len(t.Projects)
	case NETWORK_ID_LEVEL:
		return len(t.NetworkIds)
	case GROUP_LEVEL:
		return len(t.Groups)
	case SMARTCONTRACT_LEVEL:
		return len(t.Smartcontracts)
	case FULL_LEVEL:
		return len(t.Methods) + len(t.Events)
	default:
		return len(t.Organizations) + len(t.Projects) + len(t.NetworkIds) + len(t.Groups) + len(t.Smartcontracts) + len(t.Methods) + len(t.Events)
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
	if len(t.Organizations) > 0 {
		str += "o:" + list(t.Organizations) + ";"
	}
	if len(t.Projects) > 0 {
		str += "p:" + list(t.Projects) + ";"
	}
	if len(t.NetworkIds) > 0 {
		str += "n:" + list(t.NetworkIds) + ";"
	}
	if len(t.Groups) > 0 {
		str += "g:" + list(t.Groups) + ";"
	}
	if len(t.Smartcontracts) > 0 {
		str += "s:" + list(t.Smartcontracts) + ";"
	}
	if len(t.Methods) > 0 {
		str += "m:" + list(t.Methods) + ";"
	}
	if len(t.Events) > 0 {
		str += "e:" + list(t.Events) + ";"
	}

	return str
}

func ParseJSONToTopicFilter(obj map[string]interface{}) TopicFilter {
	topic := TopicFilter{
		Organizations:  []string{},
		Projects:       []string{},
		NetworkIds:     []string{},
		Groups:         []string{},
		Smartcontracts: []string{},
		Methods:        []string{},
		Events:         []string{},
	}

	if obj["n"] != nil {
		n := obj["n"].([]interface{})
		topic.NetworkIds = make([]string, len(n))
		for i, o := range n {
			topic.NetworkIds[i] = o.(string)
		}
	}

	if obj["o"] != nil {
		organizations := obj["o"].([]interface{})
		topic.Organizations = make([]string, len(organizations))
		for i, o := range organizations {
			topic.Organizations[i] = o.(string)
		}
	}

	if obj["p"] != nil {
		projects := obj["p"].([]interface{})
		topic.Projects = make([]string, len(projects))
		for i, o := range projects {
			topic.Projects[i] = o.(string)
		}
	}

	if obj["g"] != nil {
		g := obj["g"].([]interface{})
		topic.Groups = make([]string, len(g))
		for i, o := range g {
			topic.Groups[i] = o.(string)
		}
	}

	if obj["s"] != nil {
		s := obj["s"].([]interface{})
		topic.Smartcontracts = make([]string, len(s))
		for i, o := range s {
			topic.Smartcontracts[i] = o.(string)
		}
	}

	if obj["m"] != nil {
		m := obj["m"].([]interface{})
		topic.Methods = make([]string, len(m))
		for i, o := range m {
			topic.Methods[i] = o.(string)
		}
	}

	if obj["e"] != nil {
		e := obj["e"].([]interface{})
		topic.Events = make([]string, len(e))
		for i, o := range e {
			topic.Events[i] = o.(string)
		}
	}

	return topic
}
