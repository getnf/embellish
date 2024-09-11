import Gtk from "gi://Gtk";
import Adw from "gi://Adw";

export class PreviewManager {
    constructor(parent) {
        this._parent = parent;
    }

    new(font) {
        const button = new Gtk.Button({
            icon_name: "embellish-preview-symbolic",
        });
        button.add_css_class("flat");
        button.set_tooltip_text("Preview");
        button.connect("clicked", () => {
            this._showDialog(font.tarName);
        });
        return button;
    }

    _showDialog(fileName) {
        const dialog = new Adw.Dialog({
            title: fileName,
            content_width: 360,
            content_height: -1,
        });
        const page = new Adw.ToolbarView();
        page.set_extend_content_to_top_edge(true);
        const headerBar = new Adw.HeaderBar();
        headerBar.set_show_title(false);
        page.add_top_bar(headerBar);

        const previewPicture = Gtk.Picture.new_for_resource(
            `/io/github/getnf/embellish/previews/${fileName}.svg`,
        );
        previewPicture.set_can_shrink(false);
        previewPicture.set_content_fit(Gtk.ContentFit.CONTAIN);
        previewPicture.set_valign(Gtk.Align.CENTER);
        previewPicture.set_halign(Gtk.Align.CENTER);
        previewPicture.set_margin_start(12);
        previewPicture.set_margin_end(12);
        previewPicture.add_css_class("svg-preview");

        page.set_content(previewPicture);
        dialog.set_child(page);

        dialog.present(this._parent);
    }
}
