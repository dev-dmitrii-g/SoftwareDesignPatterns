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
	reason := "public eligibility or service missing"
	if approved {
		reason = "public eligibility confirmed"
	}
	return Decision{Program: "GOBMP", Approved: approved, Reason: reason}, nil
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
	reason := "active coverage, service, or referral missing"
	if approved {
		reason = "active coverage confirmed"
	}
	return Decision{Program: "OSMS", Approved: approved, Reason: reason}, nil
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
	reason := "policy expired or service missing"
	if approved {
		reason = "private policy confirmed"
	}
	return Decision{Program: "DMS", Approved: approved, Reason: reason}, nil
}
