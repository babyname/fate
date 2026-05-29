## Project Configuration

- **Project Name**: fate
- **Team Version**: v3
- **Default Flow**: dev-flow

## Paths

docs_internal: _docs/fate/
docs_external: docs/
default_flow: dev-flow

## Toolchain

### Backend (Go)
pipeline: go test ./... | go build -o bin/app
working_dir: .

### Frontend (TypeScript/React)
pipeline: bun run test | bun run build
working_dir: web

## Constraints
- No Chinese comments in code
- TDD required
