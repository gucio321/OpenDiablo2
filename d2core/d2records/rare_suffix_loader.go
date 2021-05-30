package d2records

import (
	"github.com/gucio321/d2txt"
)

func rareItemSuffixLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records, err := rareItemAffixLoader(d)
	if err != nil {
		return err
	}

	r.Debugf("Loaded %d RareSuffix records", len(records))

	r.Item.Rare.Suffix = records

	return nil
}
