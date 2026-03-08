#!/usr/bin/env bash

set -e

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 FONT [FONT ...]"
    echo "Example: $0 *.ttf"
    exit 1
fi

GLYPHS="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789{}[]()<>;:.,+-*/%&|^~!?@#\$'\"\\"

for font in "$@"; do
    if [ ! -f "$font" ]; then
        echo "Skipping '$font' (not a file)"
        continue
    fi

    filename=$(basename -- "$font")
    name="${filename%.*}"
    ext="${filename##*.}"

    output="${name}-subset.${ext}"

    echo "Subsetting $font → $output"

    pyftsubset "$font" \
        --text="$GLYPHS" \
        --output-file="$output" \
        --layout-features='*' \
        --glyph-names \
        --recommended-glyphs \
        --symbol-cmap
done

echo "Done."
