# Requirements Document

## Introduction
ブロック破壊時に **2%** の確率で「パドル拡大」アイテムが出現し、パドルで取得すると **5秒間** パドル幅が **3倍** になる。短時間の救済として機能し、効果終了後は元の幅（難易度適用後）に戻る。既存のマルチボールアイテムとは別種として共存させる。

## Requirements

### Requirement 1: パドル拡大アイテムの出現
**Objective:** As a プレイヤー, I want ブロック破壊時に一定確率でパドル拡大アイテムが落下する, so that 一時的に守りを強化できるチャンスを得られる

#### Acceptance Criteria
1. When ブロックが破壊されアイテム抽選が行われるとき, the system shall spawn a paddle-enlarge item with **2%** probability（マルチボール抽選とは独立）。
2. If アイテムがスポーンした場合, the system shall 破壊されたブロック位置（中心付近）に配置し一定速度で落下させる。
3. While アイテムが落下中, the system shall 毎フレーム位置を更新し、画面外に出たら消去する。
4. Where 複数アイテム（マルチボール含む）が同時に存在する場合, the system shall 各アイテムを独立して処理する。
5. The system shall 同時生成アイテム数の上限（MaxItems）を超える場合は新規スポーンをスキップする。

### Requirement 2: 取得時の効果適用（パドル幅3倍・5秒）
**Objective:** As a プレイヤー, I want アイテム取得でパドル幅が一定時間3倍になる, so that 球を受け止めやすくなる

#### Acceptance Criteria
1. When パドルが paddle-enlarge item と衝突したとき, the system shall パドル幅を現在のベース幅の **3倍** に設定し、**5秒** の効果タイマーを開始する。
2. While 効果が有効な間, the system shall 幅を3倍に保ち、再度取得しても3倍を超えないようにする。
3. If 効果有効中に同種アイテムを再度取得した場合, the system shall 残り時間を5秒にリセットし、幅は3倍を維持する。
4. When 効果時間が満了したとき, the system shall パドル幅を効果適用前のベース幅に戻し、効果状態をクリアする。
5. Where ゲームがポーズ中の場合, the system shall 効果残り時間のカウントダウンを停止し、再開後に継続する。

### Requirement 3: アイテムの識別と表示
**Objective:** As a プレイヤー, I want アイテムの種類と効果の有無を視覚的に区別できる, so that 状況を把握してプレイできる

#### Acceptance Criteria
1. When パドル拡大アイテムが落下中, the system shall マルチボールアイテムと識別可能な色・形で描画する。
2. When パドル拡大効果が有効になったとき, the system shall 効果中であることを示すインジケータ（例: パドル色変化など）を表示する。
3. When 効果が終了したとき, the system shall インジケータを消去し通常表示に戻す。


