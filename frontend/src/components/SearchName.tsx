import { useState } from "react";
import { searchName } from "../services/nameService";

function SearchName() {
  const [name, setName] = useState("");
  const [result, setResult] = useState<boolean | null>(null);
  const [searchedName, setSearchedName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSearch(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setError("");
    setResult(null);

    if (!name.trim()) {
      setError("Please enter a name to search");
      return;
    }

    try {
      setLoading(true);

      const response = await searchName(name);

      setResult(response.exists);
      setSearchedName(response.name);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to search name");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="rounded-xl bg-white p-6 shadow-md">
      <h2 className="mb-4 text-xl font-bold text-gray-800">Search Name</h2>

      <form onSubmit={handleSearch}>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter name to search"
          className="mb-4 w-full rounded-lg border border-gray-300 px-4 py-3 outline-none focus:border-purple-500"
        />

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-lg bg-purple-600 px-4 py-3 font-semibold text-white hover:bg-purple-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? "Searching..." : "Search Name"}
        </button>
      </form>

      {result !== null && (
        <div
          className={`mt-4 rounded-lg p-4 ${
            result
              ? "bg-green-100 text-green-700"
              : "bg-yellow-100 text-yellow-800"
          }`}
        >
          {result ? (
            <p>
              <strong>{searchedName}</strong> already exists in the database.
            </p>
          ) : (
            <p>
              <strong>{searchedName}</strong> was not found in the database.
            </p>
          )}
        </div>
      )}

      {error && (
        <p className="mt-4 rounded-lg bg-red-100 p-3 text-red-700">{error}</p>
      )}
    </div>
  );
}

export default SearchName;
