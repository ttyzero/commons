# commons

Shared building blocks for [ttyzero](https://github.com/ttyzero) companion
panes: [gitwing](https://github.com/ttyzero/gitwing),
[peek](https://github.com/ttyzero/peek),
[buscope](https://github.com/ttyzero/buscope).

```
github.com/ttyzero/commons/theme    palette, $TTYTHEME, --borderless, theme bus
github.com/ttyzero/commons/files    files channel payload parse
github.com/ttyzero/commons/connect  dial ttybus, sub files+theme, Bubble Tea waits
```

## Theme bus

All three panes subscribe to ttybus channel `theme`:

```sh
ttybus pub theme 'THEME dark'
ttybus pub theme 'THEME nord'
ttybus pub theme 'THEME catppuccin BORDERS=0'
ttybus pub theme 'BORDERS=0'
```

Looks: `auto`, `dark`, `light` (Charm), then the ricing staples —
Catppuccin, Dracula, Nord, Gruvbox, Tokyo Night, Solarized, One Dark,
Rosé Pine, Everforest, Kanagawa (each with a light variant where the
upstream palette has one). `$TTYTHEME`, `$CLITHEME`, OSC 11,
`$COLORFGBG`, `$TTYBORDERLESS` — see `theme` package docs. Switch live
with [ttythemer](https://github.com/ttyzero/ttythemer).

## License

MIT.
