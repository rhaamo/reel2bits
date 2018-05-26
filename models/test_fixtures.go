package models

import (
	"github.com/go-testfixtures/testfixtures"
)

var fixtures *testfixtures.Context

// InitFixtures initialize test fixtures for a test database
func InitFixtures(helper testfixtures.Helper, dir string) (err error) {
	testfixtures.SkipDatabaseNameCheck(true)
	fixtures, err = testfixtures.NewFolder(db.DB(), helper, dir)
	return err
}

// LoadFixtures load fixtures for a test database
func LoadFixtures() error {
	db.AutoMigrate(&User{}, &Album{}, &Track{}, &TrackInfo{})
	return fixtures.Load()
}
