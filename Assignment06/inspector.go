package main

type Inspector struct {
	store   *Store
	results Subject[Event]
}

func NewInspector(store *Store, results Subject[Event]) *Inspector {
	return &Inspector{store: store, results: results}
}

func (i *Inspector) OnEvent(e Event) {
	if e.quality >= 60 {
		return
	}
	patient := i.store.Tree(e.patient)
	result := e
	result.kind = Inspection
	result.quality = Apply[float64](patient, NewQualityVisitor())
	i.results.Publish(result)
}
