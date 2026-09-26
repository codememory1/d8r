package httpdownload

type Options struct {
	BufferSize       int64
	MaxParallelParts int
	RangeParts       int
}
