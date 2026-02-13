import Adw from "gi://Adw";
import GObject from "gi://GObject";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";

import { EmbWindow } from "./window.js";
import { EmbIconsPage } from "./iconsPage.js";

export const EmbApplication = GObject.registerClass(
    {
        GTypeName: "EmbApplication",
    },
    class extends Adw.Application {
        vfunc_startup() {
            super.vfunc_startup();
            this.#loadSettings();
            this.#setupActions();
            this.#setupAccelerators();
        }

        vfunc_activate() {
            const window = new EmbWindow({ application: this });
            window.present();
        }

        #loadSettings() {
            globalThis.settings = new Gio.Settings({
                schemaId: this.applicationId,
            });
        }

        #setupActions() {
            const quitAction = new Gio.SimpleAction({ name: "quit" });
            quitAction.connect("activate", () => this.quit());
            this.add_action(quitAction);

            const aboutAction = new Gio.SimpleAction({ name: "about" });
            aboutAction.connect("activate", () => this._openAboutDialog());
            this.add_action(aboutAction);
        }

        #setupAccelerators() {
            this.set_accels_for_action("app.quit", ["<Control>q"]);
            this.set_accels_for_action("window.close", ["<Control>w"]);
            this.set_accels_for_action("win.search", ["<Control>f"]);
        }

        _openAboutDialog() {
            const dialog = new Adw.AboutDialog({
                application_icon: "io.github.getnf.embellish",
                application_name: _("Embellish"),
                developer_name: "Ronnie Nissan",
                version: pkg.version,
                comments: _(
                    "Embellish helps you manage Nerd Fonts on your system",
                ),
                website: "https://github.com/getnf/embellish",
                issue_url: "https://github.com/getnf/embellish/issues/new",
                copyright: "© 2024 Ronnie Nissan",
                license_type: Gtk.License.GPL_3_0,
                developers: ["Ronnie Nissan <ronnie.nissan@proton.me>"],
                designers: ["Brage Fuglseth https://bragefuglseth.dev"],
                artists: [
                    "Jakub Steiner https://jimmac.eu/",
                    "Brage Fuglseth https://bragefuglseth.dev",
                ],
                // Translators: Replace "translator-credits" with your names, one name per line
                translator_credits: _("translator-credits"),
            });

            dialog.add_other_app("io.github.ronniedroid.concessio", _("Concessio"), _("Understand file permissions"))
            dialog.add_other_app("io.github.sitraorg.sitra", _("Sitra"), _("Get fonts from online sources"))

            dialog.present(this.get_active_window());
        }
    },
);
