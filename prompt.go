package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

var (
	clearMap map[string]func()
)

func init() {
	clearMap = make(map[string]func())
	clearMap["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
func ParseInput(reader *bufio.Reader) string {
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	return s
}

func ClearScreen() {
	clearFunc, ok := clearMap[runtime.GOOS]
	if !ok {
		fmt.Println("\n *** Your platform is not supported to clear the terminal screen ***")
		return
	}
	clearFunc()
}

func TemplatePrompt(prompt string, templateValues TemplateValues) string {
	tmpl, err := template.New("question").Parse(prompt)
	if err != nil {
		return ""
	}
	var newPrompt bytes.Buffer
	err = tmpl.Execute(&newPrompt, templateValues)
	if err != nil {
		return ""
	}
	return newPrompt.String()
}

func CheckScale(response string) bool {
	_, err := strconv.ParseInt(response, 10, 64)
	if err != nil {
		return false
	}
	return true
}

func CheckYesNo(response string) string {
	if strings.ToLower(response) == "yes" || strings.ToLower(response) == "y" {
		return "yes"
	}
	if strings.ToLower(response) == "no" || strings.ToLower(response) == "n" {
		return "no"
	}
	return ""
}
