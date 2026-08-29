import { fetchCars, fetchBrands, fetchTypes, deleteCar, getCarById, fetchStatus } from './api.js';
import { openEditModal } from './modal.js';
import { STATUS } from './config.js';

let carsData = [];
let allCarsData = []; // NEW: Store the original full dataset

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

// Helper function to flatten car data
function flattenCarData(car) {
    const flattened = {
        ...car,
        brandDisplay: car.merek?.merek || (car.idMerek_fk ? `[ID: ${car.idMerek_fk}]` : 'Not Specified'),
        typeDisplay: car.jenis?.jenis ||  (car.idJenis_fk ? `[ID: ${car.idJenis_fk}]` : 'Not Specified'),
        statusDisplay: car.status?.status || (car.idStatus_fk ? `[ID: ${car.idStatus_fk}]` : 'Not Specified')
    };
    console.log('Flattened car: ', flattened); // Debug log
    return flattened;
}

export function renderCarList(cars) {
    console.log('renderCarList called with: ', cars); // Debug log
    console.log('Number of cars in renderCarList: ', cars?.length); // Debug log
    console.log('Is car an array? ', Array.isArray(cars));

    carsData = cars; // Store the raw data
    const carList = document.getElementById('carList');
    carList.innerHTML = "";

    if(!cars || cars.length === 0) {
        console.log('No cars found - showing empty state');  // Debug Log
        carList.innerHTML = '<p>No cars found.</p>';
        return;
    }

    console.log('Rendering', cars.length, 'cars');

    cars.forEach(car => {
        // Use the flattened properties for dropdowns that use FKs
        const brandDisplay = car.brandDisplay || car.merek?.merek || (car.idMerek_fk ? `[ID: ${car.idMerek_fk}]` : 'Not Specified');

        const typeDisplay = car.typeDisplay || car.jenis?.jenis ||  (car.idJenis_fk ? `[ID: ${car.idJenis_fk}]` : 'Not Specified');

        const statusDisplay = car.statusDisplay || car.status?.status || (car.idStatus_fk ? `[ID: ${car.idStatus_fk}]` : 'Not Specified');

        console.log('Rendering car: ', car.nama_mobil, 'Status: ', statusDisplay); // Debug log

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
                    <!-- <img src="../public/img/yaris.jpg" style="width: 75%; max-width: 800px; height: auto; margin-left: 150px;"> -->
                    <img src="../public/img/${car.imageCar}" style="width: 75%; max-width: 800px; height: auto; margin-left: 150px;">
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
            console.log('Extracted cars array: ', cars); // Debug log

            // Flatten the data before storing and rendering
            const flattenedCars = cars.map(car => flattenCarData(car));
            console.log('Flattened cars array: ', flattenedCars); // Debug log

            // New: Store the full dataset
            allCarsData = flattenedCars;

            // renderCarList(cars);
            renderCarList(flattenedCars); // change from calling const cars to call flattenedCars
        })
        .catch(error => {
            console.log('Error in refreshCarList: ', error);
            document.getElementById('carList').innerHTML = 
                `<div class="error">Error Handling cars: ${error.message}</div>`;
        });
}


// export function getCars() {
//     console.log('getCars called, returning: ', carsData); // Debug log
//     return carsData;
// }

// MODIFIED: getCars now returns the full dataset, not the filtered
export function getCars() {
    console.log('getCars called, returning all cars count: ', allCarsData?.length);
    return allCarsData; // Return the full dataset
}

// NEW: Function to get currently displayed data (filtered)
export function getDisplayedCars() {
    return carsData;
}