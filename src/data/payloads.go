package data

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Payload represents a single payload item
type Payload struct {
	path      string
	Filename  string    `json:"filename" yaml:"filename"`
	Size      int64     `json:"size" yaml:"size"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

// Contents returns the payload contents and error if something
// goes wrong
func (p Payload) Contents() ([]byte, error) {
	return os.ReadFile(p.path)
}

// ContentDisposition returns the Content-Disposition header
// value for particular payload
func (p Payload) ContentDisposition() string {
	return fmt.Sprintf(`attachment; filename="%s"`, p.Filename)
}

// Remove deletes the payload from the storage
func (p Payload) Remove() error {
	return os.Remove(p.path)
}

// PayloadsStore is the store implementation for payloads
type PayloadsStore struct {
	log     zerolog.Logger
	content map[string]Payload
	baseDir string

	mux sync.Mutex
}

// NewPayloadsStore creates new PayloadsStore for runtime
func NewPayloadsStore(log zerolog.Logger) (*PayloadsStore, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "goosymock-payloads-*")
	if err != nil {
		return nil, errors.Wrap(err, "os.MkdirTemp")
	}

	log.Debug().Str("path", tempDir).Msg("payloads store created")

	return &PayloadsStore{
		log:     log,
		content: map[string]Payload{},
		baseDir: tempDir,
		mux:     sync.Mutex{},
	}, nil
}

// GetAll implements FileStore
func (s *PayloadsStore) GetAll() map[string]Payload {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.content
}

// Get implements FileStore
func (s *PayloadsStore) Get(id string) (Payload, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var err error
	route, ok := s.content[id]
	if !ok {
		err = fmt.Errorf("payload with given ID does not exist: %s", id)
	}

	return route, err
}

// Set implements FileStore
func (s *PayloadsStore) Set(payloads map[string]Payload) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	err := s.DeleteAll()
	s.content = payloads
	return err
}

// Insert implements FileStore
func (s *PayloadsStore) Insert(id string, payload Payload) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, exists := s.content[id]; exists {
		return fmt.Errorf("payload with given ID already exists: %s", id)
	}
	s.content[id] = payload
	return nil
}

// Update implements FileStore
func (s *PayloadsStore) Update(id string, payload Payload) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, exists := s.content[id]; !exists {
		return fmt.Errorf("payload with given ID does not exist: %s", id)
	}
	s.content[id] = payload
	return nil
}

// Upsert implements FileStore
func (s *PayloadsStore) Upsert(id string, payload Payload) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.content[id] = payload
	return nil
}

// Delete implements FileStore
func (s *PayloadsStore) Delete(id string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	payload, exists := s.content[id]
	if !exists {
		return fmt.Errorf("payload with given ID does not exist: %s", id)
	}

	if err := payload.Remove(); err != nil {
		return errors.Wrap(err, "os.Remove")
	}
	return nil
}

// DeleteAll implements FileStore
func (s *PayloadsStore) DeleteAll() error {
	errors := []string{}
	for id, item := range s.content {
		if err := item.Remove(); err != nil {
			errors = append(errors, err.Error())
		}
		delete(s.content, id)
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}
	return nil
}

// Count implements FileStore
func (s *PayloadsStore) Count() int {
	return len(s.content)
}

func getUniqueID() string {
	id, _ := gonanoid.New()
	return id
}

func (s *PayloadsStore) createTempFile(src io.ReadCloser) (string, int64, error) {
	dst, err := os.CreateTemp(s.baseDir, "*.payload")
	if err != nil {
		return "", 0, errors.Wrap(err, "os.CreateTemp")
	}

	filename := dst.Name()

	defer dst.Close()
	defer src.Close()
	size, err := io.Copy(dst, src)
	if err != nil {
		return "", 0, errors.Wrap(err, "io.Copy")
	}

	return filename, size, nil
}

// InsertFile implements FileStore
func (s *PayloadsStore) InsertFile(name string, src io.ReadCloser) (string, error) {
	id := getUniqueID()

	filename, size, err := s.createTempFile(src)
	if err != nil {
		return "", err
	}

	if err := s.Insert(id, Payload{
		path:      filename,
		Filename:  name,
		Size:      size,
		Timestamp: time.Now().UTC(),
	}); err != nil {
		return "", err
	}
	return id, nil
}

// UpdateFile implements FileStore
func (s *PayloadsStore) UpdateFile(id string, name string, src io.ReadCloser) error {
	filename, size, err := s.createTempFile(src)
	if err != nil {
		return err
	}

	if err := s.Update(id, Payload{
		path:      filename,
		Filename:  name,
		Size:      size,
		Timestamp: time.Now().UTC(),
	}); err != nil {
		return err
	}
	return nil
}

// UpsertFile implements FileStore
func (s *PayloadsStore) UpsertFile(id string, name string, src io.ReadCloser) error {
	filename, size, err := s.createTempFile(src)
	if err != nil {
		return err
	}

	if err := s.Upsert(id, Payload{
		path:      filename,
		Filename:  filename,
		Size:      size,
		Timestamp: time.Now().UTC(),
	}); err != nil {
		return err
	}
	return nil
}

// Close implements FileStore
func (s *PayloadsStore) Close() error {
	return os.RemoveAll(s.baseDir)
}

var _ FilesStore[Payload] = &PayloadsStore{}
