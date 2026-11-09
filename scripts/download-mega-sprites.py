#!/usr/bin/env python3
"""
Download missing mega evolution sprites from PokeAPI and save them
in pokedextracker/pokesprite naming format (NNN-form.png).

Uses the fit_image() approach from pokesprite-v2 to resize official art
to 68x56 box sprite format.

Usage:
    python3 scripts/download-mega-sprites.py
"""

from __future__ import annotations

import io
import json
import sys
import urllib.request
import urllib.error
from pathlib import Path
from typing import Optional

from PIL import Image

IMAGES_DIR = Path(__file__).resolve().parents[1] / "images"

# PokeAPI slug → (national_id, pokesprite form suffix)
# Only includes sprites MISSING from the current spritesheet
MISSING_SPRITES = [
    # Z-A base game megas (missing from spritesheet)
    ("meganium-mega",       154, "mega"),
    ("emboar-mega",         500, "mega"),
    ("feraligatr-mega",     160, "mega"),
    ("clefable-mega",        36, "mega"),
    ("victreebel-mega",      71, "mega"),
    ("starmie-mega",        121, "mega"),
    ("dragonite-mega",      149, "mega"),
    ("skarmory-mega",       227, "mega"),
    ("scolipede-mega",      545, "mega"),
    ("barbaracle-mega",     689, "mega"),
    ("floette-mega",        670, "mega"),
    ("pyroar-mega",         668, "mega"),
    ("eelektross-mega",     604, "mega"),
    ("scrafty-mega",        560, "mega"),
    ("chandelure-mega",     609, "mega"),
    ("greninja-mega",       658, "mega"),
    ("falinks-mega",        870, "mega"),
    ("chesnaught-mega",     652, "mega"),
    ("delphox-mega",        655, "mega"),
    ("hawlucha-mega",       701, "mega"),
    ("malamar-mega",        687, "mega"),
    ("dragalge-mega",       691, "mega"),
    ("froslass-mega",       478, "mega"),
    ("drampa-mega",         780, "mega"),
    ("zygarde-mega",        718, "mega"),
    ("excadrill-mega",      530, "mega"),
    # Z-A Mega Dimension DLC megas
    ("meowstic-mega",       678, "mega"),
    ("raichu-mega-x",        26, "mega-x"),
    ("raichu-mega-y",        26, "mega-y"),
    ("garchomp-mega-z",     445, "mega-z"),
    ("absol-mega-z",        359, "mega-z"),
    ("lucario-mega-z",      448, "mega-z"),
    ("scovillain-mega",     952, "mega"),
    ("glimmora-mega",       970, "mega"),
    ("tatsugiri-mega",      978, "mega"),
    ("baxcalibur-mega",     998, "mega"),
    ("chimecho-mega",       358, "mega"),
    ("golisopod-mega",      768, "mega"),
    ("golurk-mega",         623, "mega"),
    ("staraptor-mega",      398, "mega"),
    ("crabominable-mega",   740, "mega"),
    ("heatran-mega",        485, "mega"),
    ("darkrai-mega",        491, "mega"),
    ("magearna-mega",       801, "mega"),
    ("zeraora-mega",        807, "mega"),
    # Primals (sprites exist in CSS but let's ensure we have good ones)
    ("kyogre-primal",       382, "primal"),
    ("groudon-primal",      383, "primal"),
]


def fit_image(image: Image.Image) -> Image.Image:
    """Resize official art to 68x56 box sprite format.
    Same logic as pokesprite-v2's build_assets.py."""
    frame = Image.new("RGBA", (68, 56), (0, 0, 0, 0))
    working = image.convert("RGBA")
    working.thumbnail((64, 52), Image.Resampling.LANCZOS)
    x = (68 - working.width) // 2
    y = 56 - working.height - 2
    y = max(0, y)
    frame.alpha_composite(working, (x, y))
    return frame


def download_sprite(slug: str, shiny: bool = False) -> Optional[bytes]:
    """Download a sprite from PokeAPI, resize to box sprite format."""
    try:
        url = f"https://pokeapi.co/api/v2/pokemon/{slug}"
        req = urllib.request.Request(url, headers={"User-Agent": "PokeTrack-Sprites/1.0"})
        with urllib.request.urlopen(req, timeout=30) as resp:
            data = json.loads(resp.read().decode())
    except Exception as e:
        print(f"  ERROR fetching API data for {slug}: {e}", file=sys.stderr)
        return None

    sprites = data.get("sprites", {})
    other = sprites.get("other", {})
    home = other.get("home", {})
    artwork = other.get("official-artwork", {})

    if shiny:
        sprite_url = home.get("front_shiny") or artwork.get("front_shiny")
    else:
        sprite_url = home.get("front_default") or artwork.get("front_default")

    if not sprite_url:
        print(f"  WARNING: No {'shiny' if shiny else 'regular'} sprite URL for {slug}", file=sys.stderr)
        return None

    try:
        req = urllib.request.Request(sprite_url, headers={"User-Agent": "PokeTrack-Sprites/1.0"})
        with urllib.request.urlopen(req, timeout=30) as resp:
            raw = resp.read()
        image = Image.open(io.BytesIO(raw))
        fitted = fit_image(image)
        out = io.BytesIO()
        fitted.save(out, format="PNG", optimize=True)
        return out.getvalue()
    except Exception as e:
        print(f"  ERROR processing sprite for {slug}: {e}", file=sys.stderr)
        return None


def main():
    print(f"Downloading {len(MISSING_SPRITES)} mega/primal sprites...")
    print(f"Output directory: {IMAGES_DIR}")

    downloaded = 0
    skipped = 0
    failed = 0

    for pokeapi_slug, nat_id, form in MISSING_SPRITES:
        dex_str = f"{nat_id:03d}"

        # Regular sprite
        regular_name = f"{dex_str}-{form}.png"
        regular_path = IMAGES_DIR / regular_name

        if regular_path.exists():
            print(f"  SKIP {regular_name} (exists)")
            skipped += 1
        else:
            sprite_data = download_sprite(pokeapi_slug, shiny=False)
            if sprite_data:
                regular_path.write_bytes(sprite_data)
                print(f"  OK   {regular_name} ({len(sprite_data)} bytes)")
                downloaded += 1
            else:
                print(f"  FAIL {regular_name}")
                failed += 1

        # Shiny sprite
        shiny_name = f"{dex_str}-shiny-{form}.png"
        shiny_path = IMAGES_DIR / shiny_name

        if shiny_path.exists():
            print(f"  SKIP {shiny_name} (exists)")
            skipped += 1
        else:
            sprite_data = download_sprite(pokeapi_slug, shiny=True)
            if sprite_data:
                shiny_path.write_bytes(sprite_data)
                print(f"  OK   {shiny_name} ({len(sprite_data)} bytes)")
                downloaded += 1
            else:
                print(f"  FAIL {shiny_name}")
                failed += 1

    print(f"\nDone: {downloaded} downloaded, {skipped} skipped, {failed} failed")


if __name__ == "__main__":
    main()
