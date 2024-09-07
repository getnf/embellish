import "gi://Gdk?version=4.0";
import "gi://Gtk?version=4.0";
import "gi://Adw?version=1";

import { EmbApplication } from "./application.js";

export function main(argv) {
    return new EmbApplication({
        "application-id": pkg.name,
        resource_base_path: "/io/github/getnf/embellish",
    }).run(argv);
}
