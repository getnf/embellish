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
        button.add_css_class("image-button");
        button.set_tooltip_text(tooltip);
        const stack = new Gtk.Stack();
        stack.set_transition_type(Gtk.StackTransitionType.CROSSFADE);
        stack.set_transition_duration(150);
        const buttonIcon = Gtk.Image.new_from_icon_name(icon);
        const spinner = new Gtk.Spinner();
        stack.add_named(buttonIcon, "icon");
        stack.add_named(spinner, "spinner");
        stack.set_visible_child_name("icon");
        button.set_child(stack);

        return { button, spinner, stack };
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
