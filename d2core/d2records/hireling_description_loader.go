package d2records

import (
	"github.com/gucio321/d2txt"
)

func hirelingDescriptionLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(HirelingDescriptions)

	for d.Next() {
		record := &HirelingDescriptionRecord{
			Name:  d.String("Hireling Description"),
			Token: d.String("Code"),
		}

		records[record.Name] = record
	}

	r.Hireling.Descriptions = records

	r.Debugf("Loaded %d HirelingDescription records", len(records))

	return nil
}
