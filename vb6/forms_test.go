package vb6

import "testing"

func TestReadErrors(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		err      error
	}{
		{"nonexistent", "nonexistent", ErrFileNotExist},
		{"empty", "testdata/empty.frm", ErrFileEmpty},
		{"no_version", "testdata/no_version.frm", ErrBadVersion},
		{"unsupported_version", "testdata/unsupported_version.frm", ErrBadVersion},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Load(c.filename)
			if err != c.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	f, err := Load("testdata/valid.frm")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if f == nil {
		t.Errorf("unexpected nil form")
	}
}
