# Requirements Document

## Introduction

この仕様は、ブロックが破壊された際に一定確率でアイテムが落下し、パドルで取得すると場のボール数が倍増するマルチボール効果を定義します。ドロップ率や同時生成上限を設け、拾えなかったアイテムは画面外で消滅させることでゲームバランスを保ちます。

## Requirements

### Requirement 1: アイテムドロップの発生

**Objective:** As a プレイヤー, I want ブロック破壊時にたまにアイテムが落下する, so that プレイにランダムな楽しさが加わる

#### Acceptance Criteria
1. When ブロックが破壊される, the Block Game shall 一定確率（例:10%）でアイテムを生成する
2. If 生成済みアイテム数が上限に達している, the Block Game shall 新規アイテムを生成しない
3. While アイテムが生成される, the Block Game shall 破壊されたブロックの位置から落下を開始する
4. Where アイテムが生成される, the Block Game shall ブロック中心付近に初期座標を設定する
5. The Block Game shall 確率と上限を定数または設定値として管理できるようにする

### Requirement 2: アイテム落下と取得判定

**Objective:** As a プレイヤー, I want 落下するアイテムをパドルで拾える, so that 効果を自分の操作で発動できる

#### Acceptance Criteria
1. When アイテムが落下中, the Block Game shall 毎フレーム一定速度または加速度でY座標を更新する
2. When アイテムがパドルと衝突する, the Block Game shall アイテムを取得済みとして消去する
3. If アイテムが画面下端に到達する, the Block Game shall アイテムを消去して未取得として扱う
4. Where 当たり判定を行う, the Block Game shall アイテムとパドルのAABB判定を用いる
5. The Block Game shall 同時に複数アイテムが存在する場合でも各アイテムを独立に更新・描画する

### Requirement 3: マルチボール効果の適用

**Objective:** As a プレイヤー, I want アイテム取得でボールが倍増する, so that ラウンドが一気に賑やかになる

#### Acceptance Criteria
1. When アイテムが取得される, the Block Game shall 場にある各ボールを複製し合計数を2倍にする
2. If 複製後のボール数が上限（例:8個）を超える, the Block Game shall 上限に切り詰めて生成する
3. While 新規ボールを生成する, the Block Game shall 元のボール速度を基に反転・分岐させた速度ベクトルを割り当てる
4. Where 新規ボールの初期位置を設定する, the Block Game shall 元のボール位置に配置する
5. The Block Game shall 効果発動と同時に対象アイテムを消去し再発動を防止する

### Requirement 4: 効果の状態管理とリセット

**Objective:** As a 開発者, I want マルチボール効果がゲーム状態と整合する, so that 予期せぬ状態ずれを防げる

#### Acceptance Criteria
1. When 新しいゲームやラウンドが開始される, the Block Game shall ボール数を初期値（1個）にリセットする
2. If 全てのボールを失った場合, the Block Game shall 通常ルールに従いライフ消費や再生成処理を行う
3. While 効果が有効な間, the Block Game shall 衝突判定・スコア加算・ゲームオーバー判定を全ボールに対して適用する
4. The Block Game shall マルチボール関連の上限値や速度分岐パターンを設定可能な定数として管理する

