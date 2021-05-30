package d2records

import (
	"github.com/gucio321/d2txt"
)

// LoadMonsterAI loads MonsterAIRecords from monai.txt
func monsterAiLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(MonsterAI)

	for d.Next() {
		record := &MonsterAIRecord{
			AI: d.String("AI"),
		}
		records[record.AI] = record
	}

	r.Debugf("Loaded %d MonsterAI records", len(records))

	r.Monster.AI = records

	return nil
}
