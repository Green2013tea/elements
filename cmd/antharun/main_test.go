package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

const (
	testTimeout = 5 * time.Second
)

func makeContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())
	for _, desc := range library {
		obj := desc.Constructor()
		runner, ok := obj.(inject.Runner)
		if !ok {
			return nil, fmt.Errorf("component %q has unexpected type %T", desc.Name, obj)
		}
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name}, runner); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func runTestInput(t *testing.T, ctx context.Context, tgt *target.Target, input *executeutil.TestInput) error {
	errs := make(chan error)
	go func() {
		// HACK(ddn): Sink chdir inside goroutine to "improve" chances that
		// golang scheduler puts this goroutine on the os thread
		// corresponding to the chdir call.
		//
		// Until elements are refactored to not know their working
		// directory we can't "go test parallel" these tests
		if len(input.Dir) != 0 {
			if err := os.Chdir(input.Dir); err != nil {
				errs <- err
				return
			}
		}
		_, err := execute.Run(ctx, execute.Opt{
			Workflow: input.Workflow,
			Params:   input.Params,
			Target:   tgt,
		})
		errs <- err
	}()

	var err error
	select {
	case err = <-errs:
	case <-time.After(testTimeout):
		// TODO(ddn): reenable timeouts
		//err = fmt.Errorf("timeout after %ds", testTimeout/time.Second)
	}

	if err == nil {
		return nil
	} else if _, ok := err.(*execute.Error); ok {
		return nil
	} else {
		return fmt.Errorf("error running %s: %s", inputLabel(input), err)
	}
}

func inputLabel(input *executeutil.TestInput) string {
	if len(input.BundlePath) != 0 {
		return fmt.Sprintf("bundle %q", input.BundlePath)
	}
	return fmt.Sprintf("workflow %q with parameters %q", input.WorkflowPath, input.ParamsPath)
}

func runElements(t *testing.T, ctx context.Context, inputs []*executeutil.TestInput) {
	tgt := target.New()
	tgt.AddDevice(human.New(human.Opt{CanMix: true, CanIncubate: true, CanHandle: true}))

	odir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range inputs {
		in := input
		t.Run(inputLabel(in), func(t *testing.T) { t.Parallel(); runTestInput(t, ctx, tgt, in) })
	}

	if err := os.Chdir(odir); err != nil {
		t.Fatal(err)
	}
}

func findInputs(basePaths ...string) ([]*executeutil.TestInput, error) {
	var inputDirs []string
	for _, c := range basePaths {
		_, err := os.Stat(c)
		if err == nil {
			inputDirs = append(inputDirs, c)
		}
	}

	if len(inputDirs) == 0 {
		return nil, fmt.Errorf("could not find example inputs in %v", basePaths)
	}

	var inputs []*executeutil.TestInput
	for _, dir := range inputDirs {
		ins, err := executeutil.FindTestInputs(dir)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, ins...)
	}

	return inputs, nil
}

func TestElementsWithExampleInputs(t *testing.T) {
	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs, err := findInputs("../../workflows", "workflows", "../../defaults", "defaults")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("found %d test inputs\n", len(inputs))

	runElements(t, ctx, inputs)
}
