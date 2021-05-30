package d2records

import (
	"github.com/gucio321/d2txt"
)

func armorTypesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(ArmorTypes)

	for d.Next() {
		record := &ArmorTypeRecord{
			Name:  d.String("Name"),
			Token: d.String("Token"),
		}

		records[record.Name] = record
	}

	r.Animation.Token.Armor = records

	r.Debugf("Loaded %d ArmorType records", len(records))

	return nil
}
