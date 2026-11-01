package main

type LatencyVisitor struct {
	total float64
	count int
}

func NewLatencyVisitor() *LatencyVisitor { return &LatencyVisitor{} }

func (v *LatencyVisitor) VisitHospital(n *Hospital)       { visitChildren(n.children, v) }
func (v *LatencyVisitor) VisitWard(n *Ward)               { visitChildren(n.children, v) }
func (v *LatencyVisitor) VisitPatient(n *Patient)         { visitChildren(n.children, v) }
func (v *LatencyVisitor) VisitECG(n *ECG)                 { v.add(n.latencyMS) }
func (v *LatencyVisitor) VisitPPG(n *PPG)                 { v.add(n.latencyMS) }
func (v *LatencyVisitor) VisitEEG(n *EEG)                 { v.add(n.latencyMS) }
func (v *LatencyVisitor) VisitDevice(*Device)             {}
func (v *LatencyVisitor) VisitTemperature(n *Temperature) { v.add(n.latencyMS) }

func (v *LatencyVisitor) add(latencyMS float64) {
	v.total += latencyMS
	v.count++
}

func (v *LatencyVisitor) Result() float64 {
	if v.count == 0 {
		return 0
	}
	return v.total / float64(v.count)
}
