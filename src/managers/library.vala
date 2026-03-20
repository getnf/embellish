/* library.vala
 *
 * Copyright 2025 Ronnie Nissan Yousif
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

using Gee;

public class Embellish.Library : GLib.Object {

    private const string INSTALLED_FONTS_FILE = "fonts";

    private Soup.Session session;
    private string fonts_dir_path;
    private string installed_fonts_file_path;
    private KeyFile installed_fonts_db;

    public signal void installation_started (Font font);
    public signal void installation_completed (Font font, bool success, string? error_message);
	public signal void uninstallation_started (Font font, bool success, string? error_message);
    public signal void uninstallation_completed (Font font, bool success, string? error_message);

    public Library () {
        session = new Soup.Session ();

        fonts_dir_path = Path.build_filename (Environment.get_home_dir (), ".local", "share", "fonts");

        var config_path = Path.build_filename (Environment.get_user_config_dir (), "embellish");
        installed_fonts_file_path = Path.build_filename (config_path, INSTALLED_FONTS_FILE);

        try {
            var fonts_dir = File.new_for_path (fonts_dir_path);
            if (!fonts_dir.query_exists ()) {
                fonts_dir.make_directory_with_parents ();
            }

            var config_dir = File.new_for_path (config_path);
            if (!config_dir.query_exists ()) {
                config_dir.make_directory_with_parents ();
            }
        } catch (Error e) {
            warning ("Library: Failed to create directories: %s", e.message);
        }

        installed_fonts_db = new KeyFile ();
        try {
            installed_fonts_db.load_from_file (installed_fonts_file_path, KeyFileFlags.NONE);
        } catch (Error e) {
            warning ("Installed fonts database not found, will be created on first install");
        }
    }

    public bool is_installed (Font font) {
        return installed_fonts_db.has_group (font.archive_name);
    }

    private void track (Font font) throws Error {
        var install_path = Path.build_filename (fonts_dir_path, font.archive_name);

        installed_fonts_db.set_string (font.archive_name, "id", font.id);
        installed_fonts_db.set_string (font.archive_name, "install_path", install_path);

        try {
            installed_fonts_db.save_to_file (installed_fonts_file_path);
        } catch (Error e) {
            warning ("Failed to save installed fonts database: %s".printf (e.message));
			throw e;
        }
    }

    private void untrack (Font font) throws Error {
		try {
			installed_fonts_db.remove_group (font.archive_name);
		} catch (Error e) {
			warning ("Failed to remove font %s from installed fonts database: %s".printf (font.display_name, e.message));
			throw e;
		}
		try {
			installed_fonts_db.save_to_file(installed_fonts_file_path);
		} catch (Error e) {
			warning ("Failed to save installed fonts database: %s".printf (e.message));
			throw e;
		}
	}
}
