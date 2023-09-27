package injector

var defaultInjectContainerKey string = "__X_INJECT_CONTAINER"

func setCtxContainerKey(key string) {
	defaultInjectContainerKey = key
}
