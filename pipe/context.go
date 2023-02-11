package pipe

type Ctx struct {
	Directories struct {
		ServerConfiguration   string
		UpstreamConfiguration string
	}
	Templates struct {
		Nginx    string
		Server   string
		Upstream string
	}
}
