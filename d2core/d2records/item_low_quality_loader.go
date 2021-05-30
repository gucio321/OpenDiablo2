package d2records

import (
	"github.com/gucio321/d2txt"
)

func lowQualityLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(LowQualities, 0)

	for d.Next() {
		record := &LowQualityRecord{
			Name: d.String("Hireling Description"),
		}

		records = append(records, record)
	}

	r.Item.LowQualityPrefixes = records

	r.Debugf("Loaded %d LowQuality records", len(records))

	return nil
}
