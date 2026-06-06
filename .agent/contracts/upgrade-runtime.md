# Upgrade Runtime Dry-run

> gate_id: upgrade-runtime-dry-run
> type: hybrid
> severity: P2
> status: active

## 描述

本仓库只验证升级合同，不修改 downstream 仓库。

## 验证命令

```bash
goalcli upgrade-runtime --dry-run
```

## 通过条件

- upgrade 合同文件存在
- 版本兼容性检查通过
- 下游 adoption_claim 声明存在

## 失败条件

- 合同文件缺失
- 版本不兼容
- 缺少 adoption_claim 声明

## Evidence 要求

- 命令输出
- 合同文件 sha256
- 版本兼容性报告

## 关联规则

- RULE-UPGRADE-001
- RULE-CORE-001
