# Requirements Document

## Project Description (Input)
Goで開発されたブロック崩しゲーム。Ebitengine v2 を用い、パドルを左右に操作してボールを反射させ、画面上部のブロックを破壊するクラシックなアーケード体験を提供する。ゲームオーバーや勝利判定などの状態管理を備える。

## Requirements

### Requirement 1: 難易度選択の提供
**Objective:** As a プレイヤー, I want ゲーム開始時に難易度を選べる, so that プレイスタイルに応じて挑戦度を調整できる

#### Acceptance Criteria
1. When ゲーム開始前, the System shall イージー/ノーマル/ハードの3種から選択できるようにする
2. When 難易度が未選択, the System shall デフォルトでノーマルを適用する
3. When 難易度を選択した, the System shall 選択結果をゲームループ開始前に確定する
4. When 新しいゲームセッションを開始する, the System shall 前セッションの選択を保持しない（毎回選択可能）
5. Where 選択肢を表示する, the System shall 現在の難易度名をUIに明示する

### Requirement 2: 難易度ごとのゲームパラメータ調整
**Objective:** As a プレイヤー, I want 難易度に応じてボールとパドルの挙動が変わる, so that 体感難易度が明確に変化する

#### Acceptance Criteria
1. When イージーを選ぶ, the System shall 初期ボール速度をノーマル比で低下させる（例: -20%）
2. When ハードを選ぶ, the System shall 初期ボール速度をノーマル比で上昇させる（例: +20%）
3. When イージーを選ぶ, the System shall パドル幅または移動速度をノーマル比で増加させる（例: +10%）
4. When ハードを選ぶ, the System shall ブロックと玉を小さくしてブロックの数を多くし、同時にパドル幅または移動速度をノーマル比で減少させる（例: -10%）
5. The System shall すべての難易度でブロック総数の基本値を維持する（難易度差分は速度/操作性のみで成立）

### Requirement 3: 難易度の可視化とフィードバック
**Objective:** As a プレイヤー, I want 現在の難易度をゲーム中に確認できる, so that 設定が適用されていることを理解できる

#### Acceptance Criteria
1. When ゲームが開始する, the System shall 画面上に現在の難易度名を表示する
2. While ゲームが進行する, the System shall ポーズ／リスタート画面でも難易度名を表示し続ける
3. If 難易度設定が無効値の場合, the System shall ノーマルにフォールバックし、UIにノーマルを表示する
4. The System shall 難易度ごとに一貫した表示文言を用いる（例: EASY / NORMAL / HARD）

### Requirement 4: 設定の検証と安全性
**Objective:** As a 開発者, I want 無効な難易度や不正パラメータを防ぐ, so that ゲーム進行が破綻しない

#### Acceptance Criteria
1. When 難易度値が定義外の場合, the System shall デフォルトノーマルに強制し、エラーをロギングする
2. When 難易度ごとの速度・幅スケールが0以下や極端に大きい場合, the System shall 起動時に検証しエラーを報告する
3. The System shall 難易度ごとの設定値を定数またはコンフィグとして一元管理する
4. The System shall テストで固定パラメータを用い、挙動が決定的に再現されることを確認する


