# Famicom

Go言語で実装したファミリーコンピュータ（NES）エミュレータです。

## 概要

このプロジェクトは、[NES Emulator From Scratch](https://bugzmanov.github.io/nes_ebook/chapter_1.html)の資料を参考に、Rust版をGo言語で実装したNESエミュレータです。Go言語の学習を目的としており、画面描画にはEbitenライブラリを使用しています。

## 現在の実装状況

現在は基本的なCPU（6502）エミュレーションとシンプルなスネークゲームが動作します：

- ✅ 6502 CPUエミュレーション（基本命令セット）
- ✅ メモリマップドI/O
- ✅ バスシステム
- ✅ ROMローダー
- ✅ Ebitenを使用した画面出力
- ✅ キーボード入力（WASD）
- ✅ サンプルゲーム（スネーク）

## 必要環境

- Go 1.22.3以上
- 対応OS: Windows、macOS、Linux

## インストール・実行方法

1. リポジトリをクローン
```bash
git clone https://github.com/tabo-syu/famicom.git
cd famicom
```

2. 依存関係をインストール
```bash
go mod tidy
```

3. ゲームを実行
```bash
go run cmd/famicom/main.go
```

## 操作方法

- **W**: 上
- **A**: 左
- **S**: 下
- **D**: 右
- **Esc**: ゲーム終了

## プロジェクト構成

```
.
├── cmd/
│   ├── famicom/           # メインアプリケーション
│   └── opcode_scraper/    # 命令コード生成ツール
├── internal/
│   ├── bus/               # システムバス
│   ├── cpu/               # 6502 CPUエミュレーション
│   ├── game/              # ゲームロジック・画面描画
│   ├── memory/            # メモリ管理
│   └── rom/               # ROMローダー
├── go.mod
└── go.sum
```

## 使用ライブラリ

- [Ebiten](https://github.com/hajimehoshi/ebiten): 2Dゲームエンジン（画面描画・入力処理）
- [Colly](https://github.com/gocolly/colly): Webスクレイピング（命令コード生成用）
- [Testify](https://github.com/stretchr/testify): テストフレームワーク

## 開発・テスト

テストの実行:
```bash
go test ./...
```

命令コード生成ツールの実行:
```bash
go run cmd/opcode_scraper/main.go
```

## 今後の実装予定

- PPU（Picture Processing Unit）実装
- APU（Audio Processing Unit）実装
- マッパー対応
- 実際のROMファイル読み込み
- セーブ・ロード機能
- デバッガー機能

## 参考資料

- [NES Emulator From Scratch](https://bugzmanov.github.io/nes_ebook/chapter_1.html) - メインの実装参考資料（Rust版）
- [NESDev Wiki](https://wiki.nesdev.com/) - NESハードウェア仕様
- [6502.org](http://6502.org/) - 6502 CPU仕様

## ライセンス

MIT License

## 貢献

Issue報告やPull Requestを歓迎します。