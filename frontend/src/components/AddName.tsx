import { useState } from "react";
import { addName } from "../services/nameService";

function AddName() {
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setMessage("");
    setError("");

    if (!name.trim()) {
      setError("Please enter a name");
      return;
    }

    try {
      setLoading(true);

      const result = await addName(name);

      setMessage(result.message);
      setName("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to add name");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="rounded-xl bg-white p-6 shadow-md">
      <h2 className="mb-4 text-xl font-bold text-gray-800">Add New Name</h2>

      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter a name"
          maxLength={100}
          className="mb-4 w-full rounded-lg border border-gray-300 px-4 py-3 outline-none focus:border-blue-500"
        />

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-lg bg-blue-600 px-4 py-3 font-semibold text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? "Adding..." : "Add Name"}
        </button>
      </form>

      {message && (
        <p className="mt-4 rounded-lg bg-green-100 p-3 text-green-700">
          {message}
        </p>
      )}

      {error && (
        <p className="mt-4 rounded-lg bg-red-100 p-3 text-red-700">{error}</p>
      )}
    </div>
  );
}

export default AddName;
