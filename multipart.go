package httpclient

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type WriteMultipart func(w *multipart.Writer) error

func MultipartFile(field, filename string, data io.Reader) WriteMultipart {
	return func(w *multipart.Writer) error {
		defer w.Close()

		file, errCreate := w.CreateFormFile(field, filename)
		if errCreate != nil {
			return errCreate
		}

		_, errCopy := io.Copy(file, data)
		return errCopy
	}
}

func FormFields(fields url.Values) WriteMultipart {
	return func(w *multipart.Writer) error {
		defer w.Close()

		for name, value := range fields {
			wr, errField := w.CreateFormField(name)
			if errField != nil {
				return errField
			}

			if len(value) > 0 {
				_, errWrite := io.WriteString(wr, value[0])
				if errWrite != nil {
					return errWrite
				}
			}
		}

		return nil
	}
}

func (client *Client) PostMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, http.MethodPost, addr, writeMultipart)
}

func (client *Client) PutMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, http.MethodPut, addr, writeMultipart)
}

func (client *Client) PatchMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, http.MethodPatch, addr, writeMultipart)
}

func (client *Client) QueryMultipart(ctx context.Context, addr string, writeMultipart WriteMultipart) (*http.Response, error) {
	return client.doMultipart(ctx, MethodQuery, addr, writeMultipart)
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

		done <- writeMultipart(mwr)
	}()

	resp, errDo := client.Transport.Do(req)
	if errDo != nil {
		return resp, errDo
	}

	return resp, <-done
}
