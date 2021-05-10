package main

import "testing"

func TestTemplatePrompt(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		templateValues TemplateValues
		want           string
	}{
		{
			"success",
			"{{.DoctorName}} is my last name",
			TemplateValues{
				FirstName:  "Test",
				DoctorName: "McTesty",
			},
			"McTesty is my last name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TemplatePrompt(tt.prompt, tt.templateValues); got != tt.want {
				t.Errorf("TemplatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
