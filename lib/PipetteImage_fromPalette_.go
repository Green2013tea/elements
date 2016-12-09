// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/microArch/factory"
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"image/color"
	"strconv"
)

// Input parameters for this protocol (data)

//AvailableColours []string

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PipetteImage_fromPaletteRequirements() {

}

// Conditions to run on startup
func _PipetteImage_fromPaletteSetup(_ctx context.Context, _input *PipetteImage_fromPaletteInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_fromPaletteSteps(_ctx context.Context, _input *PipetteImage_fromPaletteInput, _output *PipetteImage_fromPaletteOutput) {

	var err error

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &_input.Palette, _input.Rotate, _input.AutoRotate)

	image.CheckAllResizealgorithms(_input.Imagefilename, _input.OutPlate, _input.Rotate, image.AllResampleFilters)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	for locationkey, colour := range positiontocolourmap {

		colourindex := strconv.Itoa(_input.Palette.Index(colour))

		component, componentpresent := _input.ColourIndextoComponentMap[colourindex]

		if componentpresent {

			if _input.LiquidType != "" {
				component.Type, err = wtype.LiquidTypeFromString(_input.LiquidType) //"DoNotMix"

				if err != nil {
					execute.Errorf(_ctx, "for component", component.CName, err.Error())
				}
			}
			if _input.OnlythisColour != "" {

				if image.Colourcomponentmap[colour] == _input.OnlythisColour {
					counter = counter + 1

					pixelSample := mixer.Sample(component, _input.VolumePerWell)
					solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
					solutions = append(solutions, solution)
				}

			} else {
				if component.CName != _input.NotthisColour {
					counter = counter + 1
					pixelSample := mixer.Sample(component, _input.VolumePerWell)
					solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
					solutions = append(solutions, solution)
				}
			}

		}

	}
	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_fromPaletteAnalysis(_ctx context.Context, _input *PipetteImage_fromPaletteInput, _output *PipetteImage_fromPaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_fromPaletteValidation(_ctx context.Context, _input *PipetteImage_fromPaletteInput, _output *PipetteImage_fromPaletteOutput) {

}
func _PipetteImage_fromPaletteRun(_ctx context.Context, input *PipetteImage_fromPaletteInput) *PipetteImage_fromPaletteOutput {
	output := &PipetteImage_fromPaletteOutput{}
	_PipetteImage_fromPaletteSetup(_ctx, input)
	_PipetteImage_fromPaletteSteps(_ctx, input, output)
	_PipetteImage_fromPaletteAnalysis(_ctx, input, output)
	_PipetteImage_fromPaletteValidation(_ctx, input, output)
	return output
}

func PipetteImage_fromPaletteRunSteps(_ctx context.Context, input *PipetteImage_fromPaletteInput) *PipetteImage_fromPaletteSOutput {
	soutput := &PipetteImage_fromPaletteSOutput{}
	output := _PipetteImage_fromPaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_fromPaletteNew() interface{} {
	return &PipetteImage_fromPaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_fromPaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_fromPaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_fromPaletteInput{},
			Out: &PipetteImage_fromPaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_fromPaletteElement struct {
	inject.CheckedRunner
}

type PipetteImage_fromPaletteInput struct {
	AutoRotate                bool
	ColourIndextoComponentMap map[string]*wtype.LHComponent
	Colourcomponents          []*wtype.LHComponent
	Imagefilename             string
	LiquidType                string
	NotthisColour             string
	OnlythisColour            string
	OutPlate                  *wtype.LHPlate
	Palette                   color.Palette
	PosterizeImage            bool
	PosterizeLevels           int
	Rotate                    bool
	VolumePerWell             wunit.Volume
}

type PipetteImage_fromPaletteOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHComponent
}

type PipetteImage_fromPaletteSOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_fromPalette",
		Constructor: PipetteImage_fromPaletteNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PipetteImage_fromPalette.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "ColourIndextoComponentMap", Desc: "", Kind: "Parameters"},
				{Name: "Colourcomponents", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "LiquidType", Desc: "", Kind: "Parameters"},
				{Name: "NotthisColour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "AvailableColours []string\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palette", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
