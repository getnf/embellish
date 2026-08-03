#!/usr/bin/env bash

set -e

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 FONT [FONT ...]"
    echo "Example: $0 *.ttf *.otf"
    exit 1
fi

GLYPHS="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789{}[]()<>;:.,+-*/%&|^~!?@#\$'\"\\"

OUTPUT_DIR="./subsetted"
mkdir -p "$OUTPUT_DIR"

FAILED_FONTS=()

for font in "$@"; do
    if [ ! -f "$font" ]; then
        echo "Skipping '$font' (not a file)"
        continue
    fi

    filename=$(basename -- "$font")
    name="${filename%.*}"
    ext="${filename##*.}"

    output="${OUTPUT_DIR}/${name}.${ext}"

    echo "Subsetting $font → $output"

    # Pass 1: Standard subsetting
    if pyftsubset "$font" \
        --text="$GLYPHS" \
        --output-file="$output" \
        --layout-features='*' \
        --glyph-names \
        --recommended-glyphs \
        --symbol-cmap 2>/dev/null; then
        continue
    fi

    echo "  --> Warning: OpenType tables failed in '$font'. Retrying without layout/BASE tables..."

    # Pass 2: Fallback explicitly dropping corrupt BASE and layout tables
    if pyftsubset "$font" \
        --text="$GLYPHS" \
        --output-file="$output" \
        --drop-tables+=BASE,GSUB,GPOS,GDEF \
        --layout-features='' \
        --glyph-names \
        --recommended-glyphs \
        --symbol-cmap 2>/dev/null; then
        echo "  --> SUCCESS (Corrupt layout/BASE tables dropped)"
        continue
    fi

    echo "  --> FAILED: Skipping '$font'"
    FAILED_FONTS+=("$font")
    rm -f "$output"
done

echo "----------------------------------------------------------------------"
echo "Done."
if [ ${#FAILED_FONTS[@]} -gt 0 ]; then
    echo "The following fonts could not be subsetted:"
    printf ' - %s\n' "${FAILED_FONTS[@]}"
fi
