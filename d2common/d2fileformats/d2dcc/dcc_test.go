package d2dcc

import (
	"os"
	"testing"
)

func TestDCC_Load(t *testing.T) {
	testFile, fileErr := os.Open("testdata/data.dcc")
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

	dcc, loadErr := Load(data)
	if loadErr != nil {
		t.Error(loadErr)
	}

	newData := dcc.Marshal()
	for i := range newData {
		if i == 11 || i == 12 {
			continue
		}
		if data[i] != newData[i] {
			t.Errorf("data aren't equal: %d != %d (idx %d)", data[i], newData[i], i)
		}
	}
	/*_, err := Load(newData)
	if err != nil {
		t.Error(err)
	}*/
}
