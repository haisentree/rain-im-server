<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { useWebSocket } from './composables/useWebSocket'
import type { ChatMessage } from './types/message'

const {
  connected,
  connecting,
  messages,
  clientId,
  error,
  myId,
  connect,
  sendMessage,
  disconnect,
} = useWebSocket()

// ---- 登录状态 ----
const loggedIn = ref(false)
const loginClientId = ref('')
const gatewayUrl = ref('ws://localhost:5173')

function doLogin() {
  if (!loginClientId.value.trim()) return
  connect({ clientId: loginClientId.value.trim(), gatewayUrl: gatewayUrl.value })
  loggedIn.value = true
}

// ---- 聊天状态 ----
const targetId = ref('')
const inputText = ref('')
const chatContainer = ref<HTMLElement | null>(null)
const contacts = ref<Contact[]>([
  { id: '550e8400-e29b-41d4-a716-446655440001', name: 'Alice' },
  { id: '550e8400-e29b-41d4-a716-446655440002', name: 'Bob' },
  { id: '550e8400-e29b-41d4-a716-446655440003', name: 'Charlie' },
])

interface Contact {
  id: string
  name: string
}

// 当前会话的消息
function filteredMessages(): ChatMessage[] {
  if (!targetId.value) return []
  return messages.value.filter(
    (m) =>
      (m.sourceId === myId.value && m.targetId === targetId.value) ||
      (m.sourceId === targetId.value && m.targetId === myId.value)
  )
}

function scrollToBottom() {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

watch(() => messages.value.length, scrollToBottom)
watch(targetId, () => scrollToBottom())

function doSend() {
  const text = inputText.value.trim()
  if (!text || !targetId.value) return
  sendMessage(targetId.value, text)
  inputText.value = ''
  scrollToBottom()
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    doSend()
  }
}

function selectContact(contact: Contact) {
  targetId.value = contact.id
}

function doDisconnect() {
  disconnect()
  loggedIn.value = false
  messages.value = []
}
</script>

