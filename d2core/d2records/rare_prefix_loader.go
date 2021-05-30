package d2records

import (
	"github.com/gucio321/d2txt"
)

func rareItemPrefixLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records, err := rareItemAffixLoader(d)
	if err != nil {
		return err
	}

	r.Item.Rare.Prefix = records

	r.Debugf("Loaded %d RarePrefix records", len(records))

	return nil
}
