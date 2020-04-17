package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/anz-bank/sysl/pkg/diagrams"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"
)

func TestDataModel(t *testing.T) {
	t.Parallel()
	p := &datamodelCmd{}
	p.Output = "whatever.svg"
	p.Project = "Project"
	p.ClassFormat = "%(classname)"
	p.Direct = false
	fs := afero.NewOsFs()
	filename := testDir + "multiple-app-datamodel.sysl"
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, fs)
	assert.Nil(t, err, "Generating the data diagrams failed")

	expected := map[string]string{"whatever.svg": `@startuml
	''''''''''''''''''''''''''''''''''''''''''
	''                                      ''
	''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
	''                                      ''
	''''''''''''''''''''''''''''''''''''''''''
	
	class "App.User" as _0 << (D,orchid) >> {
	+ beep : **Server.ftgyhb**
	+ weuiyfgwihe : **AnotherApp.AnotherType**
	+ whatever : int
	}
	class "Server.ftgyhb" as _1 << (D,orchid) >> {
	+ blah : int
	}
	_0 *-- "1..1 " _1
	@enduml
`}
	assert.Equal(t, outmap, expected)
}

func TestHTML(t *testing.T) {
	t.Parallel()
	p := &datamodelCmd{}
	p.Output = "whatever.html"
	p.Project = "Project"
	p.ClassFormat = "%(classname)"
	p.Direct = false
	fs := afero.NewOsFs()
	plantuml := os.Getenv("SYSL_PLANTUML")
	filename := testDir + "multiple-app-datamodel.sysl"
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	fs = afero.NewMemMapFs()
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, fs)
	assert.Nil(t, err)
	err = diagrams.OutputPlantuml(p.Output, plantuml, outmap["whatever.html"], fs)
	assert.Nil(t, err)
	file, err := fs.Open(p.Output)
	assert.Nil(t, err)
	html, err := ioutil.ReadAll(file)
	assert.Nil(t, err)
	expected := fmt.Sprintf(
		"<img src=\"%s/svg/UDgCaK5hmZ0OXk"+
			"_v5UzwwHRNOgyYObd2sCN25iybwgS9T3K"+
			"aKIbZ_tqudR23aFUKmvE71xvvYDh7msV7"+
			"ykInF4VIhcYzLCshNCZtMac1bqP850L4W"+
			"qMZ8CMRUhYXoTRaLhRkT0Z8QkVAPI2VGG"+
			"UqOqodU1JFUxysTlE4sGnVR8GLaSC4d7-"+
			"GFv8ljt4tc1NdR8GJ7UomtPsieL-YxlVH"+
			"a3zhhqyrEBgWHaSwaiS4je_i-o_Xd_JIKT"+
			"tfwBIv_Mx7u4M1h2hB2XdQYisftvnFzXC0"+
			"0F__jbPdCm00\" alt=\"plantuml\">\n", plantuml)
	assert.Equal(t, expected, string(html))
}

func TestLinkOutput(t *testing.T) {
	t.Parallel()
	p := &datamodelCmd{}
	p.Output = "whatever.html"
	p.Project = "Project"
	p.ClassFormat = "%(classname)"
	p.Direct = false
	plantuml := os.Getenv("SYSL_PLANTUML")
	fs := afero.NewOsFs()
	filename := testDir + "multiple-app-datamodel.sysl"
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	fs = afero.NewMemMapFs()
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, fs)
	assert.Nil(t, err)
	err = diagrams.OutputPlantuml(p.Output, plantuml, outmap["whatever.html"], fs)
	assert.Nil(t, err)
	file, err := fs.Open(p.Output)
	assert.Nil(t, err)
	link, err := ioutil.ReadAll(file)
	assert.Nil(t, err)
	expected := fmt.Sprintf(
		"%s/svg/UDgCaK5hmZ0OXk"+
			"_v5UzwwHRNOgyYObd2sCN25iybwgS9T3K"+
			"aKIbZ_tqudR23aFUKmvE71xvvYDh7msV7"+
			"ykInF4VIhcYzLCshNCZtMac1bqP850L4W"+
			"qMZ8CMRUhYXoTRaLhRkT0Z8QkVAPI2VGG"+
			"UqOqodU1JFUxysTlE4sGnVR8GLaSC4d7-"+
			"GFv8ljt4tc1NdR8GJ7UomtPsieL-YxlVH"+
			"a3zhhqyrEBgWHaSwaiS4je_i-o_Xd_JIKT"+
			"tfwBIv_Mx7u4M1h2hB2XdQYisftvnFzXC0"+
			"0F__jbPdCm00\n", plantuml)
	assert.Equal(t, expected, plantuml, string(link))
}
