export interface UUID {
  hi: string
  lo: string
}

export interface ConnectRequest {
  client_id: string
  platform: number
}

export interface RawMessage {
  type: MessageType
  data: string  // JSON string of inner message
}

export interface SingleMessage {
  source_id: UUID
  target_id: UUID
  seq: string
  content: string
}

export interface RelayMessage {
  source_id: UUID
  target_id: UUID
  seq_id: string
  content: string
}

export enum MessageType {
  MESSAGE_UNSPECIFIED = 0,
  MESSAGE_SINGLE = 1,
  MESSAGE_SINGLE_GATEWAY_RELAY = 2,
  MESSAGE_DB_SAVE = 3,
  MESSAGE_SINGLE_GROUP = 4,
  MESSAGE_PULL = 5,
  MESSAGE_PULL_SEQ = 6,
  MESSAGE_PUSH = 7,
  MESSAGE_RELAY = 8,
  MESSAGE_RELAY_PACKATE = 9,
  MESSAGE_ACTION = 10,
  MESSAGE_ACTION_MODIFY = 11,
  MESSAGE_ACTION_CONFIRM = 12,
  MESSAGE_ACTION_READED = 13,
  MESSAGE_ACTION_RECALL = 14,
}

export enum Platform {
  PLATFORM_UNSPECIFIED = 0,
  PLATFORM_SYSTEM = 1,
  PLATFORM_WEB = 2,
  PLATFORM_ANDROID = 3,
  PLATFORM_IPHONE = 4,
  PLATFORM_WINDOWS = 5,
  PLATFORM_MACOS = 6,
  PLATFORM_LINUX = 7,
  PLATFORM_MOBILE = 8,
  PLATFORM_IPAD = 9,
}

// 将标准 UUID 字符串转换为 protobuf UUID 格式（hi/lo 由 high/low 64 bits 组成）
export function uuidToProto(uuid: string): UUID {
  const hex = uuid.replace(/-/g, '')
  const hi = BigInt('0x' + hex.substring(0, 16))
  const lo = BigInt('0x' + hex.substring(16, 32))
  return { hi: hi.toString(), lo: lo.toString() }
}

// 将 protobuf UUID 转换为标准 UUID 字符串
export function protoToUuid(uuid: UUID): string {
  const hi = BigInt(uuid.hi).toString(16).padStart(16, '0')
  const lo = BigInt(uuid.lo).toString(16).padStart(16, '0')
  const hex = hi + lo
  return [
    hex.substring(0, 8),
    hex.substring(8, 12),
    hex.substring(12, 16),
    hex.substring(16, 20),
    hex.substring(20, 32),
  ].join('-')
}

export interface ChatMessage {
  id: string
  sourceId: string
  targetId: string
  content: string
  seq: number
  timestamp: number
  self: boolean
}
