package d2records

import (
	"github.com/gucio321/d2txt"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

// LoadMiscItems loads ItemCommonRecords from misc.txt
func miscItemsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := loadCommonItems(d, d2enum.InventoryItemTypeItem)

	r.Debugf("Loaded %d Misc Item records", len(records))

	r.Item.Misc = records

	return nil
}
