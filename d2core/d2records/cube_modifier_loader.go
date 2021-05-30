package d2records

import (
	"github.com/gucio321/d2txt"
)

func cubeModifierLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(CubeModifiers)

	for d.Next() {
		record := &CubeModifierRecord{
			Name:  d.String("cube modifier type"),
			Token: d.String("Code"),
		}

		records[record.Name] = record
	}

	r.Item.Cube.Modifiers = records

	r.Debugf("Loaded %d CubeModifier records", len(records))

	return nil
}
