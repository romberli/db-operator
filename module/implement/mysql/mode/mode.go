package mode

type Mode int

const (
	Standalone Mode = iota + 1
	AsyncReplication
	SemiSyncReplication
	GroupReplication
)

func (m Mode) String() string {
	switch m {
	case Standalone:
		return "standalone"
	case AsyncReplication:
		return "async-replication"
	case SemiSyncReplication:
		return "semi-sync-replication"
	case GroupReplication:
		return "group-replication"
	default:
		return "unknown"
	}
}
