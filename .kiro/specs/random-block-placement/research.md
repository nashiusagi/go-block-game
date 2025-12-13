# Research & Design Decisions
---
**Purpose**: Capture discovery findings, architectural investigations, and rationale that inform the technical design.
---

## Summary
- **Feature**: random-block-placement
- **Discovery Scope**: Extension
- **Key Findings**:
  - 既存 `initBlocks` は固定グリッド配置で乱数・重なり回避・最小距離の仕組みがない
  - パドルYは画面下部付近（約550px）に固定されているため、最小Yギャップの定数化と検証が必要
  - 無限リトライ防止のため試行上限とフォールバック戦略が必要（ランダム配置失敗時の扱い）

## Research Log

### ブロック配置拡張
- **Context**: 既存配置は等間隔グリッドで固定。ランダム配置へ拡張する必要がある。
- **Sources Consulted**: 既存 `main.go` の `initBlocks` 実装
- **Findings**:
  - 画面サイズ: 800x600、ブロックサイズ: 70x30、行列: 5x10
  - パドル位置: Y=約550px、ブロック開始Y=50px（十分離れているが定数化されていない）
  - 配置時の境界計算と重なり判定が未実装
- **Implications**:
  - ランダム生成関数で境界内配置と重なり排除を行う
  - 最小Yギャップを設定し、配置判定で使用する

### 乱数とリトライ上限
- **Context**: ランダム配置で重なりが発生する可能性があるため、リトライと上限が必要。
- **Sources Consulted**: Go標準 `math/rand` の利用想定（新規依存なし）
- **Findings**:
  - `rand.Rand` をインジェクションしてシード管理を可能にする
  - リトライ上限を `blockCount * 10` 程度に設定して無限ループを防止
- **Implications**:
  - 上限超過時はエラーを返し、フォールバック（例: 失敗をUIに表示または再シード）を設計で定義

### パドル距離制約
- **Context**: パドルとブロックの最小Y距離を要件で保証する必要がある。
- **Sources Consulted**: 現行パドルY、画面サイズ
- **Findings**:
  - パドルY=約550px、ブロック高さ30pxを考慮し、最小ギャップ例: 150–200px が妥当
  - ギャップ定数は設定で変更可能にする
- **Implications**:
  - コンフィグに `minPaddleGap` を追加し、配置時に判定
  - ギャップ未達なら再配置する

## Architecture Pattern Evaluation
| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| 既存関数拡張 | `initBlocks` 内で乱数配置と検証を直接実装 | 影響範囲が小さい、変更箇所が限定 | テストしづらい、責務肥大 | 短期対応向き |
| 生成関数新設 | `GenerateBlocks(cfg, rnd)` を新設し、純粋関数としてテスト | 責務分離、テスト容易、将来パラメータ拡張が容易 | 呼び出し側改修が必要 | 推奨 |
| ハイブリッド | 生成関数を新設し `initBlocks` から呼び出す | 互換性維持しつつ責務分離 | 若干の複雑化 | 採用案 |

## Design Decisions

### Decision: 乱数ブロック生成を純粋関数として分離
- **Context**: 配置ロジックのテスト容易性と再利用性を高めたい
- **Alternatives Considered**:
  1. 既存関数内に直書き
  2. 新規生成関数を追加し、呼び出しで注入
- **Selected Approach**: 新規生成関数を追加し、`initBlocks` から呼び出す
- **Rationale**: テスト容易、将来の難易度調整やレイアウト切替に対応しやすい
- **Trade-offs**: 呼び出し側の改修が必要
- **Follow-up**: リトライ上限とエラー処理の実装・テスト

### Decision: リトライ上限とフォールバック
- **Context**: 重なりや境界違反で無限ループのリスクがある
- **Alternatives Considered**:
  1. 上限なし（危険）
  2. 上限あり＋エラー返却
- **Selected Approach**: 上限あり、超過時はエラー返却しフォールバックを呼び出し側で判断
- **Rationale**: 安全に失敗を検出しやすい
- **Trade-offs**: 呼び出し側でのエラー処理が必要
- **Follow-up**: エラー時のUI/ログ方針を決定

## Risks & Mitigations
- リトライ上限超過で配置失敗 → 上限設定とエラー返却、必要なら再シード再試行
- 最小Yギャップ未達成 → 判定ロジックと再配置ループで検出
- ランダム性による難易度ばらつき → シード設定オプションとテストケースで分布を確認

## References
- Go `math/rand` 標準ライブラリ（乱数生成の基礎利用）

