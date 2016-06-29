package plugin

import (
	"github.com/nobu-k/sb-image"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
)

func init() {
	udf.RegisterGlobalUDF("encode_jpeg", udf.MustConvertGeneric(image.EncodeJPEG))
	udf.RegisterGlobalUDF("decode_jpeg", udf.MustConvertGeneric(image.DecodeJPEG))
}
