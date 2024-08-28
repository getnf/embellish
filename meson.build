# Define our project
project(
       'embellish',
	version: '0.4.0',
	license: ['GPL-3.0-or-later'],
	meson_version: '>= 0.59.0'
)

APP_ID = 'io.github.getnf.embellish'

# Import the modules
gnome = import('gnome')
#i18n = import('i18n')

# Load instructions from subdirectories
subdir('data')
subdir('src')
#subdir('po')

gnome.post_install(
	gtk_update_icon_cache: true,
        # update_desktop_database: true,
        glib_compile_schemas: true
)