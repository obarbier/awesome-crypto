package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func respondError(w http.ResponseWriter, allowed int, t interface{}) {

}

func parseJSONRequest(r *http.Request, w http.ResponseWriter, out interface{}) (io.ReadCloser, error) {
	// Limit the maximum number of bytes to MaxRequestSize to protect
	// against an indefinite amount of data being read.
	reader := r.Body
	ctx := r.Context()
	maxRequestSize := ctx.Value("max_request_size")
	if maxRequestSize != nil {
		max, ok := maxRequestSize.(int64)
		if !ok {
			return nil, errors.New("could not parse max_request_size from request context")
		}
		if max > 0 {
			// MaxBytesReader won't do all the internal stuff it must unless it's
			// given a ResponseWriter that implements the internal http interface
			// requestTooLarger.  So we let it have access to the underlying
			// ResponseWriter.
			inw := w
			//if myw, ok := inw.(logical.WrappingResponseWriter); ok {
			//	inw = myw.Wrapped()
			//}
			reader = http.MaxBytesReader(inw, r.Body, max)
		}
	}
	var origBody io.ReadWriter
	err := DecodeJSONFromReader(reader, out)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to parse JSON input: %w", err)
	}
	if origBody != nil {
		return ioutil.NopCloser(origBody), err
	}
	return nil, err
}

// Decodes/Unmarshals the given io.Reader pointing to a JSON, into a desired object
func DecodeJSONFromReader(r io.Reader, out interface{}) error {
	if r == nil {
		return fmt.Errorf("'io.Reader' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	dec := json.NewDecoder(r)

	// While decoding JSON values, interpret the integer values as `json.Number`s instead of `float64`.
	dec.UseNumber()

	// Since 'out' is an interface representing a pointer, pass it to the decoder without an '&'
	return dec.Decode(out)
}
