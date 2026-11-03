import AddName from "./components/AddName";
import SearchName from "./components/SearchName";

function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-blue-700 px-6 py-6 text-center text-white shadow">
        <h1 className="text-3xl font-bold">Name Management System</h1>

        <p className="mt-2 text-blue-100">
          Add and search names using PostgreSQL
        </p>
      </header>

      <main className="mx-auto grid max-w-4xl gap-6 px-6 py-12 md:grid-cols-2">
        <AddName />

        <SearchName />
      </main>
    </div>
  );
}

export default App;
