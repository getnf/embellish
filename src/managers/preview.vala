/* preview_manager.vala
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

public class Embellish.Managers.PreviewManager : Object {
    private Gtk.Widget _parent;
    private string _sample_code;

    public PreviewManager (Gtk.Widget parent) {
        _parent = parent;
        _sample_code = """// NOTE: only a small subset of the font
// and the icons are available for preview
let text = "This is text";
var symbols = "{}[]()<>;:,.";
const more = "-_+=*/%!&|^~?@#$'";
const helloNerdFont = () => {
  const icons = {
    git: "",
    branch: "",
    folder: "",
    terminal: ""
  };
  console.log(`${icons.terminal}`);
};
helloNerdFont();
""";
    }

    public Gtk.Widget create (Font font) {
        var button = new Gtk.Button ();
        button.set_icon_name ("embellish-preview-symbolic");
        button.add_css_class ("flat");
        button.set_tooltip_text (_("Preview"));
        button.clicked.connect (() => {
            show_dialog (font);
        });
        button.set_sensitive(!font.is_custom);
        return button;
    }

private void show_dialog (Font font) {
    var dialog = new Adw.Dialog ();
    dialog.title = font.display_name;
    dialog.content_width = 600;
    dialog.content_height = 600;

    var page = new Adw.ToolbarView ();
    var header_bar = new Adw.HeaderBar ();
    header_bar.set_show_title (true);
    page.add_top_bar (header_bar);

    var source_view = new GtkSource.View ();

    var provider = new Gtk.CssProvider ();
    provider.load_from_string (@"textview { font-family: \"$(font.family)\"; font-size: 13pt; }");
    Gtk.StyleContext.add_provider_for_display (
        Gdk.Display.get_default (),
        provider,
        Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
    );

    var source_buffer = source_view.get_buffer () as GtkSource.Buffer;

    source_buffer.set_text (_sample_code, -1);
    source_buffer.changed.connect (() => {
        _sample_code = source_buffer.text;
    });

    source_view.set_margin_start (12);
    source_view.set_margin_end (12);
    source_view.set_margin_top (12);
    source_view.set_margin_bottom (12);
    source_view.set_editable (true);
    source_view.set_monospace (true);
    source_view.set_show_line_numbers (true);

    var lang_manager = GtkSource.LanguageManager.get_default ();
    var lang = lang_manager.get_language ("js");
    if (lang != null) {
        source_buffer.set_language (lang);
    }

    var style_manager = Adw.StyleManager.get_default ();

    update_style_scheme (style_manager, source_buffer);
    ulong signal_id = style_manager.notify["dark"].connect (() => {
        update_style_scheme (style_manager, source_buffer);
    });

    dialog.closed.connect (() => {
        style_manager.disconnect (signal_id);
    });

    page.set_content (source_view);
    dialog.set_child (page);
    dialog.present (_parent);
}

private void update_style_scheme (Adw.StyleManager style_manager, GtkSource.Buffer source_buffer) {
    var scheme_id = style_manager.dark ? "Adwaita-dark" : "Adwaita";
    var scheme = GtkSource.StyleSchemeManager.get_default ().get_scheme (scheme_id);
    if (scheme != null) {
        source_buffer.set_style_scheme (scheme);
    }
}
}
