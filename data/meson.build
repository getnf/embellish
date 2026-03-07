blueprints = custom_target('blueprints',
  input: files(
    'shortcuts-dialog.blp',
    'ui/IconsPage.blp',
    'ui/Window.blp',
    'ui/ImportDialog.blp'
  ),
  output: '.',
  command: [find_program('blueprint-compiler'), 'batch-compile', '@OUTPUT@', '@CURRENT_SOURCE_DIR@', '@INPUT@'],
)

# Compile the resources
gnome.compile_resources(
	APP_ID + '.data',
	APP_ID + '.data.gresource.xml',
	gresource_bundle: true,
	install: true,
	install_dir: get_option('datadir') / APP_ID,
	dependencies: blueprints,
)

# Install the settings schema
install_data(
	APP_ID + '.gschema.xml',
	install_dir: get_option('datadir') / 'glib-2.0' / 'schemas'
)

# Install the application icons
install_data(
	'icons' / APP_ID + '.svg',
	install_dir: get_option('datadir') / 'icons' / 'hicolor' / 'scalable' / 'apps'
)
install_data(
	'icons' / APP_ID + '-symbolic.svg',
	install_dir: get_option('datadir') / 'icons' / 'hicolor' / 'symbolic' / 'apps'
)

# Install the D-Bus service file
configure_file(
	input: APP_ID + '.service',
	output: APP_ID + '.service',
	configuration: {
		'APP_ID': APP_ID,
		'BINDIR': get_option('prefix') / get_option('bindir')
	},
	install: true,
	install_dir: get_option('datadir') / 'dbus-1' / 'services'
)

# Translate and install the desktop file
i18n.merge_file(
	type: 'desktop',
	input: configure_file(
		input: APP_ID + '.desktop',
		output: APP_ID + '.desktop.in',
		configuration: {
			'APP_ID': APP_ID,
			'BINDIR': get_option('prefix') / get_option('bindir')
		}
	),
	output: APP_ID + '.desktop',
	po_dir: meson.project_source_root() / 'po',
	install: true,
	install_dir: get_option('datadir') / 'applications'
)

# Translate and install the AppStream metadata
i18n.merge_file(
	input: configure_file(
		input: APP_ID + '.metainfo.xml',
		output: APP_ID + '.metainfo.xml.in',
		configuration: {
			'APP_ID': APP_ID,
		}
	),
	output: APP_ID + '.metainfo.xml',
	type: 'xml',
	po_dir: meson.project_source_root() / 'po',
	install: true,
	install_dir: get_option('datadir') / 'metainfo'
)

