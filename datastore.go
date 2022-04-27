package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// Get and Delete will return NoteNotFoundError if note could not be located; this should probably not be treating notes as strings if they can be large...
type DataStore interface {
	Store(user string, name string, note string) error
	Get(user string, name string) (string, error)
	Delete(user string, name string) error
}

type NoteNotFoundError error

type cacheDataStoreImpl struct {
	c *cache.Cache
}

func (ds * cacheDataStoreImpl) Store(user string, name string, note string) error {
	return ds.c.Add(fmt.Sprintf("%s-%s", user, name), note, cache.DefaultExpiration)
}

func (ds * cacheDataStoreImpl) Get(user string, name string) (string, error) {
	note, found := ds.c.Get(fmt.Sprintf("%s-%s", user, name))
	if !found {
		return "", NoteNotFoundError(errors.New("note not found"))
	}
	return note.(string), nil
}

func (ds * cacheDataStoreImpl) Delete(user string, name string) error {
	key := fmt.Sprintf("%s-%s", user, name)
	_, found := ds.c.Get(key)
	if !found {
		return NoteNotFoundError(errors.New("note not found"))
	}
	ds.c.Delete(key)
	return nil
}


// cannot store usernames or note names with dash (-) in them (probably base64 encode the note names)
func NewCacheDataStore() DataStore {
	c := cache.New(cache.NoExpiration, time.Duration(-1))
	return &cacheDataStoreImpl{c: c}
}