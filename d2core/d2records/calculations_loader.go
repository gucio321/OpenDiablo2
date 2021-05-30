package d2records

import (
	"github.com/gucio321/d2txt"
)

func skillCalcLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	r.Calculation.Skills = loadCalculations(r, d, "Skill")

	return nil
}

func missileCalcLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	r.Calculation.Missiles = loadCalculations(r, d, "Missile")

	return nil
}

func loadCalculations(r *RecordManager, d *d2txt.DataDictionary, name string) Calculations {
	records := make(Calculations)

	for d.Next() {
		record := &CalculationRecord{
			Code:        d.String("code"),
			Description: d.String("*desc"),
		}
		records[record.Code] = record
	}

	r.Debugf("Loaded %d %s Calculation records", len(records), name)

	return records
}
