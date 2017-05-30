// Protocol for resuspending freeze dried DNA with a diluent
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Mass of DNA to Resuspend

// Target concentration to resuspend to

// Molecular weight of the DNA

// Well location of DNA

// Plate location of DNA

// If no Policy is specified the default policy will be MegaMix which mixes the sample 10 times.

// Diluent to use to resuspend the DNA

// Type of plate DNA sample is on.

func _ResuspendDNARequirements() {
}

func _ResuspendDNASetup(_ctx context.Context, _input *ResuspendDNAInput) {
}

func _ResuspendDNASteps(_ctx context.Context, _input *ResuspendDNAInput, _output *ResuspendDNAOutput) {

	targetconcgperL := _input.TargetConc.GramPerL(_input.MolecularWeight).SIValue()

	dnamassG := _input.DNAMass.SIValue()

	if _input.DNAMass.Unit().BaseSIUnit() == "kg" {
		dnamassG = dnamassG * 1000
		_output.Warnings = append(_output.Warnings, fmt.Sprintln("Base Unit correction; Base unit of mass = ", _input.DNAMass.Unit().BaseSIUnit(), " therfore multiplying by 1000 to convert to grams"))
	}

	volumetoadd := wunit.NewVolume(dnamassG/targetconcgperL, "L")

	diluentSample := mixer.Sample(_input.Diluent, volumetoadd)

	if _input.OverRideLiquidPolicy == "" {
		_input.OverRideLiquidPolicy = "MegaMix"
	}

	var err error

	diluentSample.Type, err = wtype.LiquidTypeFromString(_input.OverRideLiquidPolicy)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	_output.ResuspendedDNA = execute.MixNamed(_ctx, _input.DNAPlate.Type, _input.Well, _input.PlateName, diluentSample)

}

func _ResuspendDNAAnalysis(_ctx context.Context, _input *ResuspendDNAInput, _output *ResuspendDNAOutput) {
}

func _ResuspendDNAValidation(_ctx context.Context, _input *ResuspendDNAInput, _output *ResuspendDNAOutput) {
}
func _ResuspendDNARun(_ctx context.Context, input *ResuspendDNAInput) *ResuspendDNAOutput {
	output := &ResuspendDNAOutput{}
	_ResuspendDNASetup(_ctx, input)
	_ResuspendDNASteps(_ctx, input, output)
	_ResuspendDNAAnalysis(_ctx, input, output)
	_ResuspendDNAValidation(_ctx, input, output)
	return output
}

func ResuspendDNARunSteps(_ctx context.Context, input *ResuspendDNAInput) *ResuspendDNASOutput {
	soutput := &ResuspendDNASOutput{}
	output := _ResuspendDNARun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ResuspendDNANew() interface{} {
	return &ResuspendDNAElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ResuspendDNAInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ResuspendDNARun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ResuspendDNAInput{},
			Out: &ResuspendDNAOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ResuspendDNAElement struct {
	inject.CheckedRunner
}

type ResuspendDNAInput struct {
	DNAMass              wunit.Mass
	DNAPlate             *wtype.LHPlate
	Diluent              *wtype.LHComponent
	MolecularWeight      float64
	OverRideLiquidPolicy wtype.PolicyName
	PlateName            string
	TargetConc           wunit.Concentration
	Well                 string
}

type ResuspendDNAOutput struct {
	ResuspendedDNA *wtype.LHComponent
	Warnings       []string
}

type ResuspendDNASOutput struct {
	Data struct {
		Warnings []string
	}
	Outputs struct {
		ResuspendedDNA *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ResuspendDNA",
		Constructor: ResuspendDNANew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/ResuspendDNA/ResuspendDNA.an",
			Params: []component.ParamDesc{
				{Name: "DNAMass", Desc: "Mass of DNA to Resuspend\n", Kind: "Parameters"},
				{Name: "DNAPlate", Desc: "Type of plate DNA sample is on.\n", Kind: "Inputs"},
				{Name: "Diluent", Desc: "Diluent to use to resuspend the DNA\n", Kind: "Inputs"},
				{Name: "MolecularWeight", Desc: "Molecular weight of the DNA\n", Kind: "Parameters"},
				{Name: "OverRideLiquidPolicy", Desc: "If no Policy is specified the default policy will be MegaMix which mixes the sample 10 times.\n", Kind: "Parameters"},
				{Name: "PlateName", Desc: "Plate location of DNA\n", Kind: "Parameters"},
				{Name: "TargetConc", Desc: "Target concentration to resuspend to\n", Kind: "Parameters"},
				{Name: "Well", Desc: "Well location of DNA\n", Kind: "Parameters"},
				{Name: "ResuspendedDNA", Desc: "", Kind: "Outputs"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
