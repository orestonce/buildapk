package zipmerge

import (
	"bytes"
	"github.com/orestonce/buildapk/zipmerge/internal"
	"io"
)

func ZipMerge(zipContent []byte, m map[string][]byte) (after []byte, err error) {
	br := bytes.NewReader(zipContent)
	r, err := internal.NewReader(br, int64(len(zipContent)))
	if err != nil {
		return
	}
	w := bytes.NewBuffer(nil)
	_, err = w.Write(zipContent[:int(r.AppendOffset())])
	if err != nil {
		return
	}
	wf := r.Append(w)
	for name, content := range m {
		var inW io.Writer
		inW, err = wf.Create(name)
		if err != nil {
			return
		}
		_, err = inW.Write(content)
		if err != nil {
			return
		}
	}
	err = wf.Close()
	if err != nil {
		return
	}
	return w.Bytes(), nil
}
