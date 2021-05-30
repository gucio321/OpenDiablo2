package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadPlayerClasses loads the PlayerClassRecords into PlayerClasses
func playerClassLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(PlayerClasses)

	for d.Next() {
		record := &PlayerClassRecord{
			Name: d.String("Player Class"),
			Code: d.String("Code"),
		}

		if record.Name == expansionString {
			continue
		}

		records[record.Name] = record
	}

	r.Debugf("Loaded %d PlayerClass records", len(records))

	r.Character.Classes = records

	return nil
}
