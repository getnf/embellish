import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";
import GLib from "gi://GLib";

import { FontsManager } from "./fontsManager.js";
import { InstalledFontsManager } from "./installedFontsManager.js";
import { LicencesManager } from "./licencesManager.js";
import { PreviewManager } from "./previewManager.js";
import { Utils } from "./utils.js";
import { VersionManager } from "./versionManager.js";

export const EmbWindow = GObject.registerClass(
    {
        GTypeName: "EmbWindow",
        Template: "resource:///io/github/getnf/embellish/ui/Window.ui",
        InternalChildren: [
            "stack",
            "mainStack",
            "mainPage",
            "searchBar",
            "searchEntry",
            "searchPage",
            "statusPage",
            "searchList",
            "toastOverlay",
            "installedFontsList",
            "availableFontsList",
        ],
    },
    class extends Adw.ApplicationWindow {
        constructor(params = {}) {
            super(params);
            this.#setupActions();
            this.#setupWelcomeScreen();
            this.utils = new Utils();
            this.installedFontsManager = new InstalledFontsManager();
            this.versionManager = new VersionManager();
            this.licencesManager = new LicencesManager();
            this.previewManager = new PreviewManager(this);
            this.fontsManager = new FontsManager(
                this.installedFontsManager,
                this.versionManager,
            );

            const installedListDefaultWidget = new Adw.ActionRow({
                title: _("No Installed Fonts yet"),
            });

            const availableListDefaultWidget = new Adw.ActionRow({
                title: _("No available Fonts yet"),
            });

            this._installedFontsList.set_placeholder(
                installedListDefaultWidget,
            );
            this._availableFontsList.set_placeholder(
                availableListDefaultWidget,
            );

            this.#initialize();
        }

        async #initialize() {
            try {
                await this.versionManager.setupFontsVersion();
            } catch (error) {
                console.error("Error during version setup: ", error);
                const toast = new Adw.Toast({
                    title: _("Failed to set up font version."),
                });
                this._toastOverlay.add_toast(toast);
            }

            try {
                this.fontsManager.loadFontDirectories();
                this.fonts = this.fontsManager.loadFonts();
                this.#setupSearch();
                this.#populateFontLists();
            } catch (error) {
                console.error("Error during font initialization: ", error);
                const toast = new Adw.Toast({
                    title: _("Failed to load fonts."),
                });
                this._toastOverlay.add_toast(toast);
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

            const searchAction = new Gio.SimpleAction({ name: "search" });
            searchAction.connect("activate", () => {
                this._searchBar.search_mode_enabled =
                    !this._searchBar.search_mode_enabled;
            });
            this.add_action(searchAction);
        }

        #setupWelcomeScreen() {
            if (globalThis.settings.get_boolean("welcome-screen-shown")) {
                this._stack.set_visible_child_name("mainPage");
            } else {
                this._stack.set_visible_child_name("welcomePage");
                globalThis.settings.set_boolean("welcome-screen-shown", true);
            }
        }

        #populateFontLists() {
            this._installedFontsList.remove_all();
            this._availableFontsList.remove_all();

            this.fonts.forEach((font) => {
                const row = this._createFontRow(font);

                if (font.installed) {
                    this._installedFontsList.append(row);
                } else {
                    this._availableFontsList.append(row);
                }
            });
        }

        #setupSearch() {
            this._searchBar.connect("notify::search-mode-enabled", () => {
                if (this._searchBar.search_mode_enabled) {
                    this._mainStack.visible_child = this._searchPage;
                } else {
                    this._mainStack.visible_child = this._mainPage;
                }
            });

            this.fonts.forEach((font) => {
                const row = this._createFontRow(font);
                this._searchList.append(row);
            });

            let results_count;

            const filter = (row) => {
                const re = new RegExp(this._searchEntry.text, "i");
                const match = re.test(row.title);
                if (match) results_count++;
                return match;
            };

            this._searchList.set_filter_func((row) => filter(row));

            this._searchEntry.connect("search-changed", () => {
                results_count = -1;
                this._searchList.invalidate_filter();
                if (results_count === -1)
                    this._mainStack.visible_child = this._statusPage;
                else if (this._searchBar.search_mode_enabled)
                    this._mainStack.visible_child = this._searchPage;
            });
        }

        _createFontRow(font) {
            let title = font.name;
            if (font.patchedName !== "") {
                title = `${font.name} (${font.patchedName})`;
            }

            const row = new Adw.ActionRow({
                title: title,
                subtitle: this.utils.escapeMarkup(font.description),
            });

            const suffix = this._createFontRowSuffix(font);
            row.add_suffix(suffix);

            return row;
        }

        _createFontRowSuffix(font) {
            const box = this.utils.createBox(Gtk.Orientation.HORIZONTAL, 12);

            const licences = this.licencesManager.new(font);
            box.append(licences);

            const previewButton = this.previewManager.new(font);
            box.append(previewButton);

            const installButton = this._createInstallButton(font);
            const updateButton = this._createUpdateButton(font);
            const removeButton = this._createRemoveButton(font);

            if (font.installed && font.hasUpdate) {
                box.append(removeButton);
                box.append(updateButton);
            }
            if (font.installed && !font.hasUpdate) {
                box.append(removeButton);
            } else if (!font.installed) {
                box.append(installButton);
            }

            return box;
        }

        _createInstallButton(font) {
            const { button, spinner, stack } =
                this.utils.createSpinnerButton(
                    "embellish-download-symbolic",
                    "Install",
                );
            button.connect("clicked", async () => {
                try {
                    await this._handleInstallButton(
                        font,
                        spinner,
                        stack,
                        _("Font Installed"),
                    );
                } catch (error) {
                    this._handleError(error, _("Installation failed: %s"));
                }
            });
            return button;
        }

        _createUpdateButton(font) {
            const { button, spinner, stack } =
                this.utils.createSpinnerButton(
                    "embellish-update-symbolic",
                    "Update",
                );
            button.connect("clicked", async () => {
                try {
                    await this._handleInstallButton(
                        font,
                        spinner,
                        stack,
                        _("Font updated"),
                    );
                } catch (error) {
                    this._handleError(error, _("Updating failed: %s"));
                }
            });
            return button;
        }

        _createRemoveButton(font) {
            const { button, spinner, stack } =
                this.utils.createSpinnerButton(
                    "embellish-remove-symbolic",
                    "Remove",
                );
            button.connect("clicked", async () => {
                try {
                    await this._handleRemoveButton(
                        font,
                        spinner,
                        stack,
                    );
                } catch (error) {
                    this._handleError(error, _("Removing failed: %s"));
                }
            });
            return button;
        }

        _handleError(error, message) {
            const toast = new Adw.Toast({
                title: message.format(error),
            });
            this._toastOverlay.add_toast(toast);
            console.log(error);
        }

        async _handleFontAction(action, font, spinner, stack, message) {
            try {
                stack.set_visible_child_name("spinner");
                spinner.spinning = true;
                await action(font);
                const toast = new Adw.Toast({
                    title: message,
                });
                this._toastOverlay.add_toast(toast);
                this.fontsManager.loadFontDirectories();
                this.fonts = this.fontsManager.loadFonts();
                this._searchList.remove_all();
                this.#setupSearch();
                this.#populateFontLists();
            } catch (error) {
                spinner.spinning = false;
                stack.set_visible_child_name("icon")
                console.log("Action failed", error);
                throw error;
            }
        }

        async _handleInstallButton(font, spinner, stack, message) {
            await this._handleFontAction(
                async (font) => {
                    await this.fontsManager.downloadAndInstall(
                        font.tarName,
                        this.versionManager.get(),
                    );
                    this.installedFontsManager.update(
                        font.tarName,
                        this.versionManager.get(),
                    );
                },
                font,
                spinner,
                stack,
                message,
            );
        }

        async _handleRemoveButton(font, spinner, stack) {
            await this._handleFontAction(
                async (font) => {
                    await this.fontsManager.remove(font.tarName);
                    this.installedFontsManager.remove(font.tarName);
                },
                font,
                spinner,
                stack,
                _("Font removed"),
            );
        }

        vfunc_close_request() {
            super.vfunc_close_request();
            this.run_dispose();
        }
    },
);
