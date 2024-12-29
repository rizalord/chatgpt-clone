import { SiOpenai } from "react-icons/si"
import { TbCursorText } from "react-icons/tb"
import { Message } from '@/types/dto/chat.dto'
import Markdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

interface MessageProps {
  message: Message
}

const MessageText = ({ message }: MessageProps) => {
  const { role, content: text } = message
  const isUser = role === "user"

  return (
    <div
      className={`group w-full text-gray-800 dark:text-gray-100`}
    >
      <div className="text-base gap-4 md:gap-6 md:max-w-2xl lg:max-w-xl xl:max-w-3xl flex lg:px-0 m-auto w-full">
        <div className="flex flex-row gap-4 md:gap-6 md:max-w-2xl lg:max-w-xl xl:max-w-3xl p-4 md:py-6 lg:px-0 m-auto w-full">
          <div className="w-8 flex flex-col relative items-end">
            {
              !isUser && (
                <div className="relative h-8 w-8 p-1.5 flex items-center justify-center text-opacity-100 border border-gray-200 dark:border-gray-600 rounded-full">
                  <SiOpenai className="h-5 w-5" />
                </div>
              )
            }
            <div className="text-xs flex items-center justify-center gap-1 absolute left-0 top-2 -ml-4 -translate-x-full group-hover:visible !invisible">
              <button
                disabled
                className="text-gray-300 dark:text-gray-400"
              ></button>
              <span className="flex-grow flex-shrink-0">1 / 1</span>
              <button
                disabled
                className="text-gray-300 dark:text-gray-400"
              ></button>
            </div>
          </div>
          <div className="relative flex w-[calc(100%-50px)] flex-col gap-1 md:gap-3 lg:w-[calc(100%-115px)]">
            <div className="flex flex-grow flex-col gap-3">
              <div className="flex flex-col items-start gap-4 break-words">
                {
                  isUser ? (
                    <div className="flex flex-col w-full items-end">
                      <p className="py-2 px-5 bg-gray-100 dark:bg-gray-800 rounded-full">{text}</p>
                    </div>
                  ) : !isUser && text === null ? (
                    <TbCursorText className="h-6 w-6 animate-pulse" />
                  ) : (
                    <Markdown className="prose w-full break-words dark:prose-invert dark:text-white" remarkPlugins={[remarkGfm]}>{text}</Markdown>
                  )
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    </div >
  )
}

export default MessageText
