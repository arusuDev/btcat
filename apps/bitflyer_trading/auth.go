package bitflyer_trading

// 実際に取引する場合に利用するbitflyerのclientを作成する。
// 以下関数を持つ
import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 受け取った文字列を結合し、HMAC-SHA256署名を作成する
func generateSign(secret, method, path, timestamp, body string) string {
	message := timestamp + method + path + body

	// HMAC(Hash-based Message Authentication Code) メッセージの整合性と認証を確保するためのハッシュ関数
	// secretを鍵として、新しいHMAC-SHA256ハッシュ関数を作成する
	h := hmac.New(sha256.New, []byte(secret))

	// 作成したハッシュ関数にメッセージを書き込む
	h.Write([]byte(message))

	// h.Sum(nil) → ハッシュ計算を実施
	// hex.EncodeToString → ハッシュ計算した値を16進文字列にエンコード
	return hex.EncodeToString(h.Sum(nil))
}

func GetBitflyer() *http.Client {
	// BitflyerのAPIキーを読み取り
	apiKey := os.Getenv("BF_KEY")
	apiSecret := os.Getenv("BF_SECRET")

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	method := "GET"
	baseURL := "https://api.bitflyer.com"
	endpoint := "/v1/me/getbalance"
	body := ""

	// HMAC-SHA256証明を作成
	sign := generateSign(apiSecret, method, endpoint, timestamp, body)

	// endpointに向けたGETリクエストを作成
	req, err := http.NewRequest(method, baseURL+endpoint, nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	// Headerに必要な情報を追加
	req.Header.Add("ACCESS-KEY", apiKey)
	req.Header.Add("ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("ACCESS-SIGN", sign)
	req.Header.Add("Content-Type", "application/json")

	// APIを実行するためのクライアントを作成
	// このクライアントを通してGET/PUTメソッドを実行する
	client := &http.Client{}
	return client
}
