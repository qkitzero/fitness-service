package address

import (
	"fmt"
	"strings"
)

var prefectures = map[string]bool{
	"北海道": true, "青森県": true, "岩手県": true, "宮城県": true, "秋田県": true,
	"山形県": true, "福島県": true, "茨城県": true, "栃木県": true, "群馬県": true,
	"埼玉県": true, "千葉県": true, "東京都": true, "神奈川県": true, "新潟県": true,
	"富山県": true, "石川県": true, "福井県": true, "山梨県": true, "長野県": true,
	"岐阜県": true, "静岡県": true, "愛知県": true, "三重県": true, "滋賀県": true,
	"京都府": true, "大阪府": true, "兵庫県": true, "奈良県": true, "和歌山県": true,
	"鳥取県": true, "島根県": true, "岡山県": true, "広島県": true, "山口県": true,
	"徳島県": true, "香川県": true, "愛媛県": true, "高知県": true, "福岡県": true,
	"佐賀県": true, "長崎県": true, "熊本県": true, "大分県": true, "宮崎県": true,
	"鹿児島県": true, "沖縄県": true,
}

type Prefecture string

func (p Prefecture) String() string {
	return string(p)
}

func NewPrefecture(s string) (*Prefecture, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if !prefectures[t] {
		return nil, fmt.Errorf("invalid prefecture")
	}
	p := Prefecture(t)
	return &p, nil
}
