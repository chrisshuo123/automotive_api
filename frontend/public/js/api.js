// api.js
import { API_BASE, DEFAULT_HEADERS } from './config.js';

export async function fetchCars() {
    const response = await fetch(`${API_BASE}/api/cars`, { headers: DEFAULT_HEADERS });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return response.json();
}

export async function fetchBrands() {
    const response = await fetch(`${API_BASE}/api/brands`, { headers: DEFAULT_HEADERS });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return response.json();
}

export async function fetchTypes() {
    const response = await fetch(`${API_BASE}/api/types`, { headers: DEFAULT_HEADERS });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return response.json();
}

export async function fetchStatus() {
    const response = await fetch(`${API_BASE}/api/status`, {headers: DEFAULT_HEADERS});
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return response.json();
}

export async function createCar(data) {
    const response = await fetch(`${API_BASE}/api/cars`, {
        method: 'POST',
        headers: DEFAULT_HEADERS,
        body: JSON.stringify(data)
    });
    return response.json();
}

export async function updateCar(id, data) {
    const response = await fetch(`${API_BASE}/api/cars/${id}`, {
        method: 'PUT',
        headers: DEFAULT_HEADERS,
        body: JSON.stringify(data)
    });
    return response.json();
}

export async function deleteCar(id) {
    const response = await fetch(`${API_BASE}/api/cars/${id}`, {
        method: 'DELETE',
        headers: DEFAULT_HEADERS
    });
    return response.json();
}

export async function getCarById(id) {
    const response = await fetch(`${API_BASE}/api/cars/${id}`, { headers: DEFAULT_HEADERS });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return response.json();
}