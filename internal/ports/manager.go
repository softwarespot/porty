package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/softwarespot/porty/internal/helpers"
)

type Manager struct {
	ports *Ports

	path     string
	file     *os.File
	fileLock *helpers.Flock
	mu       sync.RWMutex
}

func Init(path string) error {
	if helpers.FileExists(path) {
		return fmt.Errorf("ports database %q already exists", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable open/create the ports database %q: %w", path, err)
	}
	defer f.Close()

	ports := New()
	return writeAsJSON(path, f, ports)
}

func Load(path string) (*Manager, error) {
	m := &Manager{
		ports: nil,

		path:     path,
		file:     nil,
		fileLock: helpers.NewFlock(path),
	}

	var err error
	if m.file, err = os.OpenFile(m.path, os.O_RDWR, 0o666); err != nil {
		return nil, fmt.Errorf("unable to open/read the ports database %q: %w", m.path, err)
	}

	if err := json.NewDecoder(m.file).Decode(&m.ports); err != nil {
		if err := m.file.Close(); err != nil {
			return nil, fmt.Errorf("unable to close the ports database %q: %w", m.path, err)
		}
		return nil, fmt.Errorf("unable to read the ports database %q: %w", m.path, err)
	}

	if err := m.fileLock.Lock(true, 1*time.Second); err != nil {
		if err := m.file.Close(); err != nil {
			return nil, fmt.Errorf("unable to close the ports database %q: %w", m.path, err)
		}
		if errors.Is(err, helpers.ErrFlockTimeout) {
			return nil, fmt.Errorf("ports database %q is currently in use. Retry again", m.path)
		}
		return nil, fmt.Errorf("unable to lock the ports database %q: %w", m.path, err)
	}
	return m, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.fileLock.Unlock(); err != nil {
		return fmt.Errorf("unable to unlock the ports database %q: %w", m.path, err)
	}
	if err := m.file.Close(); err != nil {
		return fmt.Errorf("unable to close the ports database %q: %w", m.path, err)
	}
	return nil
}

func (m *Manager) All(sortBy SortBy) []UserPort {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.ports.All(sortBy)
}

func (m *Manager) AllByUsername(username string, sortBy SortBy) ([]UserPort, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.ports.AllByUsername(username, sortBy)
}

func (m *Manager) Register(username, appName string) (UserPort, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	up, err := m.ports.Register(username, appName)
	if err != nil {
		return userPortOnError, err
	}

	if err := m.write(); err != nil {
		return userPortOnError, err
	}
	return up, nil
}

func (m *Manager) Unregister(username, appName string) (UserPort, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	up, err := m.ports.Unregister(username, appName)
	if err != nil {
		return userPortOnError, err
	}

	if err := m.write(); err != nil {
		return userPortOnError, err
	}
	return up, nil
}

func (m *Manager) Get(username, appName string) (UserPort, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	up, err := m.ports.Get(username, appName)
	if err != nil {
		return userPortOnError, err
	}

	if err := m.write(); err != nil {
		return userPortOnError, err
	}
	return up, nil
}

func (m *Manager) GetByPort(port Port) (UserPort, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ports.GetByPort(port)
}

func (m *Manager) Next() (Port, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ports.Next()
}

type Info struct {
	Path    string
	MinPort Port
	MaxPort Port
}

func (m *Manager) Info() Info {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Info{
		Path:    m.path,
		MinPort: m.ports.MinPort,
		MaxPort: m.ports.MaxPort,
	}
}

func (m *Manager) ToPort(s string) (Port, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ports.ToPort(s)
}

func (m *Manager) ToSortBy(s string) (SortBy, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ports.ToSortBy(s)
}

func (m *Manager) write() error {
	return writeAsJSON(m.path, m.file, m.ports)
}

func writeAsJSON(path string, f *os.File, v any) error {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek to the start of the ports database %q: %w", path, err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		return fmt.Errorf("unable to write to the ports database %q: %w", path, err)
	}
	return nil
}
