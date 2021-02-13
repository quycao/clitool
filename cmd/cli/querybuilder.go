package main

import (
    "bufio"
    "fmt"
    "io"
    "strings"

    "github.com/quycao/clitool/pkg/helper"

    survey "github.com/AlecAivazis/survey/v2"
    "github.com/google/uuid"
    "github.com/schollz/progressbar"
    "github.com/skratchdot/open-golang/open"
)

func createQuery(queryHead string, itemSet []string, queryFoot string, itemsEachQuery int, isFinalBatch bool) string {
    numItems := len(itemSet)
    tempQuery := strings.TrimRight(queryHead, "\n") + " ("

    for i, item := range itemSet {
        tempQuery = tempQuery + "'"
        tempQuery = tempQuery + item
        tempQuery = tempQuery + "'"
        if i != numItems-1 && (i+1)%itemsEachQuery != 0 {
            tempQuery = tempQuery + ", "
        } else {
            tempQuery = tempQuery + ") "
            tempQuery = tempQuery + queryFoot
            if len(queryFoot) == 0 {
                tempQuery = tempQuery + "\n"
            }
            if !isFinalBatch {
                tempQuery = tempQuery + "union all"
            }
        }
    }
    return tempQuery
}

func queryBuilder() {
    var qs = []*survey.Question{
        {
            Name: "query_head",
            Prompt: &survey.Multiline{
                Message: "Query content before item set:",
            },
            Validate: survey.Required,
        },
        {
            Name: "item_set",
            Prompt: &survey.Editor{
                Message:  "Item set:",
                FileName: "*.txt",
            },
            Validate: survey.Required,
        },
        {
            Name:   "query_foot",
            Prompt: &survey.Multiline{Message: "Query content after item set:"},
        },
        {
            Name: "item_each_query",
            Prompt: &survey.Input{
                Message: "Number of items in each query:",
                Default: "1000",
            },
            Validate: survey.Required,
        },
        {
            Name: "confirm",
            Prompt: &survey.Confirm{
                Message: "Do you want to build query now:",
                Default: true,
            },
        },
    }

    var errors []error

     // the answers will be written to this struct
     answers := struct {
        QueryHead     string `survey:"query_head"` // survey will match the question and field names
        ItemSet       string `survey:"item_set"`   // or you can tag fields to match a specific name
        QueryFoot     string `survey:"query_foot"` // if the types don't match, survey will convert it
        ItemEachQuery int    `survey:"item_each_query"`
        Confirm       bool   `survey:"confirm"`
     }{}

    // perform the questions
    err := survey.Ask(qs, &answers)
        if err != nil {
            fmt.Println(err.Error())
        return
    }

    if answers.Confirm == true {
        queryFile := fmt.Sprintf("C:\\Temp\\%s.txt", uuid.New().String())

        reader := strings.NewReader(answers.ItemSet)
        numOfLines, err := helper.LineCounter(reader)
        if err != nil {
            errors = append(errors, fmt.Errorf("Line counter error: %s", err))
            helper.PauseProcess(errors)
        }

        // Seek to head of file after lineCounter read to end, prepare for scanner read text
        reader.Seek(0, 0)

        // Init progress bar
        bar := progressbar.New(numOfLines)

        scanner := bufio.NewScanner(reader)

        var items []string
        itemNum := 0
        isFinalBatch := false
        for err == nil {
            items, err = helper.ReadLines(scanner, answers.ItemEachQuery)
            if err != nil && err != io.EOF {
                errors = append(errors, fmt.Errorf("Read items file error: %s", err))
                helper.PauseProcess(errors)
            }

            itemNum = itemNum + len(items)
            if itemNum >= numOfLines {
                isFinalBatch = true
            }

            if len(items) > 0 {
                query := createQuery(answers.QueryHead, items, answers.QueryFoot, answers.ItemEachQuery, isFinalBatch)
                if err := helper.WriteText(query, queryFile); err != nil {
                    errors = append(errors, fmt.Errorf("Write query to file error: %s", err))
                    helper.PauseProcess(errors)
                }

                bar.Add(len(items))
            }
        }
        // bar.Finish()

        if len(errors) == 0 {
            fmt.Printf("\n\nQueries were writen to file: %s\n", queryFile)
            open.Run(queryFile)
        }

        helper.PauseProcess(errors)
    }
}
