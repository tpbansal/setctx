# SafeCtx

**SafeCtx** is a secure, pluggable gateway for validating and filtering LLM-related JSON-RPC / MCP requests before they reach downstream AI agents or services.

---

## Features

- 🔐 **Context-Aware Policy Enforcement** (OPA / Rego or custom rules)
- 🧼 **Context Redaction & Sanitization** (passwords, API keys, etc.)
- 🧠 **Prompt Injection Detection** (regex & future ML-based checks)
- 📏 **JSON Schema Validation** for incoming requests
- 🚀 **Plug-and-Play Proxy** for existing LLM tools or MCP APIs

---

## Project Structure

```
safectx/
├── cmd/
│   └── safectx/              # Entry point (main.go)
│       └── main.go
│
├── internal/
│   ├── middleware/           # HTTP middleware: logging, auth, rate-limit
│   │   ├── logging.go
│   │   ├── auth.go
│   │   └── recovery.go
│
│   ├── policy/               # OPA/Rego or custom rule engine
│   │   ├── engine.go         # Interface
│   │   ├── opa.go            # OPA/Rego integration
│   │   └── evaluator.go      # DSL rules or wrappers
│
│   ├── rpc/                  # JSON-RPC request validation / proxy logic
│   │   ├── handler.go
│   │   ├── proxy.go
│   │   └── schema.go
│
│   ├── detection/            # Prompt injection, bad pattern detection
│   │   ├── patterns.go
│   │   └── embeddings.go     # Future: ML-based prompt classifier
│
│   └── contextfilter/        # Redaction, mutation, context shaping
│       ├── redactor.go
│       └── sanitizer.go
│
├── pkg/
│   └── schema/               # JSON schemas, Go types
│       ├── request.go
│       └── schema_loader.go
│
├── testdata/                 # JSON test payloads and policy fixtures
│   ├── valid_input.json
│   ├── invalid_input.json
│   └── policies/
│       └── block_drop.rego
│
├── go.mod
└── README.md

```

---

## Run Locally

```bash
git clone https://github.com/yourname/safectx
cd safectx
go run ./cmd/safectx
```

Server runs at `http://localhost:8080` and expects JSON-RPC/MCP-style POST payloads.

---

## Roadmap

- [ ] Implement actual reverse proxy logic to MCP endpoints
- [ ] Add fine-grained rate limits per tool
- [ ] Extend redactors with LLM-based anomaly detection
- [ ] OpenAPI/JSON Schema generation for tools
- [ ] MAYBE add support for WASM plugin rules

---

## License

MIT. Open to contributions.

---

## Credits

Originally inspired by secure agent gateway patterns and the need for robust guardrails in AI ecosystems. Designed for open LLM integrations at scale.
