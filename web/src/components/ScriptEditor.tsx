import { useForm } from "react-hook-form";
import { saveScript } from "../api";

export function ScriptEditor() {
  const { register, handleSubmit, reset } = useForm();

  const onSubmit = async (data: any) => {
    try {
      await saveScript(data);
      alert("Saved!");
      reset();
    } catch (err: any) {
      alert("Error: " + err.response?.data || err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input {...register("name")} placeholder="Name" required />
      <input {...register("slug")} placeholder="Slug" required />
      <textarea {...register("code")} placeholder="JavaScript code" required rows={10} />
      <button type="submit">Save Script</button>
    </form>
  );
}
