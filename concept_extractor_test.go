// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"github.com/getgauge/gauge/gauge_messages"
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestExtractConceptWithoutParameters(c *C) {
	STEP := "step that takes a table"
	name := "concept"
	conceptName := &gauge_messages.Step{Name: &name}
	concept, conceptText := getExtractedConcept(conceptName, []*gauge_messages.Step{&gauge_messages.Step{Name: &STEP}})

	c.Assert(concept, Equals, "# concept\n* step that takes a table\n")
	c.Assert(conceptText, Equals, "* concept")
}

func (s *MySuite) TestExtractConcept(c *C) {
	STEP := "step that takes a table \"arg\""
	name := "concept with \"arg\""
	conceptName := &gauge_messages.Step{Name: &name}
	concept, conceptText := getExtractedConcept(conceptName, []*gauge_messages.Step{&gauge_messages.Step{Name: &STEP}})

	c.Assert(concept, Equals, "# concept with <arg>\n* step that takes a table <arg>\n")
	c.Assert(conceptText, Equals, "* concept with \"arg\"")
}

func (s *MySuite) TestExtractConceptWithSkippedParameters(c *C) {
	STEP := "step that takes a table \"arg\" and \"hello again\" "
	name := "concept with \"arg\""
	conceptName := &gauge_messages.Step{Name: &name}
	concept, conceptText := getExtractedConcept(conceptName, []*gauge_messages.Step{&gauge_messages.Step{Name: &STEP}})

	c.Assert(concept, Equals, "# concept with <arg>\n* step that takes a table <arg> and \"hello again\"\n")
	c.Assert(conceptText, Equals, "* concept with \"arg\"")
}

func (s *MySuite) TestExtractConceptWithParameters(c *C) {
	STEP := "step that takes a table \"arg\" and \"hello again\" "
	name := "concept with \"arg\" \"hello again\""
	conceptName := &gauge_messages.Step{Name: &name}
	concept, conceptText := getExtractedConcept(conceptName, []*gauge_messages.Step{&gauge_messages.Step{Name: &STEP}})

	c.Assert(concept, Equals, "# concept with <arg> <hello again>\n* step that takes a table <arg> and <hello again>\n")
	c.Assert(conceptText, Equals, "* concept with \"arg\" \"hello again\"")
}

func (s *MySuite) TestExtractConceptWithTableAsArg(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "Step with inline table", lineNo: 3},
		&token{kind: tableHeader, args: []string{"id", "name"}},
		&token{kind: tableRow, args: []string{"1", "foo"}},
		&token{kind: tableRow, args: []string{"2", "bar"}},
	}
	spec, _ := new(specParser).createSpecification(tokens, new(conceptDictionary))
	step := spec.scenarios[0].steps[0]
	inlineTable := step.args[0].table
	protoTable := convertToProtoTableParam(&inlineTable)
	STEP := "step that takes a table"
	name := "concept with \"table1\""
	conceptName := &gauge_messages.Step{Name: &name}
	tableName := TABLE + "1"
	concept, conceptText := getExtractedConcept(conceptName, []*gauge_messages.Step{&gauge_messages.Step{Name: &STEP, Table: protoTable, ParamTableName: &tableName},
		&gauge_messages.Step{Name: &STEP, Table: protoTable, ParamTableName: &tableName}})

	c.Assert(concept, Equals, "# concept with <table1>\n* step that takes a table <table1>\n* step that takes a table <table1>\n")
	c.Assert(conceptText, Equals, "* concept with "+`
     |id|name|
     |--|----|
     |1 |foo |
     |2 |bar |
`)
}

func (s *MySuite) TestExtractConceptWithSkippedTableAsArg(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "Step with inline table", lineNo: 3},
		&token{kind: tableHeader, args: []string{"id", "name"}},
		&token{kind: tableRow, args: []string{"1", "foo"}},
		&token{kind: tableRow, args: []string{"2", "bar"}},
	}
	spec, _ := new(specParser).createSpecification(tokens, new(conceptDictionary))
	step := spec.scenarios[0].steps[0]
	inlineTable := step.args[0].table
	protoTable := convertToProtoTableParam(&inlineTable)
	STEP := "step that takes a table"
	name := "concept with \"table1\""
	conceptName := &gauge_messages.Step{Name: &name}
	tableName := TABLE + "1"
	concept, conceptText := getExtractedConcept(conceptName, []*gauge_messages.Step{&gauge_messages.Step{Name: &STEP, Table: protoTable, ParamTableName: &tableName},
		&gauge_messages.Step{Name: &STEP, Table: protoTable, ParamTableName: &tableName}, &gauge_messages.Step{Name: &STEP, Table: protoTable}})

	c.Assert(concept, Equals, "# concept with <table1>\n* step that takes a table <table1>\n* step that takes a table <table1>\n* step that takes a table "+`
     |id|name|
     |--|----|
     |1 |foo |
     |2 |bar |
`)
	c.Assert(conceptText, Equals, "* concept with "+`
     |id|name|
     |--|----|
     |1 |foo |
     |2 |bar |
`)
}

func (s *MySuite) TestReplaceText(c *C) {
	content := `Copyright 2015 ThoughtWorks, Inc.

	This file is part of Gauge.

	Gauge is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Gauge is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Gauge.  If not, see <http://www.gnu.org/licenses/>.`

	replacement := `* concept with
     |id|name|
     |--|----|
     |1 |foo |
     |2 |bar |
`
	five := int32(5)
	ten := int32(10)
	info := &gauge_messages.TextInfo{StartingLineNo: &five, EndLineNo: &ten}
	finalText := replaceText(content, info, replacement)

	c.Assert(finalText, Equals, `Copyright 2015 ThoughtWorks, Inc.

	This file is part of Gauge.

* concept with
     |id|name|
     |--|----|
     |1 |foo |
     |2 |bar |

	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Gauge.  If not, see <http://www.gnu.org/licenses/>.`)
}