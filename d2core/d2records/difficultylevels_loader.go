package d2records

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/gucio321/d2txt"
)

func difficultyLevelsLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(DifficultyLevels)

	for d.Next() {
		record := &DifficultyLevelRecord{
			Name:                   d.String("Name"),
			ResistancePenalty:      d.Number("ResistPenalty"),
			DeathExperiencePenalty: d.Number("DeathExpPenalty"),
			DropChanceLow:          d.Number("UberCodeOddsNormal"),
			DropChanceNormal:       d.Number("UberCodeOddsNormal"),
			DropChanceSuperior:     d.Number("UberCodeOddsNormal"),
			DropChanceExceptional:  d.Number("UberCodeOddsNormal"),
			DropChanceMagic:        d.Number("UberCodeOddsGood"),
			DropChanceRare:         d.Number("UberCodeOddsGood"),
			DropChanceSet:          d.Number("UberCodeOddsGood"),
			DropChanceUnique:       d.Number("UberCodeOddsGood"),
			MonsterSkillBonus:      d.Number("MonsterSkillBonus"),
			MonsterColdDivisor:     d.Number("MonsterColdDivisor"),
			MonsterFreezeDivisor:   d.Number("MonsterFreezeDivisor"),
			AiCurseDivisor:         d.Number("AiCurseDivisor"),
			LifeStealDivisor:       d.Number("LifeStealDivisor"),
			ManaStealDivisor:       d.Number("ManaStealDivisor"),
		}
		switch record.Name {
		case "Normal":
			records[d2enum.DifficultyNormal] = record
		case "Nightmare":
			records[d2enum.DifficultyNightmare] = record
		case "Hell":
			records[d2enum.DifficultyHell] = record
		}
	}

	r.Debugf("Loaded %d DifficultyLevel records", len(records))

	r.DifficultyLevels = records

	return nil
}
