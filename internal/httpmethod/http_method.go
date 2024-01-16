package httpmethod

type HTTPMethod string

const (
	MethodGet     HTTPMethod = "get"
	MethodHead    HTTPMethod = "head"
	MethodPost    HTTPMethod = "post"
	MethodPut     HTTPMethod = "put"
	MethodPatch   HTTPMethod = "patch"
	MethodDelete  HTTPMethod = "delete"
	MethodConnect HTTPMethod = "connect"
	MethodOptions HTTPMethod = "options"
	MethodTrace   HTTPMethod = "trace"
)

func (h HTTPMethod) Valid() bool {
	methodList := [9]HTTPMethod{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodConnect,
		MethodOptions,
	}
	for _, method := range methodList {
		if h == method {
			return true
		}
	}

	return false
}

func Valid(method string) bool {
	return HTTPMethod(method).Valid()
}
