package d2records

import (
	"strings"

	"github.com/gucio321/d2txt"
)

// LoadObjectTypes loads ObjectTypeRecords from objtype.txt
func objectTypesLoader(r *RecordManager, d *d2txt.DataDictionary) error {
	records := make(ObjectTypes, 0)

	for d.Next() {
		record := ObjectTypeRecord{
			Name:  sanitizeObjectString(d.String("Name")),
			Token: sanitizeObjectString(d.String("Token")),
		}

		records = append(records, record)
	}

	r.Debugf("Loaded %d ObjectType records", len(records))

	r.Object.Types = records

	return nil
}

func sanitizeObjectString(str string) string {
	result := strings.TrimSpace(strings.ReplaceAll(str, string(byte(0)), ""))
	result = strings.ToLower(result)

	return result
}
