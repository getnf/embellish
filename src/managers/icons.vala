/* icons.vala
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

public class Embellish.Icons : Object {
	private Gee.Map<string, Icon> icons;

	public Icons () {
		icons = new Gee.HashMap<string, Icon> ();
		try {
			load_icons_from_file ();
		} catch (Error e) {
			warning ("%s", e.message);
		}
	}

	private void load_icons_from_file () throws Error {
    try {
        var icons_file = File.new_for_uri ("resource:///io/github/getnf/embellish/icons.csv");
        uint8[] icons_data;
        icons_file.load_contents (null, out icons_data, null);

        var lines = ((string) icons_data).split ("\n");
        foreach (var line in lines) {
            var icon = Icon.from_csv_line (line);
            if (icon != null) {
                icons.set (icon.name, icon);
            }
        }
	} catch (Error e) {
        throw new FileError.NOENT ("Failed to load embedded resources: " + e.message);
	}
}

    public Icon? icon (string name) {
		return icons.get (name);
	}

	public Gee.Collection<Icon> collection () {
		return icons.values;
	}
}