package coverage

import (
	"fmt"
	"time"
)

type Request struct {
	IIN     string
	Service string
	At      time.Time
}

type Decision struct {
	Program  string
	Approved bool
	Reason   string
}

type Record struct {
	GOBMPEligible bool
	OSMSActive    bool
	DMSExpiresAt  time.Time
	DMSServices   map[string]bool
}

type Source interface {
	Lookup(iin string) (Record, error)
}

type MemorySource struct {
	Records map[string]Record
}

func (s MemorySource) Lookup(iin string) (Record, error) {
	record, ok := s.Records[iin]
	if !ok {
		return Record{}, fmt.Errorf("no coverage record for %s", iin)
	}
	return record, nil
}
