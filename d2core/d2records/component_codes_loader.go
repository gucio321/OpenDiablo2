package d2records

import (
	"github.com/gucio321/d2txt"
)

func componentCodesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(ComponentCodes)

	for d.Next() {
		record := &ComponentCodeRecord{
			Component: d.String("component"),
			Code:      d.String("code"),
		}
		records[record.Component] = record
	}

	r.Debugf("Loaded %d ComponentCode records", len(records))

	r.ComponentCodes = records

	return nil
}
