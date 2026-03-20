using Gtk;
using Adw;

[GtkTemplate (ui = "/io/github/getnf/embellish/add-font-dialog.ui")]
public class Embellish.AddFontDialog : Adw.Dialog {
    [GtkChild] private unowned Adw.EntryRow name_entry;
    [GtkChild] private unowned Adw.EntryRow description_entry;
    [GtkChild] private unowned Adw.EntryRow url_entry;
    [GtkChild] private unowned Gtk.Button add_button;

    public signal void added (Embellish.Font font);

    public AddFontDialog () {
        add_button.sensitive = false;

        name_entry.changed.connect (validate);
        description_entry.changed.connect (validate);
        url_entry.changed.connect (validate);

        add_button.clicked.connect (on_add_clicked);
    }

    private void validate () {
        bool name_ok = name_entry.text.strip ().length > 0;
        bool description_ok = description_entry.text.strip ().length > 0;
        bool url_ok = url_entry.text.strip ().length > 0 &&
                     (url_entry.text.has_suffix (".zip") || url_entry.text.has_suffix (".tar.xz"));

        add_button.sensitive = name_ok && description_ok && url_ok;
    }

    private void on_add_clicked () {
        string display_name = name_entry.text.strip ();
        string description = description_entry.text.strip ();
        string url = url_entry.text.strip ();

        string archive_name = display_name.replace (" ", "");
        string id = archive_name;

        var font = new Embellish.Font (
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

        added (font);
        this.close ();
    }
}
