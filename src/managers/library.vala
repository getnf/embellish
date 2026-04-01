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
using Archive;
using GLib;

public class Embellish.Library : GLib.Object {

	private const string INSTALLED_FONTS_FILE = "fonts";

	private Soup.Session session;
	private string fonts_dir_path;
	private string installed_fonts_file_path;
	private KeyFile installed_fonts_db;

	public signal void installation_started (Font font);
	public signal void installation_completed (Font font, bool success, string? error_message);
	public signal void uninstallation_started (Font font);
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

	public async void install (Font font) {
		installation_started (font);

		File tmp_file = null;
		try {
			tmp_file = yield download (font);

			var dest = File.new_for_path (Path.build_filename (fonts_dir_path, font.archive_name));
			if (!dest.query_exists ()) {
				dest.make_directory_with_parents ();
			}

			yield extract (tmp_file, dest);

			track (font);

			installation_completed (font, true, null);
		} catch (Error e) {
			warning ("Failed to install font %s: %s", font.display_name, e.message);
			installation_completed (font, false, e.message);
		} finally {
			if (tmp_file != null) {
				try {
					tmp_file.delete ();
				} catch (Error e) {
					// Ignore cleanup errors
				}
			}
		}
	}

	public async void uninstall (Font font) {
		uninstallation_started (font);

		try {
			var font_dir = File.new_for_path (Path.build_filename (fonts_dir_path, font.archive_name));
			if (font_dir.query_exists ()) {
				yield delete_recursive (font_dir);
			}

			untrack (font);

			uninstallation_completed (font, true, null);
		} catch (Error e) {
			warning ("Failed to uninstall font %s: %s", font.display_name, e.message);
			uninstallation_completed (font, false, e.message);
		}
	}

	private async File download (Font font) throws Error {
		var message = new Soup.Message ("GET", font.url);

		var input_stream = yield session.send_async (message, Priority.DEFAULT, null);

		if (message.status_code != 200) {
			throw new IOError.FAILED ("Download failed with status %u for %s".printf (message.status_code, font.url));
		}

		// Determine file extension from URL
		string suffix = ".tar.xz";
		if (font.url.has_suffix (".zip")) {
			suffix = ".zip";
		}

		var tmp_path = Path.build_filename (Environment.get_tmp_dir (), "embellish-%s%s".printf (font.archive_name, suffix));
		var tmp_file = File.new_for_path (tmp_path);

		var output_stream = yield tmp_file.replace_async (null, false, FileCreateFlags.REPLACE_DESTINATION, Priority.DEFAULT, null);

		// Read and write in chunks
		var buffer = new uint8[8192];
		while (true) {
			var bytes_read = yield input_stream.read_async (buffer, Priority.DEFAULT, null);

			if (bytes_read == 0)break;
			yield output_stream.write_async (buffer[0 : bytes_read], Priority.DEFAULT, null);
		}

		yield output_stream.close_async (Priority.DEFAULT, null);

		return tmp_file;
	}

	private async void extract (File archive, File destination) throws Error {
		var archive_path = archive.get_path ();
		var dest_path = destination.get_path ();

		SourceFunc callback = extract.callback;
		Error? thread_error = null;

		new Thread<void*> ("extractor", () => {
			try {
				extract_sync (archive_path, dest_path);
			} catch (Error e) {
				thread_error = e;
			}
			Idle.add (() => {
				callback ();
				return false;
			});
			return null;
		});

		yield;

		if (thread_error != null) {
			throw thread_error;
		}
	}

	private static void extract_sync (string archive_path, string dest_path) throws Error {
		var reader = new Archive.Read ();
		reader.support_filter_all ();
		reader.support_format_all ();

		var writer = new Archive.WriteDisk ();
		writer.set_options (Archive.ExtractFlags.TIME | Archive.ExtractFlags.PERM | Archive.ExtractFlags.ACL | Archive.ExtractFlags.FFLAGS);
		writer.set_standard_lookup ();

		if (reader.open_filename (archive_path, 10240) != Archive.Result.OK) {
			throw new IOError.FAILED ("Could not open archive %s: %s".printf (archive_path, reader.error_string ()));
		}

		unowned Archive.Entry entry;
		while (reader.next_header (out entry) == Archive.Result.OK) {
			var entry_path = entry.pathname ();
			var full_path = Path.build_filename (dest_path, entry_path);
			entry.set_pathname (full_path);

			var res = reader.extract2 (entry, writer);
			if (res != Archive.Result.OK) {
				throw new IOError.FAILED ("Failed to extract %s: %s".printf (entry_path, reader.error_string ()));
			}
			writer.finish_entry ();
		}

		reader.close ();
	}

	private async void delete_recursive (File dir) throws Error {
		var enumerator = dir.enumerate_children ("standard::*", FileQueryInfoFlags.NOFOLLOW_SYMLINKS);
		FileInfo info;
		while ((info = enumerator.next_file ()) != null) {
			var child = dir.get_child (info.get_name ());
			if (info.get_file_type () == GLib.FileType.DIRECTORY) {
				yield delete_recursive (child);
			} else {
				child.delete ();
			}
		}
		dir.delete ();
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
			installed_fonts_db.save_to_file (installed_fonts_file_path);
		} catch (Error e) {
			warning ("Failed to save installed fonts database: %s".printf (e.message));
			throw e;
		}
	}
}