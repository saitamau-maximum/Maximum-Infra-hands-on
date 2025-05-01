package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	"example.com/infrahandson/test/e2e"
	"github.com/stretchr/testify/assert"
)

type MessageDTO struct {
	Content string `json:"content"`
}

func TestChatWebSocketConnection(t *testing.T) {
	// Start the test server
	server := e2e.StartTestServer(t)

	// === ユーザーAの準備 ===
	clientA := e2e.NewHTTPClient(t)
	registerRespA := e2e.Register(t, clientA, server.URL, "userA", "userA@example.com", "password123")
	assert.Equal(t, http.StatusOK, registerRespA.StatusCode)

	// === ユーザーBの準備 ===
	clientB := e2e.NewHTTPClient(t)
	registerRespB := e2e.Register(t, clientB, server.URL, "userB", "userB@example.com", "password123")
	assert.Equal(t, http.StatusOK, registerRespB.StatusCode)

	// === ルーム作成（ユーザーA） ===
	createRoomReq := map[string]string{"name": "Test Room"}
	createRoomResp := e2e.PostJSON(t, clientA, server.URL+"/api/room", createRoomReq)
	assert.Equal(t, http.StatusOK, createRoomResp.StatusCode)

	var createRoomRes map[string]string
	e2e.DecodeJSON(t, createRoomResp.Body, &createRoomRes)
	roomPublicID := createRoomRes["roomPublicID"]
	assert.NotEmpty(t, roomPublicID)

	// === WebSocket接続 ===
	wsURL := "ws://" + server.Address + "/api/ws/" + roomPublicID
	connA, err := e2e.ConnectWebSocket(wsURL, &clientA.Jar)
	assert.NoError(t, err)
	defer connA.Close()

	connB, err := e2e.ConnectWebSocket(wsURL, &clientB.Jar)
	assert.NoError(t, err)
	defer connB.Close()

	// === ユーザーAがメッセージ送信 ===
	message := map[string]string{"content": "Hello, World!"}
	err = connA.WriteJSON(message)
	assert.NoError(t, err)
	fmt.Println("User A sent message:", message)

	// === ユーザーBがメッセージ受信 ===
	var receivedMessage MessageDTO
	err = connB.ReadJSON(&receivedMessage)
	fmt.Println("User B received message:", receivedMessage)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", receivedMessage.Content)
}
