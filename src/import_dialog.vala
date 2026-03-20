[GtkTemplate (ui = "/io/github/getnf/embellish/import-dialog.ui")]
public class Embellish.ImportDialog : Adw.Dialog {
    [GtkChild] private unowned GtkSource.View source_view;
    [GtkChild] private unowned Gtk.Button import_button;
    [GtkChild] private unowned Gtk.Button import_file_button;
    [GtkChild] private unowned Adw.ToastOverlay toast_overlay;

    private Gtk.Widget parent_widget;
    private Embellish.CustomFonts custom_fonts;
    private Adw.StyleManager style_manager;

    public Gee.List<Embellish.Font> imported_fonts { get; private set;
        default = new Gee.ArrayList<Embellish.Font> (); }

    public signal void imported ();

    public ImportDialog (Gtk.Widget parent, Embellish.CustomFonts custom_fonts) {
        this.parent_widget = parent;
        this.custom_fonts = custom_fonts;

        var lang_manager = GtkSource.LanguageManager.get_default ();
        var json_lang = lang_manager.get_language ("json");
        var source_buffer = source_view.get_buffer () as GtkSource.Buffer;

        if (json_lang != null && source_buffer != null) {
            source_buffer.set_language (json_lang);
        }

        source_buffer.changed.connect (() => {
            import_button.set_sensitive (source_buffer.text.strip ().length > 0);
        });

        import_button.set_sensitive (false);
        import_button.clicked.connect (on_import_clicked);
        import_file_button.clicked.connect (on_import_from_file);

        style_manager = Adw.StyleManager.get_default ();
        style_manager.notify["dark"].connect (update_style_scheme);
        update_style_scheme ();
    }

    private void update_style_scheme () {
        var source_buffer = source_view.get_buffer () as GtkSource.Buffer;
        if (source_buffer == null) return;

        var scheme_id = style_manager.dark ? "Adwaita-dark" : "Adwaita";
        var scheme = GtkSource.StyleSchemeManager.get_default ().get_scheme (scheme_id);
        if (scheme != null) {
            source_buffer.set_style_scheme (scheme);
        }
    }

    private void on_import_clicked () {
        var json = source_view.get_buffer ().text;

        try {
            imported_fonts = custom_fonts.import (json);
            imported ();
            close ();
        } catch (Error e) {
            show_toast (_("Failed to import custom fonts: %s").printf (e.message));
        }
    }

    private void on_import_from_file () {
        var dialog = new Gtk.FileDialog ();
        dialog.title = _("Import Custom Fonts");
        dialog.accept_label = _("Import");

        var filter = new Gtk.FileFilter ();
        filter.set_filter_name (_("JSON Files"));
        filter.add_suffix ("json");

        var filters = new GLib.ListStore (typeof (Gtk.FileFilter));
        filters.append (filter);
        dialog.set_filters (filters);

        dialog.open.begin (parent_widget as Gtk.Window, null, (obj, res) => {
            try {
                var file = dialog.open.end (res);
                if (file == null) return;

                file.load_contents_async.begin (null, (obj2, res2) => {
                    try {
                        uint8[] contents;
                        file.load_contents_async.end (res2, out contents, null);
                        source_view.get_buffer ().set_text ((string) contents, contents.length);
                    } catch (Error e) {
                        show_toast (_("Failed to read file: %s").printf (e.message));
                    }
                });
            } catch (Error e) {
                if (!e.matches (Gtk.DialogError.quark (), Gtk.DialogError.CANCELLED) &&
                    !e.matches (Gtk.DialogError.quark (), Gtk.DialogError.DISMISSED)) {
                    show_toast (_("Failed to open file: %s").printf (e.message));
                }
            }
        });
    }

    private void show_toast (string message) {
        var toast = new Adw.Toast (message);
        toast_overlay.add_toast (toast);
    }
}
