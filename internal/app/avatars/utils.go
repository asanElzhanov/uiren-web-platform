package avatars

import (
	"bytes"
	"io"
	"net/http"
)

func generateAvatarFilename(id string) string {
	return "avatar-" + id + ".png"
}

func readFileContentAndType(file io.Reader) ([]byte, string, error) {
	var buf bytes.Buffer

	tee := io.TeeReader(file, &buf)

	head := make([]byte, 512)
	n, err := tee.Read(head)
	if err != nil && err != io.EOF {
		return nil, "", err
	}

	contentType := http.DetectContentType(head[:n])

	rest, err := io.ReadAll(tee)
	if err != nil {
		return nil, "", err
	}

	content := append(head[:n], rest...)
	return content, contentType, nil
}
