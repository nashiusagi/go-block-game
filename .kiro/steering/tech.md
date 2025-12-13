# Technology Stack

## Architecture

シンプルなモノリシック構造。ゲームロジック、レンダリング、入力処理が単一の`main.go`ファイルに集約されています。Ebitengineのゲームループ（Update/Draw）パターンに従っています。

## Core Technologies

- **Language**: Go 1.24.0+
- **Framework**: Ebitengine v2 (2.10.0-alpha.7+)
- **Runtime**: Go標準ランタイム

## Key Libraries

- **Ebitengine v2**: 2Dゲームエンジン、レンダリングと入力処理を提供
- **標準ライブラリ**: `math`（物理演算）、`image/color`（色定義）、`fmt`（文字列フォーマット）

## Development Standards

### Type Safety

- Goの型システムを活用
- `float64`を座標と物理演算に使用（精度のため）
- 構造体でゲームオブジェクトを表現（Block, Ball, Paddle）

### Code Quality

- Goの標準的な命名規則に従う（PascalCase for exported, camelCase for unexported）
- 定数でマジックナンバーを避ける
- 構造体メソッドでゲームロジックを整理

### Testing

現在はテストコードなし。必要に応じて`*_test.go`ファイルを追加。

## Development Environment

### Required Tools

- Go 1.24.0以上
- システム依存ライブラリ（X11開発ライブラリなど、プラットフォーム依存）

### Common Commands

```bash
# 実行
go run main.go

# ビルド
go build -o block-game

# 依存関係の更新
go mod tidy
```

## Key Technical Decisions

- **Ebitengine選択**: クロスプラットフォーム対応、軽量、Goネイティブ
- **単一ファイル構造**: 小規模プロジェクトのため、シンプルさを優先
- **物理演算**: 簡易的なAABB衝突判定と反射計算を実装
- **定数ベースの設定**: ゲームパラメータ（画面サイズ、ブロック数など）を定数で管理

---
_Document standards and patterns, not every dependency_

