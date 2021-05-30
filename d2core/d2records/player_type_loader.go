package d2records

import (
	"github.com/gucio321/d2txt"
)

func playerTypeLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(PlayerTypes)

	for d.Next() {
		record := &PlayerTypeRecord{
			Name:  d.String("Player Class"),
			Token: d.String("Token"),
		}

		if record.Name == expansionString {
			continue
		}

		records[record.Name] = record
	}

	r.Debugf("Loaded %d PlayerType records", len(records))

	r.Animation.Token.Player = records

	return nil
}
