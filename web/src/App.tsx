import { FormEvent, useMemo, useState } from "react";
import { Check, Clipboard, ExternalLink, Link2, Loader2, ShieldCheck } from "lucide-react";

type ShortenResponse = {
  code: string;
  short_url: string;
  long_url: string;
};

type RequestState = "idle" | "loading" | "success" | "error";

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? "";

export function App() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<ShortenResponse | null>(null);
  const [state, setState] = useState<RequestState>("idle");
  const [message, setMessage] = useState("");
  const [copied, setCopied] = useState(false);

  const canSubmit = useMemo(() => url.trim().length > 0 && state !== "loading", [state, url]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setState("loading");
    setMessage("");
    setCopied(false);

    try {
      const response = await fetch(`${apiBaseUrl}/api/v1/shorten`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ url: url.trim() })
      });

      const data = await response.json().catch(() => null);

      if (!response.ok) {
        throw new Error(data?.message ?? "Failed to create a short URL");
      }

      setResult(data as ShortenResponse);
      setState("success");
      setMessage("Short link created");
    } catch (error) {
      setResult(null);
      setState("error");
      setMessage(error instanceof Error ? error.message : "Something went wrong");
    }
  }

  async function copyShortUrl() {
    if (!result) {
      return;
    }

    await navigator.clipboard.writeText(result.short_url);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1600);
  }

  return (
    <main className="app-shell">
      <section className="workspace" aria-labelledby="page-title">
        <div className="intro">
          <div className="brand-mark" aria-hidden="true">
            <Link2 size={28} />
          </div>
          <div>
            <p className="eyebrow">Encrypted REST API</p>
            <h1 id="page-title">URL Shortener</h1>
          </div>
        </div>

        <form className="shortener-form" onSubmit={handleSubmit}>
          <label htmlFor="url">Long URL</label>
          <div className="input-row">
            <input
              id="url"
              type="url"
              placeholder="https://example.com/articles/my-long-link"
              value={url}
              onChange={(event) => setUrl(event.target.value)}
              autoComplete="url"
              required
            />
            <button type="submit" disabled={!canSubmit}>
              {state === "loading" ? <Loader2 className="spin" size={18} /> : <Link2 size={18} />}
              Shorten
            </button>
          </div>
        </form>

        {message && (
          <p className={`status-message ${state === "error" ? "is-error" : "is-success"}`}>
            {message}
          </p>
        )}

        {result && (
          <section className="result-panel" aria-label="Created short URL">
            <div className="result-meta">
              <span>Code</span>
              <strong>{result.code}</strong>
            </div>

            <div className="result-link">
              <a href={result.short_url} target="_blank" rel="noreferrer">
                {result.short_url}
              </a>
              <div className="result-actions">
                <button type="button" className="icon-button" onClick={copyShortUrl} aria-label="Copy short URL">
                  {copied ? <Check size={18} /> : <Clipboard size={18} />}
                </button>
                <a className="icon-button" href={result.short_url} target="_blank" rel="noreferrer" aria-label="Open short URL">
                  <ExternalLink size={18} />
                </a>
              </div>
            </div>

            <p className="original-url">{result.long_url}</p>
          </section>
        )}

        <div className="security-strip">
          <ShieldCheck size={18} />
          <span>Original URLs are encrypted in PostgreSQL and API queries are parameterized.</span>
        </div>
      </section>
    </main>
  );
}
