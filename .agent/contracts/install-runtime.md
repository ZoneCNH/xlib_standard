# Install Runtime Dry-run

> gate_id: install-runtime-dry-run
> type: hybrid
> severity: P2
> status: active

## 描述

本仓库只验证安装合同，不执行全局安装或外部写入。

## 验证命令

```bash
goalcli install-runtime --dry-run
```

## 通过条件

- install 合同文件存在
- 目标路径可写
- 依赖版本兼容

## 失败条件

- 合同文件缺失
- 路径不可写
- 依赖版本冲突

## Evidence 要求

- 命令输出
- 合同文件 sha256

## 关联规则

- RULE-INSTALL-001
- RULE-CORE-001
