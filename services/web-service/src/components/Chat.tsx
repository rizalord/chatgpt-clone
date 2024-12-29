'use client'

import useAutoResizeTextArea from '@/hooks/useAutoResizeTextArea'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { FiSend } from "react-icons/fi"
import { BsPlusLg } from "react-icons/bs"
import { RxHamburgerMenu } from "react-icons/rx"
import MessageText from './MessageText'
import { CreateMessageRequest, Message, MessagePart } from '@/types/dto/chat.dto'
import { Socket } from 'socket.io-client'
import * as Yup from 'yup'
import { Form, Formik, FormikHelpers } from 'formik'
import { Session } from 'next-auth'
import { useRouter } from 'next/navigation'

interface ChatProps {
  initialMessages?: Array<Message>
  toggleComponentVisibility: () => void
  apiUrl: string
  session: Session
  socket: Socket
  chatId?: string
}

export function Chat({ initialMessages = [], toggleComponentVisibility, session, socket, chatId }: ChatProps) {
  const [isTyping, setIsTyping] = useState(false)
  const [message, setMessage] = useState("")
  const [errorMessage, setErrorMessage] = useState("")
  const [conversation, setConversation] = useState<Message[]>(initialMessages)
  const [parts, setParts] = useState<string[]>([])
  const staticParts = useMemo<string[]>(() => [], [])
  const [isConnected, setIsConnected] = useState(socket.connected)
  const textAreaRef = useAutoResizeTextArea()
  const bottomOfChatRef = useRef<HTMLDivElement>(null)
  const router = useRouter()

  const initialValues = { message: '' }

  const validationSchema = Yup.object({
    message: Yup.string().required('Message is required'),
  })

  const onConnect = useCallback(() => { setIsConnected(true) }, [])

  const onMessageStart = useCallback(() => {
    setIsTyping(true)
    setParts([])
    staticParts.length = 0
  }, [staticParts])

  const onMessage = useCallback((part: string) => {
    const data = JSON.parse(part) as MessagePart
    setParts((prev) => [...prev, data.part])
    staticParts.push(data.part)
  }, [staticParts])

  const onMessageEnd = useCallback((part: string) => {
    const data = JSON.parse(part) as MessagePart

    const newMessage: Message = {
      id: 0,
      chat_id: data.chat_id,
      user_id: Number(session.user?.id) || 0,
      role: "model",
      content: staticParts.join()
    }

    if (chatId !== data.chat_id.toString()) {
      setTimeout(() => {
        router.push(`/${data.chat_id}`)
      }, 1000)
    } else {
      setConversation((prev) => [...prev, newMessage])
      setParts([])
      staticParts.length = 0

      setIsTyping(false)
    }
  }, [chatId, router, session.user?.id, staticParts])

  const onDisconnect = useCallback(() => { setIsConnected(false) }, [])

  const onError = useCallback((error: any) => { console.log("onError", error) }, [])

  const onSubmit = async (values: typeof initialValues, { resetForm }: FormikHelpers<typeof initialValues>) => {
    setIsTyping(true)
    setErrorMessage("")

    try {
      const request: CreateMessageRequest = {
        chat_id: parseInt(chatId || "0"),
        message: values.message
      }

      socket.emit("create_message", request)

      const newMessage: Message = {
        id: 0,
        chat_id: parseInt(chatId || "0"),
        user_id: Number(session.user?.id) || 0,
        role: "user",
        content: values.message
      }

      setConversation((prev) => [...prev, newMessage])

      resetForm()
    } catch (error) {
      setErrorMessage('An error occurred. Please try again later.')
    } finally {
      setIsTyping(false)
    }
  }

  useEffect(() => {
    if (chatId) {
      socket.emit('join_chat_room', { chat_id: chatId })
    }

    socket.on('connect', onConnect)
    socket.on(chatId ? `message_chat_${chatId}` : 'message', onMessage)
    socket.on('error', onError)
    socket.on('message_start', onMessageStart)
    socket.on('message_end', onMessageEnd)
    socket.on('disconnect', onDisconnect)

    return () => {
      socket.off('connect', onConnect)
      socket.off(chatId ? `message_chat_${chatId}` : 'message', onMessage)
      socket.off('error', onError)
      socket.off('message_start', onMessageStart)
      socket.off('message_end', onMessageEnd)
      socket.off('disconnect', onDisconnect)
    }
  }, [chatId, socket, onMessage, onMessageEnd, onMessageStart, onConnect, onDisconnect, onError])

  useEffect(() => {
    if (textAreaRef.current) {
      textAreaRef.current.style.height = "24px"
      textAreaRef.current.style.height = `${textAreaRef.current.scrollHeight}px`
    }
  }, [message, textAreaRef])

  useEffect(() => {
    if (bottomOfChatRef.current) {
      bottomOfChatRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [conversation])

  return (
    <div className="flex max-w-full flex-1 flex-col dark:bg-gray-room">
      {/* Header Mobile*/}
      <div className="sticky top-0 z-10 flex flex-row items-center justify-between border-b border-white/20 dark:bg-gray-800 pl-1 py-3 dark:text-gray-200 sm:pl-3 md:hidden">
        <button
          type="button"
          className="flex h-10 w-10 items-center justify-center rounded-md hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white dark:hover:text-white"
          onClick={toggleComponentVisibility}
        >
          <span className="sr-only">Open sidebar</span>
          <RxHamburgerMenu className="h-6 w-6" />
        </button>
        <h1 className="flex-1 text-center text-base font-medium">Chatbot</h1>
        <button type="button" className="px-3 flex" onClick={() => router.push('/')}>
          <BsPlusLg className="h-6 w-6" />
        </button>
      </div>
      {/*  End of Header Mobile */}

      <div className="relative h-full w-full transition-width flex flex-col overflow-hidden items-stretch flex-1">
        {/* Header */}
        <div
          className="flex w-full items-center justify-center text-center gap-1 p-3 text-gray-500 dark:text-gray-300 text-sm font-semibold"
          onClick={() => socket.emit("MESSAGE", { chat_id: 0, message: "Hello" })}
        >
          Model: Default (Gemini-1.5 Flash)
        </div>
        {/* End of Header */}

        {/* Conversations */}
        <div className="flex-1 overflow-auto">
          <div className="h-full">
            {conversation.length > 0 ? (
              <div className="flex flex-col items-center text-sm">
                {conversation.map((message, index) => (
                  <MessageText key={index} message={message} />
                ))}

                {
                  isTyping && (
                    <MessageText message={{
                      id: 0,
                      chat_id: 0,
                      user_id: Number(session.user?.id) || 0,
                      role: "model",
                      content: parts.join()
                    }} />
                  )
                }

                <div className="h-10" />
                <div ref={bottomOfChatRef}></div>
              </div>
            ) : (
              <div className="py-10 relative w-full flex flex-col h-full items-center justify-center">
                <h1 className="text-2xl sm:text-4xl font-semibold text-center flex gap-2 items-center justify-center">
                  Chatbot Web App
                </h1>
              </div>
            )}
          </div>
        </div>
        {/* End of Conversations */}

        {/* Chat Input */}
        <div className="flex flex-col justify-center items-center w-full dark:border-white/20 md:border-transparent md:dark:border-transparent">
          <Formik
            initialValues={initialValues}
            validationSchema={validationSchema}
            onSubmit={onSubmit}
          >
            {
              (props) => (
                <Form className="mx-2 flex flex-row gap-3 w-full last:mb-2 md:mx-4 md:last:mb-6 lg:mx-auto lg:max-w-2xl xl:max-w-3xl">
                  <div className="relative flex flex-col h-full w-full flex-1 items-stretch md:flex-col px-5 md:px-5 lg:px-10 xl:px-0">
                    {errorMessage ? (
                      <div className="mb-2 md:mb-0">
                        <div className="h-full flex ml-1 md:w-full md:m-auto md:mb-2 gap-0 md:gap-2 justify-center">
                          <span className="text-red-500 text-sm">{errorMessage}</span>
                        </div>
                      </div>
                    ) : null}
                    <div className="flex flex-col w-full py-3 flex-grow md:py-3.5 md:pl-4 relative rounded-md shadow-sm dark:shadow-sm bg-gray-100 dark:bg-gray-700">
                      <textarea
                        ref={textAreaRef}
                        value={props.values.message}
                        tabIndex={0}
                        data-id="root"
                        style={{
                          height: "24px",
                          maxHeight: "200px",
                          overflowY: "hidden",
                        }}
                        rows={1}
                        name="message"
                        placeholder="Send a message..."
                        className="m-0 w-full resize-none border-0 bg-transparent p-0 pr-7 focus:ring-0 focus-visible:ring-0 dark:bg-transparent pl-2 md:pl-0 focus:outline-none dark:focus:outline-none dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
                        onChange={(e) => {
                          props.setFieldValue('message', e.target.value)
                          setMessage(e.target.value)
                        }}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter' && !e.shiftKey) {
                            e.preventDefault()
                            props.submitForm()
                          }
                        }}
                        disabled={isTyping || props.isSubmitting || !isConnected}
                      >
                      </textarea>
                      <button
                        type="submit"
                        disabled={isTyping || props.values.message.length === 0 || props.isSubmitting || !isConnected}
                        className="absolute p-2 rounded-md bottom-1.5 md:bottom-2.5 right-1 md:right-2 disabled:hidden bg-gray-950 dark:bg-white hover:opacity-60"
                      >
                        <FiSend className="h-4 w-4 mr-1 text-white disabled:text-gray-950 dark:text-gray-950" />
                      </button>
                    </div>
                  </div>
                </Form>
              )
            }

          </Formik>
          <div className="px-3 pt-2 pb-3 text-center text-xs text-black/50 dark:text-white/50 md:px-4 md:pt-3 md:pb-6">
            <span>
              Chatbot may produce inaccurate information about people,
              places, or facts.
            </span>
          </div>
        </div>
        {/* End of Chat Input */}
      </div>
    </div>
  )
}
