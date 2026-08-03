#!/usr/bin/env bash

set -euo pipefail

for cmd in curl jq tar mktemp; do
    if ! command -v "$cmd" &>/dev/null; then
        echo "Error: '$cmd' is required but not installed." >&2
        exit 1
    fi
done

OUTPUT_DIR="./nerd-fonts-regular"
mkdir -p "$OUTPUT_DIR"

echo "Fetching latest Nerd Fonts release metadata..."
RELEASE_JSON=$(curl -sL "https://api.github.com/repos/ryanoasis/nerd-fonts/releases/latest")

URLS=$(echo "$RELEASE_JSON" | jq -r '.assets[] | select(.name | endswith(".tar.xz")) | .browser_download_url')

if [ -z "$URLS" ]; then
    echo "Error: Could not find any font archives in the release." >&2
    exit 1
fi

TOTAL=$(echo "$URLS" | wc -l | tr -d ' ')
CURRENT=1

echo "Extracting 1 primary regular font (preferring NerdFontMono -> NerdFont -> OTF) per family..."
echo "----------------------------------------------------------------------"

# Words to reject when identifying a base/regular weight
EXCLUDE_PATTERN="[Bb]old|[Ii]talic|[Oo]blique|[Ll]ight|[Tt]hin|[Ee]xtra|[Uu]ltra|[Ss]emi|[Dd]emio|[Cc]ondensed|[Ee]xpanded|[Vv]ariable|[Ww]indows"

for url in $URLS; do
    archive_name=$(basename "$url")
    font_family="${archive_name%.tar.xz}"

    tmp_tar=$(mktemp)
    curl -sL "$url" -o "$tmp_tar"

    # Get all candidate font files inside the archive (excluding bold/italic/light/etc.)
    all_files=$(tar -tJ -f "$tmp_tar" 2>/dev/null | grep -E '\.(ttf|otf)$' | grep -vE "$EXCLUDE_PATTERN" || true)

    target_file=""

    # --- TTF SEARCH ---
    # 1. Look specifically for explicit NerdFontMono or NFM variant with Regular
    target_file=$(echo "$all_files" | grep -E '\.ttf$' | grep -E 'NerdFontMono|NFM' | grep -i 'regular' | head -n 1 || true)

    # 2. Any NerdFontMono / NFM TTF
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.ttf$' | grep -E 'NerdFontMono|NFM' | head -n 1 || true)
    fi

    # 3. Standard Regular TTF (fallback if no explicit Mono variant exists in package)
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.ttf$' | grep -i 'regular' | head -n 1 || true)
    fi

    # 4. Any TTF
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.ttf$' | head -n 1 || true)
    fi

    # --- OTF FALLBACK (for fonts like CommitMono, GeistMono, Hermit, etc.) ---
    # 5. Look for explicit NerdFontMono / NFM OTF with Regular
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.otf$' | grep -E 'NerdFontMono|NFM' | grep -i 'regular' | head -n 1 || true)
    fi

    # 6. Any NerdFontMono / NFM OTF
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.otf$' | grep -E 'NerdFontMono|NFM' | head -n 1 || true)
    fi

    # 7. Standard Regular OTF
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.otf$' | grep -i 'regular' | head -n 1 || true)
    fi

    # 8. Any OTF
    if [ -z "$target_file" ]; then
        target_file=$(echo "$all_files" | grep -E '\.otf$' | head -n 1 || true)
    fi

    if [ -n "$target_file" ]; then
        echo "[$CURRENT/$TOTAL] $font_family -> $(basename "$target_file")"
        tar -xJ -f "$tmp_tar" -C "$OUTPUT_DIR" --transform='s|.*/||' "$target_file" 2>/dev/null || \
            echo "  --> Warning: Failed extracting '$target_file'"
    else
        echo "[$CURRENT/$TOTAL] $font_family..."
        echo "  --> Warning: No valid font file found in $archive_name"
    fi

    rm -f "$tmp_tar"
    ((CURRENT++))
done

echo "----------------------------------------------------------------------"
echo "Done! Saved $(ls -1 "$OUTPUT_DIR"/* 2>/dev/null | wc -l | tr -d ' ') regular fonts in '$OUTPUT_DIR'."
