// Event Handlers
import { createCar, updateCar, deleteCar } from './api.js';
import { showMessage, refreshCarList, loadCarForEdit } from './ui.js';
import { closeEditModal } from './modal.js';

export function handleAddCar(e) {
    e.preventDefault();

    const carData = {
        nama_mobil: document.getElementById('nama_mobil').value.trim(),
        idMerek_fk: parseInt(document.getElementById('merek').value),
        idJenis_fk: parseInt(document.getElementById('jenis').value),
        horse_power: parseInt(document.getElementById('horse_power').value) || 0
    };

    // Validation
    if(!carData.nama_mobil) {
        showMessage('Car name is required!', 'error');
        return;
    } else if (carData.horse_power <= 0) {
        showMessage('Please enter a valid horse power!', 'error');
        return;
    }

    // Save to sessionStorage
    const pendingChanges = JSON.parse(sessionStorage.getItem('pendingChanges') || '[]');
    pendingChanges.push({
        action:'add',
        data: carData,
        timestamp: new Date().toString(),
        status: 'pending'
    });
    sessionStorage.setItem('pendingChanges', JSON.stringify(pendingChanges));

    // ✅ DEBUG: Cek apakah data tersimpan
    console.log('Data Saved: ', JSON.parse(sessionStorage.getItem('pendingChanges')));
    alert('Data Saved! Redirecting to Preview...') // Test alert

    // Redirect to preview page
    window.location.replace('preview.html');
}

export function handleUpdateCar(e) {
    e.preventDefault();

    const id = document.getElementById('edit_id').value;
    const formData = {
        nama_mobil: document.getElementById('edit_nama_mobil').value.trim(),
        idMerek_fk: parseInt(document.getElementById('edit_merek').value),
        idJenis_fk: parseInt(document.getElementById('edit_jenis').value),
        horse_power: parseInt(document.getElementById('edit_horse_power').value)
    };

    // Validation...
    if(!carData.nama_mobil) {
        showMessage('Car name is required!', 'error');
        return;
    } else if (!carData.horse_power) {
        showMessage('Horse Power is required!', 'error');
        return;
    }

    // Save to sessionStorage
    const pendingChanges = JSON.parse(sessionStorage.getItem('pendingChanges') || '[]');
    pendingChanges.push({
        action: 'edit',
        data: carData,
        editId: id,
        timestamp: new Date().toString(),
        status: 'pending'
    });
    sessionStorage.setItem('pendingChanges', JSON.stringify(pendingChanges));

    // Close edit modal and redirect
    closeEditModal();
    window.location.href = 'preview.html';
}

export function handleDeleteCar(e) {
    if(e.target.classList.contains('delete-btn')) {
        const carId = e.target.dataset.id;
        if(confirm('Are you sure you want to delete this car?')) {
            deleteCar(carId)
                .then(() => {
                    showMessage('Car Deleted Successfully!', 'success');
                    refreshCarList();
                })
                .catch(error => {
                    showMessage('Delete failed: ' + error.message, 'error');
                });
        }
    }
}