/* font.vala
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
using Json;

public class Embellish.Font : GLib.Object {
    public string id { get; set; }
    public string display_name { get; set; }
    public string patched_name { get; set; default = ""; }
    public string archive_name { get; set; }
    public Gee.List<string> licences { get; set; default = new Gee.ArrayList<string> (); }
    public string description { get; set; }
    public string family { get; set; default = "";}
    public string url { get; set; }
    public bool is_custom { get; set; default = false; }
    public bool is_installed { get; set; default = false; }

    public Font (
        string id,
        string display_name,
        string archive_name,
        string description,
        string? family = null,
        Gee.List<string>? licences = null,
        string? patched_name = null,
        string? url = null,
        bool is_custom = false,
        bool? is_installed = false
    ) {
        this.id = id;
        this.display_name = display_name;
        this.archive_name = archive_name;
        this.description = description;
		if (family != null) {
			this.family = family;
		}
        if (patched_name != null) {
            this.patched_name = patched_name;
        }
        if (licences != null) {
            this.licences = licences;
        }
        this.url = url ?? "https://github.com/ryanoasis/nerd-fonts/releases/download/%s/%s.tar.xz".printf (Config.NF_RELEASE, archive_name);
        this.is_custom = is_custom;
        if (is_installed != null) {
            this.is_installed = is_installed;
        }
    }

    public static Font? from_json (Json.Object obj) {
        if (!obj.has_member ("id") || !obj.has_member ("display_name") || !obj.has_member ("archive_name") || !obj.has_member ("description")) {
            warning ("JSON object is missing required font members");
            return null;
        }

        string id = obj.get_string_member ("id");
        string display_name = obj.get_string_member ("display_name");
        string archive_name = obj.get_string_member ("archive_name");
        string description = obj.get_string_member ("description");
        string family = obj.has_member ("family") ? obj.get_string_member ("family") : "";

        string? patched_name = obj.has_member ("patched_name")
            ? obj.get_string_member ("patched_name")
            : null;

        var licences = new Gee.ArrayList<string> ();
        if (obj.has_member ("licences")) {
            foreach (var node in obj.get_array_member ("licences").get_elements ()) {
                licences.add (node.get_string ());
            }
        }

        return new Font (id, display_name, archive_name, description, family, licences, patched_name);
    }

    public static Font? from_custom_json (Json.Object obj) {
        if (!obj.has_member ("display_name") || !obj.has_member ("url") || !obj.has_member ("description")) {
            warning ("Custom JSON object is missing required font members");
            return null;
        }

        string display_name = obj.get_string_member ("display_name");
        string description = obj.get_string_member ("description");
        string url = obj.get_string_member ("url");

        string archive_name = obj.has_member ("archive_name")
            ? obj.get_string_member ("archive_name")
            : display_name.replace (" ", "");

        string id = obj.has_member ("id") ? obj.get_string_member ("id") : archive_name;

        return new Font (
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
    }

    public Json.Object to_custom_json () {
        var obj = new Json.Object ();
        obj.set_string_member ("id", this.id);
        obj.set_string_member ("display_name", this.display_name);
        obj.set_string_member ("description", this.description);
        obj.set_string_member ("url", this.url);
        obj.set_string_member ("archive_name", this.archive_name);
        return obj;
    }
}
