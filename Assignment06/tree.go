package main

type Element interface {
	Name() string
	Accept(Dispatcher)
}

type Dispatcher interface {
	VisitHospital(*Hospital)
	VisitWard(*Ward)
	VisitPatient(*Patient)
	VisitECG(*ECG)
	VisitPPG(*PPG)
	VisitEEG(*EEG)
	VisitDevice(*Device)
	VisitTemperature(*Temperature)
}

type Visitor[R any] interface {
	Dispatcher
	Result() R
}

func Apply[R any](element Element, visitor Visitor[R]) R {
	element.Accept(visitor)
	return visitor.Result()
}

type node struct{ name string }

func (n node) Name() string { return n.name }

type container struct {
	node
	children []Element
}

func newContainer(name string, children []Element) container {
	return container{node{name}, append([]Element(nil), children...)}
}

func (c container) Children() []Element { return append([]Element(nil), c.children...) }

type Hospital struct{ container }
type Ward struct{ container }
type Patient struct{ container }

func NewHospital(name string, children ...Element) *Hospital {
	return &Hospital{newContainer(name, children)}
}
func NewWard(name string, children ...Element) *Ward {
	return &Ward{newContainer(name, children)}
}
func NewPatient(name string, children ...Element) *Patient {
	return &Patient{newContainer(name, children)}
}

func (n *Hospital) Accept(v Dispatcher) { v.VisitHospital(n) }
func (n *Ward) Accept(v Dispatcher)     { v.VisitWard(n) }
func (n *Patient) Accept(v Dispatcher)  { v.VisitPatient(n) }

type reading struct {
	node
	quality   float64
	alarms    int
	latencyMS float64
}

type ECG struct{ reading }
type PPG struct{ reading }
type EEG struct{ reading }
type Device struct{ reading }
type Temperature struct{ reading }

func NewECG(name string, quality float64, alarms int, latencyMS float64) *ECG {
	return &ECG{reading{node{name}, quality, alarms, latencyMS}}
}
func NewPPG(name string, quality float64, alarms int, latencyMS float64) *PPG {
	return &PPG{reading{node{name}, quality, alarms, latencyMS}}
}
func NewEEG(name string, quality float64, alarms int, latencyMS float64) *EEG {
	return &EEG{reading{node{name}, quality, alarms, latencyMS}}
}
func NewDevice(name string, quality float64, alarms int, latencyMS float64) *Device {
	return &Device{reading{node{name}, quality, alarms, latencyMS}}
}
func NewTemperature(name string, quality float64, alarms int, latencyMS float64) *Temperature {
	return &Temperature{reading{node{name}, quality, alarms, latencyMS}}
}

func (n *ECG) Accept(v Dispatcher)         { v.VisitECG(n) }
func (n *PPG) Accept(v Dispatcher)         { v.VisitPPG(n) }
func (n *EEG) Accept(v Dispatcher)         { v.VisitEEG(n) }
func (n *Device) Accept(v Dispatcher)      { v.VisitDevice(n) }
func (n *Temperature) Accept(v Dispatcher) { v.VisitTemperature(n) }
