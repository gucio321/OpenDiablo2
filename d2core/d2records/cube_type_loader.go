package d2records

import (
	"github.com/gucio321/d2txt"
)

func cubeTypeLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(CubeTypes)

	for d.Next() {
		record := &CubeTypeRecord{
			Name:  d.String("cube item class"),
			Token: d.String("Code"),
		}

		records[record.Name] = record
	}

	r.Item.Cube.Types = records

	r.Debugf("Loaded %d CubeType records", len(records))

	return nil
}
