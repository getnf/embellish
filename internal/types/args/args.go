package args

// Command line argumetns

type InstallCmd struct {
	Fonts    []string `arg:"positional" help:"list of space separated fonts to install"`
	KeepTars bool     `arg:"-k" help:"Keep archives in the download location"`
}

type UninstallCmd struct {
	Fonts []string `arg:"positional" help:"list of space separated fonts to uninstall"`
}

type ListCmd struct {
	Installed bool `arg:"-i" help:"list only installed fonts"`
}

type UpdateCmd struct {
	Update   bool `default:"true"`
	KeepTars bool `arg:"-k" help:"Keep archives in the download location"`
}

type Args struct {
	Install    *InstallCmd   `arg:"subcommand:install" help:"install fonts"`
	Uninstall  *UninstallCmd `arg:"subcommand:uninstall" help:"uninstall fonts"`
	List       *ListCmd      `arg:"subcommand:list" help:"list fonts"`
	Update     *UpdateCmd    `arg:"subcommand:update" help:"update installed fonts"`
	ForceCheck bool          `arg:"-f" help:"Force checking for updates"`
	Global     bool          `arg:"-g" help:"Do the operation globally, for all users"`
}

func (Args) Version() string {
	return "getnf v1.0.0"
}
