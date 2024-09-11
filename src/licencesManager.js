import Gtk from "gi://Gtk";
import GLib from "gi://GLib";
import Gio from "gi://Gio";
import Adw from "gi://Adw";

export class LicencesManager {
    constructor() {}

    new(font) {
        const licenseBox = this._createBox(Gtk.Orientation.HORIZONTAL, 12);

        const licenceButton = new Gtk.MenuButton();
        licenceButton.add_css_class("licence-button");
        licenceButton.set_tooltip_text("licence details");
        const licenceButtonLabel = new Gtk.Label();
        if (font.licences.length > 1) {
            // Translators: This means that the font has two licences
            licenceButtonLabel.set_label(_("Dual"));
        } else {
            licenceButtonLabel.set_label(font.licences[0]);
        }
        licenceButton.set_always_show_arrow(false);
        licenceButton.set_child(licenceButtonLabel);
        licenceButton.set_popover(this._createPopover(font));
        licenseBox.append(licenceButton);
        return licenseBox;
    }

    _createPopover(font) {
        const popover = new Gtk.Popover({
            "has-arrow": true,
            name: "licensePopover",
        });

        const box = this._createBox(Gtk.Orientation.VERTICAL, 4);
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
            const licenceBox = this._createLicenceBox(
                licence,
                this._getDescription(licence),
            );
            box.append(licenceBox);
        });

        const clamp = new Adw.Clamp({
            "maximum-size": 250,
        });

        clamp.set_child(box);
        popover.set_child(clamp);

        return popover;
    }

    _createLicenceBox(id, description) {
        const box = this._createBox(Gtk.Orientation.VERTICAL, 4);

        const idLabel = new Gtk.Label({ label: id });
        idLabel.add_css_class("heading");

        const descriptionLabel = new Gtk.Label({ label: description });
        descriptionLabel.set_justify(2);
        descriptionLabel.set_wrap(true);

        box.append(idLabel);
        box.append(descriptionLabel);

        return box;
    }

    _getDescription(licenceKey) {
        const resourcePath = "/io/github/getnf/embellish/licences";
        const keyFile = new GLib.KeyFile();

        try {
            const data = Gio.resources_lookup_data(
                resourcePath,
                Gio.ResourceLookupFlags.NONE,
            );
            keyFile.load_from_bytes(data, GLib.KeyFileFlags.NONE);
        } catch (e) {
            console.log(e, `Failed to load ${resourcePath}`);
            return _("No description available");
        }

        const description = keyFile.get_string("licences", licenceKey);
        return description ? description : _("No description available");
    }

    _createBox(orientation, spacing) {
        const box = new Gtk.Box({
            orientation,
            spacing,
        });
        box.set_halign(Gtk.Align.CENTER);
        box.set_valign(Gtk.Align.CENTER);
        return box;
    }
}
