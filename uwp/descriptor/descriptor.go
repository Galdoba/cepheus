package descriptor

import (
	"encoding/json"
	"fmt"
	"os"
)

type Descriptor struct {
	path     string
	Info     string            `json:"header,omitempty"`
	Category map[string]*Topic `json:"category"`
}

type Topic struct {
	Entries  int                     `json:"total entries"`
	TopicMap map[string]*Description `json:"topic"`
}

func newTopic() *Topic {
	t := Topic{}
	t.TopicMap = make(map[string]*Description)
	return &t
}

type Description struct {
	Description map[string]string `json:"description,omitempty"`
}

func newDescription() *Description {
	d := Description{}
	d.Description = make(map[string]string)
	return &d
}

func New(path string) *Descriptor {
	ds := Descriptor{}
	ds.path = path
	ds.Category = make(map[string]*Topic)
	return &ds
}

func (ds *Descriptor) Save() error {
	for _, category := range ds.Category {
		category.Entries = len(category.TopicMap)
	}
	data, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	if err := os.WriteFile(ds.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

func (ds *Descriptor) Load() error {
	data, err := os.ReadFile(ds.path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	if err := json.Unmarshal(data, ds); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return nil
}

func (ds *Descriptor) AddCategory(category string) {
	ds.Category[category] = newTopic()
}

func (ds *Descriptor) AddDescription(category, topic, description string) {
	if _, ok := ds.Category[category]; !ok {
		ds.AddCategory(category)
	}
	if _, ok := ds.Category[category].TopicMap[topic]; !ok {
		ds.Category[category].TopicMap[topic] = newDescription()
	}
	ds.Category[category].TopicMap[topic].Description["en"] = description
	ds.Category[category].TopicMap[topic].Description["ru"] = description
}
