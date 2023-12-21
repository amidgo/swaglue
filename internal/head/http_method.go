package head

type HttpMethod string

const (
	MethodGet     HttpMethod = "get"
	MethodHead    HttpMethod = "head"
	MethodPost    HttpMethod = "post"
	MethodPut     HttpMethod = "put"
	MethodPatch   HttpMethod = "patch"
	MethodDelete  HttpMethod = "delete"
	MethodConnect HttpMethod = "connect"
	MethodOptions HttpMethod = "options"
	MethodTrace   HttpMethod = "trace"
)

func (h HttpMethod) Valid() bool {
	for _, method := range []HttpMethod{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodConnect,
		MethodOptions,
		MethodTrace,
	} {
		if h == method {
			return true
		}
	}
	return false
}
