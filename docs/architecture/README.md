# Architecture

本目录是架构文档入口，位于权威顺序中 `contracts/` 之后、`AGENTS.md` 之前。

当前标准分层：

```text
xlib-standard
    ↓
L0: kernel
    ↓
L1: configx / observex / testkitx / resiliencx / schedulex
    ↓
L2: redisx / kafkax / postgresx / taosx / ossx / clickhousex / natsx
    ↓
L3+: x.go / market-data / macro-data / regime-engine / business systems
```

现有架构证据优先从以下位置恢复：

- `docs/standard/`
- `contracts/`
- `.agent/policies/layer-governance.yaml`
- `.agent/contracts/scope-locks.yaml`
- `docs/adr/`

新增架构变更应补充 ADR、影响范围、分层边界和对应 harness/evidence。
