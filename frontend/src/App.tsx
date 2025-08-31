import { useState } from "react";
import Input from "./components/Input";
import Button from "./components/Button";
import Card from "./components/Card";
import ErrorMessage from "./components/ErrorMessage";
import Loader from "./components/Loader";

function App() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function analyze(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setResult(null);
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });

      if (!res.ok) {
        const errJson = await res.json().catch(() => null);
        throw new Error(errJson?.error || `HTTP ${res.status}`);
      }

      const data = await res.json();
      setResult(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col items-center px-4 py-12">
      <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-8">
        Page Analyzer
      </h1>

      <form
        onSubmit={analyze}
        className="w-full max-w-xl flex items-center gap-3 mb-6"
      >
        <Input
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
        />
        <Button type="submit" disabled={loading || !url}>
          {loading ? <Loader size="sm" /> : "Analyze"}
        </Button>
      </form>

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
