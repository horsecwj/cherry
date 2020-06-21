package storage

type Storage interface {
	Save() (ok bool)
	DownloadUrl() (downloadUrl string)
}
