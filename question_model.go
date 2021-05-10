package main

type (
	Questions struct {
		Questions []Question `json:"questions"`
	}

	Question struct {
		Prompt                 string `json:"prompt"`
		PromptNeedsTemplate    bool   `json:"promptNeedsTemplate"`
		PromptType             string `json:"promptType"` // scale(1-10), text, yes/no
		YesFollowPrompt        string `json:"yesPrompt"`
		YesPromptNeedsTemplate bool   `json:"yesPromptNeedsTemplate"`
		NoFollowPrompt         string `json:"noPrompt"`
		NoPromptNeedsTemplate  bool   `json:"noPromptNeedsTemplate"`
	}

	TemplateValues struct {
		FirstName  string
		DoctorName string
		Diagnosis  string
	}

	Answer struct {
		Prompt   string         `json:"prompt"`
		Ans      string         `json:"answer"`
		FollowUp FollowUpAnswer `json:"followUp"`
	}

	FollowUpAnswer struct {
		Prompt string `json:"prompt"`
		Ans    string `json:"answer"`
	}
)
