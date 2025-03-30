package ini

// Package ini provides a simple INI file parser and writer.
// https://en.wikipedia.org/wiki/INI_file

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

// Ini represents an INI file with sections and key-value pairs.
// Not thread-safe by default; use IniSafe for concurrent access.
type Ini map[string]map[string]string

// IniSafe is a thread-safe wrapper around Ini.
type IniSafe struct {
	data Ini
	mu   sync.RWMutex
}

// New returns a new Ini structure.
func New() Ini {
	return make(Ini)
}

// NewThreadSafe returns a new thread-safe Ini structure.
func NewThreadSafe() *IniSafe {
	return &IniSafe{
		data: make(Ini),
	}
}

// Load will parse source and merge loaded values.
func (i Ini) Load(source io.Reader) error {
	r := bufio.NewScanner(source)
	section := "root"
	var sectionMap map[string]string
	lineNum := 0

	for r.Scan() {
		lineNum++
		line := strings.TrimSpace(r.Text())
		if len(line) == 0 {
			continue
		}

		// Handle comments
		if len(line) > 0 && (line[0] == ';' || line[0] == '#') {
			// comment line
			continue
		}

		// Handle section headers
		if len(line) >= 2 && line[0] == '[' && line[len(line)-1] == ']' {
			section = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			if section == "" {
				return fmt.Errorf("line %d: empty section name", lineNum)
			}
			sectionMap = nil
			continue
		}

		pos := strings.IndexByte(line, '=')
		if pos < 0 {
			return fmt.Errorf("line %d: invalid format, missing '='", lineNum)
		}

		k := strings.ToLower(strings.TrimSpace(line[:pos]))
		if k == "" {
			return fmt.Errorf("line %d: empty key name", lineNum)
		}
		
		v := strings.TrimSpace(line[pos+1:])
		
		// Handle quotes
		if len(v) >= 2 && (v[0] == '"' && v[len(v)-1] == '"' || v[0] == '\'' && v[len(v)-1] == '\'') {
			v = v[1 : len(v)-1]
			// TODO: handle escape sequences properly
		}

		if sectionMap == nil {
			var ok bool
			sectionMap, ok = i[section]
			if !ok {
				sectionMap = make(map[string]string)
				i[section] = sectionMap
			}
		}

		sectionMap[k] = v
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	
	return nil
}

// Write generates an ini file and writes it to the provided output.
func (i Ini) Write(d io.Writer) error {
	var builder strings.Builder
	
	// Write root section first
	if s, ok := i["root"]; ok && len(s) > 0 {
		if err := i.writeSection(&builder, s); err != nil {
			return err
		}
		builder.WriteString("\n")
	}

	// Write other sections
	for n, s := range i {
		if n == "root" || len(s) == 0 {
			continue
		}

		builder.WriteString("[")
		builder.WriteString(n)
		builder.WriteString("]\n")

		if err := i.writeSection(&builder, s); err != nil {
			return err
		}
		builder.WriteString("\n")
	}
	
	_, err := d.Write([]byte(builder.String()))
	return err
}

func (i Ini) writeSection(b *strings.Builder, s map[string]string) error {
	for k, v := range s {
		// Check if value needs quoting
		needsQuotes := strings.ContainsAny(v, " \t\n\r")
		
		b.WriteString(k)
		b.WriteString("=")
		
		if needsQuotes {
			b.WriteString("\"")
			// Escape quotes in the value
			v = strings.ReplaceAll(v, "\"", "\\\"")
		}
		
		b.WriteString(v)
		
		if needsQuotes {
			b.WriteString("\"")
		}
		
		b.WriteString("\n")
	}
	return nil
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

// GetDefault returns a value for a given key or the provided default if not found.
func (i Ini) GetDefault(section, key, defaultValue string) string {
	if v, ok := i.Get(section, key); ok {
		return v
	}
	return defaultValue
}

// Set changes a value in the ini file.
func (i Ini) Set(section, key, value string) {
	section = strings.ToLower(section)
	key = strings.ToLower(key)
	
	s, ok := i[section]
	if !ok {
		s = make(map[string]string)
		i[section] = s
	}

	s[key] = value
}

// Unset removes a value from the ini file.
func (i Ini) Unset(section, key string) {
	section = strings.ToLower(section)
	key = strings.ToLower(key)
	
	s, ok := i[section]
	if !ok {
		return
	}

	delete(s, key)

	if len(s) == 0 {
		delete(i, section)
	}
}

// HasSection checks if a section exists.
func (i Ini) HasSection(section string) bool {
	_, ok := i[strings.ToLower(section)]
	return ok
}

// Sections returns a list of all section names.
func (i Ini) Sections() []string {
	sections := make([]string, 0, len(i))
	for section := range i {
		sections = append(sections, section)
	}
	return sections
}

// Keys returns a list of all keys in a section.
func (i Ini) Keys(section string) []string {
	section = strings.ToLower(section)
	s, ok := i[section]
	if !ok {
		return nil
	}
	
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}

// Thread-safe methods

// Load parses source and merges loaded values in a thread-safe manner.
func (i *IniSafe) Load(source io.Reader) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.data.Load(source)
}

// Write generates an ini file and writes it to the provided output in a thread-safe manner.
func (i *IniSafe) Write(d io.Writer) error {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.Write(d)
}

// Get returns a value for a given key in a thread-safe manner.
func (i *IniSafe) Get(section, key string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.Get(section, key)
}

// GetDefault returns a value or default in a thread-safe manner.
func (i *IniSafe) GetDefault(section, key, defaultValue string) string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.GetDefault(section, key, defaultValue)
}

// Set changes a value in a thread-safe manner.
func (i *IniSafe) Set(section, key, value string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data.Set(section, key, value)
}

// Unset removes a value in a thread-safe manner.
func (i *IniSafe) Unset(section, key string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data.Unset(section, key)
}

// HasSection checks if a section exists in a thread-safe manner.
func (i *IniSafe) HasSection(section string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.HasSection(section)
}

// Sections returns all section names in a thread-safe manner.
func (i *IniSafe) Sections() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.Sections()
}

// Keys returns all keys in a section in a thread-safe manner.
func (i *IniSafe) Keys(section string) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.Keys(section)
}