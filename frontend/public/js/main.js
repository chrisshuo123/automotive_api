import { fetchBrands, fetchTypes, fetchStatus } from './api.js';
import { populateSelect, refreshCarList } from './ui.js';
import { setupModalListeners } from './modal.js';
import { handleAddCar, handleUpdateCar, handleDeleteCar } from './handlers.js';
import { cache } from './state.js';

document.addEventListener('DOMContentLoaded', async function() {
    console.log('DOM Fully Loaded!');

    // Setup Modal
    setupModalListeners();

    // --------- LOAD ALL DATA FIRST -----------
    try {
        const [brandsRes, typesRes, statusRes] = await Promise.all([
            fetchBrands(),
            fetchTypes(),
            fetchStatus()
        ]);

        cache.brands = brandsRes.data || brandsRes;
        cache.types = typesRes.data || typesRes;
        cache.status = statusRes.data || statusRes;

        console.log('✅ Cache ready: ', {
            brands: cache.brands.length,
            types: cache.types.length,
            status: cache.status.length
        });
    } catch(err) {
        console.error('❌ failed to load master data: ', err);
        return;     // ← jangan lanjut kalau gagal
    }

    // --------- POPULATE ALL SELECTS -----------
    populateSelect('merek', cache.brands, 'namamerek', 'idmerek');
    populateSelect('edit_merek', cache.brands, 'namamerek', 'idmerek');
    populateSelect('jenis', cache.types, 'namajenis', 'idjenis');
    populateSelect('edit_jenis', cache.types, 'namajenis', 'idjenis');
    populateSelect('status', cache.status, 'namastatus', 'idstatus');
    populateSelect('edit_status', cache.status, 'namastatus', 'idstatus');

    // --------- LOAD CAR LIST -----------
    await refreshCarList();
    // refreshCarList();

    // --------- EVENT LISTENERS ----------
    const carForm = document.getElementById('carForm');
    if (carForm) carForm.addEventListener('submit', handleAddCar);

    const editCarForm = document.getElementById('editCarForm');
    if (editCarForm) editCarForm.addEventListener('submit', handleUpdateCar);

    document.addEventListener('click', handleDeleteCar);

    console.log('✅ Init complete');
})