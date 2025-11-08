import { useEffect, useRef, useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Terminal, X } from "lucide-react"

interface Props {
  deploymentId: number
  open: boolean
  onClose: () => void
}

export const DeploymentLogsOverlay = ({ deploymentId, open, onClose }: Props) => {
  const [logs, setLogs] = useState("")
  const bottomRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  // ✅ Scroll every time logs change
  useEffect(() => {
    scrollToBottom()
  }, [logs])

  useEffect(() => {
    if (!open) return

    const ws = new WebSocket(`/api/ws/logs?id=${deploymentId}`)

    ws.onmessage = (event) => {
      if (event.data === "") return
      const cleaned = event.data
        .replace(/\r/g, "")          // remove CR if present
        .replace(/\n\s*\n+/g, "\n")  // collapse multiple newlines

      setLogs((prev) => prev + cleaned)
    }

    return () => ws.close()
  }, [open, deploymentId])

  const handleClose = () => {
    setLogs("")
    onClose()
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent
        showCloseButton={false}
        className="
          w-full
          max-w-[80vw]
          sm:max-w-[80vw]
          lg:max-w-[80vw]
          h-[90vh]
          p-0
          rounded-xl
          overflow-hidden
          border border-border
          bg-background/90
          backdrop-blur-xl
          shadow-2xl
        gap-0
        "
      >
        <DialogHeader className="px-5 py-4 border-b bg-background/60 backdrop-blur flex flex-row justify-between items-center">
          <DialogTitle className="flex items-center gap-2 text-lg">
            <Terminal className="h-5 w-5 text-primary" />
            Deployment Logs
          </DialogTitle>

          <div className="flex items-center gap-3">
            <Badge
              variant="outline"
              className="font-mono text-xs px-2 py-0.5 tracking-wider"
            >
              #{deploymentId}
            </Badge>

            <button
              onClick={handleClose}
              className="p-1.5 hover:bg-muted rounded-md transition"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        </DialogHeader>

        <div
          className="
            text-white
            font-mono
            h-full
            overflow-auto
          p-2
            whitespace-pre-wrap
            relative
          text-sm
          leading-normal
          "
        >

          {logs.length === 0 ? (
            <div className="text-center text-neutral-500 mt-12">
              Waiting for logs…
            </div>
          ) : (
            logs
          )}

          <div ref={bottomRef} />
        </div>

        {/* ✅ Footer */}
        <div className="px-5 py-3 border-t bg-background/60 backdrop-blur flex justify-end">
          <Button onClick={handleClose} className="px-6">
            Close
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
