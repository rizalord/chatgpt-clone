"use client"

import { useEffect, useState } from 'react'
import MobileSidebar from '@/components/sidebar/MobileSidebar'
import Sidebar from '@/components/sidebar/Sidebar'
import { Chat } from '@/components/Chat'
import { ChatData, Message } from '@/types/dto/chat.dto'
import { useChatStore } from '@/store/chat-store'
import { signOut, useSession } from 'next-auth/react'
import { io, Socket } from 'socket.io-client'

interface RoomProps {
    chatId?: string
    chats: ChatData[]
    initialMessages?: Array<Message>
    apiUrl: string
}

export default function Room({ chats: initialChats, apiUrl, initialMessages = [], chatId }: RoomProps) {
    const [isComponentVisible, setIsComponentVisible] = useState(false)
    const { data: session } = useSession()
    const { sets, setLoading } = useChatStore()
    const [socket, setSocket] = useState<Socket | null>(null)

    const toggleComponentVisibility = () => {
        setIsComponentVisible(!isComponentVisible)
    }

    useEffect(() => {
        sets(initialChats)
        setLoading(false)
    }, [initialChats, sets, setLoading])

    useEffect(() => {
        const handleSignOut = async () => {
            if (session?.error !== "RefreshTokenError") return
            await signOut({ redirectTo: '/login' })
        }
        handleSignOut()
    }, [session?.error])

    useEffect(() => {
        // Initialize socket
        const socketInstance = io(apiUrl, {
            transports: ['websocket'],
            auth: {
                token: session?.token.access_token
            },
            query: {
                token: session?.token.access_token
            },
            autoConnect: false,
        })

        socketInstance.connect()
        setSocket(socketInstance)

        return () => {
            socketInstance.disconnect()
        }
    }, [apiUrl, session])

    return (
        <main className="overflow-hidden w-full h-screen relative flex">
            {isComponentVisible ? (
                <MobileSidebar toggleComponentVisibility={toggleComponentVisibility} />
            ) : null}
            <div className="hidden flex-shrink-0 bg-gray-50 dark:bg-gray-sidebar md:flex md:w-[260px] md:flex-col">
                <div className="flex h-full min-h-0 flex-col ">
                    <Sidebar />
                </div>
            </div>
            {
                session && socket && (
                    <Chat
                        toggleComponentVisibility={toggleComponentVisibility}
                        initialMessages={initialMessages}
                        apiUrl={apiUrl}
                        session={session}
                        socket={socket}
                        chatId={chatId}
                    />
                )
            }
        </main>
    )
}