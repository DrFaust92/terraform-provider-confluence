package confluence

import (
	"fmt"
)

// Body is part of Content
type Body struct {
	Storage *Storage `json:"storage,omitempty"`
}

// Content is a primary resource in Confluence
type Content struct {
	Id        string        `json:"id,omitempty"`
	Type      string        `json:"type,omitempty"`
	Title     string        `json:"title,omitempty"`
	Space     *Space        `json:"space,omitempty"`
	Version   *Version      `json:"version,omitempty"`
	Body      *Body         `json:"body,omitempty"`
	Links     *ContentLinks `json:"_links,omitempty"`
	Ancestors []*Content    `json:"ancestors,omitempty"`
	MetaData  *MetaData     `json:"metadata,omitempty"`
}

// ContentLinks is part of Content
type ContentLinks struct {
	Context string `json:"context,omitempty"`
	WebUI   string `json:"webui,omitempty"`
}

// Label is part of MetaData
type Label struct {
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

// Labels is part of MetaData
type Labels struct {
	LabelResults []*Label `json:"results,omitempty"`
}

// MetaData is part of Content
type MetaData struct {
	Labels *Labels `json:"labels,omitempty"`
}

// Space is part of Content
type Space struct {
	Key string `json:"key,omitempty"`
}

// Storage is part of Body
type Storage struct {
	Value          string `json:"value,omitempty"`
	Representation string `json:"representation,omitempty"`
}

// Version is part of Content
type Version struct {
	Number int `json:"number,omitempty"`
}

func (c *Client) CreateContent(content *Content) (*Content, error) {
	var response Content
	if err := c.Post("/rest/api/content", content, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetContent(id string) (*Content, error) {
	var response Content
	path := fmt.Sprintf("/rest/api/content/%s?expand=space,body.storage,version,ancestors,metadata.labels", id)
	if err := c.Get(path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) UpdateContent(content *Content) (*Content, error) {
	var response Content
	content.Version.Number++
	path := fmt.Sprintf("/rest/api/content/%s?expand=metadata.labels", content.Id)
	if err := c.Put(path, content, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) DeleteContent(id string) error {
	path := fmt.Sprintf("/rest/api/content/%s", id)
	if err := c.Delete(path); err != nil {
		return err
	}
	return nil
}
