package navigator

// Settings is a set of settings, which apply to all Navigator's instances
type Settings struct {
	// StoreDumpOnError flag for saving screenshots and pages if any error occurs
	StoreDumpOnError bool

	// DefaultDumpLocation is a path to directory where error reports would be stored
	DefaultDumpLocation string
}

func SetSettings(s Settings) {
	settings = s
}
