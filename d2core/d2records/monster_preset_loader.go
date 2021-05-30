package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadMonPresets loads monster presets from monpresets.txt
func monsterPresetLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(MonPresets)

	for d.Next() {
		act := int32(d.Number("Act"))
		if _, ok := records[act]; !ok {
			records[act] = make([]string, 0)
		}

		records[act] = append(records[act], d.String("Place"))
	}

	r.Debugf("Loaded %d MonPreset records", len(records))

	r.Monster.Presets = records

	return nil
}
