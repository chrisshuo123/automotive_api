// Event Handlers
import { createCar, updateCar, deleteCar } from './api.js';
import { showMessage, refreshCarList, loadCarForEdit } from './ui.js';
import { closeEditModal } from './modal.js';

export function handleAddCar(e) {
    e.preventDefault();
    // const formData = {
    //     nama_mobil: document.getElementById('nama_mobil').value,
    //     idMerek_fk: parseInt(document.getElementById('merek').value),
    //     idJenis_fk: parseInt(document.getElementById('jenis').value),
    //     horse_power: parseInt(document.getElementById('horse_power').value) || 0,
    //     // idStatus_fk: parseInt(document.getElementById('status').value)
    //     imageCar: document.getElementById('image_car').value,
    //     idStatus_fk: 1
    // };
    const formData = new FormData();
    formData.append('nama_mobil', document.getElementById('nama_mobil').value);
    formData.append('idMerek_fk', document.getElementById('merek').value);
    formData.append('idJenis_fk', document.getElementById('jenis').value);
    formData.append('horse_power', document.getElementById('horse_power').value || 0);
    formData.append('idStatus_fk', 1);
    formData.append('image', document.getElementById('image_car').files[0]); // Actual File Object

    // DEBUG - lihat isi formData sebelum dikirim
    for (const [key, value] of formData.entries()) {
        console.log(key, ':', value);
    }

    createCar(formData)
        .then(apiData => {
            if(apiData.status) {
                showMessage('Car Added Successfully!', 'success');
                document.getElementById('carForm').reset();
                refreshCarList(); // Refresh the car list
            } else {
                showMessage('Failed to Add Car: ' + apiData.message, 'error');
            }
        })
        .catch(error => {
            showMessage('Error Adding Car: ' + error.message, 'error');
        });
}

export function handleUpdateCar(e) {
    e.preventDefault();
    const id = document.getElementById('edit_id').value;
    // const formData = {
    //     nama_mobil: document.getElementById('edit_nama_mobil').value,
    //     idMerek_fk: parseInt(document.getElementById('edit_merek').value),
    //     idJenis_fk: parseInt(document.getElementById('edit_jenis').value),
    //     horse_power: parseInt(document.getElementById('edit_horse_power').value),
    //     imageCar: document.getElementById('edit_image_car').value,
    //     idStatus_fk: parseInt(document.getElementById('edit_status').value)
    // };
    const formData = new FormData();
    formData.append('nama_mobil', document.getElementById('edit_nama_mobil').value);
    formData.append('idMerek_fk', document.getElementById('edit_merek').value);
    formData.append('idJenis_fk', document.getElementById('edit_jenis').value);
    formData.append('horse_power', document.getElementById('edit_horse_power').value);
    formData.append('idStatus_fk', document.getElementById('edit_status').value);
    
    const imageFile = document.getElementById('edit_image_car').files[0];
    if (imageFile) {
        formData.append('image', imageFile); // Only attach if user picked a new file
    }
    
    // formData.append('imageCar', document.getElementById('edit_image_car').files[0]); // Actual File Object

    updateCar(id, formData)
        .then(() => {
            closeEditModal();
            refreshCarList();
            showMessage('Car Updated Successfully!', 'success');
        })
        .catch(error => {
            showMessage('Update failed: ' + error.message, 'error');
        })
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