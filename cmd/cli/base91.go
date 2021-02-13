package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ekyu.moe/base91"
	"ekyu.moe/util/cli"
	"github.com/quycao/clitool/pkg/helper"

	survey "github.com/AlecAivazis/survey/v2"
)

func base91Main() {
	var qs = []*survey.Question{
	{
		Name: "mode",
		Prompt: &survey.Select{
			Message: "Choose mode:",
			Options: []string{"Encode", "Decode"},
			Default: "Encode",
		},
	},
	{
		Name: "input",
		Prompt: &survey.Input{
			Message: "In file path:",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		},
		Validate: func(val interface{}) error {
			// since we are validating an Input, the assertion will always succeed
			info, err := os.Stat(val.(string))
			if os.IsNotExist(err) {
				return err
			}
			if info.IsDir() {
				return errors.New("you have entered a directory, not a file")
			}
			return nil
		},
	},
	{
		Name: "output",
		Prompt: &survey.Input{
			Message: "Out file path:",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		},
		Validate: func(val interface{}) error {
			// since we are validating an Input, the assertion will always succeed
			info, err := os.Stat(val.(string))
			if !os.IsNotExist(err) {
				return errors.New("file already exist")
			}
			if info != nil && info.IsDir() {
				return errors.New("you have entered a directory, not a file")
			}
			return nil
		},
	},
	{
		Name: "wrap",
		Prompt: &survey.Input{
			Message: "Wrap encoded lines after COLS character. Use 0 to disable line wrapping:",
			Default: "9000",
		},
		Validate: survey.Required,
	},
	{
		Name: "confirm",
		Prompt: &survey.Confirm{
			Message: "Do you want to start now:",
			Default: true,
		},
	},
}

	var errors []error

	// the answers will be written to this struct
	answers := struct {
		Mode    string `survey:"mode"`   // survey will match the question and field names
		Input   string `survey:"input"`  // or you can tag fields to match a specific name
		Output  string `survey:"output"` // if the types don't match, survey will convert it
		Wrap    int    `survey:"wrap"`
		Confirm bool   `survey:"confirm"`
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if answers.Confirm == true {
		inFilename := answers.Input
		outFilename := answers.Output
		wrap := answers.Wrap

		if len(inFilename) == 0 {
			inFilename = "-"
		}

		if answers.Mode == "Encode" {
			if len(outFilename) == 0 && inFilename != "-" {
				outFilename = inFilename + ".asc"
			}
			err = b91encode(outFilename, inFilename, wrap)
		} else {
			if len(outFilename) == 0 && inFilename != "-" {
				if strings.HasSuffix(strings.ToLower(inFilename), ".asc") {
					outFilename = strings.TrimSuffix(inFilename, ".asc")
				} else {
					outFilename = inFilename + ".plain"
				}
				outFilename = inFilename + ".asc"
			}
			err = b91decode(outFilename, inFilename)
		}

		if err != nil {
			errors = append(errors, err)
		} else {
			fmt.Printf("\n%s complete!\n", answers.Mode)
		}

		helper.PauseProcess(errors)
	}
}

func b91encode(outFilename, inFilename string, wrap int) error {
	// validate and read in file
	inFile, _, err := cli.AccessOpenFile(inFilename)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// validate and create out file
	outFile, err := cli.PromptOverrideCreate(outFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var e io.WriteCloser
	if wrap <= 0 {
		e = base91.NewEncoder(outFile)
	} else {
		e = base91.NewLineWrapper(outFile, wrap)
	}
	defer e.Close()

	if _, err := io.Copy(e, inFile); err != nil {
		return err
	}

	return nil
}

func b91decode(outFilename, inFilename string) error {
	// validate and read in file
	inFile, _, err := cli.AccessOpenFile(inFilename)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// validate and create out file
	outFile, err := cli.PromptOverrideCreate(outFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	d := base91.NewDecoder(inFile)

	if _, err := io.Copy(outFile, d); err != nil {
		return err
	}

	return nil
}
