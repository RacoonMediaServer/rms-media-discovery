package media

// Quality represents video quality
type Quality int

const (
	QualityUnd Quality = iota
	Quality480p
	Quality720p
	Quality1080p
	Quality2160p
)

func (q Quality) String() string {
	switch q {
	case Quality480p:
		return "480p"
	case Quality720p:
		return "720p"
	case Quality1080p:
		return "1080p"
	case Quality2160p:
		return "2160p"
	default:
		return ""
	}
}
