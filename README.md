# pokesprite

Inspired by [msikma/pokesprite](https://github.com/msikma/pokesprite), this repo
is a simplified version that generates spritesheets and a stylesheet that are
meant to be used for [PokedexTracker](https://pokedextracker.com).

It currently consists of 7 scripts:

- `rename` - This renames icons from
  [msikma/pokesprite](https://github.com/msikma/pokesprite) to names that can be
  used by the other scripts. Only use this one if you're copying sprites from
  that repo. Read the comment at the top of the file for more info.
- `chop` - This takes in a JSON file explaining the details of an existing
  spritesheet, and it chops it up into individual images.
- `scale` - This takes any images in the `images` directory that are greater
  than the threshold (default 100px) in either dimension (height or width) and
  either scales it by the provided factor (default 0.5) or to the set dimensions
  passed in. This script will modify the images in place.
- `trim` - This takes all images in the `images` directory and trims any excess
  transparency from it. This is so that we can center the sprites based on
  content (non-transparent pixels) and control the padding through CSS.
- `spritesheet` - This takes all the images in the `images` directory and
  stitches them together into lossless WebP spritesheets: one for regular
  icons, one for shiny icons, and a regular and shiny one for Legends: Arceus
  icons. A page only downloads the spritesheets it shows icons from. The love
  ball is copied over on its own instead, so the home page doesn't need a
  spritesheet.
- `css` - This uses the images in the `images` directory to generate a `.css`
  file that lists classes with the correct spritesheets and positions so the
  spritesheets can be used.
- `copy` - This takes the final outputs (the spritesheets, the love ball, and
  the `.css` file) and copies them into their appropriate location in
  [pokedextracker/pokedextracker.com](https://github.com/pokedextracker/pokedextracker.com).
  By default it copies into `../pokedextracker.com`, and it takes another path
  for a different checkout (see below).

To run any of them, use `mise run`:

```sh
mise run rename
mise run chop data.json
mise run scale
mise run trim
mise run spritesheet
mise run css
mise run copy
```

`copy` puts the files into `../pokedextracker.com` by default. To copy them into
a different checkout, like a worktree, pass its path:

```sh
mise run copy ~/path/to/pokedextracker.com-worktree
```

After adding or changing images, `mise run generate` runs `scale`, `trim`,
`spritesheet`, `css`, and `copy` in order, and takes the same optional path.
Running it again with no image changes produces the same files byte for byte.

## Setup

This project uses [mise](https://mise.jdx.dev) for its tools (Go and `cwebp`)
and tasks. Once mise is installed, run this in the repo to install the tools:

```sh
mise install
```
