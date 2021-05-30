package d2records

import (
	"github.com/gucio321/d2txt"
)

func uniqueAppellationsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(UniqueAppellations)

	for d.Next() {
		record := &UniqueAppellationRecord{
			Name: d.String("Name"),
		}

		records[record.Name] = record
	}

	r.Monster.Unique.Appellations = records

	r.Debugf("Loaded %d UniqueAppellation records", len(records))

	return nil
}
