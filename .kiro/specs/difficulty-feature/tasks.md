# Implementation Plan

## Task Format
- Major+Sub タスク形式。`(P)` で並行可能を示す。
- テストが後追いなら `- [ ]*` を使用。

## Tasks
- [ ] 1. 難易度設定の定義と検証
- [x] 1.1 Difficulty 列挙と DifficultySetting 構造体を追加し、EASY/NORMAL/HARD のスケールを定義する
  - スケール例: EASY (BallSpeed 0.8, BallRadius 1.0, PaddleWidth 1.1, PaddleSpeed 1.1, BlockSize 1.0, BlockCount 1.0), NORMAL (1.0,1.0,1.0,1.0,1.0,1.0), HARD (1.2,0.9,0.9,0.9,0.9,1.3)
  - BlockCount は HARD で増やし、BlockSize/BallRadius を縮小して面積を確保
  - _Requirements: 2.1,2.2,2.3,2.4,2.5,4.3_
- [x] 1.2 DifficultyValidator を実装し、難易度キーとスケールの範囲を検証する
  - スケール範囲 (0, 10] でクランプし、0以下はエラー。BlockCount は面積に収まるよう上限をかける
  - 無効値は NORMAL にフォールバックし、ログ/エラーを出力
  - _Requirements: 3.3,4.1,4.2,4.3,4.4_

- [ ] 2. コンフィグ適用とレイアウト生成統合
- [ ] 2.1 ApplyDifficulty を実装し、基準 LayoutConfig にスケールを適用して派生コンフィグを返す
  - BallSpeed/Radius、PaddleWidth/Speed、BlockW/H、BlockCount、MaxAttempts をスケール
  - BlockCount 増加時は MaxAttempts を比例増加（例: count*10）し、面積不足なら BlockSize を再クランプ
  - _Requirements: 1.1,1.2,1.3,1.4,1.5,2.1,2.2,2.3,2.4,2.5,4.1,4.3_
- [ ] 2.2 GenerateBlocks 統合を調整し、HARD のブロック増加・サイズ縮小でも配置が成立するようにする
  - 面積チェックを追加し、成立しない場合はサイズクランプまたはフォールバック（ノーマル配置やグリッド）を選択
  - _Requirements: 2.4,2.5,4.1_

- [ ] 3. 初期化フローと状態管理
- [ ] 3.1 ゲーム開始前の難易度選択フローを追加（未選択時は NORMAL を適用）
  - 入力/設定値から Difficulty を決定し、一ゲームセッションのみ保持
  - _Requirements: 1.1,1.2,1.3,1.4_
- [ ] 3.2 NewGameState/initBlocks 周辺で ApplyDifficulty を適用し、生成したコンフィグで Blocks/Balls/Paddle を初期化
  - HARD で BlockCount 増加時も MaxBalls/Items の上限に注意（必要なら合わせて調整）
  - _Requirements: 1.3,2.1,2.2,2.3,2.4,2.5,3.3,4.1_

- [ ] 4. 表示とフォールバック
- [ ] 4.1 難易度ラベルをゲーム画面・ポーズ/リスタート表示に追加する
  - 文言は `Difficulty: EASY|NORMAL|HARD` を統一
  - _Requirements: 1.5,3.1,3.2,3.4_
- [ ] 4.2 無効設定検出時に NORMAL へフォールバックし、UI/ログに反映する
  - _Requirements: 3.3,4.1,4.2_

- [ ] 5. テスト
- [ ] 5.1 ApplyDifficulty/DifficultyValidator のユニットテストを追加
  - スケール適用、クランプ、無効値フォールバック、BlockCount 増加時の MaxAttempts 増加を検証
  - _Requirements: 2.x,3.3,4.x_
- [ ] 5.2 GenerateBlocks 統合テストを追加し、HARD 設定で配置が成功することと NORMAL で従来挙動が維持されることを確認
  - _Requirements: 2.4,2.5,3.1,3.2_
- [ ] 5.3 難易度ラベル表示とデフォルト適用のテストを追加
  - 未選択→NORMAL、無効値→NORMAL、表示文言の一貫性を確認
  - _Requirements: 1.x,3.x_



