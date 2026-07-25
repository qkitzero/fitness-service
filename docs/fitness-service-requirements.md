# fitness-service 要件・ドメイン設計書

> 自作 `fitness-service` の要件とドメイン設計をまとめた文書。既存サービス調査（[health-check-system-research.md](./health-check-system-research.md)）で抽出したドメイン理解を土台に、**マイクロサービス構成（auth-service / user-service / fitness-service）へ再配置**し、**業種非依存に汎用化**した要件へ落とし込む。
>
> - 位置づけ: 実装の指針となる要件・設計書（コードそのものではない）
> - 前提スタック: Go / DDD（domain・application・infrastructure・interface）/ buf・proto / gRPC + grpc-gateway / GORM + golang-migrate
> - 関連リポジトリ: `auth-service`, `user-service`, `fitness-service`

---

## 1. 背景と目的

調査レポートは「健康チェックシステム」を閲覧調査し、機能・画面・データ項目を抽出したもの。ただしそのままでは、本プロジェクトのマイクロサービス構成では**他サービスの責務**にあたる要件（認証・スタッフ・院/子院階層・権限区分）が混在している。

本書の目的は次の3点：

1. **サービス境界の明文化** — auth-service / user-service / fitness-service の分担を定義し、調査レポートの各要件をいずれかのサービスへ割り当てる。
2. **ドメインの汎用化** — 接骨院/介護予防に特化した概念を、業種非依存の**汎用コア**＋**業種別拡張**へ分離する。
3. **実装ロードマップ** — 既存 `customer` ドメインの DDD 構成を土台に、段階的な追加順序を示す。

### 確定方針

| # | 方針 |
|---|---|
| 1 | 要件は本書（別ファイル）に記述する。調査レポートも自作サービス向けに不要機能を除外して整理する |
| 2 | 顧客（旧「患者」＝測定対象者）は **fitness-service 独自の `customer` ドメイン**が持つ。user-service の User（＝スタッフ）とは別エンティティ。**顧客は自らログインしない**（測定対象者としてスタッフが登録・管理する）。ログインするのは院スタッフおよび企業管理者のみ |
| 3 | **汎用コア＋業種別拡張**で設計。保険証番号・診察券番号・要介護度・接骨院向け問診・診療時間/カレンダー等は拡張（profile / 設定可能マスタ）として分離する |
| 4 | **対象外（除外）機能**: ポイント（ゲーミフィケーション）／ヘルプ・資料ダウンロード／患者マイページ等の患者ログイン依存機能（患者向けお知らせ・患者申請の予約・患者によるトレーニング実施記録） |
| 5 | **子供測定は対象外**。本サービスは**健康経営**を題材とし、対象を就業者（成人）に限定する。測定項目・処方ルールに「大人/子供区分」は設けない |

---

## 2. サービス分担（境界）

### 2.1 各サービスの責務

| サービス | 責務 | 主な提供 API / ドメイン |
|---|---|---|
| **auth-service** (`auth.v1`) | 認証全般（OAuth/OIDC ベース） | Login / ExchangeCode / VerifyToken(→user_id) / RefreshToken / RevokeToken / Logout / GetM2MToken |
| **user-service** (`user.v1`) | ユーザー（＝**スタッフ**）の実体 | User{user_id, display_name, birth_date}（auth identity と紐付く）。CreateUser / GetUser / UpdateUser |
| **user-service** (`group.v1`) | **院＝テナント**（親子院階層）と**所属＋権限区分** | Group{group_id, name}＋親子階層（group_relations）／Membership{user_id, group_id, role}（AddMember / UpdateMemberRole / RemoveMember / ListMembers / ListMyGroups / ListChild・ParentGroups） |
| **fitness-service** | **フィットネス中核ドメイン** | 顧客・測定・判定・トレーニング処方・来院・情報発信・問診・各マスタ・施設プロファイル（本書3章） |

### 2.2 要件マッピング（調査レポート → 担当サービス）

