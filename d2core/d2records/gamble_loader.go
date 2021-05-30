package d2records

import (
	"github.com/gucio321/d2txt"
)

func gambleLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(Gamble)

	for d.Next() {
		record := &GambleRecord{
			Name: d.String("name"),
			Code: d.String("code"),
		}
		records[record.Name] = record
	}

	r.Debugf("Loaded %d Gamble records", len(records))

	r.Gamble = records

	return nil
}
