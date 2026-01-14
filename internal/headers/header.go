package headers

type Headers map[string]string

func NewHeaders() map[string]string {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	return 0, false, nil
}
