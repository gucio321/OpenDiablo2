package d2records

import (
	"github.com/gucio321/d2txt"
)

func hitClassLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(HitClasses)

	for d.Next() {
		record := &HitClassRecord{
			Name:  d.String("Hit Class"),
			Token: d.String("Code"),
		}

		records[record.Name] = record
	}

	r.Animation.Token.HitClass = records

	r.Debugf("Loaded %d HitClass records", len(records))

	return nil
}
