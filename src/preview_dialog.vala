using Gtk;
using Adw;
using GtkSource;

[GtkTemplate (ui = "/io/github/getnf/embellish/preview-dialog.ui")]
public class Embellish.PreviewDialog : Adw.Dialog {
    [GtkChild] private unowned GtkSource.View source_view;
    private Gtk.CssProvider provider;

    private static string sample_code = """// NOTE: only a small subset of the font
// and the icons are available for preview
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
  console.log(`${icons.terminal}`);
};
helloNerdFont();
""";

    public PreviewDialog (Embellish.Font font) {
        this.title = font.display_name;

        var source_buffer = source_view.get_buffer () as GtkSource.Buffer;
        source_buffer.set_text (sample_code, -1);

        source_buffer.changed.connect (() => {
            sample_code = source_buffer.text;
        });

        var lang_manager = GtkSource.LanguageManager.get_default ();
        var lang = lang_manager.get_language ("js");
        if (lang != null) {
            source_buffer.set_language (lang);
        }

        provider = new Gtk.CssProvider ();
        string css = """
        textview {
            font-family: "%s";
            font-size: 20px;
            }
            """.printf (font.family);
            provider.load_from_string (css);
            source_view.get_style_context ().add_provider (provider, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION);

        var style_manager = Adw.StyleManager.get_default ();
        update_style_scheme (style_manager, source_buffer);
        style_manager.notify["dark"].connect (() => {
            update_style_scheme (style_manager, source_buffer);
        });

        this.closed.connect (() => {

        });
    }

    public string get_sample_code () {
        return source_view.get_buffer ().text;
    }

    private void update_style_scheme (Adw.StyleManager style_manager, GtkSource.Buffer source_buffer) {
        var scheme_id = style_manager.dark ? "Adwaita-dark" : "Adwaita";
        var scheme = GtkSource.StyleSchemeManager.get_default ().get_scheme (scheme_id);
        if (scheme != null) {
            source_buffer.set_style_scheme (scheme);
        }
    }
}
