# Contributing

## Commits

Use Conventional Commits.

Format:

```text
type(scope): short description
```

Examples:

```text
feat(session): add NewSession with cookie-based auth
fix(parse): handle None in domain string
docs(adr): add decision for zero external dependencies
refactor(domain): rename DomainItem to DomainNode
test(parse): cover OR and AND prefix operators
```

Rules:
- keep commits small and meaningful
- write commit messages in English
- prefer one focused change per commit
- use a scope that matches the main area you changed

Common types:
- `feat`
- `fix`
- `docs`
- `refactor`
- `test`
- `chore`
