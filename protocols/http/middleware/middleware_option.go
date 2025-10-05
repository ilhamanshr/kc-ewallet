package middleware

type middlewareOption struct {
	excludedHandlers map[string]bool
	registerHandlers map[string]bool
}

type middlewareOptionFn func(*middlewareOption)

func defaultMiddlewareOption() *middlewareOption {
	return &middlewareOption{
		excludedHandlers: nil,
		registerHandlers: nil,
	}
}

// Used for excluding certain path from executing middleware
// by providing last handler names as key and set value true
func ExcludeHandlers(handlerNames map[string]bool) middlewareOptionFn {
	return func(mo *middlewareOption) {
		mo.excludedHandlers = handlerNames
	}
}

// Used for registering certain path to executing middleware
// by providing last handler names as key and set value true
func RegisterHandlers(handlerNames map[string]bool) middlewareOptionFn {
	return func(mo *middlewareOption) {
		mo.registerHandlers = handlerNames
	}
}
