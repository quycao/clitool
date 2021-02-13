package main

import (
	"errors"

	"github.com/quycao/clitool/pkg/helper"

	survey "github.com/AlecAivazis/survey/v2"
)

func main() {
    var qs = &survey.Select{
        Message: "Choose function:",
        Options: []string{"Message builder", "Query builder (IN)", "base64", "base91"},
        Default: "Message builder",
    }

    var fn string

    survey.AskOne(qs, &fn)

    if fn == "Message builder" {
        msgBuilder()
    } else if fn == "Query builder (IN)" {
        queryBuilder()
    } else if fn == "base64" {
        base64Main()
    } else if fn == "base91" {
        base91Main()
    } else {
        errors := []error{errors.New("Function does not exist")}
        helper.PauseProcess(errors)
    }
}