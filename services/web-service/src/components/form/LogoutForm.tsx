"use client"

import { signOut, useSession } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'

export default function LogoutForm() {
    const { data: session } = useSession()
    const router = useRouter()

    useEffect(() => {
        if (session && session.user) {
            signOut({ redirectTo: '/login' })
        } else {
            router.push('/login')
        }
    }, [session])

    return (<div />)
}