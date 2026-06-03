# goalkit v0.1.0 — `goal/` 迁移索引

> 本文件是审计索引，不是规范正文。规范正文见 `split/docs/standard/goalkit-runtime.md`；ADR 见 `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`；执行顺序见 `split/docs/plans/goalkit-v0.1.0-roadmap.md`。
> 生成时间：2026-06-03。
> 范围：`.worktree/goal/` 38 个文件入口，19,346 行；其中 4 组 checksum 重复。
> 状态：迁移台账；不是规范正文，也不是执行计划。

## 1. 迁移原则

- `goal/` 来源文件在迁移完成前不得删除；只能在本索引审核通过后归档到 `goal/archive/`。
- `goalkit-v0.1.0-plan.md` 是综合设计提案，不是权威替代。
- 权威职责分离：标准定义契约，ADR 记录决策，roadmap 记录执行顺序，本索引记录逐文件处置证据。
- 重复文件按 checksum 去重；重复来源仍逐行列出，便于审计。

## 2. 重复组

| canonical | duplicate | SHA-256 |
| --- | --- | --- |
| `goal/goalkit.md` | `goal/goalkit_v0_1_0_goal_runtime_complete_structural_plan.md` | `28558579d53f72f1dbb3f1f88e39ec2b1fde1a9e845b960bb2e58dc1662227c1` |
| `goal/xlib_standard_goal_runtime_v3_1_1_full_execution_pack_with_automation.md` | `goal/xlib_standard_goal_runtime_v3_1_1_full_execution_pack_with_automation (1).md` | `9167ea2205f7216948f6b2a958ce6e80d8d8d9a6a801eeb9bc2e63ddd1aa0ffc` |
| `goal/xlib_standard_pr13_goal_runtime_v3_1_1_runtime_as_code_spec_compiler_execution_pack.md` | `goal/xlib_standard_pr13_goal_runtime_v3_1_1_runtime_as_code_spec_compiler_execution_pack (1).md` | `d2b52334e09cd2224ebf57eb7c027bcb89775121de03e4f813cb7c0ba2378de6` |
| `goal/xlib_standard_pr8_goal_runtime_v3_1_1_blocking_full_mode_release_verify_execution_pack.md` | `goal/xlib_standard_pr8_goal_runtime_v3_1_1_blocking_full_mode_release_verify_execution_pack (1).md` | `46c1adca6a74ffe812df9ef96f21cb661d3d69e57614ba7513dd5310b59156d3` |

## 3. 来源文件映射

