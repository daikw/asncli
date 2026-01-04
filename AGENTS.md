# AGENTS.md

## Style

Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/):

- [Guide](https://google.github.io/styleguide/go/guide): Hard rules for formatting and structure.
- [Decisions](https://google.github.io/styleguide/go/decisions): Decisions on common style points.
- [Best Practices](https://google.github.io/styleguide/go/best-practices): Patterns that solve common problems well.

Additional:

- Keep APIs small, errors wrapped with context, names clear.
- Prefer readability over cleverness. Keep functions focused.

## After Code Changes

Run these commands after modifying Go code:

```bash
# Format (uses gofumpt: https://github.com/mvdan/gofumpt)
gofumpt -w .

# Vet
go vet ./...

# Test
go test ./...
```

Fix any issues before considering the task complete.
