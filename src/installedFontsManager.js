import GLib from "gi://GLib";

export class InstalledFontsManager {
    constructor() {
        this._keyFilePath = GLib.build_filenamev([
            GLib.get_user_config_dir(),
            "embellish",
            "fonts",
        ]);

        this.#setupInstalledFonts();
    }

    #setupInstalledFonts() {
        const keyFile = new GLib.KeyFile();
        const dirPath = GLib.path_get_dirname(this._keyFilePath);

        if (!GLib.file_test(this._keyFilePath, GLib.FileTest.EXISTS)) {
            GLib.mkdir_with_parents(dirPath, 0o755);
        }

        try {
            keyFile.load_from_file(this._keyFilePath, GLib.KeyFileFlags.NONE);
        } catch (error) {
            keyFile.save_to_file(this._keyFilePath);
        }

        this._keyFile = keyFile;
    }

    hasGroup(group) {
        return this._keyFile.has_group(group);
    }

    getVersion(group) {
        return this._keyFile.get_string(group, "version");
    }

    remove(fontName) {
        this._keyFile.remove_group(fontName);

        try {
            this._keyFile.save_to_file(this._keyFilePath);
        } catch (error) {
            throw error;
        }
    }

    update(fontName, version) {
        this._keyFile.set_string(fontName, "version", version);

        try {
            this._keyFile.save_to_file(this._keyFilePath);
        } catch (error) {
            throw error;
        }
    }
}
