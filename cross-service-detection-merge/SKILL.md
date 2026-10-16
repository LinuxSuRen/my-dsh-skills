---
name: cross-service-detection-merge
description: 多算法服务检测同一物体不同检查项时的结果合并方法论：按检查项 code 增量合并（非整体替换）、任务子集过滤防止跨服务数据混乱、checkItems 键必须用检查项编码而非 itemKey、unknown/undecidable 与 abnormal 的显示分离。提炼自机器人平台算法服务注册化实战。
whenToUse: 平台派发检测任务给多个算法服务（各自负责不同检查项/子任务），需要合并多路检测结果到统一报告且互不覆盖时。
---

# 多算法服务检测结果合并

## 1. 核心原则：按检查项 code 增量合并

多服务各自负责不同检查项时，合并结果是**按检查项 code 增量添加**，不是整体替换：

```
服务A回传 checkItems: {"44": {value: 0}}      ← 只含自己负责的检查项
服务B回传 checkItems: {"DD": {value: 1}}      ← 只含自己负责的检查项
合并后: {"44": {...}, "DD": {...}}             ← 两个 code 共存
```

- 同 code 后到覆盖先到（幂等）；不同 code 增量新增
- 合并容器用普通 `map[string]any`，不用自定义 JSON 类型（类型断言按动态类型精确匹配会失效）
- 权威结论（status/serviceId）以最新反馈为准，但 checkItems 增量不丢

## 2. dispatcher fan-out 时按服务过滤载荷检查项字段

- `tasks` 数组已按服务绑定的任务子集过滤 ✓
- **`checkItemCode` 和 `enabledCheckItems` 也必须按服务过滤**——否则服务会回报不属于自己的检查项结果，合并时键冲突互相覆盖
- 过滤依据：`algorithm_task_mappings` 反查每个检查项 code 映射的检测任务，只保留与该服务 tasks 子集有交集的编码

## 3. 服务回传 checkItems 键必须用检查项编码

- 用 `checkItemCode`（payload 携带的检查项编码），**不用 itemKey**（模板生成的兜底值如 "item-1"，平台按编码查不到）
- 多检查项步骤：`enabledCheckItems` 每个编码都生成同值条目（一次检测覆盖全部检查项）
- 本地映射缺失时用 payload 的 `checkItemCode` 兜底，不退化为任务名

## 4. 值语义（按检测类型）

| 检测类型 | value | 空/无匹配 |
|---|---|---|
| 物体检测（target_classes） | 目标类计数（数字：0/1/3…） | 0 |
| 文本识别（prompt/OCR） | 识别的文本字符串 | "" |
| 布尔判定 | true/false | — |
| 人工判定 | 人工填写的检测值（默认按判定推导） | 推导默认值 |

## 5. 迟到反馈的补充合并（不覆盖已采纳结果）

- 权重高的服务先回 → 短路定稿 → 权重低的后回 → **只存参考不覆盖**（防报告/告警抖动）
- **例外**：迟到反馈有非空 checkItems 且权威结果缺该检查项 → 升级为补充反馈（增量合并），结论不覆盖
- 仲裁器已定稿条目**长期保留**（不设宽限期淘汰），内存由 LRU 上限兜底——否则慢引擎（LLM 推理常 >10s）的迟到反馈落入"未知 requestId"路径被整体覆盖

## 6. 状态显示：unknown ≠ abnormal

| 状态 | 含义 | 显示 |
|---|---|---|
| normal | 检测标准满足 | 成功 |
| abnormal | 检测标准不满足 | 异常 |
| unknown | **未收到检测结果** | **待检测**（非异常！） |
| undecidable | 收到结果但值与标准类型不匹配 | **无法判断**（非异常！） |

- `unknown` 和 `undecidable` 必须从"异常"的显示/汇总逻辑中拆出，用中性色（琥珀）显示
- 运行列表的异常计数不应包含 unknown/undecidable
