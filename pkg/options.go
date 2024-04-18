package embedme

// Options for embedme
type Options struct {
	StripEmbedComment bool
	Stdout            bool
	Verify            bool
	DryRun            bool
	WorkingDir        string
	Base              string
}

// NewDefaultOptions returns default options for embedme
func NewDefaultOptions() Options {
	return Options{
		StripEmbedComment: false,
		Stdout:            false,
		Verify:            false,
		DryRun:            false,
		WorkingDir:        "",
		Base:              "",
	}
}
