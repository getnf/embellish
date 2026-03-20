/* custom_fonts.vala
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

public class Embellish.CustomFonts : GLib.Object {

    private const string CUSTOM_FONTS_FILE = "custom-fonts.json";
    private const string OLD_CUSTOM_FONTS_FILE = "custom-fonts";

    private string custom_fonts_file_path;
    private string old_custom_fonts_file_path;
    private Gee.ArrayList<Font> custom_fonts_list;

    public CustomFonts () {
        var config_path = Path.build_filename (Environment.get_user_config_dir (), "embellish");
        custom_fonts_file_path = Path.build_filename (config_path, CUSTOM_FONTS_FILE);
        old_custom_fonts_file_path = Path.build_filename (config_path, OLD_CUSTOM_FONTS_FILE);

        try {
            var config_dir = File.new_for_path (config_path);
            if (!config_dir.query_exists ()) {
                config_dir.make_directory_with_parents ();
            }
        } catch (Error e) {
            warning ("CustomFonts: Failed to create directories: %s", e.message);
        }

        custom_fonts_list = new Gee.ArrayList<Font> ();
        load_from_file ();
        migrate_old_system ();
    }

    private void load_from_file () {
        if (!FileUtils.test (custom_fonts_file_path, FileTest.EXISTS)) {
            return;
        }

        try {
            string json;
            FileUtils.get_contents (custom_fonts_file_path, out json);

            var parser = new Json.Parser ();
            parser.load_from_data (json);

            var root = parser.get_root ();
            if (root != null && root.get_node_type () == Json.NodeType.ARRAY) {
                var array = root.get_array ();
                for (uint i = 0; i < array.get_length (); i++) {
                    var font = Font.from_custom_json (array.get_element (i).get_object ());
                    if (font != null) {
                        custom_fonts_list.add (font);
                    }
                }
            }
        } catch (Error e) {
            warning ("CustomFonts: Failed to load from file: %s", e.message);
        }
    }

    private void migrate_old_system () {
        if (!FileUtils.test (old_custom_fonts_file_path, FileTest.EXISTS)) {
            return;
        }

        bool changed = false;
        try {
            var keyfile = new KeyFile ();
            keyfile.load_from_file (old_custom_fonts_file_path, KeyFileFlags.NONE);

            string[] groups = keyfile.get_groups ();
            foreach (string group in groups) {
                try {
                    string id = group;
                    string display_name = keyfile.get_string (group, "display_name");
                    string description = keyfile.get_string (group, "description");
                    string url = keyfile.get_string (group, "url");
                    string archive_name = keyfile.get_string (group, "dirName");

                    var font = new Font (
                        id,
                        display_name,
                        archive_name,
                        description,
                        null,
                        null,
                        null,
                        url,
                        true,
                        null
                    );

                    // Check for duplicates
                    bool exists = false;
                    foreach (var existing in custom_fonts_list) {
                        if (existing.archive_name == font.archive_name) {
                            exists = true;
                            break;
                        }
                    }

                    if (!exists) {
                        custom_fonts_list.add (font);
                        changed = true;
                    }
                } catch (Error e) {
                    warning ("CustomFonts: Failed to migrate font from group %s: %s", group, e.message);
                }
            }

            if (changed) {
                save_to_file ();
            }
            
            // Delete the old file after successful migration
            var old_file = File.new_for_path (old_custom_fonts_file_path);
            old_file.delete ();
        } catch (Error e) {
            warning ("CustomFonts: Failed to migrate old system: %s", e.message);
        }
    }

    private void save_to_file () throws Error {
        var builder = new Json.Builder ();
        builder.begin_array ();
        foreach (var font in custom_fonts_list) {
            var node = new Json.Node (Json.NodeType.OBJECT);
            node.set_object (font.to_custom_json ());
            builder.add_value (node);
        }
        builder.end_array ();

        var generator = new Json.Generator ();
        generator.pretty = true;
        generator.set_root (builder.get_root ());

        try {
            FileUtils.set_contents (custom_fonts_file_path, generator.to_data (null));
        } catch (Error e) {
            warning ("CustomFonts: Failed to save to file: %s", e.message);
            throw e;
        }
    }

    public Gee.List<Font> list () {
        return custom_fonts_list;
    }

    public void add (Font font) throws Error {
        // Avoid duplicates by archive_name
        foreach (var existing in custom_fonts_list) {
            if (existing.archive_name == font.archive_name) {
                return;
            }
        }
        custom_fonts_list.add (font);
        save_to_file ();
    }

    public void remove (Font font) throws Error {
        for (int i = 0; i < custom_fonts_list.size; i++) {
            if (custom_fonts_list.get (i).archive_name == font.archive_name) {
                custom_fonts_list.remove_at (i);
                save_to_file ();
                return;
            }
        }
    }

    public string export () throws Error {
        if (!FileUtils.test (custom_fonts_file_path, FileTest.EXISTS)) {
            return "[]";
        }
        string json;
        FileUtils.get_contents (custom_fonts_file_path, out json);
        return json;
    }

    public Gee.List<Font> import (string json_data) throws Error {
        var imported = new Gee.ArrayList<Font> ();
        var parser = new Json.Parser ();
        parser.load_from_data (json_data);

        var root = parser.get_root ();
        if (root == null || root.get_node_type () != Json.NodeType.ARRAY) {
            throw new IOError.INVALID_DATA ("Invalid import data: expected an array");
        }

        var array = root.get_array ();
        for (uint i = 0; i < array.get_length (); i++) {
            var obj = array.get_element (i).get_object ();
            if (obj == null) continue;

            var font = Font.from_custom_json (obj);
            if (font != null) {
                bool exists = false;
                foreach (var existing in custom_fonts_list) {
                    if (existing.archive_name == font.archive_name) {
                        exists = true;
                        break;
                    }
                }

                if (!exists) {
                    custom_fonts_list.add (font);
                    imported.add (font);
                }
            }
        }

        if (imported.size > 0) {
            save_to_file ();
        }

        return imported;
    }
}
