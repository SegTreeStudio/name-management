import type {
  CreateNameResponse,
  SearchNameResponse,
  Name,
} from "../types/name";

const API_URL = import.meta.env.VITE_API_URL;

async function handleResponse<T>(response: Response): Promise<T> {
  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Something went wrong");
  }

  return data as T;
}

// Add a name.
export async function addName(name: string): Promise<CreateNameResponse> {
  const response = await fetch(`${API_URL}/names`, {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      name: name,
    }),
  });

  return handleResponse<CreateNameResponse>(response);
}

// Search for a name.
export async function searchName(name: string): Promise<SearchNameResponse> {
  const params = new URLSearchParams({
    name: name,
  });

  const response = await fetch(`${API_URL}/names/exists?${params.toString()}`);

  return handleResponse<SearchNameResponse>(response);
}

// Get all names.
export async function getAllNames(): Promise<Name[]> {
  const response = await fetch(`${API_URL}/names`);

  return handleResponse<Name[]>(response);
}
