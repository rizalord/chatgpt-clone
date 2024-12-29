import NextAuth, { Session } from "next-auth"
import Credentials from "next-auth/providers/credentials"
import { Response } from './types/dto/default.dto'
import { AuthResponse } from './types/dto/auth.dto'
import moment from 'moment'
import Google from "next-auth/providers/google"

async function refreshToken(refresh_token: string): Promise<Session> {
    const apiUrl = process.env.API_URL || ''

    try {
        const response = await fetch(`${apiUrl}/auth/refresh`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ refresh_token }),
        })

        const data: Response<AuthResponse> = await response.json()

        if (!response.ok) {
            throw new Error(data.message || 'Invalid credentials')
        }

        return {
            user: {
                name: data.data.user.name,
                email: data.data.user.email,
                image: data.data.user.image_url,
            },
            token: data.data.token,
            expires: moment.unix(data.data.token.expired_at).subtract(1, 'minute').toDate().toISOString(),
        } as Session
    } catch (_) {
        throw new Error('An error occurred. Please try again later.')
    }
}

async function loginGoogle(id_token: string): Promise<Session> {
    const apiUrl = process.env.API_URL || ''

    try {
        const response = await fetch(`${apiUrl}/auth/login/google`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id_token }),
        })

        const data: Response<AuthResponse> = await response.json()

        if (!response.ok) {
            throw new Error(data.message || 'Invalid credentials')
        }

        return {
            user: {
                name: data.data.user.name,
                email: data.data.user.email,
                image: data.data.user.image_url,
            },
            token: data.data.token,
            expires: moment.unix(data.data.token.expired_at).subtract(1, 'minute').toDate().toISOString(),
        } as Session
    } catch (_) {
        throw new Error('An error occurred. Please try again later.')
    }
}

export const { handlers, signIn, signOut, auth } = NextAuth({
    providers: [
        Credentials({
            id: "register",
            credentials: {
                name: { label: "Name", type: "text" },
                email: { label: "Email", type: "email" },
                image_url: { label: "Image URL", type: "text" },
                access_token: { label: "Access Token", type: "text" },
                refresh_token: { label: "Refresh Token", type: "text" },
                expired_at: { label: "Expired At", type: "text" },
            },
            authorize: async (credentials) => {
                return {
                    user: {
                        name: credentials.name,
                        email: credentials.email,
                        image: credentials.image_url,
                    },
                    token: {
                        access_token: credentials.access_token,
                        refresh_token: credentials.refresh_token,
                        expired_at: credentials.expired_at,
                    },
                    expires: moment.unix(Number(credentials.expired_at)).subtract(1, 'minute').toDate().toISOString(),
                } as any
            },
        }),
        Credentials({
            id: "login",
            credentials: {
                email: { label: "Email", type: "email" },
                password: { label: "Password", type: "password" },
            },
            authorize: async (credentials) => {
                const apiUrl = process.env.API_URL || ''

                try {

                    const response = await fetch(`${apiUrl}/auth/login`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(credentials),
                    })

                    const data: Response<AuthResponse> = await response.json()

                    if (!response.ok) {
                        throw new Error(data.message || 'Invalid credentials')
                    }

                    const session: Session = {
                        user: {
                            name: data.data.user.name,
                            email: data.data.user.email,
                            image: data.data.user.image_url,
                        },
                        token: data.data.token,
                        expires: moment.unix(data.data.token.expired_at).subtract(1, 'minute').toDate().toISOString(),
                    }

                    return session as any
                } catch (_) {
                    throw new Error('An error occurred. Please try again later.')
                }
            }
        }),
        Google,
    ],
    callbacks: {
        async jwt({ token, user, account, session, trigger }) {
            if (account?.provider === "google" && account.id_token) {
                try {
                    user = await loginGoogle(account.id_token) as any
                } catch (error) {
                    console.error(error)
                }
            }

            if (token) {
                const session = token as unknown as Session

                // check if token is expired
                if (session.token && moment().isAfter(moment.unix(session.token.expired_at))) {
                    try {
                        // const response = await refreshToken(session.token.accessToken, session.token.refreshToken)
                        const response = await refreshToken(session.token.refresh_token)

                        user = {
                            ...user,
                            ...response
                        }

                    } catch (error) {
                        console.log(error)

                        return {
                            ...token,
                            error: 'RefreshTokenError'
                        }
                    }
                }
            }

            if (trigger === 'update' && session) {
                token = {
                    ...token,
                    ...session
                }
            }

            if (user) {
                token = {
                    ...token,
                    ...user
                }
            }

            return token
        },
        async session({ session, token }) {
            session = token as any

            return session
        },
        authorized: async ({ auth }) => {
            // Logged in users are authenticated, otherwise redirect to login page
            return !!auth
        },
    },
    pages: {
        signIn: '/login',
    }
})
