import { useChatStore } from '@/store/chat-store'
import { ChatData } from '@/types/dto/chat.dto'
import { signOut, useSession } from 'next-auth/react'
import Link from 'next/link'
import React, { useEffect, useState } from "react"
import {
  AiOutlinePlus,
} from "react-icons/ai"
import { FiMessageSquare } from "react-icons/fi"
import InfiniteScroll from 'react-infinite-scroll-component'
import { Dropdown } from "flowbite-react"
import { HiLogout, HiChevronUp, HiLightBulb, HiMoon } from "react-icons/hi"
import Image from 'next/image'
import { useTheme } from 'next-themes'

const Sidebar = () => {
  const { chats, isLoading } = useChatStore()
  const [data, setData] = useState<ChatData[]>([])
  const [hasMore, setHasMore] = useState(true)
  const { data: session } = useSession()
  const { theme, setTheme } = useTheme()

  useEffect(() => {
    setData(chats.slice(0, 20))
    setHasMore(chats.length > 20)
  }, [chats])

  const fetchMoreData = () => {
    const currentLength = data.length
    const more = chats.slice(currentLength, currentLength + 20)
    setData([...data, ...more])
    setHasMore(data.length + more.length < chats.length)
  }

  return (
    <div className="scrollbar-trigger flex h-full w-full flex-1 items-start border-white/20">
      <nav className="flex h-full flex-1 flex-col space-y-1 p-2 relative w-full overflow-hidden">
        <Link
          href={`/`}
          className="flex py-3 px-3 items-center gap-3 rounded-md transition-colors duration-200 cursor-pointer text-sm mb-1 flex-shrink-0 border text-white md:text-black md:dark:text-white hover:bg-gray-500/10 md:hover:bg-gray-500/10 md:border-black/10 md:dark:border-white/20">
          <AiOutlinePlus className="h-4 w-4" />
          New Chat
        </Link>

        <div className="flex-col h-full flex-1 overflow-y-auto border-b border-black/20 dark:border-white/20">
          {
            isLoading ? (
              <div className="flex flex-1 justify-center items-center h-full">
                <div role="status">
                  <svg aria-hidden="true" className="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-gray-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor" />
                    <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill" />
                  </svg>
                  <span className="sr-only">Loading...</span>
                </div>
              </div>
            ) : chats.length === 0 ? (
              <div className="flex flex-1 justify-center items-center h-full">
                <div className="flex flex-col items-center gap-2">
                  <span className="text-base font-semibold text-gray-100 md:dark:text-gray-100 md:text-gray-800">No chats</span>
                  <span className="text-xs text-gray-200 md:dark:text-gray-200 md:text-gray-600">Create a new chat to start chatting</span>
                </div>
              </div>
            ) : (
              <InfiniteScroll
                dataLength={data.length}
                next={fetchMoreData}
                hasMore={hasMore}
                loader={
                  <div className="flex justify-center items-center py-2">
                    <div role="status">
                      <svg aria-hidden="true" className="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-gray-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor" />
                        <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill" />
                      </svg>
                      <span className="sr-only">Loading...</span>
                    </div>
                  </div>
                }
                className={"flex flex-col gap-y-1 gap-x-2 pb-2 text-sm dark:text-gray-100 text-gray-800"}
                height={'calc(100vh - 140px)'}
              >
                {
                  data.map((chat) => (
                    <Link
                      href={`/${chat.id}`}
                      className="flex py-3 px-3 items-center gap-3 relative rounded-md cursor-pointer break-all hover:pr-4 group text-white/70 md:text-black md:dark:text-white/70 hover:bg-gray-500/10 md:hover:bg-gray-500/10"
                      key={chat.id}>
                      <FiMessageSquare className="h-4 w-4" />
                      <div className="flex-1 text-ellipsis max-h-5 overflow-hidden relative whitespace-nowrap w-40">
                        {chat.topic}
                      </div>
                    </Link>
                  ))
                }
              </InfiniteScroll>
            )
          }
        </div>

        {/* <button
          onClick={() => signOut()}
          className="flex py-3 px-3 items-center gap-3 rounded-md transition-colors duration-200 cursor-pointer text-sm text-white md:text-black md:dark:text-white hover:bg-gray-500/10 md:hover:bg-gray-500/10">
          <MdLogout className="h-4 w-4" />
          Log out
        </button> */}

        {
          session && session.user && (
            <div className="flex w-full relative">
              <Dropdown
                prefix='dropdown'
                arrowIcon={false}
                label={
                  <div className="flex flex-1 flex-row w-full max-w-[270px] md:max-w-[200px] justify-center items-center">
                    {/* Circular Avatar */}
                    {session.user.image ? (
                      <Image
                        src={session.user.image}
                        alt="User Avatar"
                        className="w-8 h-8 rounded-full"
                        width={32}
                        height={32}
                      />
                    ) : (
                      <div className="flex items-center justify-center w-8 h-8 rounded-full bg-gray-300 text-gray-900 dark:bg-gray-700 dark:text-gray-100">
                        {session.user.name?.substring(0, 2).toUpperCase()}
                      </div>
                    )}

                    {/* Email */}
                    <div className="flex-1 ml-2 overflow-hidden max-w-[240px] md:max-w-[140px]">
                      <span className="block text-sm font-medium whitespace-nowrap overflow-hidden text-ellipsis">
                        {session.user.email}
                      </span>
                    </div>

                    {/* Dropdown Arrow */}
                    <HiChevronUp className="flex-shrink-0 h-5 w-5" />
                  </div>
                }
                color="transparent"
                placement="top"
                fullSized
              >
                <Dropdown.Header className="relative w-full">
                  <div className="max-w-[270px] md:max-w-[200px">
                    <span className="block text-sm truncate">{session.user.name}</span>
                    <span className="block text-sm font-medium truncate">{session.user.email}</span>
                  </div>
                </Dropdown.Header>

                {
                  theme === 'dark' ? (
                    <Dropdown.Item icon={HiLightBulb} onClick={() => setTheme('light')}>Toggle light mode</Dropdown.Item>
                  ) : (
                    <Dropdown.Item icon={HiMoon} onClick={() => setTheme('dark')}>Toggle dark mode</Dropdown.Item>
                  )
                }

                {/* <Dropdown.Item icon={HiMoon}>Toggle dark mode</Dropdown.Item> */}
                <Dropdown.Divider />
                <Dropdown.Item icon={HiLogout} onClick={() => signOut()}>Sign out</Dropdown.Item>
              </Dropdown>
            </div>
          )
        }
      </nav >
    </div >
  )
}

export default Sidebar
