package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadEvents loads all of the event records from events.txt
func eventsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(Events, 0)

	for d.Next() {
		record := &EventRecord{
			Event: d.String("event"),
		}

		records = append(records, record)
	}

	r.Debugf("Loaded %d Event records", len(records))

	r.Character.Events = records

	return nil
}
