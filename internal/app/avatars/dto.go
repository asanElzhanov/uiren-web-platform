package avatars

import "io"

type UploadAvatarRequest struct {
	UserId string
	File   io.Reader
}

type Avatar struct {
	Content     []byte
	ContentType string
}
