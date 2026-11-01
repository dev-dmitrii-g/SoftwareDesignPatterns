package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func visitChildren(children []Element, visitor Dispatcher) {
	for _, child := range children {
		child.Accept(visitor)
	}
}

type QualityVisitor struct {
	total float64
	count int
}

func NewQualityVisitor() *QualityVisitor { return &QualityVisitor{} }

func (v *QualityVisitor) VisitHospital(n *Hospital)       { visitChildren(n.children, v) }
func (v *QualityVisitor) VisitWard(n *Ward)               { visitChildren(n.children, v) }
func (v *QualityVisitor) VisitPatient(n *Patient)         { visitChildren(n.children, v) }
func (v *QualityVisitor) VisitECG(n *ECG)                 { v.add(n.quality) }
func (v *QualityVisitor) VisitPPG(n *PPG)                 { v.add(n.quality) }
func (v *QualityVisitor) VisitEEG(n *EEG)                 { v.add(n.quality) }
func (v *QualityVisitor) VisitDevice(*Device)             {}
func (v *QualityVisitor) VisitTemperature(n *Temperature) { v.add(n.quality) }

func (v *QualityVisitor) add(quality float64) {
	v.total += quality
	v.count++
}

func (v *QualityVisitor) Result() float64 {
	if v.count == 0 {
		return 0
	}
	return v.total / float64(v.count)
}

type AlarmVisitor struct{ counts map[string]int }

func NewAlarmVisitor() *AlarmVisitor { return &AlarmVisitor{make(map[string]int)} }

func (v *AlarmVisitor) VisitHospital(n *Hospital)       { visitChildren(n.children, v) }
func (v *AlarmVisitor) VisitWard(n *Ward)               { visitChildren(n.children, v) }
func (v *AlarmVisitor) VisitPatient(n *Patient)         { visitChildren(n.children, v) }
func (v *AlarmVisitor) VisitECG(n *ECG)                 { v.counts["ecg"] += n.alarms }
func (v *AlarmVisitor) VisitPPG(n *PPG)                 { v.counts["ppg"] += n.alarms }
func (v *AlarmVisitor) VisitEEG(n *EEG)                 { v.counts["eeg"] += n.alarms }
func (v *AlarmVisitor) VisitDevice(n *Device)           { v.counts["device"] += n.alarms }
func (v *AlarmVisitor) VisitTemperature(n *Temperature) { v.counts["temperature"] += n.alarms }

func (v *AlarmVisitor) Result() map[string]int {
	result := make(map[string]int, len(v.counts))
	for kind, count := range v.counts {
		result[kind] = count
	}
	return result
}

type FilterVisitor struct {
	threshold float64
	matches   []Element
}

func NewFilterVisitor(threshold float64) *FilterVisitor {
	return &FilterVisitor{threshold: threshold, matches: make([]Element, 0)}
}

func (v *FilterVisitor) VisitHospital(n *Hospital)       { visitChildren(n.children, v) }
func (v *FilterVisitor) VisitWard(n *Ward)               { visitChildren(n.children, v) }
func (v *FilterVisitor) VisitPatient(n *Patient)         { visitChildren(n.children, v) }
func (v *FilterVisitor) VisitECG(n *ECG)                 { v.add(n, n.quality) }
func (v *FilterVisitor) VisitPPG(n *PPG)                 { v.add(n, n.quality) }
func (v *FilterVisitor) VisitEEG(n *EEG)                 { v.add(n, n.quality) }
func (v *FilterVisitor) VisitDevice(*Device)             {}
func (v *FilterVisitor) VisitTemperature(n *Temperature) { v.add(n, n.quality) }

func (v *FilterVisitor) add(element Element, quality float64) {
	if quality < v.threshold {
		v.matches = append(v.matches, element)
	}
}

func (v *FilterVisitor) Result() []Element { return append([]Element{}, v.matches...) }

type JSONVisitor struct{ output strings.Builder }

func NewJSONVisitor() *JSONVisitor { return &JSONVisitor{} }

func (v *JSONVisitor) VisitHospital(n *Hospital)       { v.container("hospital", n.container) }
func (v *JSONVisitor) VisitWard(n *Ward)               { v.container("ward", n.container) }
func (v *JSONVisitor) VisitPatient(n *Patient)         { v.container("patient", n.container) }
func (v *JSONVisitor) VisitECG(n *ECG)                 { v.leaf("ecg", n.reading) }
func (v *JSONVisitor) VisitPPG(n *PPG)                 { v.leaf("ppg", n.reading) }
func (v *JSONVisitor) VisitEEG(n *EEG)                 { v.leaf("eeg", n.reading) }
func (v *JSONVisitor) VisitDevice(n *Device)           { v.leaf("device", n.reading) }
func (v *JSONVisitor) VisitTemperature(n *Temperature) { v.leaf("temperature", n.reading) }
func (v *JSONVisitor) Result() string                  { return v.output.String() }

func (v *JSONVisitor) header(kind, name string) {
	v.output.WriteString(`{"type":`)
	v.quoted(kind)
	v.output.WriteString(`,"name":`)
	v.quoted(name)
}

func (v *JSONVisitor) container(kind string, n container) {
	v.header(kind, n.name)
	v.output.WriteString(`,"children":[`)
	for index, child := range n.children {
		if index > 0 {
			v.output.WriteByte(',')
		}
		child.Accept(v)
	}
	v.output.WriteString(`]}`)
}

func (v *JSONVisitor) leaf(kind string, n reading) {
	v.header(kind, n.name)
	fmt.Fprintf(&v.output, `,"quality":%g,"alarms":%d,"latency_ms":%g}`, n.quality, n.alarms, n.latencyMS)
}

func (v *JSONVisitor) quoted(value string) {
	encoded, _ := json.Marshal(value)
	v.output.Write(encoded)
}
