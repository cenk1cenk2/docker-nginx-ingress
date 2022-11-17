package pipe

type Ctx struct {
	Directories struct {
		ServerConfiguration   string
		UpstreamConfiguration string
	}
	Templates struct {
		Server   string
		Upstream string
	}
}
