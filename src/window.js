import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gio from "gi://Gio";
import GLib from "gi://GLib";

export const EmbWindow = GObject.registerClass(
    {
        GTypeName: "EmbWindow",
        Template: "resource:///io/github/getnf/embellish/ui/Window.ui",
        InternalChildren: ["stack"],
    },
    class extends Adw.ApplicationWindow {
        constructor(params = {}) {
            super(params);
            this.#setupActions();
            this.#setupWelcomeScreen();
        }

        #setupWelcomeScreen() {
            if (globalThis.settings.get_boolean("welcome-screen-shown")) {
                this._stack.set_visible_child_name("mainPage");
            } else {
                this._stack.set_visible_child_name("welcomePage");
                globalThis.settings.set_boolean("welcome-screen-shown", true);
            }
        }

        #setupActions() {
            const changeViewAction = new Gio.SimpleAction({
                name: "changeView",
                parameterType: GLib.VariantType.new("s"),
            });

            changeViewAction.connect("activate", (_action, params) => {
                this._stack.visibleChildName = params.unpack();
            });

            this.add_action(changeViewAction);
        }

        vfunc_close_request() {
            super.vfunc_close_request();
            this.run_dispose();
        }
    },
);
