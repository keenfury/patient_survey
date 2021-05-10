package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	gabs "github.com/Jeffail/gabs/v2"
)

func main() {
	// load appointment data
	jsonParsed, err := LoadAppointmentData()
	if err != nil {
		os.Exit(1)
	}
	// set up data
	patientBundle, err := TransformIntoPatientBundle(jsonParsed)
	if err != nil {
		os.Exit(1)
	}
	// load questions
	questions, err := LoadQuestions()
	if err != nil {
		os.Exit(1)
	}
	templateValues := BuildTemplateValues(patientBundle)
	// show questions
	answers := StartSurvey(questions, templateValues)
	// show summary
	PrintSummary(answers)
	// save results
	SaveAnswers(answers)

	fmt.Println("Have a wonderful day")
}

func LoadAppointmentData() (jsonParsed *gabs.Container, err error) {
	fileContent, err := os.ReadFile("./patient-feedback-raw-data.json")
	if err != nil {
		fmt.Printf("ReadFile AppointmentData - error: %s\n", err)
		return
	}
	jsonParsed, err = gabs.ParseJSON(fileContent)
	if err != nil {
		fmt.Printf("ParseJSON AppointmentData - error: %s\n", err)
		return
	}
	return
}

func TransformIntoPatientBundle(jsonParsed *gabs.Container) (patientBundle PatientBundle, err error) {
	// transform bundle
	patientBundle, err = TransformBundle(jsonParsed)
	if err != nil {
		fmt.Printf("TransformIntoPatientBundle - TransformBundle: %s \n", err)
		return
	}
	patientBundle, err = TransformEntry(jsonParsed, patientBundle)
	if err != nil {
		fmt.Printf("TransformIntoPatientBundle - TransformEntry: %s \n", err)
		return
	}
	return
}

func TransformBundle(jsonParsed *gabs.Container) (patientBundle PatientBundle, err error) {
	ok := false
	patientBundle.ResourceType, ok = jsonParsed.Path("resourceType").Data().(string)
	if !ok {
		err = fmt.Errorf("TransformBundle: ResourceType missing")
		return
	}
	patientBundle.ID, ok = jsonParsed.Path("id").Data().(string)
	if !ok {
		err = fmt.Errorf("TransformBundle: ID missing")
		return
	}
	timestamp := jsonParsed.Path("timestamp").Bytes()
	err = json.Unmarshal(timestamp, &patientBundle.Timestamp)
	if err != nil {
		err = fmt.Errorf("TransformBundle Timestamp missing or not correct format")
		return
	}
	return
}

func TransformEntry(jsonParsed *gabs.Container, patientBundle PatientBundle) (PatientBundle, error) {
	var err error
	entryChildren := jsonParsed.Path("entry").Children()
	for _, entry := range entryChildren {
		ok := false
		resourceType, ok := entry.Path("resource.resourceType").Data().(string)
		if !ok {
			err := fmt.Errorf("TransformPatient: ResourceType missing")
			return patientBundle, err
		}
		entryBytes := entry.Path("resource").Bytes()
		if resourceType == "Patient" {
			err := json.Unmarshal(entryBytes, &patientBundle.Patient)
			if err != nil {
				err = fmt.Errorf("TransformPatient: invalid patient format: %s", err)
				return patientBundle, err
			}
		}
		if resourceType == "Doctor" {
			err := json.Unmarshal(entryBytes, &patientBundle.Doctor)
			if err != nil {
				err = fmt.Errorf("TransformPatient: invalid doctor format: %s", err)
				return patientBundle, err
			}
		}
		if resourceType == "Appointment" {
			err := json.Unmarshal(entryBytes, &patientBundle.Appointment)
			if err != nil {
				err = fmt.Errorf("TransformPatient: invalid doctor format: %s", err)
				return patientBundle, err
			}
		}
		if resourceType == "Diagnosis" {
			err := json.Unmarshal(entryBytes, &patientBundle.Diagnosis)
			if err != nil {
				err = fmt.Errorf("TransformPatient: invalid doctor format: %s", err)
				return patientBundle, err
			}
		}
	}
	return patientBundle, err
}

