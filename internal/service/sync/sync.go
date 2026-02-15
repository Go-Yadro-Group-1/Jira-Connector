package sync

type SyncService struct {
}

func New() (*SyncService, error) {
	return &SyncService{}, nil
}
