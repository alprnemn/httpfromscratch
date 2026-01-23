package request

var validMethods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
	"PATCH":  {},
}

type parserState string

const BufferSize = 8

const (
	StateHeaders parserState = "reqline"
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateError   parserState = "error"
	StateBody    parserState = "body"
)
