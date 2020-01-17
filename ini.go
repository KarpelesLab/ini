package ini

// https://en.wikipedia.org/wiki/INI_file

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Ini map[string]map[string]string

// New returns a new Ini structure
func New() Ini {
	return make(Ini)
}

// Load will parse source and merge loaded values
func (i Ini) Load(source io.Reader) error {
	r := bufio.NewScanner(source)
	section := "root"
	var sectionMap map[string]string

	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if len(line) == 0 {
			continue
		}

		if line[0] == ';' {
			// comment line
			continue
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			section = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			sectionMap = nil
			continue
		}

		pos := strings.IndexByte(line, '=')
		if pos < 0 {
			return errors.New("failed to parse ini file: invalid line")
		}

		k := strings.ToLower(strings.TrimSpace(line[:pos]))
		line = strings.TrimSpace(line[pos+1:])

		// TODO: handle quotes, handle escape characters

		if sectionMap == nil {
			var ok bool
			sectionMap, ok = i[section]
			if !ok {
				sectionMap = make(map[string]string)
				i[section] = sectionMap
			}
		}

		sectionMap[k] = line
	}

	return r.Err()
}

// Write generates a ini file and writes it to the provided output
func (i Ini) Write(d io.Writer) error {
	if s, ok := i["root"]; ok {
		if err := i.writeSection(d, s); err != nil {
			return err
		}
	}

	for n, s := range i {
		if n == "root" {
			continue
		}

		_, err := d.Write(append(append([]byte{'['}, []byte(n)...), ']', '\n'))
		if err != nil {
			return err
		}

		if err := i.writeSection(d, s); err != nil {
			return err
		}
	}
	return nil
}

func (i Ini) writeSection(d io.Writer, s map[string]string) error {
	for k, v := range s {
		_, err := d.Write(append(append(append([]byte(k), '='), []byte(v)...), '\n'))
		if err != nil {
			return err
		}
	}
	_, err := d.Write([]byte{'\n'})
	return err
}

// Get returns a value for a given key. Use section "root" for entries at the
// beginning of the file.
func (i Ini) Get(section, key string) (string, bool) {
	s, ok := i[strings.ToLower(section)]
	if !ok {
		return "", false
	}

	r, ok := s[strings.ToLower(key)]
	return r, ok
}

// Set changes a value in the ini file
func (i Ini) Set(section, key, value string) {
	s, ok := i[strings.ToLower(section)]
	if !ok {
		s = make(map[string]string)
		i[strings.ToLower(section)] = s
	}

	s[strings.ToLower(key)] = value
}

// Unset removes a value from the ini file
func (i Ini) Unset(section, key string) {
	s, ok := i[strings.ToLower(section)]
	if !ok {
		return
	}

	delete(s, strings.ToLower(key))

	if len(s) == 0 {
		delete(i, strings.ToLower(section))
	}
}