| 調査レポートの要件 | 担当サービス | 備考 |
|---|---|---|
| 認証（院コード＋ID＋PW）/ ログアウト / トークン | **auth-service** | OAuth/OIDC へ置換。院コードによる認証は廃止 |
| スタッフの氏名・生年月日・ログインID | **user-service** `User` | display_name + birth_date。auth identity と紐付く |
| スタッフの権限区分3段階（一般(PII閲覧不可)/一般/管理者） | **user-service** `Membership.role` | role 値として定義（§2.3） |
| スタッフの付加属性（資格情報・連絡先・住所） | **未確定**（§6） | user-service 拡張か fitness-service StaffProfile か |
| 院・子院の階層 | **user-service** `Group` + `group_relations` | 院コード → group_id |
| 患者（顧客） | **fitness-service** `Customer`（+ ClinicProfile 拡張） | §3 |
| 測定・判定・運動器年齢・サルコペニア | **fitness-service** | §3 |
| トレーニングメニュー・自動処方（固定/年代別/要素別） | **fitness-service** | §3 |
| 来院/受付・会計 | **fitness-service** | §3 |
| 院内掲示板（スタッフ向け情報発信） | **fitness-service** | §3 |
| 問診項目・回答 | **fitness-service** | §3 |
| 測定項目・基準値・トレーニングマスタ | **fitness-service** | §3 |
| 施設情報の付加属性（住所/診療時間/診療カレンダー） | **fitness-service** `ClinicProfile`（group_id 参照） | user-service Group は name のみのため fitness 側が保持 |

### 2.3 権限区分（role）の定義案

user-service `Membership.role` に格納する値として、調査レポートの3段階を対応させる：

| role 値（案） | 意味 | 対応（旧） |
|---|---|---|
| `admin` | 管理者権限（全設定・全PII可） | 管理者権限 |
| `staff` | 一般権限（PII閲覧可） | 一般権限 |
| `staff_restricted` | 一般権限（個人情報閲覧不可） | 一般権限（個人情報閲覧不可） |

> fitness-service は認可判定のため `role` を user-service から取得する（Membership 参照）。PII 閲覧可否など細かな認可ポリシーの適用箇所は実装時に定義。

### 2.4 サービス間連携の原則

- fitness-service の各エンティティは、**スタッフを `user_id`、院を `group_id` として ID のみ参照**する（実体は user-service が保持）。
- fitness-service の各ユースケースは、冒頭で `authService.VerifyToken(ctx)` を呼び、`user_id` を得てから処理する（既存 `customer` usecase と同じ規約）。
- 院（テナント）スコープのデータは `group_id` で分離する。
- **fitness-service は Group の階層（`group_relations`）を解釈しない**。Group の親子階層は user-service が汎用マイクロサービスとして提供する機能であり、fitness-service のドメイン要件ではない。認可判定は `ListMyGroups` が返す**直接所属のみ**で行い、`ListChildGroups` / `ListParentGroups` は呼ばない。
  - 帰結として、親 Group の所属者は子 Group のデータを参照できない。これは仕様であり不具合ではない。
- **顧客の所属組織は fitness-service が独自に持つ**（`Organization`＝§3.3）。user-service の Group は `Membership{user_id, group_id, role}` を伴う「ログインするユーザーの所属先」であり、ログインしない顧客の所属組織を Group で表現すると、メンバーが存在しない Group が大量に生まれて `ListMyGroups` ベースの認可と噛み合わない。

---

## 3. fitness-service 汎用ドメインモデル

### 3.1 用語マッピング（業種依存 → 汎用）

| 調査レポート（業種依存） | 本書（汎用） | 実体の所在 |
|---|---|---|
| 患者 | Customer（顧客） | fitness-service |
| 院 / テナント | Group | user-service |
| 子院 | 子 Group（group_relations）※fitness-service は解釈しない（§2.4） | user-service |
| スタッフ | User | user-service |
| （調査レポートに対応概念なし） | **Organization（顧客の所属組織）** | fitness-service |
| 診察券番号・保険証番号・要介護度 等 | ClinicProfile（業種拡張） | fitness-service |
| 接骨院向け問診項目 | InterviewItem（設定可能マスタ） | fitness-service |

