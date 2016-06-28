package plugin

import (
	"github.com/nobu-k/sb-image"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
)

func init() {
	udf.RegisterGlobalUDF("to_jpeg", udf.MustConvertGeneric(image.ToJPEG))
}
