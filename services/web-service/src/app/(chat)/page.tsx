import { auth } from '@/auth'
import Room from '@/components/Room'
import { ChatData, GetChatsResponse } from '@/types/dto/chat.dto'
import { Response } from '@/types/dto/default.dto'
import { notFound, redirect } from 'next/navigation'

export default async function NewChat() {
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

    return (
        <>
            <Room chats={chats} apiUrl={publicApiUrl} />
        </>
    )
}