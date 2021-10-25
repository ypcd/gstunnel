package gspackoper

import (
	"bytes"
	"compress/flate"
	"io"
	"log"
)

var logger *log.Logger

func init() {
	logger = log.Default()
}

type compresser struct {
	fewriter *flate.Writer
	fereader io.ReadCloser
}

func NewCompresser() *compresser {
	return &compresser{}
}

func (ap *compresser) compress2(data []byte) []byte   { return data }
func (ap *compresser) uncompress2(data []byte) []byte { return data }

func (ap *compresser) compress(data []byte) []byte {
	var b bytes.Buffer

	if ap.fewriter == nil {
		zw, err := flate.NewWriter(&b, 1)
		ap.fewriter = zw
		if err != nil {
			logger.Fatalln(err)
		}
	} else {
		ap.fewriter.Reset(&b)
	}

	zw := ap.fewriter

	if _, err := io.Copy(zw, bytes.NewReader(data)); err != nil {
		logger.Fatalln(err)
	}
	if err := zw.Close(); err != nil {
		logger.Fatalln(err)
	}

	return b.Bytes()
}

func (ap *compresser) uncompress(data []byte) []byte {
	var b bytes.Buffer

	if ap.fereader == nil {
		zr := flate.NewReader(bytes.NewReader(data))
		ap.fereader = zr
	} else {
		zr := ap.fereader
		if err := zr.(flate.Resetter).Reset(bytes.NewReader(data), nil); err != nil {
			logger.Fatalln(err)
		}
	}
	zr := ap.fereader

	if _, err := io.Copy(&b, zr); err != nil {
		logger.Fatalln(err)
	}
	if err := zr.Close(); err != nil {
		logger.Fatalln(err)
	}

	return b.Bytes()
}
