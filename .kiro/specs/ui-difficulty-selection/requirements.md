# Requirements Document

## Introduction
ゲーム開始前にUI上から難易度を選択できるようにし、プレイヤーが自分に合った挑戦度でプレイを開始できるようにする。

## Requirements

### Requirement 1: UIによる難易度選択
**Objective:** As a プレイヤー, I want ゲーム開始前にUIで難易度を選択できる, so that 自分に合った挑戦度でプレイを始められる

#### Acceptance Criteria
1. When タイトル画面または新規ゲーム開始画面が表示されたとき, the system shall present at least three difficulty options (例: Easy/Normal/Hard) with one option preselected as the default.
2. If プレイヤーが難易度を変更して開始ボタンを押下した場合, the system shall persist the chosen difficulty and hand it to the game session before the firstブロック生成が行われる。
3. While キーボード操作およびマウス操作が有効な状態, the system shall allow difficulty selection via focus移動とクリック/Enterによる決定。
4. Where 各難易度オプションが表示される領域, the system shall show a short description of the expected effects（例: 初期スピード/残機/スコア倍率など）adjacent to the option.
5. The system shall prevent game start when no difficulty option is selected by auto-selecting the default and surfacing the current selection to the player.

### Requirement 2: 選択結果の反映と確認性
**Objective:** As a プレイヤー, I want 選択した難易度がゲーム設定に反映され確認できる, so that 期待した難易度でプレイしていると安心できる

#### Acceptance Criteria
1. When ゲームセッションが開始するとき, the system shall apply the selected difficulty parameters to core gameplay settings（例: ブロック落下速度/敵出現頻度/初期リソース）。
2. When ポーズメニューまたはヘルプ/情報画面を開いたとき, the system shall display the currently active difficulty level and its brief description.

