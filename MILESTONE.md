# 📌 BrickEngine Milestones

This document outlines the feature roadmap for BrickEngine across different versions.

---

- [x] Variable declarations via `let`
- [x] `return` statements
- [x] Conditional blocks: `if`, `else if`, `else`
- [x] Function definition and invocation (`fn greet(name) { ... }`)
- [x] Pipe fallback syntax (`value | default`)
- [x] Object literals (`let x = { a: 1, b: 2 }`)
- [x] Array indexing (`array[0]`)
- [x] Variable assignment (`output.hostname = "new"`)
- [x] Built-in functions support
- [x] Test infrastructure + golden output comparisons
- [x] `for item in array {}` loops
- [x] `while (condition) {}` loops
- [x] Array literals: `[1, 2, 3]`
- [x] Boolean and null constants: `true`, `false`, `null`
- [x] `try { ... } catch { ... }` exception blocks
- [ ] `import "file.bee"` support
- [ ] `step "name" {}` block structure (useful for AxisDeploy)
- [x] Memory usage limit (heap counter)
- [x] Infinite loop / deadlock detection
- [x] `examples/benchmarks/` folder with fib, heavy-loop, etc.
- [ ] Step-based execution engine (`step "deploy" {}` → for CLI/Graph UI)

---

> This milestone file helps track progress and align development goals. For feedback or suggestions, contact: hello@isaeken.com.tr
