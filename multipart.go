package httpclient

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

// WriteMultipart is a function that writes multipart data to the given writer.
type WriteMultipart func(w MultipartWriter) error

// MultipartWriter is an interface that allows writing multipart data.
type MultipartWriter interface {
	// CreateFormField calls CreatePart with a header using the
	// given field name.
	CreateFormField(name string) (io.Writer, error)

	// CreateFormFile is a convenience wrapper around CreatePart. It creates
	// a new form-data header with the provided field name and file name.
	CreateFormFile(field, filename string) (io.Writer, error)

	// WriteField calls CreateFormField and then writes the given value.
	WriteField(fieldname, value string) error

	// CreatePart creates a new multipart section with the provided
	// header. The body of the part should be written to the returned
	// Writer. After calling CreatePart, any previous part may no longer
	// be written to.
	CreatePart(header textproto.MIMEHeader) (io.Writer, error)
}

// MultiFile creates a WriteMultipart that writes a file to the given field.
func MultipartFile(field, filename string, data io.Reader) WriteMultipart {
	return func(w MultipartWriter) error {
		file, errCreate := w.CreateFormFile(field, filename)
		if errCreate != nil {
			return errCreate
		}

		_, errCopy := io.Copy(file, data)
		return errCopy
	}
}

// MultiFiles creates a WriteMultipart that applies the given multipart writers.
func WriteMultiparts(writers ...WriteMultipart) WriteMultipart {
	return func(w MultipartWriter) error {
		for _, write := range writers {
			if err := write(w); err != nil {
				return err
			}
		}

		return nil
	}
}

// MultiFields creates a WriteMultipart that writes the given fields.
func MultipartFields(fields url.Values) WriteMultipart {
	return func(w MultipartWriter) error {
		for name, value := range fields {
			var v string
			if len(value) > 0 {
				v = value[0]
			}

			if err := w.WriteField(name, v); err != nil {
				return fmt.Errorf("write field %q: %w", name, err)
			}
		}

		return nil
	}
}

// PostMultipart sends a POST request with multipart data.
func (client *Client) PostMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, http.MethodPost, addr, writeMultipart)
}

// PutMultipart sends a PUT request with multipart data.
func (client *Client) PutMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, http.MethodPut, addr, writeMultipart)
}

// PatchMultipart sends a PATCH request with multipart data.
func (client *Client) PatchMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, http.MethodPatch, addr, writeMultipart)
}

func (client *Client) QueryMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, methodQuery, addr, writeMultipart)
}

func (client *Client) doMultipart(ctx context.Context, method, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	body, writer := io.Pipe()

	req, errReq := client.newRequest(ctx, method, addr, body)
	if errReq != nil {
		return nil, errReq
	}

	mwr := multipart.NewWriter(writer)
	req.Header.Set(headerContentType, mwr.FormDataContentType())

	done := make(chan error, 1)
	go func() {
		defer close(done)
		defer writer.Close()
		defer mwr.Close()

		done <- writeMultipart(mwr)
	}()

	resp, errDo := client.Doer.Do(req)
	if errDo != nil {
		return resp, errDo
	}

	return resp, <-done
}
