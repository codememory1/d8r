package respond

// Body defines a response body that can be written by a responder.
type Body interface {
	ResponseBody() any
}
