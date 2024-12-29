import type { Metadata } from "next"
import { geistSans, geistMono } from './font'
import "./globals.css"
import "@radix-ui/themes/styles.css"
import { SessionProvider } from 'next-auth/react'
import ClientProvider from '@/components/provider/ClientProvider'
import { ThemeProvider } from 'next-themes'

export const metadata: Metadata = {
  title: "Chatbot Web App",
  description: "Created by @rizalord",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <SessionProvider>
          <ThemeProvider attribute="class">
            <ClientProvider>
              {children}
            </ClientProvider>
          </ThemeProvider>
        </SessionProvider>
      </body>
    </html>
  )
}
