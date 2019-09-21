package zipmerge

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestZipMerge(t *testing.T) {
	before := zipArchive(map[string][]byte{
		`f1`: {1, 2, 3, 4},
		`f2`: {2, 3, 4},
	})
	assert(reflect.DeepEqual(zipUnArchive(before), map[string][]byte{
		`f1`: {1, 2, 3, 4},
		`f2`: {2, 3, 4},
	}))

	after, err := ZipMerge(before, map[string][]byte{
		`f1`: {2, 3, 5},
	})
	checkError(err)
	assert(reflect.DeepEqual(zipUnArchive(after), map[string][]byte{
		`f1`: {2, 3, 5},
		`f2`: {2, 3, 4},
	}))
}

func zipUnArchive(content []byte) (m map[string][]byte) {
	m = map[string][]byte{}
	r := bytes.NewReader(content)
	zr, err := zip.NewReader(r, int64(len(content)))
	checkError(err)
	for _, f := range zr.File {
		fr, err := f.Open()
		checkError(err)
		fc, err := ioutil.ReadAll(fr)
		checkError(err)
		m[f.Name] = fc
		err = fr.Close()
		checkError(err)
	}
	return
}

func zipArchive(m map[string][]byte) []byte {
	w := bytes.NewBuffer(nil)
	zw := zip.NewWriter(w)
	for name, content := range m {
		inW, err := zw.Create(name)
		checkError(err)
		_, err = inW.Write(content)
		checkError(err)
	}
	err := zw.Close()
	checkError(err)
	return w.Bytes()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func assert(b bool) {
	if !b {
		panic(`assert failed`)
	}
}