> **Group と Organization は別概念**。混同しやすいため定義を固定する。
>
> | | Group | Organization |
> |---|---|---|
> | 所在 | user-service | fitness-service |
> | 役割 | **テナント**（データ分離の単位） | 顧客の**所属組織**（集計軸） |
> | ログイン | する（`Membership` を持つ User が所属） | しない |
> | 階層 | 持つ（fitness-service は解釈しない） | **持たない**（フラット） |
> | 例 | 測定を実施する施設・企業 | 顧客の勤務先企業 |

### 3.2 汎用化の原則

- **コアは業種非依存**：測定・判定・処方・来院は、フィットネス/ヘルスケア一般に通じる形で定義する。
- **業種特有は拡張で吸収**：接骨院/介護予防に固有の属性は `CustomerClinicProfile` や `ClinicProfile` に分離し、コアを汚さない。
- **設定可能マスタ**：`MeasurementItem` / `InterviewItem` / `TrainingMenu` はマスタとして定義し、業種差・院差を吸収する。院ごとの使用 ON/OFF は**マスタ本体ではなく院別設定テーブル**（例: `group_measurement_item_settings`）で持つ。マスタ自体は使用フラグを持たず、設定レコードが存在しない項目は ON とみなす。

### 3.3 エンティティ定義

各エンティティについて「汎用コア項目 / 業種拡張項目 / 別サービス参照」を明記する。ランク凡例などの実データは調査レポート §3（画面別詳細）・§6（ドメインモデル）を参照。

#### Customer（顧客）※既存 `customer` ドメインを拡張
- **汎用コア**: 氏名 / カナ氏名 / 性別（男性・女性・その他）/ 生年月日 / 連絡先（電話・メール）/ 住所（郵便番号・都道府県・市町村区・番地・その他）/ 緊急連絡先（氏名・続柄・電話）/ 利用停止フラグ
- **派生（集計）**: 初回・最終来院日 / 初回・最終測定日 / 来院回数 / 測定回数
- **別サービス参照**: 登録院 `group_id`
- **同一サービス内参照**: 所属組織 `organization_id`（任意。個人顧客は `null`）
- **ゲスト**: 一時的な測定対象者（氏名・性別・生年月日のみ）。後から正式 Customer へ昇格可能

#### Organization（顧客の所属組織）
- **汎用コア**: 組織名
- **別サービス参照**: 所有テナント `group_id`
- **階層を持たない**（フラット）。組織階層は user-service Group の汎用機能であり、fitness-service のドメイン要件ではない（§2.4）
- **テナント分離の帰結**: 同じ企業が複数テナントで別レコードとして重複登録される。共有マスタ化するとどのテナントが編集権を持つかで破綻するため、重複を許容する
- **組織名の重複を許容**（同名の別組織はあり得る）
- **拡張候補（後続）**: カナ組織名（五十音ソート用）/ 連絡先（担当者名・電話・メール）/ 住所 / 備考 / 利用停止フラグ

#### CustomerClinicProfile（顧客の業種拡張）
- 保険証番号 / 診察券番号 / 要介護度区分 / 既往歴 / 職業 / 身体活動レベル / 運動習慣（種類・時間・強度・頻度・継続期間）/ 嗜好品（飲酒・喫煙）/ 食習慣[多選] / 家族環境[多選] / 趣味 / 入会動機 / 備考
- Customer と 1:1（業種を使う院のみ利用）

#### Measurement（測定）＋ MeasurementValue（測定値）
- **Measurement**: 対象 `customer_id` / 測定日 / 測定者 `user_id` / 測定時年齢 / 更新日 / 更新者 `user_id` / 一時保存フラグ
- **MeasurementValue**: `measurement_item` 参照 / 値（複数試行・左右に対応：握力 左右×2回、棒反応5回 等）/ 測定不可フラグ / 備考
- 保持すべき値の個数は `MeasurementItem` のメタデータ（試行回数・左右区分・値の型）から決まる
- 形態（身長・体重）から BMI・適正体重を自動算出

