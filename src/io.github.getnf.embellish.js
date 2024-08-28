#!@GJS@ -m

import GLib from "gi://GLib";

// Initialize the package
imports.package.init({
    name: "@PACKAGE_NAME@",
    version: "@PACKAGE_VERSION@",
    prefix: "@PREFIX@",
    libdir: "@LIBDIR@",
});

imports.package.initGettext();

// Import the main module and run the main function
const loop = new GLib.MainLoop(null, false);
import("resource:///io/github/getnf/embellish/js/main.js")
    .then((main) => {
        GLib.idle_add(GLib.PRIORITY_DEFAULT_IDLE, () => {
            loop.quit();
            imports.package.run(main);
            return GLib.SOURCE_REMOVE;
        });
    })
    .catch(logError);
loop.run();
