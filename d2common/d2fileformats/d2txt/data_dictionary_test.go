package d2txt

import (
	"fmt"
	"os"
	"testing"
)

func Test_Load(t *testing.T) {
	testFile, fileErr := os.Open("testdata/data.txt")
	if fileErr != nil {
		t.Error("cannot open test data file")
		return
	}

	data := make([]byte, 0)
	buf := make([]byte, 16)

	for {
		numRead, err := testFile.Read(buf)

		data = append(data, buf[:numRead]...)

		if err != nil {
			break
		}
	}

	txt := LoadDataDictionary(data)
	newData := txt.Marshal()
	newTXT := LoadDataDictionary(newData)
	_ = fmt.Println
	var fields []string
	for key := range txt.lookup {
		fields = append(fields, key)
		if _, ok := newTXT.lookup[key]; !ok {
			t.Fatalf("key not found %s", key)
		}
	}

	for i := 0; true; i++ {
		fmt.Println(i)
		newTxtNext := newTXT.Next()
		fmt.Println(len(newTXT.records))
		if newTxtNext == true {
			if txt.records[i][0] != newTXT.String(fields[0]) {
				t.Fatalf("wrong string coded:")
			}
		} else if !newTxtNext && len(txt.records) != len(newTXT.records) {
			t.Fatalf("wrong number of records in some of files new: %v", newTxtNext)
		} else {
			fmt.Printf("broken on idx %d", i)
			break
		}
	}

	t.Fail()

}
