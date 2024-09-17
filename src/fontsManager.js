import GLib from "gi://GLib";
import Gio from "gi://Gio";
import Soup from "gi://Soup";
import Autoar from "gi://GnomeAutoar";

export class FontsManager {
    constructor(installedFonts, version) {
        this.installedFonts = installedFonts;
        this.version = version;

        Gio._promisify(
            Gio.File.prototype,
            "query_info_async",
            "query_info_finish",
        );
        Gio._promisify(
            Gio.File.prototype,
            "enumerate_children_async",
            "enumerate_children_finish",
        );

        Gio._promisify(Gio.File.prototype, "delete_async", "delete_finish");
    }

    loadFontDirectories() {
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

                    if (!this.installedFonts.hasGroup(fontName)) {
                        this.installedFonts.update(fontName, "v0");
                    }
                }
            }

            this._fontDirectories = directories;
        } catch (error) {
            console.log("Failed to enumerate font directories:", error);
            this._fontDirectories = [];
        }
    }

    loadFonts() {
        const resourcePath = "/io/github/getnf/embellish/fonts";
        const keyFile = new GLib.KeyFile();

        try {
            let data = Gio.resources_lookup_data(
                resourcePath,
                Gio.ResourceLookupFlags.NONE,
            );
            keyFile.load_from_bytes(data, GLib.KeyFileFlags.NONE);
        } catch (error) {
            throw new Error(`Failed to load ${resourcePath}`, error);
        }

        let fonts = [];
        const groups = keyFile.get_groups()[0];

        groups.forEach((group) => {
            const description = keyFile.get_string(group, "description");
            const licenceIds = keyFile.get_string_list(group, "licenceId");
            const tarName = keyFile.get_string(group, "tarName");
            const isInstalled = this._isInstalled(tarName);
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
                hasUpdate: isInstalled ? this._HasUpdate(tarName) : false,
            });
        });

        return fonts;
    }

    _isInstalled(fontName) {
        return this._fontDirectories.includes(fontName);
    }

    _HasUpdate(fontName) {
        let fontVersion;
        let latestVersion;

        try {
            fontVersion = this.installedFonts.getVersion(fontName);
            latestVersion = this.version.get();
        } catch (error) {
            console.log(error);
        }

        return fontVersion !== latestVersion ? true : false;
    }

    async downloadAndInstall(tarName, version) {
        try {
            await this._download(tarName, version);
            await this._extract(tarName);
        } catch (error) {
            throw error;
        }
    }

    async remove(tarName) {
        const fontDir = GLib.build_filenamev([
            GLib.get_home_dir(),
            ".local",
            "share",
            "fonts",
            tarName,
        ]);

        try {
            const file = Gio.File.new_for_path(fontDir);
            await this._deleteRecursively(file);
        } catch (error) {
            throw error;
        }
    }

    async _download(tarName, release) {
    try{
        const url = `https://github.com/ryanoasis/nerd-fonts/releases/download/${release}/${tarName}.tar.xz`;
        const downloadDir = GLib.build_filenamev([
            GLib.get_user_special_dir(GLib.UserDirectory.DIRECTORY_DOWNLOAD),
            "embellish",
            tarName,
        ]);
            await this._downloadTarXzFile(url, downloadDir);
        } catch (error) {
            throw error;
        }
    }

    async _extract(tarName) {
     try {
        const downloadDir = GLib.build_filenamev([
            GLib.get_user_special_dir(GLib.UserDirectory.DIRECTORY_DOWNLOAD),
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
                throw new Error(`Failed to extract ${filePath}`);
            }
        });

        try {
            extractor.start(null);
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

    async _deleteRecursively(file) {
        try {
            const info = await file.query_info_async(
                "*",
                Gio.FileQueryInfoFlags.NONE,
                GLib.PRIORITY_DEFAULT,
                null,
            );
            const fileType = info.get_file_type();

            if (fileType === Gio.FileType.DIRECTORY) {
                const enumerator = await file.enumerate_children_async(
                    "*",
                    Gio.FileQueryInfoFlags.NONE,
                    GLib.PRIORITY_DEFAULT,
                    null,
                );
                let childInfo;

                while ((childInfo = await enumerator.next_file(null))) {
                    const child = file.get_child(childInfo.get_name());
                    await this._deleteRecursively(child);
                }

                await file.delete_async(GLib.PRIORITY_DEFAULT, null);
            } else {
                await file.delete_async(GLib.PRIORITY_DEFAULT, null);
            }
        } catch (error) {
            console.error("Error while deleting files:", error);
            throw error;
        }
    }
}
