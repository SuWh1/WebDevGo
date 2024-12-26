package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

// contentType detection: type of files
func checkContentType(r io.ReadSeeker, allowedTypes []string) error { // readSeeker - reset file back to beginning, and then it again to copy the file it start at zero
	// just reader is like book mark, it saves where we read last time, readseeker everytime starts at the beginning
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}
	// reset file, put it back where supposed to be
	_, err = r.Seek(0, 0) // send back
	if err != nil {
		return fmt.Errorf("resetting file position: %w", err)
	}

	contentType := http.DetectContentType(testBytes) // check content type
	for _, t := range allowedTypes {
		if t == contentType {
			return nil // types is allowed
		}
	}
	return FileError{
		Issue: fmt.Sprintf("invalid content type: %v", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid extention: %v", filepath.Ext(filename)),
		}
	}
	return nil
}