#### MeasurementItem（測定項目マスタ）
- **属性**: 名称 / `code`（不変識別子）/ 分類（バイタル・形態・体組成・運動機能）/ 単位 / 試行回数 / 左右区分 / 値の型
- **`code`**: 判定ロジックが特定項目を名指しするための不変識別子。サルコペニア判定（歩行速度・握力・筋肉量）や BMI 算出（身長・体重）は、名称文字列や環境ごとに変わる ID ではなく `code` で項目を特定する
- **値の型**: `numeric`（単一数値）/ `paired`（上下2値＝血圧）/ `choice`（選択肢＝立ち上がりの台高さ）
- **保持する値の個数** = 試行回数 × (左右区分 ? 2 : 1) × (値の型が `paired` ? 2 : 1)
- **使用フラグは持たない**: 院ごとの ON/OFF は院別設定テーブル（§4 Phase 1'）で持つ（§3.2 の原則）
- **表示順は持たない**: 分類の並び（バイタル→形態→体組成→運動機能）は分類 enum の順序で決まる。分類内の並びはフロントで管理し、API は 分類 → `code` の順で決定的に返す
- **BMI・適正体重はマスタ項目に含めない**: 身長・体重からの派生値として算出側の責務とする
- **システム共通マスタ**とし、`group_id` は持たない。実データは調査レポート §3.14 の確定表のうち**大人測定のみ**を採用する

**seed データ（18件）**

| # | 分類 | 名称 | code | 単位 | 試行 | 左右 | 値の型 |
|---|---|---|---|---|---|---|---|
| 1 | バイタル | 血圧 | `blood_pressure` | mmHg | 1 | – | `paired` |
| 2 | バイタル | 脈拍 | `pulse_rate` | bpm | 1 | – | `numeric` |
| 3 | 形態 | 身長 | `height` | cm | 1 | – | `numeric` |
| 4 | 形態 | 体重 | `weight` | kg | 1 | – | `numeric` |
| 5 | 体組成 | 体脂肪率 | `body_fat_percentage` | percent | 1 | – | `numeric` |
| 6 | 体組成 | 筋肉量 | `muscle_mass` | kg | 1 | – | `numeric` |
| 7 | 運動機能 | 握力 | `grip_strength` | kg | 2 | ○ | `numeric` |
| 8 | 運動機能 | 立ち上がり | `stand_up_test` | cm | 1 | ○ | `choice` |
| 9 | 運動機能 | CS-30（30秒立ち座り） | `cs30` | count | 1 | – | `numeric` |
| 10 | 運動機能 | 長座体前屈 | `sit_and_reach` | cm | 2 | – | `numeric` |
| 11 | 運動機能 | 棒反応時間 | `stick_reaction` | cm | 5 | – | `numeric` |
| 12 | 運動機能 | 閉眼片足立ち | `eyes_closed_one_leg_stand` | sec | 2 | ○ | `numeric` |
| 13 | 運動機能 | 開眼片足立ち | `eyes_open_one_leg_stand` | sec | 2 ⚠ | ○ ⚠ | `numeric` |
| 14 | 運動機能 | FRT（手伸ばしテスト） | `functional_reach` | cm | 2 ⚠ | – | `numeric` |
| 15 | 運動機能 | 2ステップ | `two_step` | cm | 2 | – | `numeric` |
| 16 | 運動機能 | TUG（タイム・アップ・ゴー） | `timed_up_and_go` | sec | 2 ⚠ | – | `numeric` |
| 17 | 運動機能 | 5m歩行 | `walk_5m` | sec | 2 ⚠ | – | `numeric` |
| 18 | 運動機能 | 反復横跳び | `side_step` | count | 1 ⚠ | – | `numeric` |

> ⚠ の試行回数・左右区分は、調査レポート §3.7（測定情報入力）に入力形式の記載が無く §3.14 のマスタ一覧にのみ登場する項目のため**暫定値**。実運用で確認し後続 migration で補正する。
>
> 子供特有項目（上体起こし / 立ち幅跳び）は確定方針 #5 により除外した。

