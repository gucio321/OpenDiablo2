package d2records

import (
	"github.com/gucio321/d2txt"
)

func storePagesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(StorePages)

	for d.Next() {
		record := &StorePageRecord{
			StorePage: d.String("Store Page"),
			Code:      d.String("Code"),
		}
		records[record.StorePage] = record
	}

	r.Item.StorePages = records

	r.Debugf("Loaded %d StorePage records", len(records))

	return nil
}
