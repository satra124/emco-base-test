// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package db

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func init() {
	refSchemaFile = "test-schemas/cluster.yaml"
}

func printKeyToNameMap() {
	fmt.Printf("KEY TO NAME MAP\n")
	for k, v := range keyToName {
		fmt.Printf("  %v: %v\n", k, v)
	}
}

func printRefMap() {
	fmt.Printf("REFERENCE MAP\n")
	for name, se := range refMap {
		fmt.Printf("  Name: %v\n    Parent: %v\n    trimParent: %v\n", name, se.parent, se.trimParent)
		fmt.Printf("    Key: %v\n    Children: %v\n", se.keyMap, se.children)
		fmt.Printf("    References:\n")
		for _, re := range se.references {
			fmt.Printf("      Name: %v\n      Type: %v\n      Map: %v\n", re.Name, re.Type, re.Map)
			fmt.Printf("      FixedKv:\n")
			for k, v := range re.FixedKv {
				fmt.Printf("        %v: %v\n", k, v)
			}
			fmt.Printf("      FilterKeys:\n")
			for _, s := range re.FilterKeys {
				fmt.Printf("        %v\n", s)
			}
		}
		fmt.Printf("    ReferencedBy: %v\n    KeyId: %v\n", se.referencedBy, se.keyId)
	}
}

var _ = Describe("referential schema handler", func() {

	type testCase struct {
		schemaFile    string
		expectedError string
	}

	DescribeTable("Read Referential Schema",
		func(t testCase) {
			refSchemaFile = t.schemaFile
			err := ReadRefSchema()
			printRefMap()
			printKeyToNameMap()

			if t.expectedError == "" {
				Expect(err).To(BeNil())
			} else {
				Expect(strings.Contains(err.Error(), t.expectedError)).To(BeTrue())
			}
		},

		Entry("successful schema", testCase{
			schemaFile:    "test-schemas/cluster.yaml",
			expectedError: "",
		}),

		/*
			Entry("looping schema", testCase{
				schemaFile:    "test-schemas/loop.yaml",
				expectedError: "Circular keyMap error:",
			}),

			Entry("missing parent schema", testCase{
				schemaFile:    "test-schemas/missingparent.yaml",
				expectedError: "Database Referential Schema missing parent error:",
			}),

			Entry("duplicate resource schema", testCase{
				schemaFile:    "test-schemas/duplicateresource.yaml",
				expectedError: "Database Referential Schema duplicate resource error:",
			}),

			Entry("bad yaml schema", testCase{
				schemaFile:    "test-schemas/badyaml.yaml",
				expectedError: "Database Referential Schema YAML format error",
			}),
		*/
	)

})
