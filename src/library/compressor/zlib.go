package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
)

import (
	"library/logger"
)

const ENABLE_ZLIB = true

const (
	NoCompression      = zlib.NoCompression
	BestSpeed          = zlib.BestSpeed
	BestCompression    = zlib.BestCompression
	DefaultCompression = zlib.DefaultCompression
)

var allLevel []int = []int{NoCompression, BestSpeed, BestCompression, DefaultCompression}

func isLegalLevel(level int) bool {
	for _, v := range allLevel {
		if v == level {
			return true
		}
	}
	logger.Error("zlib level error:level:%s", level)
	return false
}

func Compress(input []byte, level int) ([]byte, bool) {
	if ENABLE_ZLIB == false {
		return input, true
	}
	if !isLegalLevel(level) {
		return nil, false
	}
	var output bytes.Buffer
	w, err := zlib.NewWriterLevel(&output, level)
	if err != nil {
		logger.Error("set write level error,err:%s", err.Error())
		return nil, false
	}
	_, err = w.Write(input)
	if err != nil {
		logger.Error("compress error,err:%s", err.Error())
		w.Close()
		return nil, false
	}

	w.Close() // close would implicitly flush. but defer close() would not
	logger.Info("compress rate,%d/%d", output.Len(), len(input))
	return output.Bytes(), true
}

func Decompress(input []byte) ([]byte, bool) {
	if ENABLE_ZLIB == false {
		return input, false
	}

	var output bytes.Buffer
	bf := bytes.NewReader(input)
	r, err := zlib.NewReader(bf)
	if err != nil {
		logger.Error("decompress error, %+v", input)
		return nil, false
	}
	io.Copy(&output, r)
	return output.Bytes(), true
}
