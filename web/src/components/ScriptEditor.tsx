import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { saveScript } from "../api";

type Script = {
  name: string;
  slug: string;
  code: string;
};

type Props = {
  onSave?: () => void;
  initial?: Script | null;
};

export function ScriptEditor({ onSave, initial }: Props) {
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { isSubmitting }
  } = useForm<Script>({
    defaultValues: {
      name: "",
      slug: "",
      code: "",
    },
  });
  

  useEffect(() => {
    console.log("ScriptEditor: received initial", initial);
    if (initial) {
      reset(initial);
    } else {
      reset();
    }
  }, [initial, reset]);

  // ensure form reflects edit target when initial changes
  useEffect(() => {
    if (initial) {
      reset(initial);
    } else {
      reset();
    }
  }, [initial, reset]);

  const onSubmit = async (data: Script) => {
    try {
      await saveScript(data);
      alert("Script saved!");
      if (onSave) onSave();
    } catch (err: any) {
      alert("Error: " + err.response?.data || err.message);
    }
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="bg-gray-800 p-6 rounded-lg shadow-md space-y-4 max-w-2xl mb-10"
    >
      <div>
        <label className="block mb-1 text-sm text-gray-300">Name</label>
        <input
          {...register("name", { required: true })}
          placeholder="Script name"
          className="w-full px-3 py-2 rounded bg-gray-700 text-white placeholder-gray-400"
        />
      </div>

      <div>
        <label className="block mb-1 text-sm text-gray-300">Slug</label>
        <input
          {...register("slug", { required: true })}
          placeholder="e.g. hello-world"
          className="w-full px-3 py-2 rounded bg-gray-700 text-white placeholder-gray-400"
        />
      </div>

      <div>
        <label className="block mb-1 text-sm text-gray-300">JavaScript Code</label>
        <textarea
          {...register("code", { required: true })}
          placeholder="// Your JS here"
          rows={10}
          className="w-full px-3 py-2 rounded bg-gray-700 text-white placeholder-gray-400 font-mono"
        />
      </div>
      <div>
        <button onClick={() => handleEdit(s)}>✏️ Edit</button>

      </div>

      <div className="text-right">
        <button
          type="submit"
          disabled={isSubmitting}
          className="bg-green-600 hover:bg-green-700 text-white font-semibold px-5 py-2 rounded disabled:opacity-50"
        >
          💾 {initial ? "Update Script" : "Save Script"}
        </button>
      </div>
    </form>
  );
}
