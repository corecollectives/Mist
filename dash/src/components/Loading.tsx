import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"

export default function Loading({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        "flex h-screen w-full flex-col items-center justify-center bg-background text-foreground",
        className
      )}
    >
      <Card className="flex flex-col items-center justify-center border-none shadow-none bg-transparent">
        <CardContent className="flex flex-col items-center">
          <div className="h-16 w-16 rounded-full border-4 border-primary border-t-transparent animate-spin"></div>

          <p className="mt-6 text-xl font-semibold">Loading...</p>
          <p className="mt-2 text-sm text-muted-foreground">
            Please wait while we make some API calls
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
