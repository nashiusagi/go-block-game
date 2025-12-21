# Implementation Plan

## Task Format
- Major+Sub タスク形式。`(P)` で並行可能を示す。
- テストが後追いなら `- [ ]*` を使用。

## Tasks

- [ ] 1. ドメインモデル拡張（Item.Type / PaddleEffect）
- [x] 1.1 ItemType 列挙と Item 構造体への Type フィールド追加
  - `ItemTypeMultiball`, `ItemTypePaddleEnlarge` を定義
  - 既存コードで Item 生成箇所に `Type: ItemTypeMultiball` を追加
  - _Requirements: 1.1, 1.2, 3.1_
- [x] 1.2 PaddleEffect 構造体を追加し GameState に組み込む (P)
  - `Active bool`, `RemainingTicks int`, `BaseWidth float64`, `Multiplier float64`
  - `GameState` に `PaddleEffect PaddleEffect` フィールド追加
  - _Requirements: 2.1, 2.2, 2.4_

- [ ] 2. Config 定数追加
- [x] 2.1 pkg/config/layout.go に PaddleEnlarge 関連定数を追加 (P)
  - `PaddleEnlargeChance = 0.02`
  - `PaddleEnlargeDuration = 300` (5秒 @ 60FPS)
  - `PaddleEnlargeMultiplier = 3.0`
  - LayoutConfig に `PaddleEnlargeChance float64` フィールド追加
  - _Requirements: 1.1, 2.1_

- [ ] 3. アイテムスポーン拡張（tryDropItem）
- [x] 3.1 tryDropItem を拡張しパドル拡大アイテムの独立抽選を追加
  - マルチボール抽選（既存 ItemDropChance）とパドル拡大抽選（PaddleEnlargeChance）を別々に実行
  - 共通の spawnItem ヘルパーで Item 生成を統一
  - MaxItems 上限チェックは全アイテム合計で適用
  - _Requirements: 1.1, 1.2, 1.5_

- [ ] 4. アイテム衝突処理拡張（updateItems）
- [ ] 4.1 updateItems でアイテム種別に応じた効果適用に分岐
  - `ItemTypeMultiball` → `applyMultiball`（既存）
  - `ItemTypePaddleEnlarge` → `applyPaddleEnlarge`（新規）
  - _Requirements: 1.3, 1.4, 2.1_

- [ ] 5. パドル拡大効果ロジック（applyPaddleEnlarge / updatePaddleEffect）
- [ ] 5.1 applyPaddleEnlarge を実装
  - 効果未適用時: BaseWidth を保存し、幅を 3倍に設定
  - 効果適用中: 幅は維持し RemainingTicks をリセット
  - _Requirements: 2.1, 2.2, 2.3_
- [ ] 5.2 updatePaddleEffect を実装し Advance から呼び出す
  - 毎フレーム RemainingTicks をデクリメント
  - 0 以下になったら幅を BaseWidth に戻し Active = false
  - _Requirements: 2.4, 2.5_

- [ ] 6. 描画対応（Renderer）
- [ ] 6.1 アイテム種別に応じた色分け描画
  - マルチボール: 既存色（黄系）
  - パドル拡大: 別色（緑系など）
  - _Requirements: 3.1_
- [ ] 6.2 パドル拡大効果中のインジケータ表示 (P)
  - 効果中はパドル色を変化（例: 白→シアン）または残り時間表示
  - _Requirements: 3.2, 3.3_

- [ ] 7. テスト
- [ ] 7.1 tryDropItem のユニットテスト追加
  - 2%確率でパドル拡大アイテムがスポーンすること
  - マルチボールと独立に抽選されること
  - _Requirements: 1.1_
- [ ] 7.2 applyPaddleEnlarge / updatePaddleEffect のユニットテスト追加 (P)
  - 幅が3倍になり Active=true, RemainingTicks=300 になること
  - 300 ticks 後に元の幅に戻ること
  - 再取得でタイマーリセット、幅維持
  - _Requirements: 2.1, 2.2, 2.3, 2.4_
- [ ] 7.3 統合テスト: ブロック破壊→取得→効果終了フロー (P)
  - _Requirements: 1.x, 2.x_