#### Judgment（判定）＋ 基準値マスタ
- **Judgment**: `measurement_id` 参照 / 項目別ランク（A〜E）/ 6要素評価（筋力・筋持久力・柔軟性・敏しょう性・バランス・移動能力）/ 運動器年齢 / サルコペニア度 / アドバイス（自由記述）
- **基準値マスタ**（システム共通 or 院共通）:
  - `AgeGroupStandard`（年代別標準値）… 同年代平均・運動器年齢の算出
  - `SarcopeniaStandard`（サルコペニア基準値）… 歩行速度・握力・筋肉量
  - `RankStandard`（ランク判定基準値）… A〜E 判定の閾値
- 時系列比較（今回 / 前回 / 初回）を提供
- **サルコペニア度は要否を再検討**（§6）。調査レポート由来の高齢者向け指標であり、就業者を対象とする健康経営では優先度が低い可能性がある

#### TrainingMenu（トレーニングメニュー マスタ）＋ PrescriptionRule（処方ルール）
- **TrainingMenu**: 名称 / 要素区分（筋力・筋持久力・柔軟性・敏しょう性・バランス・移動能力・瞬発力＝7種）/ 部位区分（上肢・下肢・全身）/ トレーニング分類 / 回数・時間 / セット数 / 手順説明 / イラスト。システム共通データ＋院独自追加（使用 ON/OFF は院別設定テーブル＝§3.2）
- **PrescriptionRule**（自動処方の元）:
  - `FixedMenu`（固定メニュー）… 常時処方
  - `AgeDecadeMenu`（年代別メニュー）… 10代〜90代ごと
  - `ElementMenu`（要素別メニュー・中核）… 要素区分 × 部位区分 × 年齢区分 × **トレーニングレベル（難易度順のはしご）**。判定した要素別ランクに応じてレベルを選び自動処方
- **PrescribedMenu**（処方結果）: `judgment_id` 参照 / 不足要素 / TrainingMenu 参照 / 回数・時間 / セット数（前回メニュー継続可）

#### Visit（来院 / 受付）
- 対象 `customer_id` / 来院日 / 入館時刻 / 退館時刻 / **会計状態（未会計・会計済）**
- 「来院者一覧」は会計済/未会計でフィルタ

#### BulletinBoard（院内掲示板）
- タイトル / 詳細メッセージ / 掲載開始日 / 掲載終了日
- 院内向け（スタッフ向け）の情報発信。掲載期間を持つ

#### InterviewItem（問診項目マスタ）＋ InterviewAnswer（問診回答）
- **InterviewItem**: 名称（使用 ON/OFF は院別設定テーブル＝§3.2）。既定13項目は調査レポート §3.15 参照
- **InterviewAnswer**: `measurement_id`（または Visit）参照 / InterviewItem 参照 / 回答

#### ClinicProfile（施設プロファイル）
- **別サービス参照**: `group_id`（院の実体は user-service Group）
- **fitness-service が保持する付加属性**: 住所 / 連絡先（担当者名・電話・FAX・メール・HP）/ 診療時間（曜日 日〜土 × 診療時間1・2＋終日休診日）/ 診療カレンダー（日別の診療/休診・時間帯の例外、名称付きパターン）/ 備考

### 3.4 集約・参照関係（俯瞰）

```
Group(user-service, テナント) ─┬─< ClinicProfile(fitness, group_id参照: 住所/診療時間/カレンダー)
      ※階層は解釈しない(§2.4)  └─< Organization(fitness, group_id参照: 顧客の所属組織・フラット)
                                        │ (organization_id 参照 / 同一サービス内なので FK あり)
User(user-service, スタッフ) ── Membership(role=権限) ─┐   │
                                                        │   │
Customer(fitness) ─┬─ CustomerClinicProfile(業種拡張)   │ ──┘
                   ├─< Visit(会計状態)                   │
                   ├─< Measurement ─< MeasurementValue ──┘(測定者/更新者 user_id)
                   │     └─ Judgment ─< PrescribedMenu
                   └─< InterviewAnswer

マスタ/基準値(システム共通 or 院共通):
  MeasurementItem / InterviewItem / TrainingMenu
  AgeGroupStandard / SarcopeniaStandard / RankStandard
  PrescriptionRule( FixedMenu / AgeDecadeMenu / ElementMenu )
```

