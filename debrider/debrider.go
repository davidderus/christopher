package debrider

import "errors"

// Debrider takes an URI and return a debrided URI
type Debrider interface {
	Init() error
	Auth(infos map[string]string) error
	Debrid(uri string, options map[string]interface{}) (string, error)

	// IsDebridable will be used to switch between multiple debriders
	IsDebridable(uri string) bool
}

// NewDebrider returns a new initialized debrider
//
// Authentication is optionnal to allow access to some methods which do not
// require it.
func NewDebrider(name string, authInfos map[string]string) (Debrider, error) {
	var debrider Debrider

	switch name {
	case "alldebrid", "AllDebrid", "Alldebrid", "ad":
		debrider = &AllDebrid{}
	default:
		return nil, errors.New("Invalid debrider given")
	}

	// Do init things
	initError := debrider.Init()
	if initError != nil {
		return nil, initError
	}

	// Authenticate the debrider
	if authInfos != nil {
		authError := debrider.Auth(authInfos)
		if authError != nil {
			return nil, authError
		}
	}

	return debrider, nil
}
