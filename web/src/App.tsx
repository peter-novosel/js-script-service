import { useEffect, useState } from "react";
import { listScripts, disableScript } from "./api";
import { ScriptEditor } from "./components/ScriptEditor";

function App() {
  const [scripts, setScripts] = useState<any[]>([]);
  const [editScript, setEditScript] = useState<any | null>(null);

  const refreshScripts = async () => {
    const data = await listScripts();
    setScripts(data);
  };

  useEffect(() => {
    refreshScripts();
  }, []);

  const handleEdit = (script: any) => setEditScript(script);

  const handleDelete = async (slug: string) => {
    if (!confirm(`Disable script "${slug}"?`)) return;
    await disableScript(slug);
    await refreshScripts();
  };

  return (
    <div className="min-h-screen bg-gray-100 text-gray-900 py-10 px-4">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-10">🛠 Script Admin</h1>

        <ScriptEditor
          onSave={async () => {
            await refreshScripts();
            setEditScript(null);
          }}
          initial={editScript}
        />

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-4">All Scripts</h2>

          <ul className="space-y-3">
            {scripts.map((s) => (
              <li
                key={s.slug}
                className="flex justify-between items-start bg-gray-100 border border-gray-200 rounded-md px-4 py-3"
              >
                <div>
                  <div className="font-semibold">{s.name}</div>
                  <div className="text-sm text-gray-500">({s.slug})</div>
                </div>
                <div className="space-x-2">
                  <button
                    onClick={() => handleEdit(s)}
                    className="text-blue-600 text-sm hover:underline"
                  >
                    ✏️ Edit
                  </button>
                  <button
                    onClick={() => handleDelete(s.slug)}
                    className="text-red-600 text-sm hover:underline"
                  >
                    🗑 Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}

export default App;
