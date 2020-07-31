package models

import (
	"errors"
	"fmt"
	"strings"

	"path/filepath"

	"github.com/jinzhu/gorm"
)

// Entry represents a path to a stored file or folder
type Entry struct {
	gorm.Model
	Path     string `gorm:"column:path;not null;unique"`
	IsFile   bool   `gorm:"column:isfile;not null"`
	Name     string `gorm:"column:name;not null"` // used for searching -- should be dirname or filename
	Ext      string `gorm:"column:ext"`
	Uploaded bool   `gorm:"column:uploaded"` // used during upload process
}

// CreateEntry creates a new entry at a given path
func CreateEntry(path string, isfile bool) (*Entry, error) {
	var matchingCount int
	err := db.Table("entries").Where("path = ?", path).Count(matchingCount).Error

	if err != nil {
		return nil, errors.New("Failed to connect to database")
	} else if matchingCount == 1 {
		return nil, errors.New("Entry already exists")
	}

	e := &Entry{Path: path, IsFile: isfile}
	getEntryInfo(e)
	db.Create(e)

	if e.ID <= 0 {
		return nil, errors.New("Failed to create entry")
	}

	return e, nil
}

// getEntryInfo takes an Enty produces its `name` and `ext` columns
func getEntryInfo(e *Entry) {
	if e.IsFile {
		e.Name = filepath.Dir(e.Path)
	} else {
		_, filename := filepath.Split(e.Path)
		ext := filepath.Ext(e.Path)

		e.Name = strings.TrimSuffix(filename, ext)
		e.Ext = ext
	}
}

// EntryExists checks if an entry matches a given path. If the first return
// value is `nil`, no match.
func EntryExists(path string, isfile bool) (*Entry, error) {
	e := &Entry{}
	err := db.Table("entries").Where("path = ?", path).Where("isfile = ?", isfile).First(e).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return e, nil
}

// SearchEntry searches for a match to a particular name query
func SearchEntry(query string) ([]*Entry, error) {
	var matches []*Entry

	err := db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query)).Find(&matches).Error

	if err != nil {
		return nil, err
	}

	return matches, nil
}

// GetEntryByID looks for a match to a particular entry ID
func GetEntryByID(id int) (*Entry, error) {
	e := &Entry{}
	err := db.Table("entries").Where("id = ?", id).First(e).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Unable to locate matching file")
		}

		return nil, err
	}

	return e, nil
}

// MarkAsUploaded marks an entry as uploaded
func MarkAsUploaded(e *Entry) {
	db.Model(e).Update("uploaded", true)
}
