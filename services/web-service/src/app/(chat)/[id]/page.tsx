import { auth } from '@/auth'
import Room from '@/components/Room'
import { ChatData, GetChatsResponse, GetMessagesResponse, Message } from '@/types/dto/chat.dto'
import { Response } from '@/types/dto/default.dto'
import { notFound, redirect } from 'next/navigation'

export default async function ExistingChat({
    params,
}: {
    params: Promise<{ id: string }>
}) {
    const id = (await params).id

    const apiUrl = process.env.API_URL || ''
    const publicApiUrl = process.env.PUBLIC_API_URL || apiUrl

    const session = await auth()

    if (!session?.user) redirect('/login')

    if (session?.error === "RefreshTokenError") redirect('/logout')

    const chats: ChatData[] = []

    const response = await fetch(`${apiUrl}/chats`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${session.token.access_token}`,
        },
    })

    const data: Response<GetChatsResponse> = await response.json()

    if (!response.ok) {
        if (response.status === 401) {
            redirect('/logout')
        }

        if (response.status === 403) throw new Error('Forbidden')
        if (response.status === 404) notFound()
        if (response.status === 500) throw new Error('Internal server error')

        throw new Error(data.message || 'Invalid credentials')
    }

    chats.push(...data.data.chats.reverse())

    const initialMessages: Message[] = []

    if (id) {
        const response = await fetch(`${apiUrl}/chats/${id}/messages`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${session.token.access_token}`,
            },
        })

        const data: Response<GetMessagesResponse> = await response.json()

        if (!response.ok) {
            if (response.status === 401) {
                redirect('/logout')
            }

            if (response.status === 403) throw new Error('Forbidden')
            if (response.status === 404) notFound()
            if (response.status === 500) throw new Error('Internal server error')

            throw new Error(data.message || 'Invalid credentials')
        }

        initialMessages.push(...data.data.messages.sort((a, b) => a.id - b.id))
    }

    return (
        <>
            <Room chats={chats} apiUrl={publicApiUrl} chatId={id} initialMessages={initialMessages} />
        </>
    )
}