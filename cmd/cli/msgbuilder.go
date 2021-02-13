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

func createMessage(msgHead string, itemSet []string, msgFoot string, itemsEachMsg int) string {
    numItems := len(itemSet)
    tempMsg := msgHead
    for i, item := range itemSet {
        tempMsg = tempMsg + item
        if i != numItems-1 && (i+1)%itemsEachMsg != 0 {
            tempMsg = tempMsg + "\n"
        } else {
            tempMsg = tempMsg + "\n"
            tempMsg = tempMsg + msgFoot
            // if i < numItems-1 {
            tempMsg = tempMsg + "\n\n"
            tempMsg = tempMsg + "sleep 5"
            tempMsg = tempMsg + "\n\n"
            // tempMsg = tempMsg + msgHead + "\n"
            // }
        }
    }
    return tempMsg
}

func msgBuilder() {
    var qs = []*survey.Question{
        {
            Name: "msg_head",
            Prompt: &survey.Multiline{
                Message: "Message content before item set:",
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
        // {
        //  Name: "item_set",
        //  Prompt: &survey.Input{
        //   Message: "Input file path that contain item set: ",
        //   Suggest: func(toComplete string) []string {
        //    files, _ := filepath.Glob(toComplete + "*")
        //    return files
        //   },
        //  },
        //  Validate: func(val interface{}) error {
        //   // since we are validating an Input, the assertion will always succeed
        //   info, err := os.Stat(val.(string))
        //   if os.IsNotExist(err) {
        //    return err
        //   }
        //   if info.IsDir() {
        //    return errors.New("you have entered a directory, not a file")
        //   }
        //   return nil
        //  },
        // },
        {
            Name:     "msg_foot",
            Prompt:   &survey.Multiline{Message: "Message content after item set:"},
            Validate: survey.Required,
        },
        {
            Name: "item_each_msg",
            Prompt: &survey.Input{
                Message: "Number of items in each msg:",
                Default: "200",
            },
            Validate: survey.Required,
        },
        {
            Name: "confirm",
            Prompt: &survey.Confirm{
                Message: "Do you want to build Message now:",
                Default: true,
            },
        },
    }

    var errors []error

    // the answers will be written to this struct
    answers := struct {
        MsgHead     string `survey:"msg_head"` // survey will match the question and field names
        ItemSet     string `survey:"item_set"` // or you can tag fields to match a specific name
        MsgFoot     string `survey:"msg_foot"` // if the types don't match, survey will convert it
        ItemEachMsg int    `survey:"item_each_msg"`
        Confirm     bool   `survey:"confirm"`
    }{}

    // perform the questions
    err := survey.Ask(qs, &answers)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if answers.Confirm == true {
        // dir := filepath.Dir(answers.ItemSet)
        // fileName := filepath.Base(answers.ItemSet)
        // newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "_result" + filepath.Ext(fileName)
        // msgFile := filepath.Join(dir, newFileName)
        msgFile := fmt.Sprintf("C:\\Temp\\%s.txt", uuid.New().String())

        // file, err := os.Open(answers.ItemSet)
        // if err != nil {
        //  errors = append(errors, fmt.Errorf("Cannot open file: %s", err))
        //  pauseProcess(errors)
        // }
        // defer file.Close()

        reader := strings.NewReader(answers.ItemSet)
        numOfLines, err := helper.LineCounter(reader)
        if err != nil {
            errors = append(errors, fmt.Errorf("Line counter error: %s", err))
            helper.PauseProcess(errors)
        }

        // Seek to head of file after lineCounter read to end, prepare for scanner read text
        // file.Seek(0, 0)
        reader.Seek(0, 0)

        // Init progress bar
        bar := progressbar.New(numOfLines)

        scanner := bufio.NewScanner(reader)

        var items []string
        for err == nil {
            items, err = helper.ReadLines(scanner, answers.ItemEachMsg)

            if err != nil && err != io.EOF {
                errors = append(errors, fmt.Errorf("Read items file error: %s", err))
                helper.PauseProcess(errors)
            }

            if len(items) > 0 {
                msg := createMessage(answers.MsgHead, items, answers.MsgFoot, answers.ItemEachMsg)
                if err := helper.WriteText(msg, msgFile); err != nil {
                    errors = append(errors, fmt.Errorf("Write msg to file error: %s", err))
                    helper.PauseProcess(errors)
                }

                bar.Add(len(items))
            }
        }
        // bar.Finish()

        if len(errors) == 0 {
            fmt.Printf("\n\nMessages were writen to file: %s\n", msgFile)
            open.Run(msgFile)
        }

        helper.PauseProcess(errors)
    }
}
