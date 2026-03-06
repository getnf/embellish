import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";
import GtkSource from "gi://GtkSource?version=5";

// Force GObject type registration for GtkSource
const _viewType = GtkSource.View;

export const EmbImportDialog = GObject.registerClass(
    {
        GTypeName: "EmbImportDialog",
        Template: "resource:///io/github/getnf/embellish/ui/ImportDialog.ui",
        InternalChildren: ["source_view", "import_button", "toast_overlay"],
        Signals: {
            "imported": {
                flags: GObject.SignalFlags.RUN_LAST,
                param_types: []
            }
        }
    },
    class extends Adw.Dialog {
        constructor(params = {}) {
            const parent = params.parent;
            const customFontsManager = params.customFontsManager;
            delete params.parent;
            delete params.customFontsManager;

            super(params);
            this._parent = parent;
            this._customFontsManager = customFontsManager;


            if (!this._source_view) {
                this._source_view = this.get_template_child(EmbImportDialog, "source_view");
            }
            if (!this._import_button) {
                this._import_button = this.get_template_child(EmbImportDialog, "import_button");
            }
            if (!this._toast_overlay) {
                this._toast_overlay = this.get_template_child(EmbImportDialog, "toast_overlay");
            }

            if (!this._source_view) {
                console.error("EmbImportDialog: source_view still not found");
                return;
            }

            const langManager = GtkSource.LanguageManager.get_default();
            const jsonLang = langManager.get_language("json");
            if (jsonLang) {
                this._source_view.get_buffer().set_language(jsonLang);
            }

            this._source_view.get_buffer().connect("changed", () => {
                const text = this._source_view.get_buffer().text.trim();
                if (this._import_button) {
                    this._import_button.set_sensitive(text.length > 0);
                }
            });

            if (this._import_button) {
                this._import_button.set_sensitive(false);
            }
        }

        onImportClicked() {
            const buffer = this._source_view.get_buffer();
            const json = buffer.text;
            try {
                this._customFontsManager.import(json);
                this.emit("imported");
                this.close();
            } catch (error) {
                this._handleError(error, _("Failed to import custom fonts: %s"));
            }
        }

        async onImportFromFile() {
            const dialog = new Gtk.FileDialog({
                title: _("Import Custom Fonts"),
                accept_label: _("Import"),
            });

            const filter = new Gtk.FileFilter();
            filter.set_name(_("JSON Files"));
            filter.add_suffix("json");
            const filters = new Gio.ListStore({ item_type: Gtk.FileFilter });
            filters.append(filter);
            dialog.set_filters(filters);

            try {
                const file = await dialog.open(this._parent, null);
                if (file) {
                    let [contents] = await file.load_contents_async(null);
                    if (contents) {
                        if (contents.toArray) {
                            contents = contents.toArray();
                        }
                        const decoder = new TextDecoder();
                        const json = decoder.decode(contents);
                        this._source_view.get_buffer().set_text(json, -1);
                    }
                }
            } catch (error) {
                if (!error.matches(Gtk.DialogError, Gtk.DialogError.CANCELLED) &&
                    !error.matches(Gtk.DialogError, Gtk.DialogError.DISMISSED)) {
                    this._handleError(error, _("Failed to read file: %s"));
                }
            }
        }

        _handleError(error, message) {
            const toast = new Adw.Toast({
                title: message.format(error),
            });
            if (this._toast_overlay) {
                this._toast_overlay.add_toast(toast);
            } else {
                console.error(message.format(error));
            }
        }
    }
);
