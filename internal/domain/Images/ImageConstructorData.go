package Images

import "io"

type ImageConstructorData struct {
	FileName    string
	File        io.Reader
	FramesCount int `default:"-1"`
}
