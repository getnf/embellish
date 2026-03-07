import Gtk from "gi://Gtk";
import Adw from "gi://Adw";
import GtkSource from "gi://GtkSource?version=5";
import Pango from "gi://Pango";

export class PreviewManager {
    constructor(parent) {
        this._parent = parent;
        // Default sample code for Nerd Fonts
        this._sampleCode = `//  Git
//  Branch
//  Folder
//  Terminal

let text = "This is text";
var symbols = "{}[]()<>;:,.";
const more = "-_+=*/%!&|^~?@#$'";

const helloNerdFont = () => {
  const icons = {
    git: "",
    branch: "",
    folder: "",
    terminal: ""
  };

 console.log(\`\${icons.terminal}\`);
};

helloNerdFont();
`;
    }

    new(font) {
        const button = new Gtk.Button({
            icon_name: "embellish-preview-symbolic",
        });
        button.add_css_class("flat");
        button.set_tooltip_text(_("Preview"));
        button.connect("clicked", () => {
            this._showDialog(font);
        });
        return button;
    }

    _showDialog(font) {
        const dialog = new Adw.Dialog({
            title: font.name,
            content_width: 600,
            content_height: 600,
        });
        const page = new Adw.ToolbarView();

        const headerBar = new Adw.HeaderBar();
        headerBar.set_show_title(true);
        page.add_top_bar(headerBar);


        const sourceView = new GtkSource.View();
        const provider = new Gtk.CssProvider();
        provider.load_from_string(`textview { font-family: "${font.family}"; font-size: 13pt; }`);
        const style_context = sourceView.get_style_context();
        style_context.add_provider(provider, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION);

        sourceView.get_buffer().set_text(this._sampleCode, -1);
        sourceView.get_buffer().connect("changed", () => {
            this._sampleCode = sourceView.get_buffer().text;
        });

        // Add padding
        sourceView.set_margin_start(12);
        sourceView.set_margin_end(12);
        sourceView.set_margin_top(12);
        sourceView.set_margin_bottom(12);

        // Settings
        sourceView.set_editable(true);
        sourceView.set_monospace(true);
        sourceView.set_show_line_numbers(true);

        // Highlight
        const langManager = GtkSource.LanguageManager.get_default();
        const lang = langManager.get_language("js");
        if (lang) {
            sourceView.get_buffer().set_language(lang);
        }

        const styleManager = Adw.StyleManager.get_default();
        const updateStyleScheme = () => {
            const schemeId = styleManager.dark ? "Adwaita-dark" : "Adwaita";
            const scheme = GtkSource.StyleSchemeManager.get_default().get_scheme(schemeId);
            if (scheme) {
                sourceView.get_buffer().set_style_scheme(scheme);
            }
        };

        const signalId = styleManager.connect("notify::dark", updateStyleScheme);
        dialog.connect("closed", () => {
            styleManager.disconnect(signalId);
        });

        updateStyleScheme();

        page.set_content(sourceView);
        dialog.set_child(page);

        dialog.present(this._parent);
    }
}
