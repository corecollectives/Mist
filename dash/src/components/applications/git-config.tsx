import { useEffect, useState } from "react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import type { App } from "@/types/app"
import { Github, Gitlab } from "lucide-react"
import { SiBitbucket, SiGitea } from "react-icons/si"
import { Skeleton } from "@/components/ui/skeleton"

interface GitProviderTabProps {
  app: App
}

export const GitProviderTab = ({ app }: GitProviderTabProps) => {
  const [provider, setProvider] = useState("github")

  // Github App State
  const [_, setGithubApp] = useState<any>(null)
  const [isInstalled, setIsInstalled] = useState<boolean>(false)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Repo + Branch
  const [repos, setRepos] = useState<any[]>([])
  const [branches, setBranches] = useState<any[]>([])
  const [selectedRepo, setSelectedRepo] = useState(app.gitRepository || "")
  const [selectedBranch, setSelectedBranch] = useState(app.gitBranch || "")
  const [isRepoLoading, setIsRepoLoading] = useState(true)
  const [isBranchLoading, setIsBranchLoading] = useState(false)
  // ✅ Fetch GitHub App installation
  const fetchApp = async () => {
    try {
      setIsLoading(true)
      setError(null)

      const response = await fetch("/api/github/app", { credentials: "include" })
      const data = await response.json()

      if (data.success) {
        setGithubApp(data.data.app)      // app metadata
        setIsInstalled(data.data.isInstalled) // boolean
      } else {
        setError(data.error || "Failed to load GitHub App details")
      }
    } catch (err) {
      setError("Failed to load GitHub App details")
    } finally {
      setIsLoading(false)
    }
  }

  const fetchRepos = async () => {
    try {
      setIsRepoLoading(true)
      const res = await fetch("/api/github/repositories", { credentials: "include" })
      const data = await res.json()
      setRepos(data || [])
    } catch {
      toast.error("Failed to load repositories")
    } finally {
      setIsRepoLoading(false)
    }
  }

  const fetchBranchList = async (repoFullName: string) => {
    try {
      setIsBranchLoading(true)
      const res = await fetch(`/api/github/branches`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ repo: repoFullName }),
        credentials: "include",
      })
      const data = await res.json()

      if (data.success) setBranches(data.data)
    } catch {
      toast.error("Failed to load branches")
    } finally {
      setIsBranchLoading(false)
    }
  }

  useEffect(() => {
    fetchApp()
    fetchRepos()
  }, [])

  useEffect(() => {
    if (selectedRepo) fetchBranchList(selectedRepo)
  }, [selectedRepo])

  useEffect(() => {
    if (app.gitRepository) {
      setRepos([{ id: 0, full_name: app.gitRepository }])
    }
    if (app.gitBranch) {
      setBranches([{ name: app.gitBranch }])
    }
  }, [app])

  const saveGitConfig = async () => {
    try {
      const res = await fetch("/api/apps/update", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          appId: app.id,
          gitRepository: selectedRepo,
          gitBranch: selectedBranch,
          gitProviderId: 1,
        }),
      })

      const data = await res.json()
      if (!data.success) throw new Error(data.error)

      toast.success("Git provider configuration saved")
    } catch {
      toast.error("Failed to save configuration")
    }
  }

  return (
    <Tabs defaultValue="github" value={provider} onValueChange={setProvider} className="w-full space-y-8">

      {/* ✅ PROVIDER LIST */}
      <TabsList className="grid w-full grid-cols-4">
        <TabsTrigger value="github" className="flex items-center gap-2">
          <Github className="h-4 w-4" />
          GitHub
        </TabsTrigger>

        <TabsTrigger value="gitlab" disabled className="flex items-center gap-2 opacity-70">
          <Gitlab className="h-4 w-4" />
          GitLab
        </TabsTrigger>

        <TabsTrigger value="gitea" disabled className="flex items-center gap-2 opacity-70">
          <SiGitea className="h-4 w-4" />
          Gitea
        </TabsTrigger>

        <TabsTrigger value="bitbucket" disabled className="flex items-center gap-2 opacity-70">
          <SiBitbucket className="h-4 w-4" />
          Bitbucket
        </TabsTrigger>
      </TabsList>

      {/* ✅ GITHUB TAB CONTENT */}
      <TabsContent value="github">
        {/* ✅ Loading state */}
        {isLoading && (
          <Card>
            <CardHeader>
              <CardTitle>Loading…</CardTitle>
            </CardHeader>
          </Card>
        )}

        {/* ✅ Error state */}
        {!isLoading && error && (
          <Card className="border-red-500">
            <CardHeader>
              <CardTitle className="text-red-500">Error Loading GitHub App</CardTitle>
              <CardDescription>{error}</CardDescription>
            </CardHeader>
          </Card>
        )}

        {/* ✅ NOT INSTALLED → Show connection card */}
        {!isLoading && !isInstalled && (
          <Card>
            <CardHeader>
              <CardTitle>GitHub App Not Connected</CardTitle>
              <CardDescription>
                You need to connect your GitHub App to enable repository syncing.
              </CardDescription>
            </CardHeader>

            <CardContent>
              <Button asChild>
                <a href="/git">Connect GitHub App</a>
              </Button>
            </CardContent>
          </Card>
        )}

        {!isLoading && isInstalled && (
          <Card>
            <CardHeader>
              <CardTitle>GitHub Repository</CardTitle>
              <CardDescription>
                Select the repository and branch to link with your app.
              </CardDescription>
            </CardHeader>

            <CardContent className="space-y-6">

              <div className="flex flex-col md:flex-row gap-6">

                {/* ✅ Repo select or skeleton */}
                <div className="flex-1">
                  <Label className="text-muted-foreground">Repository</Label>

                  {isRepoLoading ? (
                    <Skeleton className="w-full h-10 mt-2" />
                  ) : (
                    <Select value={selectedRepo} onValueChange={setSelectedRepo}>
                      <SelectTrigger className="mt-2 w-full">
                        <SelectValue placeholder="Select a repository" />
                      </SelectTrigger>

                      <SelectContent>
                        {repos.map((repo) => (
                          <SelectItem key={repo.id} value={repo.full_name}>
                            {repo.full_name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                </div>

                {/* ✅ Branch select or skeleton */}
                <div className="flex-1">
                  <Label className="text-muted-foreground">Branch</Label>

                  {isBranchLoading ? (
                    <Skeleton className="w-full h-10 mt-2" />
                  ) : (
                    <Select value={selectedBranch} onValueChange={setSelectedBranch}>
                      <SelectTrigger className="mt-2 w-full">
                        <SelectValue placeholder="Select a branch" />
                      </SelectTrigger>

                      <SelectContent>
                        {branches.map((branch) => (
                          <SelectItem key={branch.name} value={branch.name}>
                            {branch.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                </div>

              </div>

              <Button onClick={saveGitConfig} className="w-fit">
                Save Configuration
              </Button>
            </CardContent>
          </Card>
        )}
      </TabsContent>

      {/* ✅ OTHER PROVIDERS (Disabled) */}
      <TabsContent value="gitlab">
        <Card>
          <CardHeader>
            <CardTitle>GitLab Support</CardTitle>
            <CardDescription>GitLab integration is coming soon.</CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

      <TabsContent value="gitea">
        <Card>
          <CardHeader>
            <CardTitle>Gitea Support</CardTitle>
            <CardDescription>Gitea integration is coming soon.</CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

      <TabsContent value="bitbucket">
        <Card>
          <CardHeader>
            <CardTitle>Bitbucket Support</CardTitle>
            <CardDescription>Bitbucket integration is coming soon.</CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

    </Tabs>
  )
}
