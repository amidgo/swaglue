package httpmethod

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

var (
	HttpMethodList = [9]HttpMethod{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodConnect,
		MethodOptions,
	}
)

func (h HttpMethod) Valid() bool {
	for _, method := range HttpMethodList {
		if h == method {
			return true
		}
	}
	return false
}

func Valid(method string) bool {
	return HttpMethod(method).Valid()
}
