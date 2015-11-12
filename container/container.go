package container

const buildBotUrl string = "http://10.0.0.5/"

// ContainerBag carries all the instantiated dependencies necessary to the handlers work
type ContainerBag struct {
	BuildBotUrl string
}

// NewContainerBag return a new instance of the ContainerBag with the instantiated dependencies for the given config
func NewContainerBag() *ContainerBag {
	return &ContainerBag{
		BuildBotUrl: buildBotUrl,
	}
}
