import { ref, computed } from 'vue'
import type { ConnectRequest, RawMessage, SingleMessage, ChatMessage, UUID } from '../types/message'
import { MessageType, Platform, protoToUuid, uuidToProto } from '../types/message'

export function useWebSocket() {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const connecting = ref(false)
  const messages = ref<ChatMessage[]>([])
  const clientId = ref('')
  const error = ref('')

  const myId = computed(() => clientId.value)

  function connect(params: { clientId: string; gatewayUrl: string }) {
    if (ws.value) {
      ws.value.close()
    }

    connecting.value = true
    error.value = ''
    clientId.value = params.clientId

    const token: ConnectRequest = {
      client_id: params.clientId,
      platform: Platform.PLATFORM_WEB,
    }

    const tokenStr = JSON.stringify(token)
    const url = `${params.gatewayUrl}/gateway?token=${encodeURIComponent(tokenStr)}`

    try {
      const socket = new WebSocket(url)
      ws.value = socket

      socket.onopen = () => {
        connected.value = true
        connecting.value = false
      }

      socket.onmessage = (event) => {
        // 忽略服务端测试回显 "send"
        if (event.data === 'send') return

        try {
          const raw: RawMessage = JSON.parse(event.data)
          handleRawMessage(raw)
        } catch {
          // 非 JSON 消息忽略
        }
      }

      socket.onclose = () => {
        connected.value = false
        connecting.value = false
      }

      socket.onerror = () => {
        error.value = 'WebSocket 连接失败'
        connecting.value = false
      }
    } catch (e) {
      error.value = '创建连接失败: ' + (e as Error).message
      connecting.value = false
    }
  }

  function handleRawMessage(raw: RawMessage) {
    // 处理单聊消息（服务端 relay 时会转回 MESSAGE_SINGLE）
    if (raw.type === MessageType.MESSAGE_SINGLE || raw.type === MessageType.MESSAGE_SINGLE_GATEWAY_RELAY) {
      try {
        const sm: SingleMessage = JSON.parse(raw.data)
        const fromId = protoToUuid(sm.source_id)
        const toId = protoToUuid(sm.target_id)

        messages.value.push({
          id: `${fromId}-${sm.seq}`,
          sourceId: fromId,
          targetId: toId,
          content: sm.content,
          seq: Number(sm.seq),
          timestamp: Date.now(),
          self: fromId === myId.value,
        })
      } catch {
        // 解析失败忽略
      }
    }
  }

  function sendMessage(targetUuid: string, content: string) {
    if (!ws.value || !connected.value) return

    const sourceId: UUID = uuidToProto(myId.value)
    const targetId: UUID = uuidToProto(targetUuid)

    const singleMsg: Record<string, unknown> = {
      source_id: sourceId,
      target_id: targetId,
      content: content,
    }

    const rawMsg: Record<string, unknown> = {
      type: MessageType.MESSAGE_SINGLE,
      data: JSON.stringify(singleMsg),
    }

    ws.value.send(JSON.stringify(rawMsg))

    messages.value.push({
      id: `${myId.value}-${Date.now()}`,
      sourceId: myId.value,
      targetId: targetUuid,
      content,
      seq: 0,
      timestamp: Date.now(),
      self: true,
    })
  }

  function disconnect() {
    ws.value?.close()
    ws.value = null
    connected.value = false
    connecting.value = false
  }

  return {
    ws,
    connected,
    connecting,
    messages,
    clientId,
    error,
    myId,
    connect,
    sendMessage,
    disconnect,
  }
}
