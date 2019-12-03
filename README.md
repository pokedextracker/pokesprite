# pokesprite

Inspired by [msikma/pokesprite](https://github.com/msikma/pokesprite), this repo
is a simplified version that generates a spritesheet and stylesheet that are
meant to be used for [PokedexTracker](https://pokedextracker.com).

It currently consists of 3 scripts:

- `rename` - This renames icons from
  [msikma/pokesprite](https://github.com/msikma/pokesprite) to names that can be
  used by the other scripts. Only use this one if you're copying over the
  `icons` and `data` directories from that repo.
- `spritesheet` - This takes all the images in the `images` directory and
  stitches them together into a single image.
- `scss` - This uses the images in the `images` directory to generate a `.scss`
  file that lists classes with the correct positions so the spritesheet can be
  used.

To run any of them, it's a simple `make` target:

```sh
make rename
make spritesheet
make scss
```
