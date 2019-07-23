package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// EntriesService servıce
type EntriesService service

//Entry model
type Entry struct {
	locale string
	Sys    *Sys `json:"sys"`
	Fields map[string]interface{}
}

// GetVersion returns entity version
func (entry *Entry) GetVersion() int {
	version := 1
	if entry.Sys != nil {
		version = entry.Sys.Version
	}

	return version
}

// GetEntryKey returns the entry's keys
func (service *EntriesService) GetEntryKey(entry *Entry, key string) (*EntryField, error) {
	ef := EntryField{
		value: entry.Fields[key],
	}

	col, err := service.c.ContentTypes.List(entry.Sys.Space.Sys.ID).Next()
	if err != nil {
		return nil, err
	}

	for _, ct := range col.ToContentType() {
		if ct.Sys.ID != entry.Sys.ContentType.Sys.ID {
			continue
		}

		for _, field := range ct.Fields {
			if field.ID != key {
				continue
			}

			ef.dataType = field.Type
		}
	}

	return &ef, nil
}

// List returns entries collection
func (service *EntriesService) List(spaceID string) *Collection {
	path := service.c.newPath(spaceID, "/entries")
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single entry
func (service *EntriesService) Get(spaceID, entryID string) (*Entry, error) {
	path := service.c.newPath(spaceID, fmt.Sprintf("/entries/%s", entryID))
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Entry{}, err
	}

	var entry Entry
	if ok := service.c.do(req, &entry); ok != nil {
		return nil, err
	}

	return &entry, err
}

// Delete the entry
func (service *EntriesService) Delete(spaceID string, entryID string) error {
	path := service.c.newPath(spaceID, fmt.Sprintf("/entries/%s", entryID))
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Update the entry
func (service *EntriesService) Update(spaceID string, entry *Entry) error {
	path := service.c.newPath(spaceID, fmt.Sprintf("/entries/%s", entry.Sys.ID))
	method := "PUT"

	byt, err := json.Marshal(map[string]interface{}{"fields": entry.Fields})
	if err != nil {
		return err
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(byt))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Content-Type", entry.Sys.ContentType.Sys.ID)

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Publish the entry
func (service *EntriesService) Publish(spaceID string, entry *Entry) error {
	path := service.c.newPath(spaceID, fmt.Sprintf("/entries/%s/published", entry.Sys.ID))
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Unpublish the entry
func (service *EntriesService) Unpublish(spaceID string, entry *Entry) error {
	path := service.c.newPath(spaceID, fmt.Sprintf("/entries/%s/published", entry.Sys.ID))
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
