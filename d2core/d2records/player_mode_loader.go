package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadPlayerModes loads PlayerModeRecords into PlayerModes
func playerModesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(PlayerModes)

	for d.Next() {
		record := &PlayerModeRecord{
			Name:  d.String("Name"),
			Token: d.String("Token"),
		}

		records[record.Name] = record
	}

	r.Character.Modes = records

	r.Debugf("Loaded %d PlayerMode records", len(records))

	return nil
}
