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

// Input parameters for this protocol Lesson(data)

// PCRprep parameters:

/*
	// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume

	//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...
	//FullTemplatesequence string // better to use Sid's type system here after proof of concept
	//FullTemplatelength int	// clearly could be calculated from the sequence... Sid will have a method to do this already so check!
	//TargetTemplatesequence string // better to use Sid's type system here after proof of concept
	//TargetTemplatelengthinBP int
*/
// Reaction parameters: (could be a entered as thermocycle parameters type possibly?)

//Denaturationtemp Temperature

// Should be calculated from primer and template binding
// should be calculated from template length and polymerase rate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol Lessonwith types

// e.g. DMSO

// Physical outputs from this protocol Lessonwith types

func _Lesson0_PCRRequirements() {
}

// Conditions to run on startup
func _Lesson0_PCRSetup(_ctx context.Context, _input *Lesson0_PCRInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Lesson0_PCRSteps(_ctx context.Context, _input *Lesson0_PCRInput, _output *Lesson0_PCROutput) {

	// rename components

	_input.Template.CName = _input.TemplateName
	_input.FwdPrimer.CName = _input.FwdPrimerName
	_input.RevPrimer.CName = _input.RevPrimerName

	bufferVolume := (wunit.CopyVolume(_input.ReactionVolume))
	bufferVolume.DivideBy(float64(_input.BufferConcinX))

	// Make a mastermix
	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.Sample(_input.Water, _input.WaterVolume)
	bufferSample := mixer.Sample(_input.Buffer, bufferVolume)
	samples = append(samples, waterSample, bufferSample)

	dntpSample := mixer.Sample(_input.DNTPS, _input.DNTPVol)
	samples = append(samples, dntpSample)

	if len(_input.Additives) != len(_input.AdditiveVols) {
		execute.Errorf(_ctx, "Bad things are going to happen if you have different numbers of additives and additivevolumes")
	}

	for i := range _input.Additives {
		additiveSample := mixer.Sample(_input.Additives[i], _input.AdditiveVols[i])
		samples = append(samples, additiveSample)
	}

	if _input.Hotstart == false {
		polySample := mixer.Sample(_input.PCRPolymerase, _input.PolymeraseVolume)
		samples = append(samples, polySample)
	}

	// if this is true do stuff inside {}
	if _input.AddPrimerstoMasterMix {

		FwdPrimerSample := mixer.Sample(_input.FwdPrimer, _input.FwdPrimerVol)
		samples = append(samples, FwdPrimerSample)
		RevPrimerSample := mixer.Sample(_input.RevPrimer, _input.RevPrimerVol)
		samples = append(samples, RevPrimerSample)

	}

	// pipette out to make mastermix
	var mastermix *wtype.LHComponent
	for j := range samples {
		if j == 0 {
			mastermix = execute.MixInto(_ctx, _input.OutPlate, _input.WellPosition, samples[j])
		} else {
			mastermix = execute.Mix(_ctx, mastermix, samples[j])
		}
	}

	// reset samples to zero
	samples = make([]*wtype.LHComponent, 0)

	// if this is false do stuff inside {}
	if !_input.AddPrimerstoMasterMix {

		FwdPrimerSample := mixer.Sample(_input.FwdPrimer, _input.FwdPrimerVol)
		samples = append(samples, FwdPrimerSample)
		RevPrimerSample := mixer.Sample(_input.RevPrimer, _input.RevPrimerVol)
		samples = append(samples, RevPrimerSample)

	}

	templateSample := mixer.Sample(_input.Template, _input.Templatevolume)
	samples = append(samples, templateSample)

	for j := range samples {
		mastermix = execute.Mix(_ctx, mastermix, samples[j])
	}
	reaction := mastermix

	// this needs to go after an initial denaturation!
	if _input.Hotstart {
		polySample := mixer.Sample(_input.PCRPolymerase, _input.PolymeraseVolume)

		reaction = execute.Mix(_ctx, reaction, polySample)
	}

	// thermocycle parameters called from enzyme lookup:

	polymerase := _input.PCRPolymerase.CName

	extensionTemp := enzymes.DNApolymerasetemps[polymerase]["extensiontemp"]
	meltingTemp := enzymes.DNApolymerasetemps[polymerase]["meltingtemp"]

	// initial Denaturation

	r1 := execute.Incubate(_ctx, reaction, execute.IncubateOpt{
		Temp: meltingTemp,
		Time: _input.InitDenaturationtime,
	})

	for i := 0; i < _input.Numberofcycles; i++ {

		// Denature

		r1 = execute.Incubate(_ctx, r1, execute.IncubateOpt{
			Temp: meltingTemp,
			Time: _input.Denaturationtime,
		})

		// Anneal
		r1 = execute.Incubate(_ctx, r1, execute.IncubateOpt{
			Temp: _input.AnnealingTemp,
			Time: _input.Annealingtime,
		})

		//extensiontime := TargetTemplatelengthinBP/PCRPolymerase.RateBPpers // we'll get type issues here so leave it out for now

		// Extend
		r1 = execute.Incubate(_ctx, r1, execute.IncubateOpt{
			Temp: extensionTemp,
			Time: _input.Extensiontime,
		})

	}
	// Final Extension
	r1 = execute.Incubate(_ctx, r1, execute.IncubateOpt{
		Temp: extensionTemp,
		Time: _input.Finalextensiontime,
	})

	// all done
	_output.Reaction = r1

	_output.Reaction.CName = _input.ReactionName

	_output.Status = "Success"
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Lesson0_PCRAnalysis(_ctx context.Context, _input *Lesson0_PCRInput, _output *Lesson0_PCROutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Lesson0_PCRValidation(_ctx context.Context, _input *Lesson0_PCRInput, _output *Lesson0_PCROutput) {
}
func _Lesson0_PCRRun(_ctx context.Context, input *Lesson0_PCRInput) *Lesson0_PCROutput {
	output := &Lesson0_PCROutput{}
	_Lesson0_PCRSetup(_ctx, input)
	_Lesson0_PCRSteps(_ctx, input, output)
	_Lesson0_PCRAnalysis(_ctx, input, output)
	_Lesson0_PCRValidation(_ctx, input, output)
	return output
}

func Lesson0_PCRRunSteps(_ctx context.Context, input *Lesson0_PCRInput) *Lesson0_PCRSOutput {
	soutput := &Lesson0_PCRSOutput{}
	output := _Lesson0_PCRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Lesson0_PCRNew() interface{} {
	return &Lesson0_PCRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Lesson0_PCRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Lesson0_PCRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Lesson0_PCRInput{},
			Out: &Lesson0_PCROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Lesson0_PCRElement struct {
	inject.CheckedRunner
}

type Lesson0_PCRInput struct {
	AddPrimerstoMasterMix bool
	AdditiveVols          []wunit.Volume
	Additives             []*wtype.LHComponent
	AnnealingTemp         wunit.Temperature
	Annealingtime         wunit.Time
	Buffer                *wtype.LHComponent
	BufferConcinX         int
	DNTPS                 *wtype.LHComponent
	DNTPVol               wunit.Volume
	Denaturationtime      wunit.Time
	Extensiontime         wunit.Time
	Finalextensiontime    wunit.Time
	FwdPrimer             *wtype.LHComponent
	FwdPrimerName         string
	FwdPrimerVol          wunit.Volume
	Hotstart              bool
	InitDenaturationtime  wunit.Time
	Numberofcycles        int
	OutPlate              *wtype.LHPlate
	PCRPolymerase         *wtype.LHComponent
	PolymeraseVolume      wunit.Volume
	ReactionName          string
	ReactionVolume        wunit.Volume
	RevPrimer             *wtype.LHComponent
	RevPrimerName         string
	RevPrimerVol          wunit.Volume
	Template              *wtype.LHComponent
	TemplateName          string
	Templatevolume        wunit.Volume
	Water                 *wtype.LHComponent
	WaterVolume           wunit.Volume
	WellPosition          string
}

type Lesson0_PCROutput struct {
	Reaction *wtype.LHComponent
	Status   string
}

type Lesson0_PCRSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Lesson0_PCR",
		Constructor: Lesson0_PCRNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/AnthaAcademy/Lesson0_Examples/AutoPCR/PCR.an",
			Params: []component.ParamDesc{
				{Name: "AddPrimerstoMasterMix", Desc: "", Kind: "Parameters"},
				{Name: "AdditiveVols", Desc: "", Kind: "Parameters"},
				{Name: "Additives", Desc: "e.g. DMSO\n", Kind: "Inputs"},
				{Name: "AnnealingTemp", Desc: "Should be calculated from primer and template binding\n", Kind: "Parameters"},
				{Name: "Annealingtime", Desc: "Denaturationtemp Temperature\n", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferConcinX", Desc: "", Kind: "Parameters"},
				{Name: "DNTPS", Desc: "", Kind: "Inputs"},
				{Name: "DNTPVol", Desc: "", Kind: "Parameters"},
				{Name: "Denaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "Extensiontime", Desc: "should be calculated from template length and polymerase rate\n", Kind: "Parameters"},
				{Name: "Finalextensiontime", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimer", Desc: "", Kind: "Inputs"},
				{Name: "FwdPrimerName", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "Hotstart", Desc: "", Kind: "Parameters"},
				{Name: "InitDenaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "Numberofcycles", Desc: "\t\t// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume\n\n\t\t//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...\n\t\t//FullTemplatesequence string // better to use Sid's type system here after proof of concept\n\t\t//FullTemplatelength int\t// clearly could be calculated from the sequence... Sid will have a method to do this already so check!\n\t\t//TargetTemplatesequence string // better to use Sid's type system here after proof of concept\n\t\t//TargetTemplatelengthinBP int\n\nReaction parameters: (could be a entered as thermocycle parameters type possibly?)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PCRPolymerase", Desc: "", Kind: "Inputs"},
				{Name: "PolymeraseVolume", Desc: "", Kind: "Parameters"},
				{Name: "ReactionName", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimer", Desc: "", Kind: "Inputs"},
				{Name: "RevPrimerName", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "Template", Desc: "", Kind: "Inputs"},
				{Name: "TemplateName", Desc: "", Kind: "Parameters"},
				{Name: "Templatevolume", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "WaterVolume", Desc: "PCRprep parameters:\n", Kind: "Parameters"},
				{Name: "WellPosition", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}

/*type Polymerase struct {
	LHComponent
	Rate_BPpers float64
	Fidelity_errorrate float64 // could dictate how many colonies are checked in validation!
	Extensiontemp Temperature
	Hotstart bool
	StockConcentration Concentration // this is normally in U?
	TargetConcentration Concentration
	// this is also a glycerol solution rather than a watersolution!
}
*/