<template>
  <!-- 登录界面 -->
  <div v-if="!loggedIn" class="login-wrapper">
    <div class="login-card">
      <h1 class="login-title">Rain IM</h1>
      <p class="login-subtitle">即时通讯系统</p>

      <div class="login-field">
        <label>网关地址</label>
        <input v-model="gatewayUrl" placeholder="ws://localhost:5173" />
      </div>

      <div class="login-field">
        <label>用户 ID (UUID)</label>
        <input
          v-model="loginClientId"
          placeholder="例如: 550e8400-e29b-41d4-a716-446655440001"
          @keydown.enter="doLogin"
        />
      </div>

      <p v-if="error" class="login-error">{{ error }}</p>

      <button class="login-btn" :disabled="connecting" @click="doLogin">
        {{ connecting ? '连接中...' : '连接' }}
      </button>
    </div>
  </div>

  <!-- 聊天界面 -->
  <div v-else class="chat-wrapper">
    <!-- 左侧: 联系人列表 -->
    <div class="sidebar">
      <div class="sidebar-header">
        <span class="my-id">{{ clientId }}</span>
        <button class="disconnect-btn" @click="doDisconnect">断开</button>
      </div>
      <div class="sidebar-list">
        <div
          v-for="c in contacts"
          :key="c.id"
          class="contact-item"
          :class="{ active: targetId === c.id }"
          @click="selectContact(c)"
        >
          <div class="contact-avatar">{{ c.name.charAt(0) }}</div>
          <span class="contact-name">{{ c.name }}</span>
        </div>
      </div>
    </div>

    <!-- 右侧: 聊天区域 -->
    <div class="chat-area">
      <template v-if="targetId">
        <div class="chat-header">
          <span>{{ contacts.find((c) => c.id === targetId)?.name || targetId }}</span>
          <span class="connection-status" :class="connected ? 'online' : 'offline'">
            {{ connected ? '在线' : '离线' }}
          </span>
        </div>
        <div ref="chatContainer" class="chat-messages">
          <div
            v-for="msg in filteredMessages()"
            :key="msg.id"
            class="message-row"
            :class="msg.self ? 'self' : 'other'"
          >
            <div class="message-bubble" :class="msg.self ? 'self' : 'other'">
              <div class="message-text">{{ msg.content }}</div>
              <div class="message-meta">
                <span class="message-time">{{ new Date(msg.timestamp).toLocaleTimeString() }}</span>
                <span v-if="msg.seq > 0" class="message-seq">#{{ msg.seq }}</span>
              </div>
            </div>
          </div>
          <div v-if="filteredMessages().length === 0" class="empty-chat">
            暂无消息，发送第一条消息吧
          </div>
        </div>
        <div class="chat-input-area">
          <input
            v-model="inputText"
            class="chat-input"
            placeholder="输入消息，按 Enter 发送..."
            @keydown="onKeydown"
          />
          <button class="send-btn" :disabled="!inputText.trim()" @click="doSend">
            发送
          </button>
        </div>
      </template>
      <div v-else class="no-contact">
        选择一个联系人开始聊天
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ---- 登录页 ---- */
.login-wrapper {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  background: #fff;
  padding: 40px;
  border-radius: 12px;
  width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.login-title {
  text-align: center;
  font-size: 28px;
  color: #1677ff;
}

.login-subtitle {
  text-align: center;
  color: #999;
  margin-bottom: 32px;
}

.login-field {
  margin-bottom: 16px;
}

.login-field label {
  display: block;
  margin-bottom: 4px;
  font-size: 13px;
  color: #666;
}

.login-field input {
  width: 100%;
  padding: 10px 12px;
}

.login-error {
  color: #e74c3c;
  font-size: 13px;
  margin-bottom: 12px;
}

.login-btn {
  width: 100%;
  padding: 10px;
  background: #1677ff;
  color: #fff;
  border-radius: 6px;
  font-size: 15px;
  margin-top: 8px;
  transition: background 0.2s;
}

.login-btn:hover {
  background: #4096ff;
}

.login-btn:disabled {
  background: #a0c4ff;
  cursor: not-allowed;
}

/* ---- 聊天布局 ---- */
.chat-wrapper {
  height: 100%;
  display: flex;
}

.sidebar {
  width: 260px;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.my-id {
  font-size: 12px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 160px;
}

.disconnect-btn {
  font-size: 12px;
  color: #999;
  background: transparent;
  padding: 4px 8px;
  border-radius: 4px;
}

.disconnect-btn:hover {
  color: #e74c3c;
  background: #fff1f0;
}

.sidebar-list {
  flex: 1;
  overflow-y: auto;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
}

.contact-item:hover {
  background: #f5f5f5;
}

.contact-item.active {
  background: #e6f4ff;
}

.contact-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #1677ff;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.contact-name {
  font-size: 14px;
}

/* ---- 聊天区 ---- */
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.chat-header {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 15px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.connection-status {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 400;
}

.connection-status.online {
  color: #52c41a;
  background: #f6ffed;
}

.connection-status.offline {
  color: #999;
  background: #f5f5f5;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #fafafa;
}

.message-row {
  display: flex;
  margin-bottom: 16px;
}

.message-row.self {
  justify-content: flex-end;
}

.message-bubble {
  max-width: 70%;
  padding: 10px 14px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.5;
}

.message-bubble.self {
  background: #1677ff;
  color: #fff;
  border-bottom-right-radius: 4px;
}

.message-bubble.other {
  background: #fff;
  border: 1px solid #f0f0f0;
  border-bottom-left-radius: 4px;
}

.message-text {
  word-break: break-word;
}

.message-meta {
  display: flex;
  gap: 6px;
  margin-top: 4px;
  font-size: 11px;
  opacity: 0.7;
}

.message-bubble.self .message-meta {
  justify-content: flex-end;
}

.empty-chat {
  text-align: center;
  color: #bbb;
  margin-top: 120px;
  font-size: 14px;
}

.chat-input-area {
  display: flex;
  gap: 12px;
  padding: 12px 20px;
  border-top: 1px solid #f0f0f0;
}

.chat-input {
  flex: 1;
  padding: 10px 14px;
}

.send-btn {
  padding: 10px 24px;
  background: #1677ff;
  color: #fff;
  border-radius: 6px;
  font-size: 14px;
  transition: background 0.2s;
}

.send-btn:hover {
  background: #4096ff;
}

.send-btn:disabled {
  background: #d9d9d9;
  cursor: not-allowed;
}

.no-contact {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #bbb;
  font-size: 15px;
}
</style>
