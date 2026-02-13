package sync

//nolint:revive
type SyncService struct{}

func New() (*SyncService, error) {
	return &SyncService{}, nil
}
