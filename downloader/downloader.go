package downloader

import "errors"

// Downloader takes uri and download them
type Downloader interface {
	Auth(infos map[string]interface{}) error
	Download(uri string, options map[string]interface{}) (string, error)
	DownloadStatus(downloadID string) (map[string]interface{}, error)
}

// NewDownloader returns a new authenticated downloader
func NewDownloader(name string, authInfos map[string]interface{}) (Downloader, error) {
	var downloader Downloader

	switch name {
	case "Aria2", "aria2", "aria":
		downloader = &Aria2{}
	default:
		return nil, errors.New("Invalid downloader given")
	}

	if authInfos != nil {
		authError := downloader.Auth(authInfos)
		if authError != nil {
			return nil, authError
		}
	}

	return downloader, nil
}
