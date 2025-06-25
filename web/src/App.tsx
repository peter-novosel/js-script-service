import { useEffect, useState } from "react";
import { listScripts } from "./api";
import { ScriptEditor } from "./components/ScriptEditor";

function App() {
  const [scripts, setScripts] = useState<any[]>([]);

  useEffect(() => {
    listScripts().then(setScripts);
  }, []);

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Script Admin</h1>
      <ScriptEditor />
      <hr />
      <h2>All Scripts</h2>
      <ul>
        {scripts.map((s) => (
          <li key={s.slug}>
            <strong>{s.name}</strong> ({s.slug}) - {s.enabled ? "✅" : "❌"}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
