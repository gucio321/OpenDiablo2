package d2records

import (
	"github.com/gucio321/d2txt"
)

func bodyLocationsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(BodyLocations)

	for d.Next() {
		location := &BodyLocationRecord{
			Name: d.String("Name"),
			Code: d.String("Code"),
		}
		records[location.Code] = location
	}

	r.Debugf("Loaded %d BodyLocation records", len(records))

	r.BodyLocations = records

	return nil
}
