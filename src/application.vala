/* application.vala
 *
 * Copyright 2025 Ronnie Nissan Yousif
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

public class Embellish.Application : Adw.Application {
    public Application () {
        Object (
            application_id: "io.github.getnf.embellish",
            flags: ApplicationFlags.DEFAULT_FLAGS,
            resource_base_path: "/io/github/getnf/embellish"
        );
    }

    construct {
        ActionEntry[] action_entries = {
            { "about", this.on_about_action },
            { "quit", this.quit }
        };
        this.add_action_entries (action_entries, this);
        this.set_accels_for_action ("app.quit", {"<primary>q"});
        this.set_accels_for_action("window.close", {"<Primary>w"});
        this.set_accels_for_action ("win.search", {"<primary>f"});
    }

    public override void activate () {
        base.activate ();
        var win = this.active_window ?? new Embellish.Window (this);
        win.present ();
    }

    private void on_about_action () {
        string[] developers = { "Ronnie Nissan https://ronnienissan.pages.dev/" };
		string[] designers = {"Brage Fuglseth https://bragefuglseth.dev"};
		string[] artists = {"Jakub Steiner https://jimmac.eu/",
                    "Brage Fuglseth https://bragefuglseth.dev",};
        var about = new Adw.AboutDialog () {
			application_name = _("Embellish"),
			comments = _(
                "Embellish helps you manage Nerd Fonts on your system"
            ),
			issue_url = "https://github.com/getnf/embellish/issues/new",
            license_type = Gtk.License.GPL_3_0,
            application_icon = "io.github.getnf/embellish",
			developer_name = "Ronnie Nissan",
            translator_credits = _("translator-credits"),
            version = Config.PACKAGE_VERSION,
            developers = developers,
			designers = designers,
			artists = artists,
            copyright = "© 2026 Ronnie Nissan",
        };

		about.add_other_app("io.github.ronniedroid.concessio", _("Concessio"), _("Understand file permissions"));
		about.add_other_app("io.github.sitraorg.sitra", _("Sitra"), _("Get fonts from online sources"));

        about.present (this.active_window);
    }
}
