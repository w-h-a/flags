package writereader

import (
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/clients/writer"
)

type WriteReader interface {
	writer.Writer
	reader.Reader
}
