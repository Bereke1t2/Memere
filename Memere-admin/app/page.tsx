import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

export default function Home() {
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6 p-8">
      <div className="flex items-center justify-between w-full max-w-lg">
        <h1 className="text-2xl font-semibold tracking-tight">Memere Admin</h1>
        <ThemeToggle />
      </div>

      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Phase 1 · Skill 1 — Scaffold</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <p className="text-muted-foreground text-sm">
            Next.js 15 + TypeScript strict + Tailwind v4 + shadcn/ui shell.
            Business screens are built in subsequent skills.
          </p>
          <div className="flex flex-wrap gap-2">
            <Badge>Next.js</Badge>
            <Badge variant="secondary">TypeScript</Badge>
            <Badge variant="outline">Tailwind v4</Badge>
            <Badge variant="outline">shadcn/ui</Badge>
            <Badge variant="outline">TanStack Query</Badge>
          </div>
          <Button variant="default">Sample Button</Button>
        </CardContent>
      </Card>
    </div>
  );
}
