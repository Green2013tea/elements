// This protocol is intended to design assembly parts using the MoClo assembly standard.
// Overhangs for a part are chosen according to the designated class of each part (e.g. promoter).
// The MoClo standard is hierarchical so the enzyme is chosen based on the level of assembly.
// i.e. first level 0 parts are made which may comprise of a promoter, 5prime upstream part, coding sequene, and terminator.
// Level 0 parts can then be assembled together by using level 1 enzymes and overhangs.
// currently this protocol only supports level 0 steps.
// see http://journals.plos.org/plosone/article?id=10.1371/journal.pone.0016765
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strings"
)

// Input parameters for this protocol (data)

//MoClo
//Level						string // of assembly standard
//PartMoClotypesinorder		[]string // labels e.g. pro = promoter

//string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSetup(_ctx context.Context, _input *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSteps(_ctx context.Context, _input *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput, _output *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput) {
	//var msg string
	// set warnings reported back to user to none initially
	warnings := make([]string, 1)
	warnings[0] = "none"
	var err error

	/* find sequence data from keyword; looking it up by a given name in an inventory
	   or by biobrick ID from iGem parts registry */
	partsinorder := make([]wtype.DNASequence, 0)

	partsinorder = append(partsinorder, _input.Promoter, _input.Rbs, _input.CDS, _input.Terminator)

	vectordata := _input.Vector

	Level := "Level0"

	//lookup restriction enzyme
	restrictionenzyme := enzymes.Enzymelookup[_input.AssemblyStandard][Level]

	// (1) Add standard overhangs using chosen assembly standard
	_output.PartswithOverhangs, err = enzymes.MakeStandardTypeIIsassemblyParts(partsinorder, _input.AssemblyStandard, Level, []string{"Pro", "5U + NT1", "CDS1", "3U + Ter"})
	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
		execute.Errorf(_ctx, err.Error())
	}

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
	status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

	endreport := "Only run in the event of assembly failure"
	_output.NewDNASequence = newDNASequence
	if err == nil && numberofassemblies == 1 {

		_output.Simulationpass = true
	} else {
		warnings = append(warnings, status)
		// perform mock digest to test fragement overhangs (fragments are hidden by using _, )
		_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

		allends := make([]string, 0)
		ends := ""

		ends = text.Print(vectordata.Nm+" 5 Prime end: ", stickyends5)
		allends = append(allends, ends)
		ends = text.Print(vectordata.Nm+" 3 Prime end: ", stickyends3)
		allends = append(allends, ends)

		for _, part := range _output.PartswithOverhangs {
			_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(part, restrictionenzyme)
			ends = text.Print(part.Nm+" 5 Prime end: ", stickyends5)
			allends = append(allends, ends)
			ends = text.Print(part.Nm+" 3 Prime end: ", stickyends3)
			allends = append(allends, ends)
		}
		endreport = strings.Join(allends, " ")
	}

	_output.Warnings = strings.Join(warnings, ";")

	partsummary := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {
		partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
	}

	partstoorder := text.Print("PartswithOverhangs: ", partsummary)

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Warnings:", _output.Warnings),
			text.Print("Simulationpass=", _output.Simulationpass),
			text.Print("NewDNASequence: ", _output.NewDNASequence),
			partstoorder,
		)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORAnalysis(_ctx context.Context, _input *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput, _output *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORValidation(_ctx context.Context, _input *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput, _output *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput) {
}
func _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORRun(_ctx context.Context, input *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput) *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput {
	output := &MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput{}
	_MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSetup(_ctx, input)
	_MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSteps(_ctx, input, output)
	_MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORAnalysis(_ctx, input, output)
	_MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORValidation(_ctx, input, output)
	return output
}

func MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORRunSteps(_ctx context.Context, input *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput) *MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSOutput {
	soutput := &MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSOutput{}
	output := _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORNew() interface{} {
	return &MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput{},
			Out: &MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORElement struct {
	inject.CheckedRunner
}

type MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORInput struct {
	AssemblyStandard string
	CDS              wtype.DNASequence
	Constructname    string
	Partsinorder     []string
	Promoter         wtype.DNASequence
	Rbs              wtype.DNASequence
	Terminator       wtype.DNASequence
	Vector           wtype.DNASequence
}

type MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOROutput struct {
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           string
}

type MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORSOutput struct {
	Data struct {
		NewDNASequence     wtype.DNASequence
		PartswithOverhangs []wtype.DNASequence
		Simulationpass     bool
		Status             string
		Warnings           string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTOR",
		Constructor: MoCloLevel0_PRO_RBS_CDS_TERMINATOR_VECTORNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol is intended to design assembly parts using the MoClo assembly standard.\nOverhangs for a part are chosen according to the designated class of each part (e.g. promoter).\nThe MoClo standard is hierarchical so the enzyme is chosen based on the level of assembly.\ni.e. first level 0 parts are made which may comprise of a promoter, 5prime upstream part, coding sequene, and terminator.\nLevel 0 parts can then be assembled together by using level 1 enzymes and overhangs.\ncurrently this protocol only supports level 0 steps.\nsee http://journals.plos.org/plosone/article?id=10.1371/journal.pone.0016765\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/TypeIISAssembly_design/MoClo_level0.an",
			Params: []component.ParamDesc{
				{Name: "AssemblyStandard", Desc: "MoClo\n", Kind: "Parameters"},
				{Name: "CDS", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "Partsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Promoter", Desc: "Level\t\t\t\t\t\tstring // of assembly standard\nPartMoClotypesinorder\t\t[]string // labels e.g. pro = promoter\n", Kind: "Parameters"},
				{Name: "Rbs", Desc: "", Kind: "Parameters"},
				{Name: "Terminator", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "string\n", Kind: "Parameters"},
				{Name: "NewDNASequence", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "PartswithOverhangs", Desc: "parts to order\n", Kind: "Data"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
