/* icon.vala
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

public class Embellish.Icon : GLib.Object {
	public string name {get; set;}
	public string unicode {get; set;}
	public string glyph {get; set;}
	public string category {get; set;}

	public Icon (
				 string name,
				 string unicode,
				 string glyph,
				 string category
		) {
		this.name = name;
		this.unicode = unicode;
		this.glyph = glyph;
		this.category = category;
	}

public static Icon? from_csv_line (string line) {
    string[] parts = line.strip ().split (",");

    if (parts.length < 2 || parts[0] == "") {
        return null;
    }

    string name    = parts[0];
    string unicode = parts[1];
    string glyph   = parts[2];

    // extract "nf-cod" from "nf-cod-book"
    int second_dash = name.index_of ("-", name.index_of ("-") + 1);
    string prefix   = second_dash >= 0 ? name.substring (0, second_dash) : name;
    string category = category_from_prefix (prefix);

    return new Icon (name, unicode, glyph, category);
}

private static string category_from_prefix (string prefix) {
    switch (prefix) {
        case "nf-fa":          return "Font Awesome";
        case "nf-md":          return "Material Design";
        case "nf-cod":         return "Codicons";
        case "nf-dev":         return "Devicons";
        case "nf-oct":         return "Octicons";
        case "nf-fae":         return "Font Awesome Extension";
        case "nf-linux":       return "Linux";
        case "nf-weather":     return "Weather";
        case "nf-seti":        return "Seti UI";
        case "nf-custom":      return "Custom";
        case "nf-pl":          return "Powerline";
        case "nf-ple":         return "Powerline Extra";
        case "nf-pom":         return "Pomicons";
        case "nf-iec":         return "IEC Power";
        case "nf-extra":       return "Extra";
        case "nf-indent":      return "Indentation";
        case "nf-indentation": return "Indentation";
        default:               return "Unknown";
    }
}
}