---

## 4. 実装ロードマップ

既存 `customer` ドメインの DDD 規約を踏襲し、各ドメインを段階的に追加する。各フェーズは同一パターンで実装する：

> **実装パターン（1ドメインあたり）**
> `internal/domain/<ctx>`（値オブジェクト＋集約 interface/struct＋Repository IF＋error）→ `internal/application/<ctx>`（usecase、冒頭で `VerifyToken`）→ `internal/infrastructure/<ctx>`（GORM model＋repository）＋ `internal/infrastructure/db/migrations`（up/down）→ `internal/interface/grpc/<ctx>`（handler）→ `proto/<ctx>/v1/<ctx>.proto`（buf generate）→ `mocks/` ＋ 各層 `*_test.go`

| Phase | 内容 | 主な依存 |
|---|---|---|
| **0** | 既存 `customer` を汎用コアへ拡充＋ `CustomerClinicProfile`（業種拡張）を追加 | なし（既存土台） |
| **1a** | `MeasurementItem`（システム共通マスタ・参照専用。公開 API は `ListMeasurementItems` のみ） | なし（Customer 非依存＝**Phase 0 と並行可**） |
| **1b** | `Measurement` / `MeasurementValue` | Customer + MeasurementItem |
| **1'** | 院別使用設定（`group_measurement_item_settings`）＋ ON/OFF 切替 API | 1a ＋ role 取得（§6） |
| **0'a** | `Organization`（顧客の所属組織・フラット・CRUD） | なし（Customer 非依存＝**Phase 0 と並行可**） |
| **0'b** | `Customer` に `organization_id` を追加 | Customer ＋ Organization |
| **2** | 基準値マスタ（AgeGroupStandard / SarcopeniaStandard / RankStandard）→ `Judgment` | Measurement + Item |
| **3** | `TrainingMenu`（マスタ）→ `PrescriptionRule`（固定/年代別/要素別）→ 自動処方 `PrescribedMenu` | TrainingMenu + Judgment |
| **4** | `Visit`（会計状態） | Customer |
| **5** | `BulletinBoard` / `InterviewItem`・`InterviewAnswer` | Customer / Measurement |
| **6** | `ClinicProfile`（施設・診療時間・診療カレンダー） | Group(user-service) 参照 |

- 各フェーズは独立して価値が出る順に並べている。
- クリティカルな依存: `Judgment ← Measurement ← MeasurementItem`、`自動処方 ← TrainingMenu ＋ Judgment`。
- **並行実装の指針**: マスタ系（`MeasurementItem` / `TrainingMenu` / 基準値マスタ / `InterviewItem`）と `group_id` 参照のみのドメイン（`Organization` / `BulletinBoard` / `ClinicProfile`）は Customer に依存しないため、Phase 順を待たず並行して着手できる。
- **並行時に必ず競合する箇所**: `cmd/fitness/main.go`（DI 配線・`RegisterXxxServiceServer`・`SetServingStatus`）/ `cmd/gateway/main.go`（`RegisterXxxServiceHandler`）/ `Makefile` の `mock-gen`。migration と `gen/` はドメインごとに別ファイル・別ディレクトリのため競合しない。
- **マスタの ON/OFF 切替 API を作る際の注意**: システム共通マスタのフラグを院スタッフが変更すると全院に波及するため、切替 API は院別設定テーブル（1'）の導入と admin ロール限定の認可が揃ってから追加する。

---

## 5. 命名・規約メモ（既存 `customer` 実装に準拠）

