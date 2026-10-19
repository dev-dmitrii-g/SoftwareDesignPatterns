package coverage

import (
	"errors"
	"strings"
)

type Selector struct {
	Source Source
}

func (s Selector) Manual(name string) (Strategy, error) {
	switch strings.ToUpper(name) {
	case "GOBMP":
		return GOBMP{Source: s.Source}, nil
	case "OSMS":
		return OSMS{Source: s.Source}, nil
	case "DMS":
		return DMS{Source: s.Source}, nil
	default:
		return nil, errors.New("unknown coverage program")
	}
}

func (s Selector) Automatic(req Request) (Strategy, error) {
	for _, candidate := range []Strategy{
		DMS{Source: s.Source},
		OSMS{Source: s.Source},
		GOBMP{Source: s.Source},
	} {
		decision, err := candidate.Check(req)
		if err != nil {
			return nil, err
		}
		if decision.Approved {
			return candidate, nil
		}
	}
	return nil, errors.New("no coverage program approves this service")
}
