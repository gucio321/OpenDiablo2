package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadElemTypes loads ElemTypeRecords into ElemTypes
func elemTypesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(ElemTypes)

	for d.Next() {
		record := &ElemTypeRecord{
			ElemType: d.String("Elemental Type"),
			Code:     d.String("Code"),
		}
		records[record.ElemType] = record
	}

	r.Debugf("Loaded %d ElemType records", len(records))

	r.ElemTypes = records

	return nil
}
