package d2records

import (
	"github.com/gucio321/d2txt"
)

func monsterTypesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(MonsterTypes)

	for d.Next() {
		record := &MonTypeRecord{
			Type:      d.String("type"),
			Equiv1:    d.String("equiv1"),
			Equiv2:    d.String("equiv2"),
			Equiv3:    d.String("equiv3"),
			StrSing:   d.String("strsing"),
			StrPlural: d.String("strplur"),
		}
		records[record.Type] = record
	}

	r.Debugf("Loaded %d MonType records", len(records))

	r.Monster.Types = records

	return nil
}
