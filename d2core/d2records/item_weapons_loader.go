package d2records

import (
	"github.com/gucio321/d2txt"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

// LoadWeapons loads weapon records
func weaponsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := loadCommonItems(d, d2enum.InventoryItemTypeWeapon)

	r.Debugf("Loaded %d Weapon records", len(records))

	r.Item.Weapons = records

	return nil
}
