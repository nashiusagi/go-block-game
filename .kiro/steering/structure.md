# Project Structure

## Organization Philosophy

シンプルなモノリシック構造。小規模なゲームプロジェクトのため、複雑な階層構造は避け、単一ファイルにゲームロジックを集約しています。

## Directory Patterns

### Root Level
**Location**: `/`  
**Purpose**: プロジェクトのルートディレクトリ。メインのソースコードと設定ファイルを配置  
**Example**: `main.go`, `go.mod`, `README.md`

### Game Logic
**Location**: `main.go`  
**Purpose**: ゲームの全ロジック（初期化、更新、描画、衝突判定）を含む  
**Pattern**: 
- `Game`構造体がゲーム状態を保持
- `Update()`メソッドでゲームロジックを更新
- `Draw()`メソッドでレンダリング
- `Layout()`メソッドで画面サイズを定義

## Naming Conventions

- **Files**: `snake_case.go`（Go標準）
- **Types**: `PascalCase`（例: `Game`, `Block`, `Ball`, `Paddle`）
- **Functions**: `PascalCase` for exported, `camelCase` for unexported
- **Constants**: `camelCase`（例: `screenWidth`, `ballSpeed`）
- **Variables**: `camelCase`

## Import Organization

```go
import (
    // 標準ライブラリ
    "fmt"
    "image/color"
    "math"
    
    // 外部ライブラリ
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)
```

**Import順序**: 標準ライブラリ → 外部ライブラリ（空行で区切る）

## Code Organization Principles

- **構造体定義**: ファイル上部に型定義を配置
- **初期化関数**: `NewGame()`, `initBlocks()`, `initBall()`, `initPaddle()`で初期化ロジックを分離
- **ゲームループ**: Ebitengineの`Update()`/`Draw()`パターンに従う
- **定数定義**: ファイル上部でゲームパラメータを定義

## Game State Management

- `Game`構造体がすべてのゲーム状態を保持
- 状態変更は`Update()`メソッド内で行う
- 描画は`Draw()`メソッド内で行う（状態を変更しない）

---
_Document patterns, not file trees. New files following patterns shouldn't require updates_

