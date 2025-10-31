import { useEffect, useState } from "react"
import type { User } from "@/lib/types"
import { useAuth } from "@/context/AuthContext"
import { toast } from "react-toastify"
import Loading from "@/components/Loading"

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { FormModal } from "@/components/FormModal"

const getRoleStyles = (role: string) => {
  switch (role) {
    case "owner":
      return "bg-purple-500/20 text-purple-400"
    case "admin":
      return "bg-blue-500/20 text-blue-400"
    default:
      return "bg-muted text-muted-foreground"
  }
}

export function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { user } = useAuth()

  const fetchUsers = async () => {
    try {
      const response = await fetch("/api/users/getAll")
      const data = await response.json()

      if (!data.success) throw new Error(data.error || "Failed to fetch users")

      const updatedUsers: User[] = data.data.map((u: User) => ({
        ...u,
        isAdmin: u.role === "admin" || u.role === "owner",
      }))

      setUsers(updatedUsers)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch users")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [])

  const handleCreateUser = async (userData: {
    username: string
    email: string
    password: string
    role: "admin" | "user"
  }) => {
    try {
      const response = await fetch("/api/users/create", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userData),
      })

      const data = await response.json()
      if (!data.success) toast.error(data.error || "Failed to create user")

      toast.success(data.message || "User created successfully")
      fetchUsers()
      setIsModalOpen(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create user")
    }
  }

  if (loading)
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <Loading />
      </div>
    )

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Users
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage users and their permissions
          </p>
        </div>
        <Button
          onClick={() => setIsModalOpen(true)}
          disabled={!user?.isAdmin}
          className="transition-colors"
        >
          Add User
        </Button>
      </div>

      {/* Error */}
      {error ? (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {users.map((u) => (
            <Card
              key={u.id}
              className="cursor-pointer border-border bg-card hover:border-primary transition-colors"
            >
              <CardHeader className="pb-2">
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted text-foreground">
                      {u.username[0].toUpperCase()}
                    </div>
                    <div>
                      <CardTitle className="text-lg font-semibold text-foreground">
                        {u.username}
                      </CardTitle>
                      <CardDescription className="text-sm text-muted-foreground">
                        {u.email}
                      </CardDescription>
                    </div>
                  </div>
                  <Badge
                    variant="secondary"
                    className={`capitalize ${getRoleStyles(u.role)}`}
                  >
                    {u.role}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="border-t border-border pt-3">
                <p className="text-sm text-muted-foreground font-mono break-all">
                  User ID: {u.id}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateUser}
        title="Create New User"
        fields={[
          { name: "username", label: "Username", type: "text", required: true },
          { name: "email", label: "Email", type: "email", required: true },
          { name: "password", label: "Password", type: "password", required: true },
          {
            name: "role",
            label: "Role",
            type: "select",
            options: [
              { label: "User", value: "user" },
              ...(user?.isAdmin ? [{ label: "Admin", value: "admin" }] : []),
            ],
            required: true,
          },
        ]}
      />
    </div>
  )
}
