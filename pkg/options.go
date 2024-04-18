package embedme

// Options is used to configure embedme
type Options struct {
	StripEmbedComment bool
	Stdout            bool
	Verify            bool
	DryRun            bool
	WorkingDir        string
	Base              string
}

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
