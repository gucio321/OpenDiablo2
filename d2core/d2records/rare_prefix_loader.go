package d2records

import (
	"github.com/gucio321/d2txt"
)

func rareItemPrefixLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	r.Item.Rare.Prefix = rareItemAffixLoader(d)

	r.Debugf("Loaded %d RarePrefix records", len(r.Item.Rare.Prefix))

	return nil
}
