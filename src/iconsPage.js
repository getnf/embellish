import GObject from "gi://GObject";
import Adw from "gi://Adw";
import Gtk from "gi://Gtk";
import Gio from "gi://Gio";
import Pango from "gi://Pango";
import Gdk from "gi://Gdk";
import GLib from "gi://GLib";

export const EmbIconsPage = GObject.registerClass(
    {
        GTypeName: "EmbIconsPage",
        Template: "resource:///io/github/getnf/embellish/ui/IconsPage.ui",
        InternalChildren: [
            "searchBar",
            "searchEntry",
            "iconsStack",
            "sectionsContainer",
        ],
    },
    class extends Adw.Bin {
        constructor(params = {}) {
            super(params);
            this.icons = [];
            this.filteredIcons = [];
            this.iconSections = new Map();
            this.visibleSections = new Set();
            this.searchTimeout = null;
            this.observedSections = new Map();
            this.#setupSearch();
            this.#setupIntersectionObserver();
            this.#loadIcons();

            // Clean up timeout when dialog is destroyed
            this.connect("destroy", () => {
                this.#clearSearchTimeout();
                this.#cleanupObserver();
            });
        }

        #setupSearch() {
            this._searchEntry.connect("search-changed", () => {
                // Clear existing timeout if it exists
                this.#clearSearchTimeout();

                // Debounce search to improve performance
                this.searchTimeout = GLib.timeout_add(
                    GLib.PRIORITY_DEFAULT,
                    150,
                    () => {
                        // Mark timeout as completed so we don't try to remove it
                        const currentTimeout = this.searchTimeout;
                        this.searchTimeout = null;

                        this.#filterIcons();
                        this.#populateAllSections();
                        return GLib.SOURCE_REMOVE;
                    },
                );
            });

            // Enable search mode by default
            this._searchBar.search_mode_enabled = true;
        }

        #clearSearchTimeout() {
            if (this.searchTimeout !== null) {
                try {
                    GLib.source_remove(this.searchTimeout);
                } catch (error) {
                    // Timeout may have already been removed or completed
                    console.debug(
                        "Search timeout already removed:",
                        error.message,
                    );
                }
                this.searchTimeout = null;
            }
        }

        #setupIntersectionObserver() {
            // Create intersection observer to implement lazy loading
            // Since GTK doesn't have IntersectionObserver, we'll use scroll events
            const scrolledWindow = this._sectionsContainer.get_ancestor(
                Gtk.ScrolledWindow,
            );
            if (scrolledWindow) {
                this.scrolledWindow = scrolledWindow;
                const adjustment = scrolledWindow.get_vadjustment();
                adjustment.connect("value-changed", () => {
                    this.#checkVisibleSections();
                });
            }
        }

        #cleanupObserver() {
            // Cleanup any remaining observers or timeouts
            this.observedSections.clear();
            this.visibleSections.clear();
            this.scrolledWindow = null;
        }

        async #loadIcons() {
            try {
                // Show loading state
                this._iconsStack.set_visible_child_name("loadingState");

                // Load the cheatsheet CSV from resources
                const resource = Gio.resources_lookup_data(
                    "/io/github/getnf/embellish/cheatsheet.csv",
                    Gio.ResourceLookupFlags.NONE,
                );

                if (!resource) {
                    throw new Error("Could not load cheatsheet.csv resource");
                }

                const csvData = new TextDecoder().decode(resource.get_data());
                this.#parseCSV(csvData);
                this.#organizeBySections();
                this.#filterIcons(); // Initialize filtered icons
                this.#populateAllSections();
            } catch (error) {
                console.error("Failed to load cheatsheet:", error);
                this._iconsStack.set_visible_child_name("emptyState");
            }
        }

        #parseCSV(csvContent) {
            const lines = csvContent.trim().split("\n");
            this.icons = [];

            for (const line of lines) {
                const parts = line.split(",");
                if (parts.length >= 3) {
                    const name = parts[0].trim();
                    const unicode = parts[1].trim();
                    const icon = parts[2]; // The actual icon character

                    if (name && unicode) {
                        const category = this.#getCategoryFromName(name);
                        this.icons.push({
                            name: name,
                            unicode: unicode.toUpperCase(),
                            icon: icon || "",
                            category: category,
                        });
                    }
                }
            }
        }

        #getCategoryFromName(name) {
            // Extract category from icon name prefix (e.g., "nf-fa" from "nf-fa-home")
            const parts = name.split("-");
            if (parts.length >= 2) {
                return `${parts[0]}-${parts[1]}`;
            }
            return "other";
        }

        #organizeBySections() {
            this.iconSections.clear();

            // Group icons by category
            for (const icon of this.icons) {
                if (!this.iconSections.has(icon.category)) {
                    this.iconSections.set(icon.category, []);
                }
                this.iconSections.get(icon.category).push(icon);
            }

            // Sort sections by name and sort icons within each section
            for (const [category, icons] of this.iconSections) {
                icons.sort((a, b) => a.name.localeCompare(b.name));
            }
        }

        #filterIcons() {
            const searchText = this._searchEntry.text.toLowerCase();

            if (searchText === "") {
                this.filteredIcons = [...this.icons];
            } else {
                this.filteredIcons = this.icons.filter(
                    (icon) =>
                        icon.name.toLowerCase().includes(searchText) ||
                        icon.unicode.toLowerCase().includes(searchText) ||
                        this.#getCategoryTitle(icon.category)
                            .toLowerCase()
                            .includes(searchText),
                );
            }

            // Sort filtered results by relevance
            if (searchText !== "") {
                this.filteredIcons.sort((a, b) => {
                    const aNameMatch = a.name
                        .toLowerCase()
                        .startsWith(searchText)
                        ? 1
                        : 0;
                    const bNameMatch = b.name
                        .toLowerCase()
                        .startsWith(searchText)
                        ? 1
                        : 0;

                    if (aNameMatch !== bNameMatch) {
                        return bNameMatch - aNameMatch; // Prioritize name matches
                    }

                    return a.name.localeCompare(b.name); // Alphabetical fallback
                });
            }
        }

        #populateAllSections() {
            // Clear existing children
            this.#clearSections();

            if (this.filteredIcons.length === 0) {
                this._iconsStack.set_visible_child_name("emptyState");
                return;
            }

            this._iconsStack.set_visible_child_name("iconGrid");

            // Group filtered icons by category
            const filteredSections = new Map();
            for (const icon of this.filteredIcons) {
                if (!filteredSections.has(icon.category)) {
                    filteredSections.set(icon.category, []);
                }
                filteredSections.get(icon.category).push(icon);
            }

            // Create sections in sorted order
            const sortedCategories = Array.from(filteredSections.keys()).sort();

            for (const category of sortedCategories) {
                const icons = filteredSections.get(category);
                this.#createSection(category, icons);
            }
        }

        #clearSections() {
            let child = this._sectionsContainer.get_first_child();
            while (child) {
                const next = child.get_next_sibling();
                this._sectionsContainer.remove(child);
                child = next;
            }
            this.visibleSections.clear();
            this.observedSections.clear();
        }

        #createSection(category, icons) {
            // Create section container
            const sectionBox = new Gtk.Box({
                orientation: Gtk.Orientation.VERTICAL,
                spacing: 12,
            });

            // Create section header
            const headerBox = new Gtk.Box({
                orientation: Gtk.Orientation.HORIZONTAL,
                spacing: 8,
            });
            headerBox.add_css_class("icon-section-header");

            const categoryTitle = this.#getCategoryTitle(category);
            const headerLabel = new Gtk.Label({
                label: `${categoryTitle}`,
                halign: Gtk.Align.START,
                hexpand: true,
            });
            headerLabel.add_css_class("heading");

            headerBox.append(headerLabel);
            sectionBox.append(headerBox);

            // Create flow box for icons in this section
            const flowBox = new Gtk.FlowBox({
                selection_mode: Gtk.SelectionMode.NONE,
                homogeneous: true,
                column_spacing: 6,
                row_spacing: 6,
                min_children_per_line: 3,
                max_children_per_line: 6,
            });

            // For better performance, load different amounts based on section size
            let initialCount = Math.min(30, icons.length);
            if (icons.length > 100) {
                initialCount = 20; // Smaller initial load for large sections
            }

            const initialIcons = icons.slice(0, initialCount);
            const remainingIcons = icons.slice(initialCount);

            // Add initial icons
            for (const icon of initialIcons) {
                const iconWidget = this.#createIconWidget(icon);
                flowBox.append(iconWidget);
            }

            // Add "Load More" button if there are remaining icons
            if (remainingIcons.length > 0) {
                const loadMoreButton = new Gtk.Button({
                    label: _("Load %d more icons").format(remainingIcons.length),
                    halign: Gtk.Align.CENTER,
                    margin_top: 12,
                });
                loadMoreButton.add_css_class("pill");

                loadMoreButton.connect("clicked", () => {
                    // Load remaining icons in batches for better performance
                    const batchSize = 50;
                    let loadedCount = 0;

                    const loadBatch = () => {
                        const batch = remainingIcons.slice(
                            loadedCount,
                            loadedCount + batchSize,
                        );
                        for (const icon of batch) {
                            const iconWidget = this.#createIconWidget(icon);
                            flowBox.append(iconWidget);
                        }
                        loadedCount += batch.length;

                        if (loadedCount >= remainingIcons.length) {
                            // All icons loaded, remove button
                            sectionBox.remove(loadMoreButton);
                        } else {
                            // Update button text
                            const remaining =
                                remainingIcons.length - loadedCount;
                            loadMoreButton.set_label(
                                _("Load %d more icons...").format(remaining),
                            );

                            // Schedule next batch with a small delay to keep UI responsive
                            GLib.timeout_add(GLib.PRIORITY_DEFAULT, 10, () => {
                                loadBatch();
                                return GLib.SOURCE_REMOVE;
                            });
                        }
                    };

                    loadBatch();
                });

                sectionBox.append(flowBox);
                sectionBox.append(loadMoreButton);
            } else {
                sectionBox.append(flowBox);
            }

            this._sectionsContainer.append(sectionBox);
        }

        #getCategoryTitle(category) {
            const categoryMap = {
                "nf-fa": "Font Awesome",
                "nf-md": "Material Design",
                "nf-cod": "Codicons",
                "nf-dev": "Devicons",
                "nf-oct": "Octicons",
                "nf-fae": "Font Awesome Extension",
                "nf-linux": "Linux",
                "nf-weather": "Weather",
                "nf-seti": "Seti UI",
                "nf-custom": "Custom",
                "nf-pl": "Powerline",
                "nf-ple": "Powerline Extra",
                "nf-pom": "Pomicons",
                "nf-iec": "IEC Power",
                "nf-extra": "Extra",
                "nf-indent": "Indentation",
                "nf-indentation": "Indentation",
            };
            return categoryMap[category] || category.toUpperCase();
        }

        #checkVisibleSections() {
            // Simple implementation for scroll position tracking
            // This could be enhanced to automatically load sections when they come into view
            if (!this.scrolledWindow) return;

            const adjustment = this.scrolledWindow.get_vadjustment();
            const scrollPosition = adjustment.get_value();
            const windowHeight = adjustment.get_page_size();
            const totalHeight = adjustment.get_upper();

            // Scroll to top functionality - could be used for a scroll-to-top button
            this.scrollPosition = {
                current: scrollPosition,
                total: totalHeight,
                visible: windowHeight,
                atTop: scrollPosition < 50,
                atBottom: scrollPosition > totalHeight - windowHeight - 50,
            };
        }

        #createIconWidget(icon) {
            const button = new Gtk.Button({
                label: icon.icon,
                tooltip_text: `${icon.name}\nUnicode: ${icon.unicode}`,
                // Translated version, disabled as it makes tooltip alignment weird.
                // tooltip_text: _("%s\nUnicode: %s").format(icon.name, icon.unicode),
            });

            button.add_css_class("flat");
            button.add_css_class("icon-button");

            button.connect("clicked", () => {
                this.#copyToClipboard(icon.icon, button);
            });

            // Add keyboard support for copying
            button.connect("activate", () => {
                this.#copyToClipboard(icon.icon, button);
            });

            const buttonContainer = new Adw.Bin({
                child: button,
                halign: Gtk.Align.CENTER,
                width_request: 60,
                height_request: 60,
            });

            return buttonContainer;
        }

        #copyToClipboard(text, button = null) {
            const display = this.get_display
                ? this.get_display()
                : Gdk.Display.get_default();
            const clipboard = display.get_clipboard();
            clipboard.set(text);

            // Add visual feedback to the button
            if (button) {
                GLib.timeout_add(GLib.PRIORITY_DEFAULT, 600, () => {
                    return GLib.SOURCE_REMOVE;
                });
            }

            // Try to show a toast notification
            this.#showToast(_('Copied "%s" to clipboard').format(text));
        }

        #showToast(message) {
            const toast = new Adw.Toast({
                title: message,
                timeout: 2,
            });

            // Try to find a toast overlay in the application window
            let parent = this.get_parent();
            while (parent) {
                if (parent._toastOverlay) {
                    parent._toastOverlay.add_toast(toast);
                    return;
                }
                parent = parent.get_parent();
            }

            // Try to find the main window and its toast overlay
            const application = this.get_root ? this.get_root() : null;
            if (application && application._toastOverlay) {
                application._toastOverlay.add_toast(toast);
                return;
            }
        }
    },
);


