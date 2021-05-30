package d2records

import (
	"github.com/gucio321/d2txt"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

func armorLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	if r.Item.Armors != nil {
		return nil // already loaded
	}

	records := loadCommonItems(d, d2enum.InventoryItemTypeArmor)

	r.Debugf("Loaded %d Armor Item records", len(records))

	r.Item.Armors = records

	return nil
}
