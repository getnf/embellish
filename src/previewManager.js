import Gtk from "gi://Gtk";
import Adw from "gi://Adw";
import GtkSource from "gi://GtkSource?version=5";
import Pango from "gi://Pango";

export class PreviewManager {
    constructor(parent) {
        this._parent = parent;
    }

    new(font) {
        const button = new Gtk.Button({
            icon_name: "embellish-preview-symbolic",
        });
        button.add_css_class("flat");
        button.set_tooltip_text(_("Preview"));
        button.connect("clicked", () => {
            this._showDialog(font.tarName);
        });
        return button;
    }

    _showDialog(fileName) {
        const dialog = new Adw.Dialog({
            title: fileName,
            content_width: 500,
            content_height: 400,
        });
        const page = new Adw.ToolbarView();
        
        const headerBar = new Adw.HeaderBar();
        headerBar.set_show_title(false);
        page.add_top_bar(headerBar);

        const sourceView = new GtkSource.View();
        
        // Add some good sample code for Nerd Fonts
        const sampleCode = `#!/bin/bash
#  Git
#  Branch
#  Folder
#  Terminal

function hello_nerdfont() {
    echo "This is    "
}

hello_nerdfont
`;

        sourceView.get_buffer().set_text(sampleCode, -1);
        
        // Add padding
        sourceView.set_margin_start(12);
        sourceView.set_margin_end(12);
        sourceView.set_margin_top(12);
        sourceView.set_margin_bottom(12);
        
        // Settings
        sourceView.set_editable(false);
        sourceView.set_monospace(true);
        sourceView.set_show_line_numbers(true);
        
        // Highlight
        const langManager = GtkSource.LanguageManager.get_default();
        const bashLang = langManager.get_language("sh");
        if (bashLang) {
            sourceView.get_buffer().set_language(bashLang);
        }

        page.set_content(sourceView);
        dialog.set_child(page);

        dialog.present(this._parent);
    }
}
