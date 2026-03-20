/* fonts.vala
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

public class Embellish.Fonts : Object {
    private Gee.Map<string, Font> fonts;
	private CustomFonts custom_fonts;

    public Fonts () {
        fonts = new Gee.HashMap<string, Font> ();
		custom_fonts = new CustomFonts();

        try {
            load_fonts_from_files ();
        } catch (Error e) {
            warning (@"Failed to load fonts: $(e.message)");
        }

        merge_custom_fonts ();
    }

    private void load_fonts_from_files () throws Error {
        try {
            var fonts_file = File.new_for_uri ("resource:///io/github/getnf/embellish/fonts.json");

            uint8[] fonts_data;

            fonts_file.load_contents (null, out fonts_data, null);

			load_fonts ((string) fonts_data, fonts_data.length);
        } catch (Error e) {
            throw new FileError.NOENT ("Failed to load embedded resources: " + e.message);
        }
    }

    private void load_fonts (string json, int len) throws Error {
        var parser = new Json.Parser ();
        parser.load_from_data (json, len);

        foreach (var node in parser.get_root ().get_array ().get_elements ()) {
            var font_info = Font.from_json (node.get_object ());

			fonts.set (font_info.id, font_info);
        }
    }

	private void merge_custom_fonts () {
		foreach (var font in custom_fonts.list ()) {
			fonts.set (font.id, font);
		}
	}

    public Font? font (string id) {
        return fonts.get (id);
    }

    public Gee.Collection<Font> collection () {
        return fonts.values;
    }
}
