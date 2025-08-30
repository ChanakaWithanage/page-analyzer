import { useState } from "react";

function App() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  async function analyze(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const res = await fetch("http://localhost:8080/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      const data = await res.json();
      setResult(data);
    } catch (err: any) {
      setError(err.message);
    }
  }

  return (
    <div style={{ margin: "2rem", fontFamily: "sans-serif" }}>
      <h1>Page Analyzer</h1>
      <form onSubmit={analyze}>
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
          style={{ width: "400px", padding: "8px" }}
        />
        <button type="submit" style={{ marginLeft: "8px", padding: "8px 16px" }}>
          Analyze
        </button>
      </form>

      {error && <p style={{ color: "red" }}>Error: {error}</p>}

      {result && (
        <pre
          style={{
            marginTop: "1rem",
            background: "#1e1e1e",  // dark background
            color: "#0f0",          // bright green text
            padding: "1rem",
            borderRadius: "8px",
            overflowX: "auto"
          }}
        >
          {JSON.stringify(result, null, 2)}
        </pre>
      )}
    </div>
  );
}

export default App;
