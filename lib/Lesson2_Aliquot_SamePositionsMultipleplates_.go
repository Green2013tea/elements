// example protocol Lessonshowing how the MixTo command can be used to specify different plates of the same type  i.e. plate 1 ,2, 3 of type greiner384
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol Lesson(data)

// positions on each plate to add aliquots
// number of plates to fill aliquots into

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol Lessonwith types

// Physical outputs from this protocol Lessonwith types

func _Lesson2_Aliquot_SamePositionsMultipleplatesRequirements() {

}

// Conditions to run on startup
func _Lesson2_Aliquot_SamePositionsMultipleplatesSetup(_ctx context.Context, _input *Lesson2_Aliquot_SamePositionsMultipleplatesInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Lesson2_Aliquot_SamePositionsMultipleplatesSteps(_ctx context.Context, _input *Lesson2_Aliquot_SamePositionsMultipleplatesInput, _output *Lesson2_Aliquot_SamePositionsMultipleplatesOutput) {

	aliquots := make([]*wtype.LHComponent, 0)

	// this time round we're going to change the platenumber with each loop until we reach NumberofPlates specified in parameters
	// note that we're starting counting from 1 instead of zero since a platenumber of 0 is an invalid entry into MixTo
	// to ensure we reach the correct number or plates to aliquot to we also therefore need to change the evaluation condition to platenumber < (NumberofPlates +1) rather than platenumber < NumberofPlates
	// alternatively we could have changed the evaluation condition to platenumber <= NumberofPlates
	for platenumber := 1; platenumber < (_input.NumberofPlates + 1); platenumber++ {

		// for each plate we will aliquot to every position specified in the parameters
		// this introduces an alternative syntax for looping through an array using the range keyword
		// a position in the slice Positions can alternatively be accessed with the range command by
		// (i) using the index of which position is required, e.g. Positions[0],Positions[1],Positions[i]
		// using the range syntax this would look like this:
		// for i := range Positions {
		// aliquot := MixTo(OutPlate,Positions[i],platenumber,aliqiotSample)
		// }
		// in that case i starts at 0 and increases by 1 with each loop finishing at len(Positions)
		//(ii) the method as shown below where we use a temporary variable name position for each value of the slice and ignore the index by using the underscore _,
		for _, position := range _input.Positions {
			if _input.Solution.TypeName() == "dna" {
				_input.Solution.Type = wtype.LTDoNotMix
			}
			aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
			// position and platenumber are termporary variables filled in and updated per loop
			aliquot := execute.MixTo(_ctx, _input.OutPlate, position, platenumber, aliquotSample)
			aliquots = append(aliquots, aliquot)
		}
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Lesson2_Aliquot_SamePositionsMultipleplatesAnalysis(_ctx context.Context, _input *Lesson2_Aliquot_SamePositionsMultipleplatesInput, _output *Lesson2_Aliquot_SamePositionsMultipleplatesOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Lesson2_Aliquot_SamePositionsMultipleplatesValidation(_ctx context.Context, _input *Lesson2_Aliquot_SamePositionsMultipleplatesInput, _output *Lesson2_Aliquot_SamePositionsMultipleplatesOutput) {

}
func _Lesson2_Aliquot_SamePositionsMultipleplatesRun(_ctx context.Context, input *Lesson2_Aliquot_SamePositionsMultipleplatesInput) *Lesson2_Aliquot_SamePositionsMultipleplatesOutput {
	output := &Lesson2_Aliquot_SamePositionsMultipleplatesOutput{}
	_Lesson2_Aliquot_SamePositionsMultipleplatesSetup(_ctx, input)
	_Lesson2_Aliquot_SamePositionsMultipleplatesSteps(_ctx, input, output)
	_Lesson2_Aliquot_SamePositionsMultipleplatesAnalysis(_ctx, input, output)
	_Lesson2_Aliquot_SamePositionsMultipleplatesValidation(_ctx, input, output)
	return output
}

func Lesson2_Aliquot_SamePositionsMultipleplatesRunSteps(_ctx context.Context, input *Lesson2_Aliquot_SamePositionsMultipleplatesInput) *Lesson2_Aliquot_SamePositionsMultipleplatesSOutput {
	soutput := &Lesson2_Aliquot_SamePositionsMultipleplatesSOutput{}
	output := _Lesson2_Aliquot_SamePositionsMultipleplatesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Lesson2_Aliquot_SamePositionsMultipleplatesNew() interface{} {
	return &Lesson2_Aliquot_SamePositionsMultipleplatesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Lesson2_Aliquot_SamePositionsMultipleplatesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Lesson2_Aliquot_SamePositionsMultipleplatesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Lesson2_Aliquot_SamePositionsMultipleplatesInput{},
			Out: &Lesson2_Aliquot_SamePositionsMultipleplatesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Lesson2_Aliquot_SamePositionsMultipleplatesElement struct {
	inject.CheckedRunner
}

type Lesson2_Aliquot_SamePositionsMultipleplatesInput struct {
	NumberofPlates   int
	OutPlate         string
	Positions        []string
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Lesson2_Aliquot_SamePositionsMultipleplatesOutput struct {
	Aliquots []*wtype.LHComponent
}

type Lesson2_Aliquot_SamePositionsMultipleplatesSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Lesson2_Aliquot_SamePositionsMultipleplates",
		Constructor: Lesson2_Aliquot_SamePositionsMultipleplatesNew,
		Desc: component.ComponentDesc{
			Desc: "example protocol Lessonshowing how the MixTo command can be used to specify different plates of the same type  i.e. plate 1 ,2, 3 of type greiner384\n",
			Path: "src/github.com/antha-lang/elements/AnthaAcademy/Lesson2_mix/D_AliquotTo_samepositionmultipleplates.an",
			Params: []component.ParamDesc{
				{Name: "NumberofPlates", Desc: "number of plates to fill aliquots into\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Parameters"},
				{Name: "Positions", Desc: "positions on each plate to add aliquots\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
