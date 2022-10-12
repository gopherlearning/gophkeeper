package conf

type Args struct {
	Version VersionCmd `cmd:"" help:"Показать информацию о версии"`
	Verbose bool       `name:"verbose" short:"v" help:"Включить расширенное логирование"`
}

func (a *Args) Debug() bool {
	return a.Verbose
}

type Context struct {
	Verbose bool
}
