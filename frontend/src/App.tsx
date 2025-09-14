import { useState, type ReactNode } from "react";
import Input from "./components/Input";
import Button from "./components/Button";
import Card from "./components/Card";
import ErrorMessage from "./components/ErrorMessage";
import Loader from "./components/Loader";
import { isLikelyUrl } from "./utils/url";


type Headings = Partial<Record<"h1" | "h2" | "h3" | "h4" | "h5" | "h6", number>>;

interface AnalyzedPage {
  url?: string;
  html_version?: string;
  title?: string;
  headings?: Headings;
  links_internal?: number;
  links_external?: number;
  links_inaccessible?: number;
  login_form_present?: boolean;
  // Allow unknown keys without breaking
  [k: string]: unknown;
}


function Badge({ children }: { children: ReactNode }) {
  return (
    <span className="inline-flex items-center rounded-full border border-slate-600/60 bg-slate-800/60 px-2.5 py-0.5 text-xs font-medium text-indigo-200">
      {children}
    </span>
  );
}

function StatCard({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="rounded-xl border border-slate-700/60 bg-slate-900/50 p-4">
      <div className="text-xs uppercase tracking-wide text-slate-400">{label}</div>
      <div className="mt-1 text-2xl font-semibold text-slate-100">{value}</div>
    </div>
  );
}

function HeadingsTable({ headings }: { headings?: Headings }) {
  const keys: Array<keyof Headings> = ["h1", "h2", "h3", "h4", "h5", "h6"];
  return (
    <div className="rounded-xl border border-slate-700/60 bg-slate-900/50 p-4">
      <div className="mb-2 text-sm font-semibold text-slate-200">Headings</div>
      <div className="grid grid-cols-2 gap-y-2 text-sm">
        <div className="text-slate-400">Level</div>
        <div className="text-slate-400">Count</div>
        {keys.map((k) => (
          <div className="contents" key={k as string}>
            <div className="text-slate-200">{String(k).toUpperCase()}</div>
            <div className="text-slate-100">{Number(headings?.[k] ?? 0)}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

function ResultView({ data }: { data: AnalyzedPage }) {
  const [showRaw, setShowRaw] = useState(false);
  const [copied, setCopied] = useState(false);

  const title =
    (data.title && String(data.title).trim()) ||
    (data.url && String(data.url)) ||
    "Analyzed Page";

  const url = data.url ?? "";
  const htmlVersion = data.html_version ?? "—";

  async function copyJSON() {
    try {
      await navigator.clipboard.writeText(JSON.stringify(data, null, 2));
      setCopied(true);
      setTimeout(() => setCopied(false), 1200);
    } catch {
      // non-blocking
    }
  }

  return (
    <Card>
      {/* Header */}
      <div className="mb-4">
        <div className="flex flex-wrap items-center gap-2">
          <h2 className="text-xl font-semibold text-slate-100">{title}</h2>
          <Badge>{htmlVersion}</Badge>
        </div>
        {url && (
          <a
            href={url}
            target="_blank"
            rel="noreferrer"
            className="mt-1 block max-w-full break-words text-sm text-sky-300 hover:underline"
          >
            {url}
          </a>
        )}
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard label="Internal Links" value={Number(data.links_internal ?? 0)} />
        <StatCard label="External Links" value={Number(data.links_external ?? 0)} />
        <StatCard
          label="Inaccessible Links"
          value={Number(data.links_inaccessible ?? 0)}
        />
        <StatCard
          label="Login Form"
          value={
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                data.login_form_present
                  ? "bg-emerald-900/60 text-emerald-200"
                  : "bg-rose-900/60 text-rose-200"
              }`}
            >
              {data.login_form_present ? "Yes" : "No"}
            </span>
          }
        />
      </div>

      {/* Headings */}
      <div className="mt-4">
        <HeadingsTable headings={data.headings} />
      </div>

      {/* Raw JSON (collapsible) */}
      <div className="mt-4 rounded-xl border border-slate-700/60 bg-slate-900/50 p-4">
        <div className="mb-2 flex items-center justify-between">
          <div className="text-sm font-semibold text-slate-200">Raw JSON</div>
          <div className="flex items-center gap-2">
            <Button
              type="button"
              aria-expanded={showRaw}
              onClick={() => setShowRaw((s) => !s)}
            >
              {showRaw ? "Hide" : "View"} Raw
            </Button>
            <Button
              type="button"
              variant="ghost"
              aria-label="Copy JSON"
              onClick={copyJSON}
            >
              {copied ? "Copied!" : "Copy"}
            </Button>
          </div>
        </div>
        {showRaw && (
          <pre
            className="max-h-96 overflow-auto rounded-lg border border-slate-700/60 bg-slate-950 p-3 text-xs text-sky-200"
            role="region"
            aria-label="Analysis JSON"
          >
            {JSON.stringify(data, null, 2)}
          </pre>
        )}
      </div>
    </Card>
  );
}


function App() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const valid = isLikelyUrl(url);

  async function analyze(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setResult(null);

    if (!isLikelyUrl(url)) {
      setError("Please enter a valid http(s) URL.");
      return;
    }

    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });

      // Read the body ONCE
      const text = await res.text();
      let data: any = null;
      try {
        data = text ? JSON.parse(text) : null;
      } catch {
        /* ignore JSON parse errors */
      }

      if (!res.ok) {
        const msg =
          (data && (data.error as string)) ||
          (Array.isArray(data?.errors) && data.errors.join("; ")) ||
          `${res.status} ${res.statusText || ""}`.trim();
        throw new Error(msg);
      }

      setResult(data);
    } catch (err: any) {
      setError(err?.message ?? "Request failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col items-center px-4 py-12">
      <h1 className="mb-8 text-4xl font-bold text-gray-900 dark:text-white">
        Page Analyzer
      </h1>

      <form onSubmit={analyze} className="mb-2 flex w-full max-w-3xl items-center gap-3">

        <Input
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
          aria-invalid={!valid && url.length > 0}
        />
        <Button type="submit" disabled={loading || !valid}>
          {loading ? <Loader size="sm" /> : "Analyze"}
        </Button>
      </form>

      {/* Inline hint for UX; server still validates authoritatively */}
      {!valid && url.length > 0 && (
        <p className="mb-4 w-full max-w-3xl text-sm text-amber-600 dark:text-amber-400">
          Enter a valid URL (http or https). You can omit the scheme; we’ll assume https://.
        </p>
      )}

      {error && (
        <div className="w-full max-w-3xl mx-auto mb-2">
          <ErrorMessage message={error} />
        </div>
      )}

      {result && (
        <div className="w-full max-w-3xl">
          <ResultView data={result as AnalyzedPage} />
        </div>
      )}
    </div>
  );
}

export default App;
