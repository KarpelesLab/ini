package ini

import (
	"io"
	"sync"
)

// IniSafe is a thread-safe wrapper around Ini.
type IniSafe struct {
	data Ini
	mu   sync.RWMutex
}

// NewThreadSafe returns a new thread-safe Ini structure.
func NewThreadSafe() *IniSafe {
	return &IniSafe{
		data: make(Ini),
	}
}

// Load parses source and merges loaded values in a thread-safe manner.
// Deprecated: Use ReadFrom instead which implements io.ReaderFrom interface.
func (i *IniSafe) Load(source io.Reader) error {
	_, err := i.ReadFrom(source)
	return err
}

// ReadFrom implements the io.ReaderFrom interface with thread safety.
func (i *IniSafe) ReadFrom(source io.Reader) (int64, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.data.ReadFrom(source)
}

// Write generates an ini file and writes it to the provided output in a thread-safe manner.
// Deprecated: Use WriteTo instead which implements io.WriterTo interface.
func (i *IniSafe) Write(d io.Writer) error {
	_, err := i.WriteTo(d)
	return err
}

// WriteTo implements the io.WriterTo interface with thread safety.
func (i *IniSafe) WriteTo(d io.Writer) (int64, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.WriteTo(d)
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
