embellish_sources += files(
  'main.vala',
  'application.vala',
  'modals/font.vala',
  'managers/fonts.vala',
  'managers/library.vala',
  'managers/custom_fonts.vala',
  'window.vala',
  'import_dialog.vala',
  'add_font_dialog.vala',
  'managers/licences.vala',
  'preview_dialog.vala'
)

embellish_deps = [
  config_dep,
  dependency('glib-2.0'),
  dependency('gio-2.0'),
  dependency('gobject-2.0'),
  dependency('gtk4'),
  dependency('libadwaita-1', version: '>= 1.4'),
  dependency('gee-0.8'),
  dependency('libsoup-3.0'),
  dependency('json-glib-1.0'),
  dependency('gtksourceview-5')
]

blueprints = custom_target('blueprints',
  input: files(
    'shortcuts-dialog.blp',
    'window.blp',
    'import-dialog.blp',
    'add-font-dialog.blp',
    'preview-dialog.blp',
    'icons-page.blp'
  ),
  output: '.',
  command: [find_program('blueprint-compiler'), 'batch-compile', '@OUTPUT@', '@CURRENT_SOURCE_DIR@', '@INPUT@'],
)

# Compile the resources
embellish_sources += gnome.compile_resources(
	'io.github.getnf.embellish',
	'io.github.getnf.embellish.src.gresource.xml',
	c_name: 'embellish',
	dependencies: blueprints,
)

executable('embellish', embellish_sources,
  dependencies: embellish_deps,
  include_directories: config_inc,
  install: true,
)
