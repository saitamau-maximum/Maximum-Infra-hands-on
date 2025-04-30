package e2e_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"example.com/infrahandson/test/e2e"
)

func TestRegisterAndLogin(t *testing.T) {
	// テストサーバーとHTTPクライアントの準備
	server := e2e.StartTestServer(t)
	client := e2e.NewHTTPClient(t)

	// テストデータ
	name := "Test User"
	email := "test@example.com"
	password := "securepassword"

	// --- Register ---
	resp := e2e.Register(t, client, server.URL, name, email, password)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200 on register, got %d: %s", resp.StatusCode, string(body))
	}

	// --- 認証確認（/api/user/me）---
	meResp, err := client.Get(server.URL + "/api/user/me")
	if err != nil {
		t.Fatalf("failed to call /api/user/me: %v", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(meResp.Body)
		t.Fatalf("expected status 200 from /api/user/me after register, got %d: %s", meResp.StatusCode, string(body))
	}

	// --- Cookie を初期化してログインからやり直し（別クライアント）---
	client2 := e2e.NewHTTPClient(t)
	// client2 に client のクッキーをコピー
	for _, cookie := range client.Jar.Cookies(resp.Request.URL) {
		client2.Jar.SetCookies(resp.Request.URL, []*http.Cookie{cookie})
	}
	resp = e2e.Login(t, client2, server.URL, email, password)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200 on login, got %d: %s", resp.StatusCode, string(body))
	}

	// --- 再度 /api/user/me で認証確認 ---
	meResp, err = client2.Get(server.URL + "/api/user/me")
	if err != nil {
		t.Fatalf("failed to call /api/me after login: %v", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(meResp.Body)
		t.Fatalf("expected status 200 from /api/user/me after login, got %d: %s", meResp.StatusCode, string(body))
	}

	// --- レスポンスの内容を軽く検証 ---
	var me map[string]interface{}
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil {
		t.Fatalf("failed to decode /api/user/me response: %v", err)
	}
	emailValue, ok := me["email"].(string)
	if !ok {
		t.Errorf("expected email to be string, got %T", me["email"])
	}
	if emailValue != email {
		t.Errorf("expected email %s in /api/user/me, got %v", email, me["email"])
	}
}

