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
