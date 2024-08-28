# Configure the entry point
configure_file(
	input : APP_ID + '.js',
	output : APP_ID,
	configuration: {
		'GJS': find_program('gjs').full_path(),
		'PACKAGE_NAME': APP_ID,
		'PACKAGE_VERSION': meson.project_version(),
		'PREFIX': get_option('prefix'),
		'LIBDIR': get_option('prefix') / get_option('libdir')
	},
	install_dir: get_option('bindir'),
	install_mode: 'rwxr-xr-x'
)

# Compile the resources
app_resource = gnome.compile_resources(
	APP_ID + '.src',
	APP_ID + '.src.gresource.xml',
	gresource_bundle: true,
	install: true,
	install_dir : get_option('datadir') / APP_ID
)
