package d2records

import "github.com/gucio321/d2txt"

// recordLoader represents something that can load a data dictionary and
// handles placing it in the record manager exports
type recordLoader func(r *RecordManager, dictionary *d2txt.DataDictionary) error
