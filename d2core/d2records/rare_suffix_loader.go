package d2records

import (
	"github.com/gucio321/d2txt"
)

func rareItemSuffixLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := rareItemAffixLoader(d)

	r.Debugf("Loaded %d RareSuffix records", len(records))

	r.Item.Rare.Suffix = records

	return nil
}
