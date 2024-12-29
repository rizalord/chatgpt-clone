export interface GetChatsRequest {
    user_id: number
}

export interface GetChatRequest {
    chat_id: number
    user_id: number
}

export interface CreateMessageRequest {
    chat_id?: number
    message: string
}

export interface GetMessagesRequest {
    chat_id: number
    user_id: number
}

export interface ChatData {
    id: number
    user_id: number
    topic: string
}

export interface GetChatsResponse {
    chats: ChatData[]
}

export interface Message {
    id: number
    chat_id: number
    user_id: number
    role: "user" | "model"
    content: string
}

export interface GetMessagesResponse {
    messages: Message[]
}

export enum MessageStatus {
    STARTED = 1,
    ON_PROGRESS = 2,
    FINISHED = 3,
}

export interface MessagePart {
    chat_id: number
    part: string
    status: MessageStatus
}
