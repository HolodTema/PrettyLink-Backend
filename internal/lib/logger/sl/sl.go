package sl

//we set _ before import to handle unused import error
//because of sqlite driver is unused, but we import it to init it
import (
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
