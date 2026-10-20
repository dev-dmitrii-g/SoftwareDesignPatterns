package main

import (
	"fmt"
	"time"

	"hospitaladmission/internal/admission"
	"hospitaladmission/internal/coverage"
	"hospitaladmission/internal/flow"
)

func showFlow(label string, context *flow.Context, a admission.Admission) {
	result, err := context.Begin(a)
	if err != nil {
		fmt.Printf("%s: %v\n", label, err)
		return
	}
	fmt.Printf("%s: %s, steps=%v", label, result.Status, result.Steps)
	if !result.ScheduledAt.IsZero() {
		fmt.Printf(", at=%s", result.ScheduledAt.Format(time.RFC3339))
	}
	fmt.Println()
}

func showCoverage(label string, context *coverage.Context, req coverage.Request) {
	decision, err := context.Check(req)
	if err != nil {
		fmt.Printf("%s: %v\n", label, err)
		return
	}
	fmt.Printf("%s: %s, approved=%t\n", label, decision.Program, decision.Approved)
}

func main() {
	a := admission.Admission{
		IIN: "123456789123", Service: "consultation",
		RequestedMedicine: "Activated Charcoal", RequestedEquipment: "IV Stand",
		ScheduledAt: time.Now().Add(24 * time.Hour),
	}
	flowSelector := flow.Selector{}
	fmt.Println("Admission flow: same input with every strategy")
	for _, name := range []string{"automatic", "manual", "scheduled"} {
		selected, err := flowSelector.Manual(name)
		if err != nil {
			panic(err)
		}
		showFlow(name, flow.NewContext(selected), a)
	}
	flowContext := flow.NewContext(flow.Manual{})
	showFlow("manual choice", flowContext, a)
	flowContext.SetStrategy(flowSelector.Automatic(a))
	showFlow("automatic choice after SetStrategy", flowContext, a)
	urgent := a
	urgent.Urgent = true
	showFlow("automatic urgent rule", flow.NewContext(flowSelector.Automatic(urgent)), urgent)

	source := coverage.MemorySource{Records: map[string]coverage.Record{
		a.IIN: {
			GOBMPEligible: true,
			OSMSActive:    true,
			DMSExpiresAt:  time.Now().Add(365 * 24 * time.Hour),
			DMSServices:   map[string]bool{"consultation": true},
		},
	}}
	req := coverage.Request{IIN: a.IIN, Service: a.Service, At: time.Now()}
	coverageSelector := coverage.Selector{Source: source}
	fmt.Println("\nCoverage: same input with every strategy")
	for _, name := range []string{"GOBMP", "OSMS", "DMS"} {
		selected, err := coverageSelector.Manual(name)
		if err != nil {
			panic(err)
		}
		showCoverage(name, coverage.NewContext(selected), req)
	}
	selected, err := coverageSelector.Manual("OSMS")
	if err != nil {
		panic(err)
	}
	coverageContext := coverage.NewContext(selected)
	showCoverage("manual choice", coverageContext, req)
	selected, err = coverageSelector.Automatic(req)
	if err != nil {
		panic(err)
	}
	coverageContext.SetStrategy(selected)
	showCoverage("automatic choice after SetStrategy", coverageContext, req)
}
