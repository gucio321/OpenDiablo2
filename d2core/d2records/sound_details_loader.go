package d2records

import (
	"github.com/gucio321/d2txt"
)

// Loadrecords loads SoundEntries from sounds.txt
func soundDetailsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(SoundDetails)

	for d.Next() {
		entry := &SoundDetailRecord{
			Handle:    d.String("Sound"),
			Index:     d.Number("Index"),
			FileName:  d.String("FileName"),
			Volume:    d.Number("Volume"),
			GroupSize: d.Number("Group Size"),
			Loop:      d.Bool("Loop"),
			FadeIn:    d.Number("Fade In"),
			FadeOut:   d.Number("Fade Out"),
			DeferInst: d.Bool("Defer Inst"),
			StopInst:  d.Bool("Stop Inst"),
			Duration:  d.Number("Duration"),
			Compound:  d.Number("Compound"),
			Reverb:    d.Number("Reverb"),
			Falloff:   d.Number("Falloff"),
			Cache:     d.Bool("Cache"),
			AsyncOnly: d.Bool("Async Only"),
			Priority:  d.Number("Priority"),
			Stream:    d.Bool("Stream"),
			Stereo:    d.Bool("Stereo"),
			Tracking:  d.Bool("Tracking"),
			Solo:      d.Bool("Solo"),
			MusicVol:  d.Bool("Music Vol"),
			Block1:    d.Number("Block 1"),
			Block2:    d.Number("Block 2"),
			Block3:    d.Number("Block 3"),
		}

		records[entry.Handle] = entry
	}

	r.Sound.Details = records

	r.Debugf("Loaded %d SoundDetail records", len(records))

	return nil
}
