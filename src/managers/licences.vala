/* licences_manager.vala
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

public class Embellish.Managers.LicencesManager : Object {

    public static Gtk.Widget create (Font font) {
        var license_box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 12);

        var licence_button = new Gtk.MenuButton ();
        licence_button.add_css_class ("licence-button");
        licence_button.set_tooltip_text (_("licence details"));

        var licence_button_label = new Gtk.Label (null);
        if (font.licences.size > 1) {
            licence_button_label.set_label (_("Dual"));
        } else if (font.licences.size == 1) {
            licence_button_label.set_label (font.licences[0]);
        } else {
            licence_button_label.set_label (_("Unknown"));
            licence_button.set_sensitive(false);
        }

        licence_button.set_always_show_arrow (false);
        licence_button.set_child (licence_button_label);
        licence_button.set_popover (create_popover (font));
        license_box.append (licence_button);

        return license_box;
    }

    private static Gtk.Popover create_popover (Font font) {
        var popover = new Gtk.Popover ();
        popover.has_arrow = true;
        popover.name = "licensePopover";

        var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 4);
        box.set_margin_top (12);
        box.set_margin_start (24);
        box.set_margin_end (24);
        box.set_margin_bottom (12);

        var licence_logo = new Gtk.Image.from_resource ("/io/github/getnf/embellish/licence.svg");
        licence_logo.set_pixel_size (64);
        box.append (licence_logo);

        foreach (var licence in font.licences) {
            var licence_box = create_licence_box (licence, get_description (licence));
            box.append (licence_box);
        }

        var clamp = new Adw.Clamp ();
        clamp.maximum_size = 250;
        clamp.set_child (box);
        popover.set_child (clamp);

        return popover;
    }

    private static Gtk.Widget create_licence_box (string id, string description) {
        var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 4);

        var id_label = new Gtk.Label (id);
        id_label.add_css_class ("heading");

        var description_label = new Gtk.Label (description);
        description_label.set_justify (Gtk.Justification.CENTER);
        description_label.set_wrap (true);

        box.append (id_label);
        box.append (description_label);

        return box;
    }

    private static string get_description (string licence_key) {
        var resource_path = "/io/github/getnf/embellish/licences";
        var key_file = new GLib.KeyFile ();

        try {
            var data = GLib.resources_lookup_data (resource_path, GLib.ResourceLookupFlags.NONE);
            key_file.load_from_bytes (data, GLib.KeyFileFlags.NONE);
        } catch (Error e) {
            warning ("Failed to load %s: %s", resource_path, e.message);
            return _("No description available");
        }

        try {
            var description = key_file.get_string ("licences", licence_key);
            return description != null ? _(description) : _("No description available");
        } catch (Error e) {
            return _("No description available");
        }
    }
}
