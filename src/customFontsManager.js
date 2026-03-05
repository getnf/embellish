import GLib from "gi://GLib";

export class CustomFontsManager {
    constructor() {
        this._keyFilePath = GLib.build_filenamev([
            GLib.get_user_config_dir(),
            "embellish",
            "custom-fonts",
        ]);

        this.#setup();
    }

    #setup() {
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

    getAll() {
        const fonts = [];
        const [groups] = this._keyFile.get_groups();

        groups.forEach((group) => {
            try {
                const description = this._keyFile.get_string(group, "description");
                const url = this._keyFile.get_string(group, "url");
                const dirName = this._keyFile.get_string(group, "dirName");

                fonts.push({
                    name: group,
                    description,
                    url,
                    dirName,
                });
            } catch (error) {
                console.error(`Failed to read custom font "${group}":`, error);
            }
        });

        return fonts;
    }

    add(name, description, url) {
        const dirName = name.replace(/[^a-zA-Z0-9]/g, "");

        this._keyFile.set_string(name, "description", description);
        this._keyFile.set_string(name, "url", url);
        this._keyFile.set_string(name, "dirName", dirName);

        try {
            this._keyFile.save_to_file(this._keyFilePath);
        } catch (error) {
            throw error;
        }

        return dirName;
    }

    export() {
        return JSON.stringify(this.getAll(), null, 2);
    }

    import(jsonData) {
        try {
            const fonts = JSON.parse(jsonData);
            if (!Array.isArray(fonts)) {
                throw new Error("Invalid import data: expected an array");
            }

            fonts.forEach((font) => {
                if (font.name && font.url) {
                    this.add(font.name, font.description || "", font.url);
                }
            });
        } catch (error) {
            throw error;
        }
    }

    remove(name) {
        this._keyFile.remove_group(name);

        try {
            this._keyFile.save_to_file(this._keyFilePath);
        } catch (error) {
            throw error;
        }
    }
}
