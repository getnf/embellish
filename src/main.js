import "gi://Gdk?version=4.0";
import "gi://Gtk?version=4.0";
import "gi://Adw?version=1";

import { EmbApplication } from "./application.js";

String.prototype.format = function (args) {
    let str = this;
    if (arguments.length <= 0) return str;

    // Simple implementation of %s and %d replacement
    for (let i = 0; i < arguments.length; i++) {
        str = str.replace(/%[sd]/, arguments[i]);
    }
    return str;
};

export function main(argv) {
    return new EmbApplication({
        "application-id": pkg.name,
        resource_base_path: "/io/github/getnf/embellish",
    }).run(argv);
}
