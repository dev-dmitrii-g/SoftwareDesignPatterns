package coverage

type Strategy interface {
	Check(Request) (Decision, error)
}

type GOBMP struct {
	Source Source
}

func (s GOBMP) Check(req Request) (Decision, error) {
	record, err := s.Source.Lookup(req.IIN)
	if err != nil {
		return Decision{}, err
	}
	approved := record.GOBMPEligible && (req.Service == "emergency" || req.Service == "basic")
	return Decision{Program: "GOBMP", Approved: approved}, nil
}

type OSMS struct {
	Source Source
}

func (s OSMS) Check(req Request) (Decision, error) {
	record, err := s.Source.Lookup(req.IIN)
	if err != nil {
		return Decision{}, err
	}
	approved := record.OSMSActive && (req.Service == "consultation" || (req.Service == "rehabilitation" && record.HasReferral))
	return Decision{Program: "OSMS", Approved: approved}, nil
}

type DMS struct {
	Source Source
}

func (s DMS) Check(req Request) (Decision, error) {
	record, err := s.Source.Lookup(req.IIN)
	if err != nil {
		return Decision{}, err
	}
	approved := req.At.Before(record.DMSExpiresAt) && record.DMSServices[req.Service]
	return Decision{Program: "DMS", Approved: approved}, nil
}
