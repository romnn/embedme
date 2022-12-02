package embedme

// Options is used to configure embedme
type Options struct {
	StripEmbedComment bool
	Stdout            bool
	Verify            bool
	DryRun            bool
	Cwd               string
	Base              string
}