# Install the font files
font_files = files(
  'font-files/0xProtoNerdFontMono-Regular.ttf',
  'font-files/3270NerdFontMono-Regular.ttf',
  'font-files/AdwaitaMonoNerdFontMono-Regular.ttf',
  'font-files/AgaveNerdFontMono-Regular.ttf',
  'font-files/AnonymiceProNerdFontMono-Regular.ttf',
  'font-files/ArimoNerdFont-Regular.ttf',
  'font-files/AtkynsonMonoNerdFontMono-Regular.otf',
  'font-files/AurulentSansMNerdFontMono-Regular.otf',
  'font-files/BigBlueTerm437NerdFontMono-Regular.ttf',
  'font-files/BitstromWeraNerdFontMono-Regular.ttf',
  'font-files/BlexMonoNerdFontMono-Regular.ttf',
  'font-files/CaskaydiaCoveNerdFontMono-Regular.ttf',
  'font-files/CaskaydiaMonoNerdFontMono-Regular.ttf',
  'font-files/CodeNewRomanNerdFontMono-Regular.otf',
  'font-files/ComicShannsMonoNerdFontMono-Regular.otf',
  'font-files/CommitMonoNerdFontMono-Regular.otf',
  'font-files/CousineNerdFontMono-Regular.ttf',
  'font-files/D2CodingLigatureNerdFontMono-Regular.ttf',
  'font-files/DaddyTimeMonoNerdFontMono-Regular.ttf',
  'font-files/DejaVuSansMNerdFontMono-Regular.ttf',
  'font-files/DepartureMonoNerdFontMono-Regular.otf',
  'font-files/DroidSansMNerdFontMono-Regular.otf',
  'font-files/EnvyCodeRNerdFontMono-Regular.ttf',
  'font-files/FantasqueSansMNerdFontMono-Regular.ttf',
  'font-files/FiraCodeNerdFontMono-Regular.ttf',
  'font-files/FiraMonoNerdFontMono-Regular.otf',
  'font-files/GeistMonoNerdFontMono-Regular.otf',
  'font-files/GohuFont11NerdFontMono-Regular.ttf',
  'font-files/GoMonoNerdFontMono-Regular.ttf',
  'font-files/HackNerdFontMono-Regular.ttf',
  'font-files/HasklugNerdFontMono-Regular.otf',
  'font-files/HeavyDataNerdFont-Regular.ttf',
  'font-files/HurmitNerdFontMono-Regular.otf',
  'font-files/iMWritingMonoNerdFontMono-Regular.ttf',
  'font-files/InconsolataGoNerdFontMono-Regular.ttf',
  'font-files/InconsolataLGCNerdFontMono-Regular.ttf',
  'font-files/InconsolataNerdFontMono-Regular.ttf',
  'font-files/IntoneMonoNerdFontMono-Regular.ttf',
  'font-files/IosevkaNerdFontMono-Regular.ttf',
  'font-files/IosevkaTermNerdFontMono-Regular.ttf',
  'font-files/IosevkaTermSlabNerdFontMono-Regular.ttf',
  'font-files/JetBrainsMonoNerdFontMono-Regular.ttf',
  'font-files/LektonNerdFontMono-Regular.ttf',
  'font-files/LilexNerdFontMono-Regular.ttf',
  'font-files/LiterationMonoNerdFontMono-Regular.ttf',
  'font-files/M+1CodeNerdFontMono-Regular.ttf',
  'font-files/MartianMonoNerdFontMono-Regular.ttf',
  'font-files/MesloLGLDZNerdFontMono-Regular.ttf',
  'font-files/MonaspiceArNerdFontMono-Regular.otf',
  'font-files/MonofurNerdFontMono-Regular.ttf',
  'font-files/MonoidNerdFontMono-Regular.ttf',
  'font-files/MononokiNerdFontMono-Regular.ttf',
  'font-files/NotoMonoNerdFontMono-Regular.ttf',
  'font-files/OpenDyslexicMNerdFontMono-Regular.otf',
  'font-files/OverpassMNerdFontMono-Regular.otf',
  'font-files/ProFontIIxNerdFontMono-Regular.ttf',
  'font-files/ProggyCleanCENerdFontMono-Regular.ttf',
  'font-files/RecMonoCasualNerdFontMono-Regular.ttf',
  'font-files/RobotoMonoNerdFontMono-Regular.ttf',
  'font-files/SauceCodeProNerdFontMono-Regular.ttf',
  'font-files/ShureTechMonoNerdFontMono-Regular.ttf',
  'font-files/SpaceMonoNerdFontMono-Regular.ttf',
  'font-files/SymbolsNerdFontMono-Regular.ttf',
  'font-files/TerminessNerdFontMono-Regular.ttf',
  'font-files/TinosNerdFont-Regular.ttf',
  'font-files/UbuntuMonoNerdFontMono-Regular.ttf',
  'font-files/UbuntuNerdFont-Regular.ttf',
  'font-files/UbuntuSansMonoNerdFontMono-Regular.ttf',
  'font-files/VictorMonoNerdFontMono-Regular.ttf',
  'font-files/ZedMonoNerdFontMono-Regular.ttf'
)

install_data(
  font_files,
  install_dir: join_paths(
    get_option('prefix'),
    'share',
    'fonts',
    'io.github.getnf.embellish'
  )
)
