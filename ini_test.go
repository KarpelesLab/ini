package ini_test

import (
	"bytes"
	"testing"

	"github.com/KarpelesLab/ini"
)

func TestIni(t *testing.T) {
	f := `; this is a test ini file

var1=value1
var2=value2

[section]
var1=value3
var2=value4`

	ini := ini.New()
	err := ini.Load(bytes.NewReader([]byte(f)))
	if err != nil {
		t.Errorf("failed to parse ini: %s", err)
		return
	}

	if v, ok := ini.Get("root", "var1"); !ok || v != "value1" {
		t.Errorf("failed to get value root/var1, read %#v %#v", v, ok)
	}

	if v, ok := ini.Get("section", "var2"); !ok || v != "value4" {
		t.Errorf("failed to get value section/var2, read %#v %#v", v, ok)
	}
}
