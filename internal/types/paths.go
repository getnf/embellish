package types

type Paths struct {
	Download string
	Install  string
	Db       string
}

func (p *Paths) GetDownloadPath() string {
	return p.Download
}

func (p *Paths) GetInstallPath() string {
	return p.Install
}

func (p *Paths) GetDbPath() string {
	return p.Db
}
