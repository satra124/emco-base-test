// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	pkgerrors "github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var refSchemaFile string = "ref-schemas/emco-db-ref-schema.yaml"

type ReferenceEntry struct {
	Key   Key
	KeyId string
}

// ReferenceSchema defines the structure of a reference entry in the referential schema
type ReferenceSchema struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`       // can be not present or "map" or "many"
	CommonKey  string            `yaml:"commonKey"`  // optional to disambiguate - part of reference key that is common
	Map        string            `yaml:"map"`        // if Type is "map", this is JSON tag of the Map object
	FixedKv    map[string]string `yaml:"fixedKv"`    // Key/Values of referenced resource that are known at compile time
	FilterKeys []string          `yaml:"filterKeys"` // if "map" type, list of keys to filter (not count as references)
}

// ResourceSchema defines the structure of a data resource in the referential schema
type ResourceSchema struct {
	Name       string            `yaml:"name"`
	Parent     string            `yaml:"parent"`
	References []ReferenceSchema `yaml:"references"` // if present, overrides default resource id
}

// DbSchema is the top level structure for the referential schema
type DbSchema struct {
	Resources []ResourceSchema `yaml:"resources"`
}

// SchemaMapElement defines the structure that will be prepared
// for each element of the referential schema map
type SchemaMapElement struct {
	parent       string              // Name of the parent resource
	trimParent   bool                // Special case for combined key resource - e.g. compositeApp.compositeAppVersion
	keyId        string              // the keyId string identifier for this resource
	children     map[string]struct{} // list of children
	keyMap       map[string]struct{} // map structure for looking up this item
	references   []ReferenceSchema   // list of references
	referencedBy map[string]struct{} // list of resource that may reference this resource
}

// refMap is a map of the referential schema setup for convenient reference
// by the DB interface functions to identify and manage references.
var refMap map[string]SchemaMapElement

// keyToName maps the keyId to the name of a resource
var keyToName map[string]string

// GetKeyMap returns a key structure (map) for a given resource name
func GetKeyMap(name string) (map[string]struct{}, error) {
	e, ok := refMap[name]
	if !ok {
		return nil, pkgerrors.New("Missing reference map element error: " + name)
	}

	keyMap := make(map[string]struct{})
	id := name
	for true {
		if _, ok := keyMap[id]; ok {
			return nil, pkgerrors.New("Circular keyMap error: " + name)
		}
		keyMap[id] = struct{}{}
		if len(e.parent) == 0 {
			break
		}
		id = e.parent
		if e, ok = refMap[e.parent]; !ok {
			return nil, pkgerrors.New("Missing reference map parent element error: " + e.parent)
		}
	}

	return keyMap, nil
}

// ReadRefSchema reads in the referential schema and creates the refMap
func ReadRefSchema() error {
	var emcoRefSchema DbSchema

	// Read the Referential Schema File
	if _, err := os.Stat(refSchemaFile); err != nil {
		if os.IsNotExist(err) {
			err = pkgerrors.New("Database Referential Schema File " + refSchemaFile + " not found")
		} else {
			err = pkgerrors.Wrap(err, "Database Referential Schema Stat file error")
		}
		return err
	}
	rawBytes, err := ioutil.ReadFile(refSchemaFile)
	if err != nil {
		return pkgerrors.Wrap(err, "Database Referential Schema Read file error")
	}

	err = yaml.Unmarshal(rawBytes, &emcoRefSchema)
	if err != nil {
		log.Error("DB Ref Map", log.Fields{"error": err})
		return pkgerrors.Wrap(err, "Database Referential Schema YAML format error")
	}

	// Pass 1 - put all the resource elements into the referential schema map
	refMap = make(map[string]SchemaMapElement)
	for _, e := range emcoRefSchema.Resources {
		// special case for a name with two elements - e.g. name1.name2
		// this turns into two entries where name1 is the parent of name2
		names := strings.Split(e.Name, ".")
		if len(names) > 2 {
			log.Error("DB Ref Map - Invalid ref-schema resource name", log.Fields{"resource": e.Name})
			return pkgerrors.New("Database Referential Schema invalid resource name: " + e.Name)
		}
		for i, name := range names {
			if _, ok := refMap[name]; !ok {
				se := SchemaMapElement{
					children:     make(map[string]struct{}),
					keyMap:       make(map[string]struct{}),
					referencedBy: make(map[string]struct{}),
					trimParent:   false,
				}
				if len(names) == 1 {
					se.parent = e.Parent
					se.references = e.References
				} else {
					if i == 0 {
						se.parent = e.Parent
					} else {
						se.parent = names[0]
						se.trimParent = true
						se.references = e.References
					}
				}
				refMap[name] = se
			} else {
				log.Error("DB Ref Map", log.Fields{"error": err})
				return pkgerrors.New("Database Referential Schema duplicate resource error: " + name)
			}
		}
	}

	// Pass 2 - fill out children entries, populate keyMaps
	keyToName = make(map[string]string)
	for name, m := range refMap {
		// Add to parents child list
		if m.parent != "" {
			if p, ok := refMap[m.parent]; ok {
				p.children[name] = struct{}{}
			} else {
				log.Error("DB Ref Map", log.Fields{"Missing parent": name})
				return pkgerrors.New("Database Referential Schema missing parent error: " + name)
			}
		}

		// Update referencedBy list of references
		// Check if referenced by:
		for _, r := range m.references {
			// get the referenced refMap entry
			re, ok := refMap[r.Name]
			if !ok {
				return pkgerrors.New("Database Referential Schema missing reference entry error: " + r.Name)
			}

			// add 'id' to the referencedBy list of the referenced refMap entry
			re.referencedBy[name] = struct{}{}
			refMap[r.Name] = re
		}

		// check keyMap
		km, err := GetKeyMap(name)
		if err != nil {
			return err
		}
		m.keyMap = km

		keyStr, err := createKeyId(km)
		if err != nil {
			return pkgerrors.New("Database Referential Schema error making keyStr: " + name)
		}
		keyToName[keyStr] = name

		m.keyId = keyStr
		refMap[name] = m
	}

	return nil
}

func createKeyId(key interface{}) (string, error) {

	var n map[string]struct{}
	st, err := json.Marshal(key)
	if err != nil {
		return "", pkgerrors.Errorf("Error Marshalling key: %s", err.Error())
	}
	err = json.Unmarshal([]byte(st), &n)
	if err != nil {
		return "", pkgerrors.Errorf("Error Unmarshalling key to Json Map: %s", err.Error())
	}
	var keys []string
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	s := "{"
	for _, k := range keys {
		s = s + k + ","
	}
	s = s + "}"
	return s, nil
}
