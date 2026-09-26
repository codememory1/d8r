package inspect

// Options contains configuration used while inspecting HTTP resources.
type Options struct {
	MinRangeProbeSize int64
	ReservedHeaders   []string
}
