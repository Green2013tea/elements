package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _TypeIISConstructAssemblyMMX_forscreen_transformRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblyMMX_forscreen_transformSetup(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreen_transformInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblyMMX_forscreen_transformSteps(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreen_transformInput, _output *TypeIISConstructAssemblyMMX_forscreen_transformOutput) {
	var err error

	samples := make([]*wtype.LHComponent, 0)
	_output.ConstructName = _input.OutputConstructName

	last := len(_input.PartSeqs) - 1
	output, count, _, seq, err := enzymes.Assemblysimulator(enzymes.Assemblyparameters{
		Constructname: _output.ConstructName,
		Enzymename:    _input.EnzymeName,
		Vector:        _input.PartSeqs[last],
		Partsinorder:  _input.PartSeqs[:last],
	})
	_output.Output = output
	if err != nil {
		//  Errorf("%s: %s", output, err)
	}
	if count != 1 {
		//  Errorf("no successful assembly")
	}

	_output.Sequence = seq

	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	for k, part := range _input.Parts {
		part.Type, err = wtype.LiquidTypeFromString(_input.LHPolicyName)

		if err != nil {
			execute.Errorf(_ctx, "cannot find liquid type: %s", err)
		}

		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartSeqs[k].Nm
		samples = append(samples, partSample)
	}

	mmxSample := mixer.Sample(_input.MasterMix, _input.MasterMixVolume)
	samples = append(samples, mmxSample)

	// ensure the last step is mixed
	samples[len(samples)-1].Type = wtype.LTDNAMIX
	_output.Reaction = execute.MixTo(_ctx, _input.OutPlate.Type, _input.OutputLocation, _input.OutputPlateNum, samples...)
	_output.Reaction.Extra["label"] = _output.ConstructName

	dnaSample := mixer.Sample(_output.Reaction, _input.TransformationVolume)

	execute.Incubate(_ctx, dnaSample, _input.ReactionTemp, _input.ReactionTime, false)

	transformation := execute.MixTo(_ctx, _input.PlateWithCompetentCells.Type, _input.CompetentCellPlateWell, 1, dnaSample)

	execute.Incubate(_ctx, transformation, _input.PostPlasmidTemp, _input.PostPlasmidTime, false)

	transformationSample := mixer.Sample(transformation, _input.CompetentCellTransferVolume)

	_output.Recovery = execute.MixTo(_ctx, _input.PlatewithRecoveryMedia.Type, _input.RecoveryPlateWell, _input.RecoveryPlateNumber, transformationSample)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	execute.Incubate(_ctx, _output.Recovery, _input.RecoveryTemp, _input.RecoveryTime, true)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyMMX_forscreen_transformAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreen_transformInput, _output *TypeIISConstructAssemblyMMX_forscreen_transformOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyMMX_forscreen_transformValidation(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreen_transformInput, _output *TypeIISConstructAssemblyMMX_forscreen_transformOutput) {
}
func _TypeIISConstructAssemblyMMX_forscreen_transformRun(_ctx context.Context, input *TypeIISConstructAssemblyMMX_forscreen_transformInput) *TypeIISConstructAssemblyMMX_forscreen_transformOutput {
	output := &TypeIISConstructAssemblyMMX_forscreen_transformOutput{}
	_TypeIISConstructAssemblyMMX_forscreen_transformSetup(_ctx, input)
	_TypeIISConstructAssemblyMMX_forscreen_transformSteps(_ctx, input, output)
	_TypeIISConstructAssemblyMMX_forscreen_transformAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyMMX_forscreen_transformValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyMMX_forscreen_transformRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyMMX_forscreen_transformInput) *TypeIISConstructAssemblyMMX_forscreen_transformSOutput {
	soutput := &TypeIISConstructAssemblyMMX_forscreen_transformSOutput{}
	output := _TypeIISConstructAssemblyMMX_forscreen_transformRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyMMX_forscreen_transformNew() interface{} {
	return &TypeIISConstructAssemblyMMX_forscreen_transformElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyMMX_forscreen_transformInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyMMX_forscreen_transformRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyMMX_forscreen_transformInput{},
			Out: &TypeIISConstructAssemblyMMX_forscreen_transformOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyMMX_forscreen_transformElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyMMX_forscreen_transformInput struct {
	CompetentCellPlateWell      string
	CompetentCellTransferVolume wunit.Volume
	EnzymeName                  string
	InactivationTemp            wunit.Temperature
	InactivationTime            wunit.Time
	LHPolicyName                wtype.PolicyName
	MasterMix                   *wtype.LHComponent
	MasterMixVolume             wunit.Volume
	OutPlate                    *wtype.LHPlate
	OutputConstructName         string
	OutputLocation              string
	OutputPlateNum              int
	OutputReactionName          string
	PartSeqs                    []wtype.DNASequence
	PartVols                    []wunit.Volume
	Parts                       []*wtype.LHComponent
	PlateWithCompetentCells     *wtype.LHPlate
	PlatewithRecoveryMedia      *wtype.LHPlate
	PostPlasmidTemp             wunit.Temperature
	PostPlasmidTime             wunit.Time
	ReactionTemp                wunit.Temperature
	ReactionTime                wunit.Time
	ReactionVolume              wunit.Volume
	RecoveryPlateNumber         int
	RecoveryPlateWell           string
	RecoveryTemp                wunit.Temperature
	RecoveryTime                wunit.Time
	TransformationVolume        wunit.Volume
	Water                       *wtype.LHComponent
}

type TypeIISConstructAssemblyMMX_forscreen_transformOutput struct {
	ConstructName string
	Output        string
	Reaction      *wtype.LHComponent
	Recovery      *wtype.LHComponent
	Sequence      wtype.DNASequence
}

type TypeIISConstructAssemblyMMX_forscreen_transformSOutput struct {
	Data struct {
		ConstructName string
		Output        string
		Sequence      wtype.DNASequence
	}
	Outputs struct {
		Reaction *wtype.LHComponent
		Recovery *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssemblyMMX_forscreen_transform",
		Constructor: TypeIISConstructAssemblyMMX_forscreen_transformNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX_forscreen/TypeIISConstructAssemblyMMX_transform.an",
			Params: []component.ParamDesc{
				{Name: "CompetentCellPlateWell", Desc: "", Kind: "Parameters"},
				{Name: "CompetentCellTransferVolume", Desc: "", Kind: "Parameters"},
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "MasterMixVolume", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputConstructName", Desc: "", Kind: "Parameters"},
				{Name: "OutputLocation", Desc: "", Kind: "Parameters"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "", Kind: "Parameters"},
				{Name: "PartSeqs", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "PlateWithCompetentCells", Desc: "", Kind: "Inputs"},
				{Name: "PlatewithRecoveryMedia", Desc: "", Kind: "Inputs"},
				{Name: "PostPlasmidTemp", Desc: "", Kind: "Parameters"},
				{Name: "PostPlasmidTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryPlateNumber", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryPlateWell", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryTemp", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryTime", Desc: "", Kind: "Parameters"},
				{Name: "TransformationVolume", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "ConstructName", Desc: "", Kind: "Data"},
				{Name: "Output", Desc: "", Kind: "Data"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
				{Name: "Recovery", Desc: "", Kind: "Outputs"},
				{Name: "Sequence", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
