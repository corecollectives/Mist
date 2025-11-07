import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { MoreVertical, GitBranch, Github, Server, Globe } from "lucide-react";
import type { App } from "@/types/app";

interface AppCardProps {
  app: App;
  onClick: () => void;
  onEdit?: (app: App) => void;
  onDelete?: (appId: number) => void;
  canEdit?: boolean;
  canDelete?: boolean;
}

export const AppCard: React.FC<AppCardProps> = ({
  app,
  onClick,
  onEdit,
  onDelete,
  canEdit = true,
  canDelete = true,
}) => {
  return (
    <Card
      className="relative transition-colors hover:border-primary cursor-pointer group"
      onClick={onClick}
    >
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div className="min-w-0 flex-1">
          <CardTitle className="truncate">{app.name}</CardTitle>
          <CardDescription className="truncate">
            {app.description || "No description provided"}
          </CardDescription>
        </div>

        {(canEdit || canDelete) && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => e.stopPropagation()}
              >
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {canEdit && (
                <DropdownMenuItem
                  onClick={(e) => {
                    e.stopPropagation();
                    onEdit?.(app);
                  }}
                >
                  Edit
                </DropdownMenuItem>
              )}
              {canDelete && (
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={(e) => {
                    e.stopPropagation();
                    onDelete?.(app.id);
                  }}
                >
                  Delete
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </CardHeader>

      <CardContent>
        {/* Tags or metadata */}
        <div className="flex flex-wrap gap-2 mb-3">
          {app.status && (
            <Badge
              variant={
                app.status === "running"
                  ? "default"
                  : app.status === "error"
                    ? "destructive"
                    : "secondary"
              }
            >
              {app.status}
            </Badge>
          )}
          {app.deploymentStrategy && (
            <Badge variant="secondary">{app.deploymentStrategy}</Badge>
          )}
        </div>

        {/* Git Info */}
        {(app.gitRepository || app.gitBranch) && (
          <div className="flex flex-col gap-2 text-xs text-muted-foreground mb-3">
            {app.gitRepository && (
              <div className="flex items-center gap-2">
                <Github className="w-4 h-4 text-muted-foreground" />
                <span className="truncate">{app.gitRepository}</span>
              </div>
            )}
            {app.gitBranch && (
              <div className="flex items-center gap-2">
                <GitBranch className="w-4 h-4 text-muted-foreground" />
                <span>{app.gitBranch}</span>
              </div>
            )}
          </div>
        )}

        <Separator className="my-3" />

        {/* Footer info */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between text-xs text-muted-foreground gap-2">
          <div className="flex items-center gap-2">
            <Server className="w-4 h-4 text-muted-foreground" />
            <span>
              Port:{" "}
              <span className="text-foreground font-medium">
                {app.port || "—"}
              </span>
            </span>
          </div>
          {app.createdAt && (
            <div className="flex items-center gap-2">
              <Globe className="w-4 h-4 text-muted-foreground" />
              <span>
                Created:{" "}
                {new Date(app.createdAt).toLocaleDateString(undefined, {
                  day: "numeric",
                  month: "short",
                  year: "numeric",
                })}
              </span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
