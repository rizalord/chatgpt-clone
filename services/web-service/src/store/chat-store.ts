import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type { } from '@redux-devtools/extension'
import { ChatData } from '@/types/dto/chat.dto'

interface ChatState {
    chats: ChatData[]
    isLoading: boolean
    sets: (chats: ChatData[]) => void
    clear: () => void
    toggleLoading: () => void
    setLoading: (loading: boolean) => void
    push: (chat: ChatData) => void
}

export const useChatStore = create<ChatState>()(
    devtools(
        persist(
            (set) => ({
                chats: [],
                isLoading: true,
                sets: (chats: ChatData[]) => set({ chats }),
                clear: () => set({ chats: [] }),
                toggleLoading: () => set((state) => ({ isLoading: !state.isLoading })),
                setLoading: (loading: boolean) => set({ isLoading: loading }),
                push: (chat: ChatData) => set((state) => ({ chats: [...state.chats, chat] })),
            }),
            {
                name: 'chat-storage',
            },
        ),
    ),
)