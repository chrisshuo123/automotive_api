import { fetchCars, fetchBrands, fetchTypes, deleteCar, getCarById } from './api.js';
import { openEditModal } from './modal.js';

export function showMessage(message, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.textContent = message;
    messageDiv.className = type;

    // Auto-hide for 5 seconds
    setTimeout(() => {
        messageDiv.textContent = '';
        messageDiv.className = '';
    }, 5000);
}

export async function loadCarForEdit(id) {
    try {
        const response = await getCarById(id);
        console.log('Full API Response: ', response);   // Debug

        let car = response;
        if(response.o) car = response.o;
        else if (response.data) car = response.data;
        else if (Array.isArray(response) && response.length > 0) car = response[0];

        console.log('Extracted car: ', car);
        // console.log('Car Fields: ', Object.keys(car));

        document.getElementById('edit_id').value = car.idCars;
        document.getElementById('edit_nama_mobil').value = car.nama_mobil;
        document.getElementById('edit_merek').value = car.idMerek_fk;
        document.getElementById('edit_jenis').value = car.idJenis_fk;
        document.getElementById('edit_horse_power').value = car.horse_power;
        document.getElementById('edit_status').value = car.idStatus_fk;
        openEditModal();
    } catch(error) {
        console.error('Edit Error: ', error);
        alert('Failed to load car: ' + error.message);
    }
}

export function renderCarList(cars) {
    const carList = document.getElementById('carList');
    carList.innerHTML = '';

    if(!cars || cars.length === 0) {
        carList.innerHTML = '<p>No cars found.</p>';
        return;
    }

    cars.forEach(car => {
        const brandDisplay = car.merek?.merek || (car.idMerek_fk ? `[ID: ${car.idMerek_fk}]` : 'Not Specified');

        const typeDisplay = car.jenis?.jenis ||  (car.idJenis_fk ? `[ID: ${car.idJenis_fk}]` : 'Not Specified');

        const statusDisplay = car.status?.status || (car.idStatus_fk ? `[ID: ${car.idStatus_fk}]` : 'Not Specified');

        // Determine status color
        let statusColor = 'yellow'; // default
        if(statusDisplay.toLowerCase() === 'approved') {
            statusColor = 'green';
        } else if(statusDisplay.toLowerCase() === 'needs preview') {
            statusColor = 'yellow';
        }

        const carItem = document.createElement('div');
        carItem.className = 'car-item';
        carItem.innerHTML = `
            <div class="flex-crud">
                <div style="width: 200px;">
                    <h3>${car.nama_mobil}</h3>
                    <p><b>Brand: </b> ${brandDisplay}</p>
                    <p><b>Type: </b> ${typeDisplay}</p>
                    <p><b>Horse Power: </b> ${car.horse_power ?? 'N/A'}</p>
                    <p><b>Status: </b><br><span style="background-color:${statusColor}; color: ${statusColor === 'yellow' ? 'black' : 'white'}; padding: 5px; border-radius: 5px;">${statusDisplay}</span></p>
                    <!-- This is for Edit Button in Panel Update Menu -->
                    <div class="flex-crud-button">
                        <button class="edit-btn" data-id="${car.idCars}">Edit</button>
                        <button class="delete-btn" data-id="${car.idCars}">Delete</button>
                    </div>
                </div>
                <div>
                    <img src="../public/img/yaris.jpg" style="width: 75%; max-width: 800px; height: auto; margin-left: 150px;">
                </div>
            </div>
        `;
        carList.appendChild(carItem);
    });

    // Attach Edit Listeners
    document.querySelectorAll('.edit-btn').forEach(btn => {
        btn.addEventListener('click', () => loadCarForEdit(btn.dataset.id));
    });
}

export function populateSelect(selectId, items, labelKey, valueKey, defaultText = 'Select an option') {
    const select = document.getElementById(selectId);
    select.innerHTML = `<option value="">${defaultText}</option>`;
    items.forEach(item => {
        select.add(new Option(item[labelKey], item[valueKey]));
    });
}

export function refreshCarList() {
    console.log("refreshCarList called");
    fetchCars()
        .then(apiData => {
            console.log('API Cars Response: ', apiData);
            const cars = apiData.data || apiData;
            renderCarList(cars);
        })
        .catch(error => {
            console.log('Error in refreshCarList: ', error);
            document.getElementById('carList').innerHTML = 
                `<div class="error">Error Handling cars: ${error.message}</div>`;
        });
}