package d2records

import (
	"github.com/gucio321/d2txt"
)

func weaponClassesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(WeaponClasses)

	for d.Next() {
		record := &WeaponClassRecord{
			Name:  d.String("Weapon Class"),
			Token: d.String("Code"),
		}

		records[record.Name] = record
	}

	r.Animation.Token.Weapon = records

	r.Debugf("Loaded %d WeaponClass records", len(records))

	return nil
}
