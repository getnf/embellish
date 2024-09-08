import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";
import GLib from "gi://GLib";
import Soup from "gi://Soup";
import Autoar from "gi://GnomeAutoar";

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
            "installedFonts",
            "availableFonts",
        ],
    },
    class extends Adw.ApplicationWindow {
        constructor(params = {}) {
            super(params);
            this.#setupActions();
            this.#setupWelcomeScreen();
            this.#loadFontDirectories();
            this.#loadFonts();
            this.#setupSearch();
            this.#populateFontLists();

            Gio._promisify(
                Soup.Session.prototype,
                "send_and_read_async",
                "send_and_read_finish",
            );

            this.#setupFontsVersion();
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

        #loadFontDirectories() {
            const fontDir = GLib.build_filenamev([
                GLib.get_home_dir(),
                ".local",
                "share",
                "fonts",
            ]);

            try {
                const fontDirectoryFile = Gio.File.new_for_path(fontDir);
                const enumerator = fontDirectoryFile.enumerate_children(
                    "standard::name,standard::type",
                    Gio.FileQueryInfoFlags.NONE,
                    null,
                );

                const directories = [];
                let fileInfo;

                while ((fileInfo = enumerator.next_file(null)) !== null) {
                    if (fileInfo.get_file_type() === Gio.FileType.DIRECTORY) {
                        const fontName = fileInfo.get_name();
                        directories.push(fontName);

                        let keyFile;

                        try {
                            keyFile = this._getInstalledFontsKeyFile();
                        } catch (error) {
                            console.error(`Error getting key file: ${error}`);
                            continue;
                        }

                        const hasGroup = keyFile.has_group(fontName);

                        if (!hasGroup) {
                            this._updateInstalledFonts(fontName, "v0");
                        }
                    }
                }

                this._fontDirectories = directories;
            } catch (error) {
                console.log("Failed to enumerate font directories:", error);
                this._fontDirectories = [];
            }
        }

        #loadFonts() {
            const resourcePath = "/io/github/getnf/embellish/fonts";
            const keyFile = new GLib.KeyFile();

            try {
                let data = Gio.resources_lookup_data(
                    resourcePath,
                    Gio.ResourceLookupFlags.NONE,
                );
                keyFile.load_from_bytes(data, GLib.KeyFileFlags.NONE);
            } catch (error) {
                console.log(`Failed to load ${resourcePath}`, error);
                return;
            }

            let fonts = [];
            const groups = keyFile.get_groups()[0];

            groups.forEach((group) => {
                const description = keyFile.get_string(group, "description");
                const licenceIds = keyFile.get_string_list(group, "licenceId");
                const tarName = keyFile.get_string(group, "tarName");
                const isInstalled = this._isFontInstalled(tarName);
                let patchedName = "";
                try {
                    patchedName = keyFile.get_string(group, "patchedName");
                } catch (e) {
                    patchedName = "";
                }

                fonts.push({
                    name: group,
                    patchedName,
                    tarName,
                    description,
                    licences: licenceIds,
                    installed: isInstalled,
                    hasUpdate: isInstalled
                        ? this._fontHasUpdate(tarName)
                        : false,
                });
            });

            this.fonts = fonts;
        }

        _isFontInstalled(fontName) {
            return this._fontDirectories.includes(fontName);
        }

        _fontHasUpdate(fontName) {
            let fontVersion;
            let latestVersion;

            try {
                fontVersion = this._getInstalledFontVersion(fontName);
                latestVersion = this._getVersion();
            } catch (error) {
                console.log(error);
            }

            return fontVersion !== latestVersion ? true : false;
        }

        _getInstalledFontsKeyFile() {
            const keyFilePath = GLib.build_filenamev([
                GLib.get_user_config_dir(),
                "embellish",
                "fonts",
            ]);
            const keyFile = new GLib.KeyFile();

            if (GLib.file_test(keyFilePath, GLib.FileTest.EXISTS)) {
                try {
                    keyFile.load_from_file(keyFilePath, GLib.KeyFileFlags.NONE);
                } catch (error) {
                    throw error;
                }
            } else {
                GLib.mkdir_with_parents(
                    GLib.path_get_dirname(keyFilePath),
                    0o755,
                );

                try {
                    keyFile.save_to_file(keyFilePath);
                    keyFile.load_from_file(keyFilePath, GLib.KeyFileFlags.NONE);
                } catch (error) {
                    throw error;
                }
            }

            return keyFile;
        }

        _getInstalledFontVersion(group) {
            let keyFile;
            try {
                keyFile = this._getInstalledFontsKeyFile();
            } catch (error) {
                throw error;
            }

            return keyFile.get_string(group, "version");
        }

        _removeInstalledFont(fontName) {
            const keyFilePath = GLib.build_filenamev([
                GLib.get_user_config_dir(),
                "embellish",
                "fonts",
            ]);
            let keyFile;
            try {
                keyFile = this._getInstalledFontsKeyFile();
            } catch (error) {
                throw error;
            }

            keyFile.remove_group(fontName);
            try {
                keyFile.save_to_file(keyFilePath);
            } catch (error) {
                throw error;
            }
        }

        _updateInstalledFonts(fontName, version) {
            const keyFilePath = GLib.build_filenamev([
                GLib.get_user_config_dir(),
                "embellish",
                "fonts",
            ]);
            let keyFile;
            try {
                keyFile = this._getInstalledFontsKeyFile();
            } catch (error) {
                throw error;
            }

            keyFile.set_string(fontName, "version", version);

            try {
                keyFile.save_to_file(keyFilePath);
            } catch (error) {
                throw error;
            }
        }

        #populateFontLists() {
            this._installedFonts.remove_all();
            this._availableFonts.remove_all();

            this.fonts.forEach((font) => {
                let title = font.name;
                if (font.patchedName !== "") {
                    title = `${font.name} (${font.patchedName})`;
                }

                const row = new Adw.ActionRow({
                    title: title,
                    subtitle: this._escapeMarkup(font.description),
                });

                const suffix = this._makeRowSuffix(font);
                row.add_suffix(suffix);

                if (font.installed) {
                    this._installedFonts.append(row);
                } else {
                    this._availableFonts.append(row);
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
                let title = font.name;
                if (font.patchedName !== "") {
                    title = `${font.name} (${font.patchedName})`;
                }

                const row = new Adw.ActionRow({
                    title: title,
                    subtitle: this._escapeMarkup(font.description),
                });

                const suffix = this._makeRowSuffix(font);
                row.add_suffix(suffix);

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

        _makeRowSuffix(font) {
            const box = new Gtk.Box({
                orientation: "vertical",
                spacing: 12,
            });
            box.set_halign(3);
            box.set_valign(3);

            const licenseBox = new Gtk.Box({
                orientation: "horizontal",
                spacing: 12,
            });
            licenseBox.set_halign(3);
            licenseBox.set_valign(3);

            const licenceButton = new Gtk.MenuButton();
            licenceButton.add_css_class("licence-button");
            const licenceButtonLabel = new Gtk.Label();
            if (font.licences.length > 1) {
                // Translators: This means that the font has two licences
                licenceButtonLabel.set_label(_("Dual"));
            } else {
                licenceButtonLabel.set_label(font.licences[0]);
            }
            licenceButton.set_always_show_arrow(false);
            licenceButton.set_child(licenceButtonLabel);
            licenceButton.set_popover(this._makeLicencesPopover(font));
            licenseBox.append(licenceButton);
            box.append(licenseBox);

            const previewButton = new Gtk.Button({
                icon_name: "embellish-preview-symbolic",
            });
            previewButton.add_css_class("flat");
            previewButton.connect("clicked", () => {
                this._showPreviewDialog(font.tarName);
            });

            box.append(previewButton);

            const installButton = new Gtk.Button();
            installButton.add_css_class("flat");
            const installButtonBox = new Gtk.Box({
                orientation: Gtk.Orientation.HORIZONTAL,
            });
            const installButtonIcon = Gtk.Image.new_from_resource(
                "/io/github/getnf/embellish/icons/scalable/actions/embellish-download-symbolic.svg",
            );
            const installButtonSpinner = new Gtk.Spinner();
            installButtonSpinner.set_visible(false);
            installButtonBox.append(installButtonIcon);
            installButtonBox.append(installButtonSpinner);
            installButton.set_child(installButtonBox);
            installButton.connect("clicked", async () => {
                try {
                    await this._handleInstallButton(
                        font,
                        installButtonSpinner,
                        installButtonIcon,
                        _("Font Intalled"),
                    );
                } catch (error) {
                    const toast = new Adw.Toast({
                        title: _(`Installation failed: ${error}`),
                    });
                    this._toastOverlay.add_toast(toast);
                    console.log(error);
                }
            });

            const updateButton = new Gtk.Button();
            updateButton.add_css_class("flat");
            const updateButtonBox = new Gtk.Box({
                orientation: Gtk.Orientation.HORIZONTAL,
            });
            const updateButtonIcon = Gtk.Image.new_from_resource(
                "/io/github/getnf/embellish/icons/scalable/actions/embellish-update-symbolic.svg",
            );
            const updateButtonSpinner = new Gtk.Spinner();
            updateButtonSpinner.set_visible(false);
            updateButtonBox.append(updateButtonIcon);
            updateButtonBox.append(updateButtonSpinner);
            updateButton.set_child(updateButtonBox);
            updateButton.connect("clicked", async () => {
                try {
                    await this._handleInstallButton(
                        font,
                        updateButtonSpinner,
                        updateButtonIcon,
                        _("Font updated"),
                    );
                } catch (error) {
                    const toast = new Adw.Toast({
                        title: _(`Updating failed: ${error}`),
                    });
                    this._toastOverlay.add_toast(toast);
                    console.log(error);
                }
            });

            const removeButton = new Gtk.Button({
                icon_name: "embellish-remove-symbolic",
            });
            removeButton.add_css_class("flat");

            removeButton.connect("clicked", () => {
                this._handleRemoveButton(font);
            });

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

        async _handleInstallButton(font, spinner, icon, message) {
            try {
                icon.set_visible(false);
                spinner.set_visible(true);
                spinner.spinning = true;
                await this._downloadAndInstallFont(font.tarName);
                spinner.spinning = false;
                spinner.set_visible(false);
                icon.set_visible(true);
                const toast = new Adw.Toast({
                    title: message,
                });
                this._toastOverlay.add_toast(toast);
                this._updateInstalledFonts(font.tarName, this._getVersion());
                this.#loadFontDirectories();
                this.#loadFonts();
                this._searchList.remove_all();
                this.#setupSearch();
                this.#populateFontLists();
            } catch (error) {
                spinner.spinning = false;
                spinner.set_visible(false);
                icon.set_visible(true);
                throw error;
            }
        }

        _handleRemoveButton(font) {
            try {
                this._removeFonts(font.tarName);
                const toast = new Adw.Toast({
                    title: _("Font removed"),
                });
                this._toastOverlay.add_toast(toast);
                this._removeInstalledFont(font.tarName);
                this.#loadFontDirectories();
                this.#loadFonts();
                this._searchList.remove_all();
                this.#setupSearch();
                this.#populateFontLists();
            } catch (error) {
                const toast = new Adw.Toast({
                    title: _(`Removal failed: ${error}`),
                });
                this._toastOverlay.add_toast(toast);
                console.log("Font removal failed: ", error);
            }
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

            previewPicture.set_halign(3);
            previewPicture.set_valign(3);
            previewPicture.set_can_shrink(false);
            previewPicture.set_margin_start(12);
            previewPicture.set_margin_end(12);
            previewPicture.add_css_class("svg-preview");

            page.set_content(previewPicture);
            dialog.set_child(page);

            dialog.present(this);
        }

        _makeLicencesPopover(font) {
            const popover = new Gtk.Popover({
                "has-arrow": true,
                name: "licensePopover",
            });

            const box = new Gtk.Box({
                orientation: Gtk.Orientation.VERTICAL,
                spacing: 12,
            });
            box.set_halign(3);
            box.set_valign(3);
            box.set_margin_top(12);
            box.set_margin_start(24);
            box.set_margin_end(24);
            box.set_margin_bottom(12);

            const licenceLogo = Gtk.Image.new_from_resource(
                `/io/github/getnf/embellish/licence.svg`,
            );
            licenceLogo.set_pixel_size(64);
            box.append(licenceLogo);

            font.licences.forEach((licence) => {
                const lice = this._makeLicenceBox(
                    licence,
                    this._getLicenceDescription(licence),
                );
                box.append(lice);
            });

            const clamp = new Adw.Clamp({
                "maximum-size": 250,
            });

            clamp.set_child(box);

            popover.set_child(clamp);
            return popover;
        }

        _makeLicenceBox(id, description) {
            const box = new Gtk.Box({
                orientation: Gtk.Orientation.VERTICAL,
                spacing: 4,
            });
            box.set_halign(3);
            box.set_valign(3);

            const idLabel = new Gtk.Label({ label: id });
            idLabel.add_css_class("heading");
            const descriptionLabel = new Gtk.Label({ label: description });
            descriptionLabel.set_justify(2);
            descriptionLabel.set_wrap(true);

            box.append(idLabel);
            box.append(descriptionLabel);

            return box;
        }

        _getLicenceDescription(licenceKey) {
            const resourcePath = "/io/github/getnf/embellish/licences";
            const keyFile = new GLib.KeyFile();

            try {
                let data = Gio.resources_lookup_data(
                    resourcePath,
                    Gio.ResourceLookupFlags.NONE,
                );
                keyFile.load_from_bytes(data, GLib.KeyFileFlags.NONE);
            } catch (e) {
                console.log(e, `Failed to load ${resourcePath}`);
                return null;
            }
            const description = keyFile.get_string("licences", licenceKey);

            return description ? description : _("No description available");
        }

        async #setupFontsVersion() {
            let latestVersion;
            let currentVersion;

            try {
                currentVersion = this._getVersion();
            } catch (error) {
                console.log(error);
            }

            try {
                latestVersion = await this._getLatestRelease();
            } catch (error) {
                console.log("Failed to fetch the latest release: ", error);
                const toast = new Adw.Toast({
                    title: _(
                        `Failed to fetch the latest release version ${error}`,
                    ),
                });
                this._toastOverlay.add_toast(toast);
                return;
            }

            if (latestVersion !== currentVersion.current) {
                try {
                    this._updateVersion(latestVersion);
                } catch (error) {
                    console.log(error);
                }
            }
        }

        _getVersionKeyFile() {
            const keyFilePath = GLib.build_filenamev([
                GLib.get_user_config_dir(),
                "embellish",
                "version",
            ]);
            const keyFile = new GLib.KeyFile();

            if (GLib.file_test(keyFilePath, GLib.FileTest.EXISTS)) {
                try {
                    keyFile.load_from_file(keyFilePath, GLib.KeyFileFlags.NONE);
                } catch (error) {
                    throw error;
                }
            } else {
                keyFile.set_string("NerdFonts", "version", "v0");

                GLib.mkdir_with_parents(
                    GLib.path_get_dirname(keyFilePath),
                    0o755,
                );

                try {
                    keyFile.save_to_file(keyFilePath);
                    console.log(
                        "Version Keyfile initialized with default value.",
                    );
                    keyFile.load_from_file(keyFilePath, GLib.KeyFileFlags.NONE);
                } catch (error) {
                    throw error;
                }
            }

            return keyFile;
        }

        _getVersion() {
            let keyFile;
            try {
                keyFile = this._getVersionKeyFile();
            } catch (error) {
                throw error;
            }

            return keyFile.get_string("NerdFonts", "version");
        }

        _updateVersion(version) {
            const keyFilePath = GLib.build_filenamev([
                GLib.get_user_config_dir(),
                "embellish",
                "version",
            ]);
            let keyFile;
            try {
                keyFile = this._getVersionKeyFile();
            } catch (error) {
                throw error;
            }

            keyFile.set_string("NerdFonts", "version", version);

            try {
                keyFile.save_to_file(keyFilePath);
            } catch (error) {
                throw error;
            }
        }

        async _getLatestRelease() {
            const session = new Soup.Session();

            const request = Soup.Message.new(
                "GET",
                "https://api.github.com/repos/ryanoasis/nerd-fonts/releases/latest",
            );

            request.request_headers.append("User-Agent", "Embellish/0.4");

            try {
                const bytes = await session.send_and_read_async(
                    request,
                    GLib.PRIORITY_DEFAULT,
                    null,
                );

                if (request.get_status() !== Soup.Status.OK) {
                    throw new Error(
                        `HTTP request failed with status: ${request.get_status()}`,
                    );
                }

                const textDecoder = new TextDecoder("utf-8");
                const responseText = textDecoder.decode(bytes.toArray());
                const jsonResponse = JSON.parse(responseText);
                const release = jsonResponse.tag_name;

                return release;
            } catch (error) {
                throw error;
            }
        }

        _removeFonts(tarName) {
            const fontDir = GLib.build_filenamev([
                GLib.get_home_dir(),
                ".local",
                "share",
                "fonts",
                tarName,
            ]);

            try {
                let file = Gio.File.new_for_path(fontDir);
                this._deleteRecursively(file);
            } catch (error) {
                throw error;
            }
        }

        async _downloadAndInstallFont(tarName) {
            try {
                await this._downloadFont(tarName);
                await this._extractFont(tarName);
            } catch (error) {
                throw error;
            }
        }

        async _downloadFont(tarName) {
            const release = this._getVersion();
            const url = `https://github.com/ryanoasis/nerd-fonts/releases/download/${release}/${tarName}.tar.xz`;
            const downloadDir = GLib.build_filenamev([
                GLib.get_user_special_dir(
                    GLib.UserDirectory.DIRECTORY_DOWNLOAD,
                ),
                "embellish",
                tarName,
            ]);

            try {
                await this._downloadTarXzFile(url, downloadDir);
            } catch (error) {
                throw error;
            }
        }

        async _extractFont(tarName) {
            const downloadDir = GLib.build_filenamev([
                GLib.get_user_special_dir(
                    GLib.UserDirectory.DIRECTORY_DOWNLOAD,
                ),
                "embellish",
                tarName,
            ]);

            const extractDir = GLib.build_filenamev([
                GLib.get_home_dir(),
                ".local",
                "share",
                "fonts",
                tarName,
            ]);

            try {
                await this._extractTarXz(downloadDir, extractDir);
            } catch (error) {
                throw error;
            }
        }

        async _downloadTarXzFile(url, destinationPath) {
            const session = new Soup.Session();

            try {
                const request = Soup.Message.new("GET", url);
                const bytes = await session.send_and_read_async(
                    request,
                    GLib.PRIORITY_DEFAULT,
                    null,
                );

                if (request.get_status() !== Soup.Status.OK) {
                    throw new Error(
                        `HTTP request failed with status: ${request.get_status()}`,
                    );
                }

                try {
                    this._saveBytesToFile(bytes, destinationPath);
                } catch (error) {
                    throw error;
                }
            } catch (error) {
                throw error;
            }
        }

        _saveBytesToFile(bytes, filePath) {
            try {
                const file = Gio.File.new_for_path(filePath);
                const outputStream = file.replace(
                    null,
                    false,
                    Gio.FileCreateFlags.NONE,
                    null,
                );
                outputStream.write_all(bytes.get_data(), null);
                outputStream.close(null);
            } catch (error) {
                throw error;
            }
        }

        async _extractTarXz(filePath, fontsDir) {
            const file = Gio.File.new_for_path(filePath);

            const destination = Gio.File.new_for_path(fontsDir);

            const extractor = new Autoar.Extractor({
                source_file: file,
                output_file: destination,
            });

            extractor.set_output_is_dest(true);

            extractor.connect("error", (error) => {
                if (error) {
                    throw new Error(`failed to extract ${filePath}`);
                }
            });

            try {
                extractor.start(null);
            } catch (error) {
                throw error;
            }
        }

        _deleteRecursively(file) {
            try {
                if (
                    file.query_file_type(Gio.FileQueryInfoFlags.NONE, null) ===
                    Gio.FileType.DIRECTORY
                ) {
                    let enumerator = file.enumerate_children(
                        "*",
                        Gio.FileQueryInfoFlags.NONE,
                        null,
                    );
                    let info;
                    while ((info = enumerator.next_file(null))) {
                        let child = file.get_child(info.get_name());
                        this._deleteRecursively(child);
                    }
                    file.delete(null);
                } else {
                    file.delete(null);
                }
            } catch (e) {
                throw error;
            }
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