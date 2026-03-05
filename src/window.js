import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";
import GLib from "gi://GLib";

import { FontsManager } from "./fontsManager.js";
import { InstalledFontsManager } from "./installedFontsManager.js";
import { CustomFontsManager } from "./customFontsManager.js";
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
            "noFontsBanner",
            "searchBar",
            "searchEntry",
            "searchPage",
            "statusPage",
            "searchList",
            "toastOverlay",
            "scroller",
            "installedFontsList",
            "customFontsBox",
            "customFontsList",
            "availableFontsList",
        ],
    },
    class extends Adw.ApplicationWindow {
        constructor(params = {}) {
            super(params);
            this.#setupActions();
            this.#setupWelcomeScreen();
            this.utils = new Utils();
            this.scrollValue = 0;
            this.installedFontsManager = new InstalledFontsManager();
            this.versionManager = new VersionManager();
            this.licencesManager = new LicencesManager();
            this.previewManager = new PreviewManager(this);
            this.fontsManager = new FontsManager(
                this.installedFontsManager,
                this.versionManager,
            );
            this.customFontsManager = new CustomFontsManager();

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
                this.#populateCustomFontsList();
                this.#setupNoFontsBanner();
            } catch (error) {
                console.error("Error during font initialization: ", error);
                const toast = new Adw.Toast({
                    title: _("Failed to load fonts."),
                });
                this._toastOverlay.add_toast(toast);
            }
        }

        #setupNoFontsBanner() {
            this._mainStack.connect("notify::visible-child", () => {
                this._updateNoFontsBanner();
            });

            this._noFontsBanner.connect("button-clicked", () => {
                this._mainStack.set_visible_child_name("fontsPage");
            });
        }

        _updateNoFontsBanner() {
            const visible_child = this._mainStack.get_visible_child_name();
            const installedFonts = this.fonts
                ? this.fonts.filter((font) => font.installed)
                : [];

            if (visible_child === "iconsPage" && installedFonts.length === 0) {
                this._noFontsBanner.set_revealed(true);
            } else {
                this._noFontsBanner.set_revealed(false);
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

            const addCustomFontAction = new Gio.SimpleAction({ name: "addCustomFont" });
            addCustomFontAction.connect("activate", () => {
                this.#showAddCustomFontDialog();
            });
            this.add_action(addCustomFontAction);
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

            this.fonts.forEach((font) => {
                const row = this._createFontRow(font);

                if (font.installed) {
                    this._installedFontsList.append(row);
                } else {
                    this._availableFontsList.append(row);
                }
            });
        }

        #populateCustomFontsList() {
            this._customFontsList.remove_all();

            const customFonts = this.customFontsManager.getAll();

            this._customFontsBox.set_visible(customFonts.length > 0);

            customFonts.forEach((font) => {
                const row = new Adw.ActionRow({
                    title: font.name,
                    subtitle: this.utils.escapeMarkup(font.description),
                });

                const box = this.utils.createBox(Gtk.Orientation.HORIZONTAL, 12);
                const isInstalled = this.fontsManager._isInstalled(font.dirName);

                if (isInstalled) {
                    const { button, spinner, stack } = this.utils.createSpinnerButton(
                        "embellish-remove-symbolic",
                        _("Remove"),
                    );

                    button.connect("clicked", async () => {
                        try {
                            stack.set_visible_child_name("spinner");
                            spinner.spinning = true;
                            await this.fontsManager.remove(font.dirName);
                            this.customFontsManager.remove(font.name);
                            this.fontsManager.loadFontDirectories();
                            this.fonts = this.fontsManager.loadFonts();
                            this._searchList.remove_all();
                            this.#setupSearch();
                            this.#populateFontLists();
                            this.#populateCustomFontsList();
                            const toast = new Adw.Toast({
                                title: _("Custom font removed"),
                            });
                            this._toastOverlay.add_toast(toast);
                        } catch (error) {
                            spinner.spinning = false;
                            stack.set_visible_child_name("icon");
                            this._handleError(error, _("Removing failed: %s"));
                        }
                    });
                    box.append(button);
                } else {
                    const { button, spinner, stack } = this.utils.createSpinnerButton(
                        "embellish-download-symbolic",
                        _("Install"),
                    );

                    button.connect("clicked", async () => {
                        try {
                            stack.set_visible_child_name("spinner");
                            spinner.spinning = true;
                            await this.fontsManager.downloadAndInstallFromUrl(font.url, font.dirName);
                            this.fontsManager.loadFontDirectories();
                            this.fonts = this.fontsManager.loadFonts();
                            this._searchList.remove_all();
                            this.#setupSearch();
                            this.#populateFontLists();
                            this.#populateCustomFontsList();
                            const toast = new Adw.Toast({
                                title: _("Custom font installed"),
                            });
                            this._toastOverlay.add_toast(toast);
                        } catch (error) {
                            spinner.spinning = false;
                            stack.set_visible_child_name("icon");
                            this._handleError(error, _("Installation failed: %s"));
                        }
                    });
                    box.append(button);
                }

                row.add_suffix(box);
                this._customFontsList.append(row);
            });
        }

        #setupSearch() {
            this._searchBar.connect("notify::search-mode-enabled", () => {
                if (this._searchBar.search_mode_enabled) {
                    this._mainStack.set_visible_child_name("searchPage");
                } else {
                    this._mainStack.set_visible_child_name("fontsPage");
                }
            });

            this.fonts.forEach((font) => {
                const row = this._createFontRow(font);
                this._searchList.append(row);
            });

            let results_count;

            const filter = (row) => {
                const query = this._searchEntry.text.trim().toLowerCase();
                const title = row.title.toLowerCase();
                const description = (row.subtitle || "").toLowerCase();

                if (query === "") {
                    results_count++;
                    return true;
                }

                const match = title.includes(query) || description.includes(query);
                if (match) results_count++;
                return match;
            };

            this._searchList.set_filter_func((row) => filter(row));

            this._searchEntry.connect("search-changed", () => {
                const query = this._searchEntry.text.trim();
                results_count = 0;
                this._searchList.invalidate_filter();

                if (query !== "" && results_count === 0) {
                    this._mainStack.set_visible_child_name("statusPage");
                } else if (this._searchBar.search_mode_enabled) {
                    this._mainStack.set_visible_child_name("searchPage");
                } else {
                    this._mainStack.set_visible_child_name("fontsPage");
                }
            });

            // Start with the search bar closed by default
            this._searchBar.search_mode_enabled = false;
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
            const { button, spinner, stack } = this.utils.createSpinnerButton(
                "embellish-download-symbolic",
                _("Install"),
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
            const { button, spinner, stack } = this.utils.createSpinnerButton(
                "embellish-update-symbolic",
                _("Update"),
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
            const { button, spinner, stack } = this.utils.createSpinnerButton(
                "embellish-remove-symbolic",
                _("Remove"),
            );
            button.connect("clicked", async () => {
                try {
                    await this._handleRemoveButton(font, spinner, stack);
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

        _handleScrollPosition(scrollValue) {
            GLib.idle_add(GLib.PRIORITY_DEFAULT_IDLE, () => {
                const adjustment = this._scroller.get_vadjustment();
                if (adjustment) {
                    adjustment.set_value(scrollValue);
                }
                return GLib.SOURCE_REMOVE;
            });
        }

        async _handleFontAction(action, font, spinner, stack, message) {
            const adjustment = this._scroller.get_vadjustment();
            this.scrollValue = adjustment ? adjustment.get_value() : 0;
            try {
                stack.set_visible_child_name("spinner");
                spinner.spinning = true;
                await this.versionManager.setupFontsVersion();
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
                this.#populateCustomFontsList();
                this._handleScrollPosition(this.scrollValue);
                this._updateNoFontsBanner();
            } catch (error) {
                spinner.spinning = false;
                stack.set_visible_child_name("icon");
                console.log("Action failed", error);
                throw error;
            }
        }

        async _handleInstallButton(font, spinner, stack, message) {
            try {
                await this._handleFontAction(
                    async (font) => {
                        try {
                            await this.fontsManager.downloadAndInstall(
                                font.tarName,
                                this.versionManager.get(),
                            );
                            this.installedFontsManager.update(
                                font.tarName,
                                this.versionManager.get(),
                            );
                        } catch (error) {
                            throw error;
                        }
                    },
                    font,
                    spinner,
                    stack,
                    message,
                );
            } catch (error) {
                throw error;
            }
        }

        async _handleRemoveButton(font, spinner, stack) {
            try {
                await this._handleFontAction(
                    async (font) => {
                        try {
                            await this.fontsManager.remove(font.tarName);
                            this.installedFontsManager.remove(font.tarName);
                        } catch (error) {
                            throw error;
                        }
                    },
                    font,
                    spinner,
                    stack,
                    _("Font removed"),
                );
            } catch (error) {
                throw error;
            }
        }

        vfunc_close_request() {
            super.vfunc_close_request();
            this.run_dispose();
        }

        #showAddCustomFontDialog() {
            const dialog = new Adw.Dialog({
                title: _("Add Custom Font"),
                content_width: 400,
            });

            const toolbarView = new Adw.ToolbarView();

            const headerBar = new Adw.HeaderBar();
            toolbarView.add_top_bar(headerBar);

            const clamp = new Adw.Clamp({
                maximum_size: 400,
                margin_top: 24,
                margin_bottom: 24,
                margin_start: 12,
                margin_end: 12,
            });

            const box = new Gtk.Box({
                orientation: Gtk.Orientation.VERTICAL,
                spacing: 24,
            });

            const listBox = new Gtk.ListBox({
                selection_mode: Gtk.SelectionMode.NONE,
            });
            listBox.add_css_class("boxed-list");

            const nameRow = new Adw.EntryRow({
                title: _("Font name"),
            });

            const descriptionRow = new Adw.EntryRow({
                title: _("Description"),
            });

            const urlRow = new Adw.EntryRow({
                title: _("Zip file URL"),
            });

            listBox.append(nameRow);
            listBox.append(descriptionRow);
            listBox.append(urlRow);

            const addButton = new Gtk.Button({
                label: _("Add"),
                halign: Gtk.Align.CENTER,
            });
            addButton.add_css_class("suggested-action");
            addButton.add_css_class("pill");
            addButton.set_sensitive(false);

            const updateButtonSensitivity = () => {
                const hasName = nameRow.text.trim() !== "";
                const hasUrl = urlRow.text.trim() !== "";
                addButton.set_sensitive(hasName && hasUrl);
            };

            nameRow.connect("changed", updateButtonSensitivity);
            urlRow.connect("changed", updateButtonSensitivity);

            addButton.connect("clicked", async () => {
                const name = nameRow.text.trim();
                const description = descriptionRow.text.trim();
                const url = urlRow.text.trim();

                dialog.close();

                try {
                    this.customFontsManager.add(name, description, url);

                    this.fontsManager.loadFontDirectories();
                    this.fonts = this.fontsManager.loadFonts();
                    this._searchList.remove_all();
                    this.#setupSearch();
                    this.#populateFontLists();
                    this.#populateCustomFontsList();

                    const toast = new Adw.Toast({
                        title: _("Custom font added to list"),
                    });
                    this._toastOverlay.add_toast(toast);
                } catch (error) {
                    this._handleError(error, _("Failed to add custom font: %s"));
                }
            });

            box.append(listBox);
            box.append(addButton);
            clamp.set_child(box);
            toolbarView.set_content(clamp);
            dialog.set_child(toolbarView);
            dialog.present(this);
        }
    },
);
