package d2records

import (
	"github.com/gucio321/d2txt"
)

func compositeTypeLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(CompositeTypes)

	for d.Next() {
		record := &CompositeTypeRecord{
			Name:  d.String("Name"),
			Token: d.String("Token"),
		}

		records[record.Name] = record
	}

	r.Animation.Token.Composite = records

	r.Debugf("Loaded %d CompositeType records", len(records))

	return nil
}
