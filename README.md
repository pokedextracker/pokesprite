# pokesprite

Inspired by [msikma/pokesprite](https://github.com/msikma/pokesprite), this repo
is a simplified version that generates a spritesheet and stylesheet that are
meant to be used for [PokedexTracker](https://pokedextracker.com).

It currently consists of 6 scripts:

- `rename` - This renames icons from
  [msikma/pokesprite](https://github.com/msikma/pokesprite) to names that can be
  used by the other scripts. Only use this one if you're copying over the
  `icons` and `data` directories from that repo.
- `scale` - This takes any images in the `images` directory that are greater
  than 100px in either dimension (height or width) and scales it by factor or
  0.5. This script will modify the images in place.
- `trim` - This takes all images in the `images` directory and trims any excess
  transparency from it. This is so that we can center the sprites based on
  content (non-transparent pixels) and control the padding through CSS.
- `spritesheet` - This takes all the images in the `images` directory and
  stitches them together into a single image.
- `scss` - This uses the images in the `images` directory to generate a `.scss`
  file that lists classes with the correct positions so the spritesheet can be
  used.
- `copy` - This takes the final outputs (the spritesheet and the `.scss` file)
  and copies them into their appropriate location in
  [pokedextracker/pokedextracker.com](https://github.com/pokedextracker/pokedextracker.com).
  It assumes that this repo and that repo are both cloned in the same parent
  directory. If that is not the case, this script will err.

To run any of them, it's a simple `task` command:

```sh
task rename
task scale
task trim
task spritesheet
task scss
task copy
```

## Setup

### Task

Instead of `make`, this project uses [`task`](https://taskfile.dev/#/). It seems
to be a bit cleaner for some specific things that we want to do.

You can find instructions on how to install it
[here](https://taskfile.dev/#/installation).

### Go

To have everything working as expected, you need to have a module-aware version
of Go installed (v1.11 or greater) and `pngcrush`.

To install Go, you can do it any way you prefer. We recommend using
[`goenv`](https://github.com/syndbg/goenv) so that you can use the correct
version of Go for different projects depending on `.go-version` files. In its
current state, the v2 beta of `goenv` can't be installed through `brew`
normally, so you need to fetch from `HEAD` using the following command:

```sh
brew install --HEAD goenv
```

**Note**: If you already have a v1 version of `goenv` installed, you need to
uninstall it first.

Once installed, you can go into this projects directory and run the following to
install the correct version of Go:

```sh
goenv install
```

### `pngcrush`

`pngcrush` is required for the `spritesheet` command. To install it, you can
just run the following command:

```sh
task setup
```