| 来源文件 | SHA-256 | 迁移目标 | 处置 |
|---|---|---|---|
| `goal/goal_runtime_v3_1_1_structural_refactor_plan.md` | `f16f697e2bee88090955be7d52c9618c36b5f7882c674cb19c7bc47f4618345b` | `split/docs/standard/goalkit-runtime.md`; `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md` | 保留参考：结构重构决策来源。 |
| `goal/goal_runtime_v3_1_1_structural_refactor_plan_v2_harness_runtime.md` | `07896823848c2e358957cb6464d8fd60210adf340c95bd17e2484700ec80d219` | `split/docs/standard/goalkit-runtime.md`; `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md` | 保留参考：结构重构决策来源。 |
| `goal/goalkit-analysis-20260603.md` | `63141ae403a809d2fdb6c783a037b0fcef2e3b3637c2e8aa69a3dd05c58eed74` | `.omx/reports/goalkit-v0.1.0-structural-issues-20260603.md`; `split/docs/plans/goalkit-v0.1.0-roadmap.md` | 保留参考：结构性问题与修复候选已迁移到报告和路线图。 |
| `goal/goalkit.md` | `28558579d53f72f1dbb3f1f88e39ec2b1fde1a9e845b960bb2e58dc1662227c1` | `split/docs/standard/goalkit-runtime.md`; `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`; `split/docs/plans/goalkit-v0.1.0-roadmap.md` | 归档：综合方案来源；重复内容已拆分，保留用于审计。 |
| `goal/goalkit_v0_1_0_goal_runtime_complete_structural_plan.md` | `28558579d53f72f1dbb3f1f88e39ec2b1fde1a9e845b960bb2e58dc1662227c1` | `split/docs/standard/goalkit-runtime.md`; `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`; `split/docs/plans/goalkit-v0.1.0-roadmap.md` | 归档：综合方案来源；重复内容已拆分，保留用于审计。 |
| `goal/xlib_standard_goal_runtime_v3_1_1_full_execution_pack_with_automation (1).md` | `9167ea2205f7216948f6b2a958ce6e80d8d8d9a6a801eeb9bc2e63ddd1aa0ffc` | `split/docs/standard/goalkit-runtime.md`; `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-12 | 保留参考：Full execution pack；重复文件按 checksum 去重。 |
| `goal/xlib_standard_goal_runtime_v3_1_1_full_execution_pack_with_automation.md` | `9167ea2205f7216948f6b2a958ce6e80d8d8d9a6a801eeb9bc2e63ddd1aa0ffc` | `split/docs/standard/goalkit-runtime.md`; `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-12 | 保留参考：Full execution pack；重复文件按 checksum 去重。 |
| `goal/xlib_standard_issue_pr_commit_release_automation_patch.md` | `c9b872fa553dc2cfb11037a207a62bc9347d16ec2882f4f1c6e401d225fb4dd1` | `split/docs/standard/goalkit-runtime.md` Automation Surface; `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-12 | 保留参考：自动化范围降级为后置 PR-12。 |
| `goal/xlib_standard_pr10_goal_runtime_v3_1_1_full_execution_pack_generator_execution_pack.md` | `87ce012ebd7600998b53b063210c2db1528fc0c242080262f5061fe7b19e5e9e` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-1 and PR-3 | 保留参考：generator/index 内容前移到基础文档阶段。 |
| `goal/xlib_standard_pr11_goal_runtime_v3_1_1_runtime_pack_ci_release_evidence_integration_execution_pack.md` | `c6b1d54568273911bba23d7749b0505b4fece1108fdcc1b6af02efe05e7792f2` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-5 and PR-8 | 保留参考：CI release evidence 与 runtime pack。 |
| `goal/xlib_standard_pr12_goal_runtime_v3_1_1_runtime_operations_lifecycle_governance_execution_pack.md` | `37d1283797e7cbd3de4efd141cfd838df81808e766d46d06438fb16f9d5173cf` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-8 | 保留参考：operations lifecycle governance。 |
| `goal/xlib_standard_pr13_goal_runtime_v3_1_1_runtime_as_code_spec_compiler_execution_pack (1).md` | `d2b52334e09cd2224ebf57eb7c027bcb89775121de03e4f813cb7c0ba2378de6` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-7 | 保留参考：runtime-as-code spec compiler；重复文件按 checksum 去重。 |
| `goal/xlib_standard_pr13_goal_runtime_v3_1_1_runtime_as_code_spec_compiler_execution_pack.md` | `d2b52334e09cd2224ebf57eb7c027bcb89775121de03e4f813cb7c0ba2378de6` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-7 | 保留参考：runtime-as-code spec compiler；重复文件按 checksum 去重。 |
| `goal/xlib_standard_pr14_goal_runtime_v3_1_1_downstream_adoption_orchestrator_execution_pack.md` | `a1f49a22d60d574967125d09e5a361c0abdd3be597454fcd777d9909824a52a3` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-9 | 保留参考：downstream adoption orchestrator。 |
| `goal/xlib_standard_pr15_goal_runtime_v3_1_1_ecosystem_adoption_dashboard_release_queue_execution_pack.md` | `e6acd1e9d9cc1d15023a69e3f7661e33a3a8e818a0f5d12039a2fe2a763b781b` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-9 and PR-10 | 保留参考：ecosystem dashboard/release queue 分拆到生态与观测阶段。 |
| `goal/xlib_standard_pr16_goal_runtime_v3_1_1_agent_team_orchestration_execution_pack.md` | `66442a0984eefcafeb126ecf3f4509978bf36176ad0c3a93eaa00a1ef490c826` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-11 | 保留参考：agent team orchestration 归入成熟化/DX。 |
| `goal/xlib_standard_pr17_goal_runtime_v3_1_1_self_improving_autoresearch_runtime_execution_pack.md` | `e2d558a030df5f98a14510c3ef2e6ad2006d10e0f38025af3c70dc63882e41fd` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-11 | 保留参考：autoresearch runtime 归入成熟化。 |
| `goal/xlib_standard_pr18_goal_runtime_v3_1_1_policy_as_code_runtime_execution_pack.md` | `03523f46c2c0dadff026f3d64522bab43c88850a0e9cacea26fc67b36ee283b3` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-7 | 保留参考：policy-as-code。 |
| `goal/xlib_standard_pr19_goal_runtime_v3_1_1_runtime_test_harness_golden_fixtures_execution_pack.md` | `ffed0cd92864f07360824473b11e9bb2df07766568259ad6df1191f5d62b013f` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-8 | 保留参考：test harness/golden fixtures。 |
| `goal/xlib_standard_pr1_goal_runtime_v3_1_1_templates_docs_execution_pack.md` | `0fbeaa7b25dbf244f2587d91e04c018fcae129c15a90400e59b7ea6edae68554` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-1; `split/docs/standard/goalkit-runtime.md` templates/docs sections | 保留参考：PR-1 执行包。 |
| `goal/xlib_standard_pr20_goal_runtime_v3_1_1_trust_root_evidence_attestation_execution_pack.md` | `3f92afabcfa582bf010af1ea4ff6529dc2f05dc9b69a93e23356f123ce048abd` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr21_goal_runtime_v3_1_1_runtime_observability_slo_execution_pack.md` | `20d30c2e1c27a874c0ef3398bd370ac30258ede7b59d495e713535ae1402900f` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr22_goal_runtime_v3_1_1_runtime_budget_anti_bloat_execution_pack.md` | `8958781b6ab47700af5039063b34a8f0962b9f438fe54f5418cf4faabad82034` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr23_goal_runtime_v3_1_1_runtime_dx_onboarding_execution_pack.md` | `89c9018348028a01c1b51add3558eb0970a506023947be8de02fc5076e9d9626` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr24_goal_runtime_v3_1_1_runtime_versioning_release_train_execution_pack.md` | `55c6d19b97aa846ffcf6dbd80c2eb2e71d992cbfb7ebbc9908884db5055843cd` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr25_goal_runtime_v3_1_1_runtime_conformance_certification_benchmark_execution_pack.md` | `e6322f4a393105e93d0f2522010b75941cb3c86fbe0c8bc3fd4bb21c58504f70` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr26_goal_runtime_v3_1_1_standard_publishing_adoption_kit_execution_pack.md` | `a5f76844f3e0a7015e25072857992394ac7a16743bc4c760998e11e731c08213` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr27_goal_runtime_v3_1_1_runtime_constitution_final_stop_conditions_execution_pack.md` | `c0883f04751a44371e96b455c79ede4fdc8cba221544eb41c735825115f41f11` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-10 and PR-11; `split/docs/standard/goalkit-runtime.md` maturity constraints | 保留参考：observability/budget/DX/versioning/conformance/publishing/constitution 成熟化内容已收敛。 |
| `goal/xlib_standard_pr28_goal_runtime_v3_1_1_issue_pr_commit_release_automation_execution_pack.md` | `2285f284cc11997692d99978dfa5e5fd435703a33e1a5ece33345cb35a2ab8a7` | `split/docs/standard/goalkit-runtime.md` Automation Surface; `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-12 | 保留参考：自动化范围降级为后置 PR-12。 |
| `goal/xlib_standard_pr2_goal_runtime_v3_1_1_schemas_output_contract_execution_pack.md` | `66fe8a2afa678727f967b3cfa854c153726b05a18fba28fd1c4824f6193dd054` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-2; `split/docs/standard/goalkit-runtime.md` output contract | 保留参考：PR-2 执行包。 |
| `goal/xlib_standard_pr3_goal_runtime_v3_1_1_runtime_index_compatibility_adr_execution_pack.md` | `86ea9e5e64820072a15f156d699eeaa2014547d7f8a9b45b63eb155d9b4e50ec` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-3; `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md` | 保留参考：PR-3 执行包。 |
| `goal/xlib_standard_pr4_goal_runtime_v3_1_1_makefile_harness_command_registry_execution_pack.md` | `5dc4328e53e8a5941849dea2143cbd90a8e2b52c165571ff2a067f50379e8acb` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-4 | 保留参考：PR-4 Makefile/harness slice。 |
| `goal/xlib_standard_pr5_goal_runtime_v3_1_1_xlibgate_commands_fixtures_tests_execution_pack.md` | `995e6d82b8052a0ab2c734d2539274eab035590c30085e7a2ae18360b20f48b1` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-4 | 保留参考：xlibgate 命令/fixture slice 合并到 PR-4。 |
| `goal/xlib_standard_pr6_goal_runtime_v3_1_1_release_manifest_generated_artifact_policy_execution_pack.md` | `aba4ce095c91eeef6671ed6c023f82b917f63aa8407e33f1c473a8ce3447b661` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-5 | 保留参考：artifact policy/blocking slice。 |
| `goal/xlib_standard_pr7_goal_runtime_v3_1_1_ci_artifacts_scorecard_execution_pack.md` | `d1cea6204cc8ca7c7ebc951fd068a6e915baa1862a39c855a3c0839e13e95f0d` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-8 | 保留参考：CI artifacts/scorecard 并入可信治理阶段。 |
| `goal/xlib_standard_pr8_goal_runtime_v3_1_1_blocking_full_mode_release_verify_execution_pack (1).md` | `46c1adca6a74ffe812df9ef96f21cb661d3d69e57614ba7513dd5310b59156d3` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-5 and PR-8 | 保留参考：blocking 与 test harness/trust root 拆分。 |
| `goal/xlib_standard_pr8_goal_runtime_v3_1_1_blocking_full_mode_release_verify_execution_pack.md` | `46c1adca6a74ffe812df9ef96f21cb661d3d69e57614ba7513dd5310b59156d3` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-5 and PR-8 | 保留参考：blocking 与 test harness/trust root 拆分。 |
| `goal/xlib_standard_pr9_goal_runtime_v3_1_1_runtime_freeze_pack_check_drift_control_execution_pack.md` | `af3099e4d3de54e7e6c781565d380c80882ae3d3d3b25c17934f5558dc89fb76` | `split/docs/plans/goalkit-v0.1.0-roadmap.md` PR-6 | 保留参考：freeze/drift control。 |

## 4. 审计状态

- 当前状态：文档级迁移索引已建立；尚未执行归档或删除。
- 归档前检查：确认本表每个来源文件均有 checksum、迁移目标和处置说明。
- 禁止事项：不得仅凭综合提案声明作为唯一处置证据；不得在 PR-4 之前声明 G12-G16 命令已可执行。
