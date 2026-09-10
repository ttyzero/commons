# Agents guide for commons

Shared library for ttyzero companion TUIs.

- Go 1.25, module `github.com/ttyzero/commons`
- `theme` — look (palette, borders, `theme` bus commands)
- `files` — `files` payload parse
- `connect` — ttybus dial + files/theme subscriptions

```sh
go test ./...
```

Do not add app-specific UI here. New shared look or bus helpers belong
in this module so gitwing, peek, and buscope stay in sync.
