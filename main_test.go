package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	gabs "github.com/Jeffail/gabs/v2"
)

func TestTransformBundle(t *testing.T) {
	parseJsonCorrect, _ := gabs.ParseJSON([]byte(`{"resourceType":"Bundle", "id":"my_bundle_id", "timestamp":"2021-05-09T00:00:00Z"}`))
	parseJsonResourseType, _ := gabs.ParseJSON([]byte(`{"id":"my_bundle_id", "timestamp":"2021-05-09T00:00:00Z"}`))
	timestamp := time.Date(2021, time.Month(5), 9, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name              string
		jsonParsed        *gabs.Container
		wantPatientBundle PatientBundle
		wantErr           bool
		wantedError       error
	}{
		{
			"Successful Transformation - Bundle",
			parseJsonCorrect,
			PatientBundle{ResourceType: "Bundle", ID: "my_bundle_id", Timestamp: timestamp},
			false,
			nil,
		},
		{
			"Missing ResourseType Transformation - Bundle",
			parseJsonResourseType,
			PatientBundle{},
			true,
			fmt.Errorf("TransformBundle: ResourceType missing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patientBundle, err := TransformBundle(tt.jsonParsed)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformBundle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.wantPatientBundle, patientBundle) {
				t.Errorf("TransformBundle() = %v, want %v", tt.wantPatientBundle, patientBundle)
			}
		})
	}
}

func TestTransformEntry(t *testing.T) {
	parseJsonCorrect, _ := gabs.ParseJSON([]byte(`{
		"resourceType":"Bundle",
		"id":"my_bundle_id",
		"timestamp":"2021-05-09T00:00:00Z",
		"entry": [{
			"resource": {
				"resourceType": "Patient",
				"id": "my_patient_id",
				"active": true,
				"name": [{
						"text": "Test McTesty",
						"family": "McTesty",
						"given": [
							"Test"
						]
					}
				],
				"contact": [{
						"system": "phone",
						"value": "555-555-5555",
						"use": "mobile"
					},
					{
						"system": "email",
						"value": "test@mctesty.com",
						"use": "work"
					}
				],
				"gender": "male",
				"birthDate": "1901-09-12",
				"address": [{
						"use": "home",
						"line": [
							"123 Some St",
							"Some City"
						]
					}
				]
			}
		}]
	}`))
	timestamp := time.Date(2021, time.Month(5), 9, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name              string
		jsonParsed        *gabs.Container
		wantPatientBundle PatientBundle
		wantErr           bool
		wantedError       error
	}{
		{
			"Successful Transformation - Entry(Patient)",
			parseJsonCorrect,
			PatientBundle{
				ResourceType: "Bundle",
				ID:           "my_bundle_id",
				Timestamp:    timestamp,
				Patient: Patient{
					ResourceType: "Patient",
					ID:           "my_patient_id",
					Active:       true,
					Names: []Name{
						{Text: "Test McTesty", Family: "McTesty", Given: []string{"Test"}},
					},
					Contacts: []Contact{
						{System: "phone", Value: "555-555-5555", Use: "mobile"},
						{System: "email", Value: "test@mctesty.com", Use: "work"},
					},
					Gender:    "male",
					BirthDate: "1901-09-12",
					Addresses: []Address{
						{Use: "home", Line: []string{"123 Some St", "Some City"}},
					},
				},
			},
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patientBundle, err := TransformEntry(tt.jsonParsed, PatientBundle{ResourceType: "Bundle", ID: "my_bundle_id", Timestamp: timestamp})
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.wantPatientBundle, patientBundle) {
				t.Errorf("TransformEntry() = %v, want %v", patientBundle, tt.wantPatientBundle)
			}
		})
	}
}
