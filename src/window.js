import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";
import GLib from "gi://GLib";

import { FontsManager } from "./fontsManager.js";
import { InstalledFontsManager } from "./installedFontsManager.js";
import { LicencesManager } from "./licencesManager.js";
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
            this.installedFontsManager = new InstalledFontsManager();
            this.versionManager = new VersionManager();
            this.licencesManager = new LicencesManager();
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
                subtitle: this._escapeMarkup(font.description),
            });

            const suffix = this._createRowSuffix(font);
            row.add_suffix(suffix);

            return row;
        }

        _createRowSuffix(font) {
            const box = this._createBox("horizontal", 12);

            const licences = this.licencesManager.new(font);
            box.append(licences);

            const previewButton = this._createPreviewButton(font);
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

        _createBox(orientation, spacing) {
            const box = new Gtk.Box({
                orientation,
                spacing,
            });
            box.set_halign(3);
            box.set_valign(3);
            return box;
        }

        _createPreviewButton(font) {
            const button = new Gtk.Button({
                icon_name: "embellish-preview-symbolic",
            });
            button.add_css_class("flat");
            button.set_tooltip_text("Preview");
            button.connect("clicked", () => {
                this._showPreviewDialog(font.tarName);
            });
            return button;
        }

        _createButton(icon, tooltip) {
            const button = new Gtk.Button();
            button.add_css_class("flat");
            button.set_tooltip_text(tooltip);
            const buttonBox = new Gtk.Box({
                orientation: Gtk.Orientation.HORIZONTAL,
            });
            const buttonIcon = Gtk.Image.new_from_resource(
                `/io/github/getnf/embellish/icons/scalable/actions/${icon}.svg`,
            );
            const buttonSpinner = new Gtk.Spinner();
            buttonSpinner.set_visible(false);
            buttonBox.append(buttonIcon);
            buttonBox.append(buttonSpinner);
            button.set_child(buttonBox);

            return { button, buttonIcon, buttonSpinner };
        }

        _createInstallButton(font) {
            const { button, buttonIcon, buttonSpinner } = this._createButton(
                "embellish-download-symbolic",
                "Install",
            );
            button.connect("clicked", async () => {
                try {
                    await this._handleInstallButton(
                        font,
                        buttonSpinner,
                        buttonIcon,
                        _("Font Installed"),
                    );
                } catch (error) {
                    this._handleError(error, _("Installation failed: %s"));
                }
            });
            return button;
        }

        _createUpdateButton(font) {
            const { button, buttonIcon, buttonSpinner } = this._createButton(
                "embellish-update-symbolic",
                "Update",
            );
            button.connect("clicked", async () => {
                try {
                    await this._handleInstallButton(
                        font,
                        buttonSpinner,
                        buttonIcon,
                        _("Font updated"),
                    );
                } catch (error) {
                    this._handleError(error, _("Updating failed: %s"));
                }
            });
            return button;
        }

        _createRemoveButton(font) {
            const { button, buttonIcon, buttonSpinner } = this._createButton(
                "embellish-remove-symbolic",
                "Remove",
            );
            button.connect("clicked", async () => {
                try {
                    await this._handleRemoveButton(
                        font,
                        buttonSpinner,
                        buttonIcon,
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

        async _handleFontAction(action, font, spinner, icon, message) {
            try {
                icon.set_visible(false);
                spinner.set_visible(true);
                spinner.spinning = true;

                // Execute the action (install or remove)
                await action(font);

                spinner.spinning = false;
                spinner.set_visible(false);
                icon.set_visible(true);

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
                spinner.set_visible(false);
                icon.set_visible(true);
                console.log("Action failed", error);
                throw error;
            }
        }

        async _handleInstallButton(font, spinner, icon, message) {
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
                icon,
                message,
            );
        }

        async _handleRemoveButton(font, spinner, icon) {
            await this._handleFontAction(
                async (font) => {
                    await this.fontsManager.remove(font.tarName);
                    this.installedFontsManager.remove(font.tarName);
                },
                font,
                spinner,
                icon,
                _("Font removed"),
            );
        }

        _showPreviewDialog(fileName) {
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
            previewPicture.set_margin_start(12);
            previewPicture.set_margin_end(12);
            previewPicture.add_css_class("svg-preview");

            page.set_content(previewPicture);
            dialog.set_child(page);

            dialog.present(this);
        }

        _escapeMarkup(text) {
            return text
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;");
        }

        vfunc_close_request() {
            super.vfunc_close_request();
            this.run_dispose();
        }
    },
);
