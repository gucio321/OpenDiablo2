package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadMonModes loads monster records
func monsterModeLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(MonModes)

	for d.Next() {
		record := &MonModeRecord{
			Name:  d.String("name"),
			Token: d.String("token"),
			Code:  d.String("code"),
		}
		records[record.Name] = record
	}

	r.Debugf("Loaded %d MonMode records", len(records))

	r.Monster.Modes = records

	return nil
}
