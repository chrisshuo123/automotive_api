import { fetchBrands, fetchTypes, fetchStatus } from './api.js';
import { populateSelect, refreshCarList } from './ui.js';
import { setupModalListeners } from './modal.js';
import { handleAddCar, handleUpdateCar, handleDeleteCar } from './handlers.js';

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM Fully Loaded!');

    // Setup Modal
    setupModalListeners();

    // Load Initial Data
    refreshCarList();

    // Load Brands, Types, & Status for selects
    fetchBrands()
        .then(apiData => {
            const brands = apiData.data || apiData;
            populateSelect('merek', brands, 'merek', 'idMerek');
            populateSelect('edit_merek', brands, 'merek', 'idMerek');
        })
        .catch(error => console.error('Brands error: ', error));
    
    fetchTypes()
        .then(apiData => {
            const types = apiData.data || apiData;
            populateSelect('jenis', types, 'jenis', 'idJenis');
            populateSelect('edit_jenis', types, 'jenis', 'idJenis');
        })
        .catch(error => console.error('Types Error: ', error));

    fetchStatus()
        .then(apiData => {
            const status = apiData.data || apiData;
            // populateSelect('status', status, 'status', 'idStatus');
            populateSelect('edit_status', status, 'status', 'idStatus');
        })
        .catch(error => console.error('Status Error: ', error));

    // EVENT LISTENERS
    // Add a Car
    document.getElementById('carForm').addEventListener('submit', handleAddCar);
    // Update a Car (Edit)
    document.getElementById('editCarForm').addEventListener('submit', handleUpdateCar);
    // Delete a Car
    document.addEventListener('click', handleDeleteCar);
})