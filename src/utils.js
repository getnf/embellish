import Gtk from "gi://Gtk";

export class Utils {
    constructor() {}

    createBox(orientation, spacing) {
        const box = new Gtk.Box({
            orientation,
            spacing,
        });
        box.set_halign(Gtk.Align.CENTER);
        box.set_valign(Gtk.Align.CENTER);
        return box;
    }

    createSpinnerButton(icon, tooltip) {
        const button = new Gtk.Button();
        button.add_css_class("flat");
        button.set_tooltip_text(tooltip);
        const buttonBox = this.createBox(Gtk.Orientation.HORIZONTAL, 0);
        const buttonIcon = Gtk.Image.new_from_resource(
            `/io/github/getnf/embellish/icons/scalable/actions/${icon}.svg`,
        );
        const buttonSpinner = new Gtk.Spinner();
        buttonSpinner.set_visible(false);
        buttonBox.append(buttonIcon);
        buttonBox.append(buttonSpinner);
        button.set_child(buttonBox);

        return { button, buttonIcon, buttonSpinner };
    }

    escapeMarkup(text) {
        return text
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}
