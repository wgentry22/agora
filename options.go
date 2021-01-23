package agora

type Options struct {
	configFilePath string
	configFileName string
}

func defaultOptions() *Options {
	return &Options{
		configFilePath: ".",
		configFileName: "app.toml",
	}
}

func New(confs ...Configuration) {
	opts := defaultOptions()

	for _, conf := range confs {
		conf(opts)
	}
}
