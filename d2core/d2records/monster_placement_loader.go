package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadMonsterPlacements loads the MonsterPlacementRecords into MonsterPlacements.
func monsterPlacementsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(MonsterPlacements, 0)

	for d.Next() {
		records = append(records, MonsterPlacementRecord(d.String("code")))
	}

	r.Monster.Placements = records

	r.Debugf("Loaded %d MonsterPlacement records", len(records))

	return nil
}