- **proto**: `package <ctx>.v1;` / `option go_package = "github.com/qkitzero/fitness-service/gen/go/<ctx>/v1";` / `google.api.http` アノテーションで REST 併設（grpc-gateway）
- **レイヤ配置**: `internal/domain|application|infrastructure|interface/<ctx>/`
- **集約**: `interface`（getter＋振る舞い）＋非公開 `struct`＋`New<Entity>(...)` コンストラクタ。更新は集約メソッド（例: `Update(name)`）
- **値オブジェクト**: ID・Name 等を型化（例: `CustomerID`, `Name`）。ID 生成関数（例: `NewCustomerID()`）
- **Repository**: `internal/domain/<ctx>/repository.go` に IF、`internal/infrastructure/<ctx>/repository.go` に実装（GORM）
- **usecase**: `authService.VerifyToken(ctx)` を冒頭で呼ぶ。依存は IF で注入
- **migration**: `internal/infrastructure/db/migrations/<timestamp>_create_<table>_table.up.sql` / `.down.sql`。ID は `VARCHAR(36) PRIMARY KEY`、`created_at` / `updated_at TIMESTAMP NOT NULL`
- **外部参照**: staff=`user_id`、院=`group_id` は文字列 ID で保持（FK は張らない＝別サービス境界）
- **同一サービス内参照**: fitness-service 内のエンティティ間（例: `customers.organization_id` → `organizations.id`）は**FK を張る**。上記「FK を張らない」は別サービス境界に対する規約であり、サービス内には適用しない
- **テスト/モック**: `mocks/` に生成モック、各層に `*_test.go`

---

## 6. 未確定事項（実装時に判断）

- **スタッフの付加属性**（資格情報・連絡先・住所）を user-service User の拡張とするか、fitness-service 側に `StaffProfile`（user_id 参照）として持つか。
- **要素別ランク → トレーニングレベルの対応閾値**（ElementMenu のはしごから、どのランクでどのレベルを選ぶか）。調査レポート未確認のため、`RankStandard` と併せて定義が必要。
- **認可ポリシーの適用箇所**（`role=staff_restricted` での PII マスキング等）。現状 `UserService` は `ListMyGroups` のみを提供し `Membership.role` を取得していないため、role を使う認可（マスタ切替の admin 限定など）には user-service 連携の追加が必要。
- **立ち上がりテストの選択肢定義**（両足 50〜10cm / 片足 40〜0cm）をどこに持たせるか。`MeasurementItem` の値の型は `choice` として型だけ用意済み。選択肢の実体は Phase 1b（`MeasurementValue`）の設計時に確定させる。
- **測定項目と6要素（筋力・筋持久力・柔軟性・敏しょう性・バランス・移動能力）の対応**。1項目が複数要素に寄与しうるため多対多になる可能性がある。Phase 2 で `RankStandard` と併せて定義する。
- **サルコペニア度の要否**（§3.3 Judgment）。高齢者向け指標のため、就業者を対象とする健康経営で採用するかを判断する。
- **`MeasurementItem` の暫定メタデータ**（§3.3 の ⚠ 5項目の試行回数・左右区分）。実運用で確認し後続 migration で補正する。
- **業種特化した命名の見直し**。健康経営を題材とする（確定方針 #5）ため、「院 / 子院 / 院スタッフ / 来院」や `ClinicProfile` / `CustomerClinicProfile`（保険証番号・診察券番号・要介護度）といった医療特化の語彙が実態と合わない。`ClinicProfile` → `GroupProfile`、「診療時間 / 診療カレンダー」→「営業時間 / 稼働カレンダー」等への一括見直しを検討する。`Group` / `group_id` は user-service の語彙なので維持する（境界を跨ぐ翻訳を避けるため）。
- **`Visit` の「会計状態（未会計・会計済）」の要否**。接骨院の保険診療を前提とした属性であり、企業が費用を負担する健康経営では不要になる可能性がある。
- **`Organization` の管理者ログイン**。顧客の所属企業の担当者が自社従業員の集計を参照する要求が出た場合、`Organization.linked_group_id`（user-service Group への任意参照）で橋渡しする。先に作ると二重管理になるため、要求が出るまで作らない。

---

*本書は調査レポートのドメイン理解を、自作 fitness-service のマイクロサービス構成へ再配置・汎用化した設計指針である。実データ・画面の詳細は [health-check-system-research.md](./health-check-system-research.md) を参照。*
