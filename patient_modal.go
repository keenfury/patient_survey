package main

import "time"

type (
	PatientBundle struct {
		ResourceType string    `json:"resourceType"`
		ID           string    `json:"id"`
		Timestamp    time.Time `json:"timestamp"`
		Patient
		Doctor
		Appointment
		Diagnosis
	}

	Patient struct {
		ResourceType string    `json:"resourceType"`
		ID           string    `json:"id"`
		Active       bool      `json:"active"`
		Names        []Name    `json:"name"`
		Contacts     []Contact `json:"contact"`
		Gender       string    `json:"gender"`
		BirthDate    string    `json:"birthDate"`
		Addresses    []Address `json:"address"`
	}

	Doctor struct {
		ResourceType string `json:"resourceType"`
		ID           string `json:"id"`
		Names        []Name `json:"name"`
	}

	Appointment struct {
		ResourceType     string            `json:"resourceType"`
		ID               string            `json:"id"`
		Status           string            `json:"status"`
		AppointmentTypes []AppointmentType `json:"type"`
		Subject          Reference         `json:"subject"`
		Actor            Reference         `json:"actor"`
		Period           `json:"period"`
	}

	Diagnosis struct {
		ResourceType  string `json:"resourceType"`
		ID            string `json:"id"`
		Meta          `json:"meta"`
		Status        string    `json:"status"`
		DiagnosisCode Code      `json:"code"`
		Appointment   Reference `json:"appointment"`
	}

	Name struct {
		Text   string   `json:"text"`
		Family string   `json:"family"`
		Given  []string `json:"given"`
	}

	Contact struct {
		System string `json:"system"`
		Value  string `json:"value"`
		Use    string `json:"use"`
	}
	Address struct {
		Use  string   `json:"use"`
		Line []string `json:"line"`
	}

	AppointmentType struct {
		Text string `json:"text"`
	}

	Reference struct {
		Reference string `json:"reference"`
	}

	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	Meta struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}

	Code struct {
		Coding []Coding `json:"coding"`
	}

	Coding struct {
		System string `json:"system"`
		Code   string `json:"code"`
		Name   string `json:"name"`
	}
)
