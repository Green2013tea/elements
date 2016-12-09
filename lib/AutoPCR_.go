package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

// PCRprep parameters

// e.g. ["left homology arm"]:"templatename"
// e.g. ["left homology arm"]:"fwdprimer","revprimer"

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoPCRRequirements() {
}

// Conditions to run on startup
func _AutoPCRSetup(_ctx context.Context, _input *AutoPCRInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCRSteps(_ctx context.Context, _input *AutoPCRInput, _output *AutoPCROutput) {

	// set up a counter to use as an index for increasing well position
	var counter int

	// set up some empty slices to fill as we iterate through the reactions
	_output.Reactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	// range through the Reaction to template map
	for reactionname, templatename := range _input.Reactiontotemplate {

		// use counter to find next available well position in plate
		wellposition := _input.Plate.AllWellPositions(wtype.BYCOLUMN)[counter]

		// Run PCR_vol element
		result := PCR_volRunSteps(_ctx, &PCR_volInput{WaterVolume: wunit.NewVolume(10, "ul"),
			ReactionVolume:        wunit.NewVolume(25, "ul"),
			BufferConcinX:         5,
			FwdPrimerName:         _input.Reactiontoprimerpair[reactionname][0],
			RevPrimerName:         _input.Reactiontoprimerpair[reactionname][1],
			TemplateName:          templatename,
			ReactionName:          reactionname,
			FwdPrimerVol:          wunit.NewVolume(1, "ul"),
			RevPrimerVol:          wunit.NewVolume(1, "ul"),
			AdditiveVols:          []wunit.Volume{wunit.NewVolume(5, "ul")},
			Templatevolume:        wunit.NewVolume(1, "ul"),
			PolymeraseVolume:      wunit.NewVolume(1, "ul"),
			DNTPVol:               wunit.NewVolume(1, "ul"),
			Numberofcycles:        30,
			InitDenaturationtime:  wunit.NewTime(30, "s"),
			Denaturationtime:      wunit.NewTime(5, "s"),
			Annealingtime:         wunit.NewTime(10, "s"),
			AnnealingTemp:         wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
			Extensiontime:         wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
			Finalextensiontime:    wunit.NewTime(180, "s"),
			Hotstart:              false,
			AddPrimerstoMasterMix: false,
			WellPosition:          wellposition,

			FwdPrimer:     _input.FwdPrimertype,
			RevPrimer:     _input.RevPrimertype,
			DNTPS:         factory.GetComponentByType("DNTPs"),
			PCRPolymerase: factory.GetComponentByType("Q5Polymerase"),
			Buffer:        factory.GetComponentByType("Q5buffer"),
			Water:         factory.GetComponentByType("water"),
			Template:      _input.Templatetype,
			Additives:     []*wtype.LHComponent{factory.GetComponentByType("GCenhancer")},
			OutPlate:      _input.Plate},
		)

		// add result to reactions slice
		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
		volumes = append(volumes, result.Outputs.Reaction.Volume())
		welllocations = append(welllocations, wellposition)
		// increase counter by 1 ready for next iteration of loop
		counter++

	}

	// once all values of loop have been completed, export the plate contents as a csv file
	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCRAnalysis(_ctx context.Context, _input *AutoPCRInput, _output *AutoPCROutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCRValidation(_ctx context.Context, _input *AutoPCRInput, _output *AutoPCROutput) {
}
func _AutoPCRRun(_ctx context.Context, input *AutoPCRInput) *AutoPCROutput {
	output := &AutoPCROutput{}
	_AutoPCRSetup(_ctx, input)
	_AutoPCRSteps(_ctx, input, output)
	_AutoPCRAnalysis(_ctx, input, output)
	_AutoPCRValidation(_ctx, input, output)
	return output
}

func AutoPCRRunSteps(_ctx context.Context, input *AutoPCRInput) *AutoPCRSOutput {
	soutput := &AutoPCRSOutput{}
	output := _AutoPCRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCRNew() interface{} {
	return &AutoPCRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCRInput{},
			Out: &AutoPCROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoPCRElement struct {
	inject.CheckedRunner
}

type AutoPCRInput struct {
	FwdPrimertype        *wtype.LHComponent
	Plate                *wtype.LHPlate
	Projectname          string
	Reactiontoprimerpair map[string][]string
	Reactiontotemplate   map[string]string
	RevPrimertype        *wtype.LHComponent
	Templatetype         *wtype.LHComponent
}

type AutoPCROutput struct {
	Error     error
	Reactions []*wtype.LHComponent
}

type AutoPCRSOutput struct {
	Data struct {
		Error error
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPCR",
		Constructor: AutoPCRNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PCR/AutoPCR.an",
			Params: []component.ParamDesc{
				{Name: "FwdPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "Projectname", Desc: "PCRprep parameters\n", Kind: "Parameters"},
				{Name: "Reactiontoprimerpair", Desc: "e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "Reactiontotemplate", Desc: "e.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "RevPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Templatetype", Desc: "", Kind: "Inputs"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
