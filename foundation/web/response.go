package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type NoResponse struct{}

func NewNoResponse() NoResponse {
	return NoResponse{}
}

// Encode implements the Encoder interface.
func (NoResponse) Encode() ([]byte, string, error) {
	return nil, "", nil
}

// Response represents a standardized API response envelope
// @Description Standard API response wrapper
type Response[T any] struct {
	// The actual payload of the response
	Data T `json:"data"`
}

// NewResponse creates a new Response with the given data
func NewResponse[T any](data T) Response[T] {
	return Response[T]{
		Data: data,
	}
}

// Encode implements the Encoder interface
func (r Response[T]) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}

type httpStatus interface {
	HTTPStatus() int
}

// Respond sends a response to the client.
func Respond(ctx context.Context, w http.ResponseWriter, resp Encoder) error {
	if _, ok := resp.(NoResponse); ok {
		return nil
	}

	// If the context has been canceled, it means the client is no longer
	// waiting for a response.
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client disconnected, do not send response")
		}
	}

	statusCode := http.StatusOK

	switch v := resp.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()

	case error:
		statusCode = http.StatusInternalServerError

	default:
		if resp == nil {
			statusCode = http.StatusNoContent
		}
	}

	//_, span := addSpan(ctx, "web.send.response", attribute.Int("status", statusCode))
	//defer span.End()

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, contentType, err := resp.Encode()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("respond: encode: %w", err)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("respond: write: %w", err)
	}

	return nil
}
