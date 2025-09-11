import { useState } from "react";
import Input from "./components/Input";
import Button from "./components/Button";
import Card from "./components/Card";
import ErrorMessage from "./components/ErrorMessage";
import Loader from "./components/Loader";
import { isLikelyUrl } from "./utils/url";

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
      try { data = text ? JSON.parse(text) : null; } catch { /* ignore */ }

      if (!res.ok) {
        const msg =
          data?.error ||
          (Array.isArray(data?.errors) && data.errors.join("; ")) ||
          `${res.status} ${res.statusText || ""}`.trim();
        throw new Error(msg);
      }

      setResult(data);
    } catch (err: any) {
      setError(err.message ?? "Request failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col items-center px-4 py-12">
      <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-8">
        Page Analyzer
      </h1>

      <form onSubmit={analyze} className="w-full max-w-xl flex items-center gap-3 mb-2">
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
        <p className="w-full max-w-xl text-sm text-amber-600 dark:text-amber-400 mb-4">
          Enter a valid URL (http or https). You can omit the scheme; we’ll assume https://.
        </p>
      )}

      {error && <ErrorMessage message={error} />}

      {result && (
        <Card>
          <pre className="text-sm text-gray-800 dark:text-green-400 whitespace-pre-wrap">
            {JSON.stringify(result, null, 2)}
          </pre>
        </Card>
      )}
    </div>
  );
}

export default App;
