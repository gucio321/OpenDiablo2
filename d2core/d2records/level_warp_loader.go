package d2records

import (
	"github.com/gucio321/d2txt"
)

func levelWarpsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(LevelWarps)

	for d.Next() {
		record := &LevelWarpRecord{
			Name:       d.String("Name"),
			ID:         d.Number("Id"),
			SelectX:    d.Number("SelectX"),
			SelectY:    d.Number("SelectY"),
			SelectDX:   d.Number("SelectDX"),
			SelectDY:   d.Number("SelectDY"),
			ExitWalkX:  d.Number("ExitWalkX"),
			ExitWalkY:  d.Number("ExitWalkY"),
			OffsetX:    d.Number("OffsetX"),
			OffsetY:    d.Number("OffsetY"),
			LitVersion: d.Bool("LitVersion"),
			Tiles:      d.Number("Tiles"),
			Direction:  d.String("Direction"),
		}
		records[record.ID] = record
	}

	r.Debugf("Loaded %d LevelWarp records", len(records))

	r.Level.Warp = records

	return nil
}
