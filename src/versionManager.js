import GLib from "gi://GLib";
import Gio from "gi://Gio";
import Soup from "gi://Soup";

export class VersionManager {
    constructor() {
        this._keyFilePath = GLib.build_filenamev([
            GLib.get_user_config_dir(),
            "embellish",
            "version",
        ]);

        this._lastCheck = globalThis.settings.get_string("last-check");

        Gio._promisify(
            Soup.Session.prototype,
            "send_and_read_async",
            "send_and_read_finish",
        );

        this.#setupVersionKeyFile();
    }

    #setupVersionKeyFile() {
        const keyFile = new GLib.KeyFile();
        const dirPath = GLib.path_get_dirname(this._keyFilePath);

        if (!GLib.file_test(this._keyFilePath, GLib.FileTest.EXISTS)) {
            GLib.mkdir_with_parents(dirPath, 0o755);
        }

        try {
            keyFile.load_from_file(this._keyFilePath, GLib.KeyFileFlags.NONE);
        } catch (error) {
            keyFile.set_string("NerdFonts", "version", "v0");
            keyFile.save_to_file(this._keyFilePath);
            console.log("Version Keyfile initialized with default value.");
        }

        this._keyFile = keyFile;
    }

    get() {
        return this._keyFile.get_string("NerdFonts", "version");
    }

    update(version) {
        this._keyFile.set_string("NerdFonts", "version", version);

        try {
            this._keyFile.save_to_file(this._keyFilePath);
        } catch (error) {
            throw error;
        }
    }

    async setupFontsVersion() {
        const lastCheckDate = new Date(this._lastCheck);
        const currentDate = new Date();
        const daysDifference =
            (currentDate - lastCheckDate) / (1000 * 3600 * 24);

        let latestVersion;
        let currentVersion;

        try {
            currentVersion = this.get();
        } catch (error) {
            console.log(error);
        }

        if (daysDifference < 7 && currentVersion !== "v0") {
            console.log(
                "Version check skipped. Last checked: " + this._lastCheck,
            );
            return;
        }

        globalThis.settings.set_string("last-check", currentDate.toISOString());

        try {
            latestVersion = await this._getLatestRelease();
        } catch (error) {
            console.log("Failed to fetch the latest release: ", error);
            return;
        }

        if (latestVersion !== currentVersion) {
            try {
                this.update(latestVersion);
            } catch (error) {
                console.log(error);
            }
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
}

