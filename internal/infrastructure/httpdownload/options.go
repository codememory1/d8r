package httpdownload

// Options contains configuration used by the HTTP downloader.
type Options struct {
	BufferSize       int64
	MaxParallelParts int
	RangeParts       int
}
