# Upgrade Standard Dry-run

> gate_id: upgrade-standard-dry-run
> type: hybrid
> severity: P2
> status: active

## 描述

`goalcli upgrade-standard --dry-run` 只验证合同文件存在，不提升版本或写 downstream。

## 验证命令

```bash
goalcli upgrade-standard --dry-run
```

## 通过条件

- standard 合同文件存在
- 文件格式符合 schema
- 版本号未降级

## 失败条件

- 合同文件缺失
- 格式不符合 schema
- 版本号降级

## Evidence 要求

- 命令输出
- 合同文件 sha256
- schema 校验结果

## 关联规则

- RULE-UPGRADE-STD-001
- RULE-CORE-001
