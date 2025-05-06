package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/infrastructure/server"
	"github.com/gorilla/websocket"
)

type TestServer struct {
	Address string // "localhost:8080"
	URL     string // "http://localhost:8080"
}

func StartTestServer(t *testing.T) *TestServer {
	t.Helper()

	// 固定ポート番号（Cookieスコープを安定させるため）
	const port = 8080

	// ポートが使用中でないか確認
	addr := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("port %d is already in use: %v", port, err)
	}
	listener.Close() // Echoで使用するため一旦閉じる

	// 設定読み込み
	cfg := config.LoadConfig()
	cfg.Port = fmt.Sprintf("%d", port)

	e, db := server.ServerStart(cfg)
	t.Cleanup(func() { db.Close() })

	// Echoサーバー起動

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			panic("failed to start echo server: " + err.Error())
		}
	}()

	// 起動待機
	time.Sleep(100 * time.Millisecond)

	// テスト終了時にサーバー停止とDB削除
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = e.Shutdown(ctx)

		if err := os.Remove("./database.db"); err != nil && !os.IsNotExist(err) {
			t.Logf("failed to remove database.db: %v", err)
		}
	})

	return &TestServer{
		Address: addr,
		URL:     "http://" + addr,
	}
}

func NewHTTPClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to create cookiejar: %v", err)
	}
	return &http.Client{
		Jar: jar,
	}
}

func Register(t *testing.T, client *http.Client, baseURL, name, email, password string) *http.Response {
	t.Helper()

	// リクエストボディの作成
	body := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal register body: %v", err)
	}

	// POSTリクエストの作成
	req, err := http.NewRequest("POST", baseURL+"/api/user/register", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("failed to create register request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:5173")

	// クッキー付きでリクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send register request: %v", err)
	}

	return resp
}

func Login(t *testing.T, client *http.Client, baseURL, email, password string) *http.Response {
	t.Helper()

	body := map[string]string{
		"email":    email,
		"password": password,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal login body: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/api/user/login", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("failed to create login request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send login request: %v", err)
	}
	return resp
}

// PostJSON sends a POST request with a JSON body.
func PostJSON(t *testing.T, client *http.Client, url string, body interface{}) *http.Response {
	t.Helper()

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal JSON body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}

	return resp
}

// DecodeJSON decodes a JSON response body into the provided interface.
func DecodeJSON(t *testing.T, body io.ReadCloser, v interface{}) {
	t.Helper()
	defer body.Close()

	if err := json.NewDecoder(body).Decode(v); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
}

// ConnectWebSocket establishes a WebSocket connection.
func ConnectWebSocket(url string, jar *http.CookieJar) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		Jar: *jar,
	}
	conn, _, err := dialer.Dial(url, nil)
	return conn, err
}