func LoadQuestions() (questions Questions, err error) {
	fileContent, err := os.ReadFile("./questions.json")
	if err != nil {
		fmt.Printf("ReadFile Questions - error: %s\n", err)
		return
	}
	err = json.Unmarshal(fileContent, &questions)
	if err != nil {
		fmt.Printf("ParseJSON Questions - error: %s\n", err)
		return
	}
	return
}

func BuildTemplateValues(patientBundle PatientBundle) (templateValues TemplateValues) {
	// TODO: check for array len
	templateValues.DoctorName = patientBundle.Doctor.Names[0].Family
	templateValues.Diagnosis = patientBundle.Diagnosis.DiagnosisCode.Coding[0].Name
	templateValues.FirstName = patientBundle.Patient.Names[0].Given[0]
	return
}

func StartSurvey(questions Questions, templateValues TemplateValues) []Answer {
	ClearScreen()
	fmt.Printf("Welcome to our patient survey\n\n")
	answers := []Answer{}
	reader := bufio.NewReader(os.Stdin)
	for _, question := range questions.Questions {
		prompt := question.Prompt
		if question.PromptNeedsTemplate {
			prompt = TemplatePrompt(prompt, templateValues)
		}
		response := ""
		for {
			breakOut := false
			followUp := false
			fmt.Printf("%s:\n", prompt)
			fmt.Print("Answer: ")
			response = ParseInput(reader)
			switch question.PromptType {
			case "scale":
				ok := CheckScale(response)
				if !ok {
					fmt.Print("Invalid number, press 'enter' to continue")
					ParseInput(reader)
					ClearScreen()
					continue
				}
			case "yes/no":
				yesNo := CheckYesNo(response)
				switch yesNo {
				case "yes":
					response = "yes"
					followUp = true
				case "no":
					response = "no"
					followUp = true
				default:
					fmt.Print("Invalid input (yes/no), press 'enter' to continue")
					ParseInput(reader)
					ClearScreen()
					continue
				}
			}
			breakOut = true
			answer := Answer{Prompt: prompt, Ans: response}
			if followUp {
				if response == "yes" {
					yesPrompt := question.YesFollowPrompt
					if question.YesPromptNeedsTemplate {
						yesPrompt = TemplatePrompt(yesPrompt, templateValues)
					}
					fmt.Printf("%s:\n", yesPrompt)
					fmt.Print("Answer: ")
					response = ParseInput(reader)
					answer.FollowUp.Prompt = yesPrompt
					answer.FollowUp.Ans = response
				}
				if response == "no" {
					noPrompt := question.NoFollowPrompt
					if question.YesPromptNeedsTemplate {
						noPrompt = TemplatePrompt(noPrompt, templateValues)
					}
					fmt.Printf("%s:\n", noPrompt)
					fmt.Print("Answer: ")
					response = ParseInput(reader)
					answer.FollowUp.Prompt = noPrompt
					answer.FollowUp.Ans = response
				}
			}
			answers = append(answers, answer)
			ClearScreen()
			if breakOut {
				break
			}
		}
	}
	return answers
}

func PrintSummary(answers []Answer) {
	fmt.Printf("Thanks again! Here’s what we heard:\n\n")
	for _, answer := range answers {
		fmt.Printf("%s\n", answer.Prompt)
		fmt.Printf("-- %s\n", answer.Ans)
		if answer.FollowUp.Prompt != "" {
			fmt.Println("")
			fmt.Printf("\t%s\n", answer.FollowUp.Prompt)
			fmt.Printf("\t-- %s\n", answer.FollowUp.Ans)
		}
		fmt.Println("")
	}
}

func SaveAnswers(answers []Answer) {
	answerBytes, err := json.MarshalIndent(answers, "", "    ")
	if err != nil {
		fmt.Println("unable to save file")
		return
	}
	err = os.WriteFile("./answers.json", answerBytes, 0600)
	if err != nil {
		fmt.Println("unable to save file")
		return
	}
}
