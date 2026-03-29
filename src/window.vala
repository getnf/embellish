/* window.vala
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

[GtkTemplate (ui = "/io/github/getnf/embellish/window.ui")]
public class Embellish.Window : Adw.ApplicationWindow {
    [GtkChild] private unowned Gtk.Stack stack;
	[GtkChild] private unowned Adw.ViewStack main_stack;
	[GtkChild] private unowned Gtk.ListBox installed_fonts;
	[GtkChild] private unowned Gtk.ListBox available_fonts;
	[GtkChild] private unowned Adw.ToastOverlay toast_overlay;
	[GtkChild] private unowned Gtk.SearchBar search_bar;
	[GtkChild] private unowned Gtk.SearchEntry search_entry;
	[GtkChild] private unowned Gtk.GridView icons_grid;
	[GtkChild] private unowned Adw.WrapBox categories;

	private Embellish.Fonts fonts_manager;
    private Embellish.Library library;
	private Embellish.CustomFonts custom_fonts;
	private Embellish.Icons icons_manager;
	private ListStore store;
	private ListStore icons_store;
	private Gtk.FilterListModel filtered_model;
    private Gtk.CustomFilter filter;
	private Gtk.CustomFilter icons_filter;
	private Gee.HashMap<string, Gtk.ToggleButton> category_toggles;

    public Window (Adw.Application app) {
        Object (application: app);
        fonts_manager = new Embellish.Fonts ();
        library = new Embellish.Library ();
        custom_fonts = new Embellish.CustomFonts ();
        icons_manager = new Embellish.Icons ();
        setup_actions ();
        setup_fonts ();
        setup_icons ();
		setup_categories ();

		search_entry.search_changed.connect (() => {
            this.filter.changed (Gtk.FilterChange.DIFFERENT);
			this.icons_filter.changed (Gtk.FilterChange.DIFFERENT);
        });

		main_stack.notify["visible-child"].connect (() => {
			var page = main_stack.visible_child_name;
			if (page == "icons_page") {
				search_entry.placeholder_text = _("Search icons");
			} else {
				search_entry.placeholder_text = _("Search fonts");
			}
		});
    }

    private void setup_fonts () {
        store = new ListStore (typeof (Embellish.Font));
		
        foreach (var font in fonts_manager.collection ()) {
            font.is_installed = library.is_installed (font);
            store.append (font);
        }
		
        var sorter = new Gtk.CustomSorter ((a, b) => {
            var fa = (Embellish.Font) a;
            var fb = (Embellish.Font) b;
            if (fa.is_custom != fb.is_custom)
                return fa.is_custom ? -1 : 1;
            return fa.family.collate (fb.family);
        });
		
        var sorted_model = new Gtk.SortListModel (store, sorter);
		
        this.filter = new Gtk.CustomFilter ((obj) => {
            var font = obj as Embellish.Font;
            if (font == null)
                return false;
            string q = this.search_entry.text.strip ().down ();
            if (q.length == 0)
                return true;
            if (font.display_name.down ().contains (q))
                return true;
            if (font.description != null && font.description.down ().contains (q))
                return true;
            return false;
        });
		
        filtered_model = new Gtk.FilterListModel (sorted_model, this.filter);
		
        var installed_filter = new Gtk.CustomFilter ((obj) => {
            var font = obj as Embellish.Font;
            return font.is_installed;
        });
		
        var available_filter = new Gtk.CustomFilter ((obj) => {
            var font = obj as Embellish.Font;
            return !font.is_installed;
        });
		
        var installed_model = new Gtk.FilterListModel (filtered_model, installed_filter);
        var available_model = new Gtk.FilterListModel (filtered_model, available_filter);
		
        installed_fonts.bind_model (installed_model, (obj) => {
            return create_row (obj as Embellish.Font);
        });
		
        available_fonts.bind_model (available_model, (obj) => {
            return create_row (obj as Embellish.Font);
        });
		
        installed_fonts.set_placeholder (create_placeholder (_("No installed fonts match your search")));
        available_fonts.set_placeholder (create_placeholder (_("No available fonts match your search")));
    }

private void setup_icons () {
    icons_store = new ListStore (typeof (Embellish.Icon));
    foreach (var icon in icons_manager.collection ()) {
        icons_store.append (icon);
    }

    var sorter = new Gtk.CustomSorter ((a, b) => {
        var ia = (Embellish.Icon) a;
        var ib = (Embellish.Icon) b;
        int cat = ia.category.collate (ib.category);
        if (cat != 0) return cat;
        return ia.name.collate (ib.name);
    });
    var sorted_model = new Gtk.SortListModel (icons_store, sorter);

    this.icons_filter = new Gtk.CustomFilter ((obj) => {
        var icon = obj as Embellish.Icon;
        if (icon == null) return false;

        bool category_match = true;
        if (this.category_toggles != null) {
            bool has_active = false;
            bool is_in_active = false;
            foreach (var entry in this.category_toggles.entries) {
                if (entry.value.active) {
                    has_active = true;
                    if (icon.category == entry.key) {
                        is_in_active = true;
                    }
                }
            }
            if (has_active && !is_in_active) {
                category_match = false;
            }
        }

        if (!category_match) return false;

        string q = this.search_entry.text.strip ().down ();
        if (q.length == 0) return true;
        if (icon.name.down ().contains (q)) return true;
        if (icon.category.down ().contains (q)) return true;
        if (icon.unicode.down ().contains (q)) return true;
        return false;
    });
    var filtered_model = new Gtk.FilterListModel (sorted_model, this.icons_filter);

    var factory = new Gtk.SignalListItemFactory ();
    factory.setup.connect ((obj) => {
        var item = obj as Gtk.ListItem;
        var button = new Gtk.Button ();
        button.add_css_class ("flat");
        button.add_css_class ("icon-button");
        button.set_size_request (56, 56);
        item.child = button;
    });
    factory.bind.connect ((obj) => {
        var item = obj as Gtk.ListItem;
        var icon = item.item as Embellish.Icon;
        var button = item.child as Gtk.Button;
        button.label = icon.glyph;
        button.tooltip_text = "%s\nU+%s\n%s".printf (icon.name, icon.unicode.up (), icon.category);
        button.clicked.connect (() => {
            on_icon_clicked (icon);
        });
    });
    factory.unbind.connect ((obj) => {
        var item = obj as Gtk.ListItem;
        var button = item.child as Gtk.Button;
        SignalHandler.disconnect_matched (
            button,
            SignalMatchType.ID,
            Signal.lookup ("clicked", button.get_type ()),
            0, null, null, null
        );
    });

    icons_grid.model = new Gtk.NoSelection (filtered_model);
    icons_grid.factory = factory;
}

 private void setup_categories () {
	 category_toggles = new Gee.HashMap<string, Gtk.ToggleButton> ();
	     string[] categories_list = {
        "Font Awesome", "Font Awesome Extension", "Material Design",
        "Codicons", "Devicons", "Octicons", "Linux", "Weather",
        "Seti UI", "Custom", "Powerline", "Powerline Extra",
        "Pomicons", "IEC Power", "Extra", "Indentation"
    };

	 foreach (string category in categories_list) {
		 var toggle = new Gtk.ToggleButton.with_label (category);
		 toggle.set_css_classes ({ "category", category });
		 categories.append (toggle);
		 category_toggles[category] = toggle;
		 
		 toggle.toggled.connect (() => {
				 this.icons_filter.changed (Gtk.FilterChange.DIFFERENT);
			 });
	 }
 }

private void on_icon_clicked (Embellish.Icon icon) {
    // copy glyph to clipboard, show a toast, etc.
    var clipboard = get_clipboard ();
    clipboard.set_text (icon.glyph);
    toast_overlay.add_toast (new Adw.Toast (_("Copied %s").printf (icon.name)));
}

    private Gtk.Widget create_placeholder (string message) {
        var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 12);
        box.halign = Gtk.Align.CENTER;
        box.valign = Gtk.Align.CENTER;
        box.margin_top = 24;
        box.margin_bottom = 24;

        var icon = new Gtk.Image.from_icon_name ("edit-find-symbolic");

        var label = new Gtk.Label (message);
        label.add_css_class ("dim-label");

        box.append (icon);
        box.append (label);

        return box;
    }

    private Gtk.Widget create_row (Embellish.Font font) {
				var row = new Adw.ActionRow ();

                string label;
                if (font.is_installed && font.family != null) {
                    label = "<span font_family=\"%s\">%s</span>".printf(font.family, font.display_name);
                    row.set_use_markup(true);
                    } else {
                        label = font.display_name;
                    }

                if (font.patched_name != null && font.patched_name.length > 0) {
                    label += " (%s)".printf (font.patched_name);
                }

				row.set_title (label);

				if (font.description != null)
					row.set_subtitle (Markup.escape_text(font.description));
				var suffix = create_suffix (font);
                row.add_suffix(suffix);

				return row;
    }

	private Gtk.Box create_suffix (Font font) {
            var box = new Gtk.Box(Gtk.Orientation.HORIZONTAL, 12);
            box.set_halign(Gtk.Align.CENTER);
            box.set_valign(Gtk.Align.CENTER);

            var licence_button = Embellish.Managers.LicencesManager.create (font);

            var preview_button = new Gtk.Button ();
            preview_button.set_icon_name ("embellish-preview-symbolic");
            preview_button.add_css_class ("flat");
            preview_button.set_tooltip_text (_("Preview"));
            preview_button.clicked.connect (() => {
                var dialog = new Embellish.PreviewDialog (font);
                dialog.present (this);
            });
            preview_button.set_sensitive(!font.is_custom);

            box.append(licence_button);
            box.append(preview_button);

            if (font.is_custom) {
                var remove_button = new Gtk.Button.from_icon_name ("embellish-remove-custom-symbolic");
                remove_button.add_css_class ("flat");
                remove_button.set_tooltip_text (_("Remove from list"));
                remove_button.clicked.connect (() => {
                    this.remove_custom_font (font);
                });
                box.append (remove_button);
            }

            return box;
    }

    private void remove_custom_font (Font font) {
        try {
            custom_fonts.remove (font);
            for (uint i = 0; i < store.get_n_items (); i++) {
                var f = (Font) store.get_item (i);
                if (f.id == font.id) {
                    store.remove (i);
                    break;
                }
            }
            toast_overlay.add_toast (new Adw.Toast (_("Custom font removed")));
        } catch (Error e) {
            toast_overlay.add_toast (new Adw.Toast (_("Failed to remove custom font: %s").printf (e.message)));
        }
    }

    private void import_custom_fonts () {
    var dialog = new Embellish.ImportDialog (this, custom_fonts);
    dialog.imported.connect (() => {
        int added = 0;
        foreach (var font in dialog.imported_fonts) {
            bool exists = false;
            for (uint i = 0; i < store.get_n_items (); i++) {
                var existing = (Embellish.Font) store.get_item (i);
                if (existing.id == font.id) {
                    exists = true;
                    break;
                }
            }
            if (exists) continue;
            font.is_installed = library.is_installed (font);
            store.append (font);
            added++;
        }
        if (added > 0) {
            toast_overlay.add_toast (new Adw.Toast (ngettext (
                "%d custom font imported",
                "%d custom fonts imported",
                (ulong) added).printf (added)));
        } else {
            toast_overlay.add_toast (new Adw.Toast (_("All fonts already exist")));
        }
    });
    dialog.present (this);
}

private void add_custom_font () {
    var dialog = new Embellish.AddFontDialog ();
    dialog.added.connect ((font) => {
        bool exists = false;
        for (uint i = 0; i < store.get_n_items (); i++) {
            var existing = (Embellish.Font) store.get_item (i);
            if (existing.id == font.id) {
                exists = true;
                break;
            }
        }

        if (exists) {
            toast_overlay.add_toast (new Adw.Toast (_("Font already exists")));
            return;
        }

        try {
            custom_fonts.add (font);
            font.is_installed = library.is_installed (font);
            store.append (font);
            toast_overlay.add_toast (new Adw.Toast (_("Custom font added")));
        } catch (Error e) {
            toast_overlay.add_toast (new Adw.Toast (_("Failed to add custom font: %s").printf (e.message)));
        }
    });
    dialog.present (this);
}

private void export_custom_fonts () {
    var dialog = new Gtk.FileDialog ();
    dialog.title = _("Export Custom Fonts");
    dialog.initial_name = "custom-fonts.json";

    var filter = new Gtk.FileFilter ();
    filter.set_filter_name (_("JSON Files"));
    filter.add_suffix ("json");

    var filters = new GLib.ListStore (typeof (Gtk.FileFilter));
    filters.append (filter);
    dialog.set_filters (filters);

    dialog.save.begin (this, null, (obj, res) => {
        try {
            var file = dialog.save.end (res);
            if (file == null) return;

            var json = custom_fonts.export ();
            var bytes = new Bytes (json.data);
            file.replace_contents_bytes_async.begin (
                bytes,
                null,
                false,
                FileCreateFlags.REPLACE_DESTINATION,
                null,
                (obj2, res2) => {
                    try {
                        file.replace_contents_bytes_async.end (res2, null);
                        toast_overlay.add_toast (new Adw.Toast (_("Custom fonts exported")));
                    } catch (Error e) {
                        toast_overlay.add_toast (new Adw.Toast (_("Failed to save file: %s").printf (e.message)));
                    }
                }
            );
        } catch (Error e) {
            if (!e.matches (Gtk.DialogError.quark (), Gtk.DialogError.CANCELLED) &&
                !e.matches (Gtk.DialogError.quark (), Gtk.DialogError.DISMISSED)) {
                toast_overlay.add_toast (new Adw.Toast (_("Failed to export: %s").printf (e.message)));
            }
        }
    });
}


	private void setup_actions () {
    var change_view_action =
        new SimpleAction ("change_view",
                              new GLib.VariantType ("s"));

    change_view_action.activate.connect ((action, param) => {
        if (param == null)
            return;

        var view = param.get_string();
        this.stack.visible_child_name =
            view;
    });

        var search_action = new SimpleAction("search", null);
            search_action.activate.connect(() => {
                this.search_bar.search_mode_enabled =
                    !this.search_bar.search_mode_enabled;
            });


        var import_action = new SimpleAction("import_custom_fonts", null);
            import_action.activate.connect(() => {
                this.import_custom_fonts();
            });

        var export_action = new SimpleAction("export_custom_fonts", null);
            export_action.activate.connect(() => {
                this.export_custom_fonts();
            });

        var add_font_action = new SimpleAction("add_custom_font", null);
            add_font_action.activate.connect(() => {
                this.add_custom_font();
            });

    this.add_action (change_view_action);
    this.add_action (search_action);
    this.add_action (import_action);
    this.add_action (export_action);
    this.add_action (add_font_action);
    }
}
