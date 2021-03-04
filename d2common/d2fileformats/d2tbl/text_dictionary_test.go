package d2tbl

import (
	"os"

	"testing"
)

func exampleData() TextDictionary {
	result := TextDictionary{
		"abc":             "def",
		"someStr":         "Some long string",
		"teststring":      "TeStxwsas123 long strin122*8:wq",
		"multilinestring": "this is come string\nwith more then one line\n123",
	}

	return result
}

func TestTBL_Marshal_testdata(t *testing.T) {
	testFile, fileErr := os.Open("testdata/testdata.tbl")
	if fileErr != nil {
		t.Error("cannot open test data file")
		return
	}

	fileData := make([]byte, 0)
	buf := make([]byte, 16)

	for {
		numRead, err := testFile.Read(buf)

		fileData = append(fileData, buf[:numRead]...)

		if err != nil {
			break
		}
	}

	tbl, err := LoadTextDictionary(fileData)
	if err != nil {
		t.Error(err)
	}

	data := tbl.Marshal()

	newTbl, err := LoadTextDictionary(data)
	if err != nil {
		t.Error(err)
	}

	for key, value := range tbl {
		newValue, ok := newTbl[key]

		if !ok {
			t.Fatalf("string %s wasn't encoded to table", key)
		}

		if newValue != value {
			t.Fatal("unexpected value set")
		}
	}
}

func TestTBL_Marshal(t *testing.T) {
	tbl := exampleData()

	data := tbl.Marshal()

	newTbl, err := LoadTextDictionary(data)
	if err != nil {
		t.Error(err)
	}

	for key, value := range tbl {
		newValue, ok := newTbl[key]

		if !ok {
			t.Fatalf("string %s wasn't encoded to table", key)
		}

		if newValue != value {
			t.Fatal("unexpected value set")
		}
	}
}

func TestTBL_MarshalNoNameString(t *testing.T) {
	tbl := &TextDictionary{
		"#0": "OKEY",
	}

	data := tbl.Marshal()

	newTbl, err := LoadTextDictionary(data)
	if err != nil {
		t.Error(err)
	}

	for key, value := range *tbl {
		newValue, ok := newTbl[key]

		if !ok {
			t.Fatalf("string %s wasn't encoded to table", key)
		}

		if newValue != value {
			t.Fatal("unexpected value set")
		}
	}
}
