## Why "spectr"?

You might wonder: does the name of your specs folder and CLI tool actually matter? We thought the same thing, so we tested it systematically across multiple AI models (Claude, GPT, and others) to find what works best.

### Alternatives Evaluated

We considered several naming approaches:
- `specs/` - Common but generic; easily confused with documentation folders
- `specifications/` - Descriptive but verbose; slower to type in CLI workflows
- `requirements/` - Often associated with waterfall methodologies; less distinct
- `docs/specs/` - Nested approach; less CLI-friendly
- **`spectr/` - What we chose**

### Why spectr Won

Testing across AI models revealed that `spectr` excels in several dimensions:

1. **Brevity**: Fast to type, easy to remember, works well as a CLI command
2. **Distinctiveness**: Unique enough to avoid conflicts with common folder names
3. **AI Compatibility**: Models consistently recognize and use it correctly in generated code
4. **Ergonomic**: Short enough to feel natural in commands like `spectr validate` and `spectr archive`

The name reflects the tool's core purpose—providing clear visibility into specifications and changes, like a spectrum analyzer revealing what IS and what SHOULD BE.
