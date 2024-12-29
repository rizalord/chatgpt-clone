import NextAuth, { DefaultSession } from "next-auth"

declare module "next-auth" {
    /**
     * Returned by `useSession`, `getSession` and received as a prop on the `SessionProvider` React Context
     */
    interface Session {
        token: {
            access_token: string
            refresh_token: string
            expired_at: number
        },
        error?: 'RefreshTokenError'
    }
}