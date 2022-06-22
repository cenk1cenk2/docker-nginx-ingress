package pipe

type Ctx struct {
	NginxConfiguration Configuration
	Directories        struct {
		ServerConfiguration   string
		UpstreamConfiguration string
	}
	Templates struct {
		Server   string
		Upstream string
	}
}
