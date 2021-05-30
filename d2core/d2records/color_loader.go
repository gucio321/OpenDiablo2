package d2records

import (
	"github.com/gucio321/d2txt"
)

func colorsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(Colors)

	for d.Next() {
		record := &ColorRecord{
			TransformColor: d.String("Transform Color"),
			Code:           d.String("Code"),
		}

		records[record.TransformColor] = record
	}

	r.Colors = records

	r.Debugf("Loaded %d Color records", len(records))

	return nil
}
