package mydict

import (
	"errors"
)

type Dictionary map[string]string

var errNotFound = errors.New("Error: Not Found")
var errAlreadyExists = errors.New("Error: That word already exists")

func (d Dictionary) Search(word string) (string, error) {
	definition, exists := d[word]

	if exists {
		return definition, nil
	}

	return "", errNotFound
}

func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)

	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errAlreadyExists
	}

	return nil
}

func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	if err == errNotFound {
		return errNotFound
	}
	d[word] = def
	return nil
}

func (d Dictionary) Delete(word string) {
	delete(d, word)
}
