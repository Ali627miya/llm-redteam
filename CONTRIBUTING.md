# Contributing

Thanks for helping improve **LLM Red Team**.

## How to contribute

1. **Issues** — bug reports, provider examples, attack-pack ideas, doc fixes.  
2. **Pull requests** — keep changes focused; match existing Go style and tests (`cd redteam && go test ./...`).  
3. **Attack packs** — add YAML under `redteam/pkg/attacks/library/` or share external packs via `-attacks-dir` (document the category stem).

## Attack YAML format

```yaml
category: my_category
attacks:
  - name: Short title
    description: Why this matters
    prompt: |
      Multi-line user prompt sent to the target
    tags: [injection]
```

## Code of conduct

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/) — see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## Security

Do not open public issues for vulnerabilities — see [SECURITY.md](SECURITY.md).
