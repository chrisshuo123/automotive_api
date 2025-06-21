// Global configuration (only what's needed)
const API_BASE = 'http://localhost:8000';
const DEFAULT_HEADERS = {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
};

console.log("Before DOMContentLoaded"); // Debug 1

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM fully loaded!"); // Debug 2
    fetchBrands();
    fetchTypes();
    fetchCars();

    // === 1 - UPDATE MODAL CONTROL FUNCTIONS ===
    function OpenEditModal() {
        document.getElementById('editModal').style.display = 'flex';
    }
    function CloseEditModal() {
        document.getElementById('editModal').style.display = 'none';
    }

    // Close Modal when clicking X or outside
    document.querySelector('.close-modal').addEventListener('click', CloseEditModal);
    window.addEventListener('click', (e) => {
        if(e.target === document.getElementById('edit-modal')) {
            CloseEditModal();
        }
    })

    // Modify your loadCarForEdit function to open the modal
    function loadCarForEdit(id) {
        fetch(`${API_BASE}/api/cars/${id}`)
            .then(response => {
                if(!response.ok) {
                    throw new Error(`HTTP Error! Status: ${response.status}`);
                }
                return response.json();
            })
            .then(response => { // car => {...
                // Handle both response formats
                const car = response.data || response;

                if (!car || !car.idCars) {
                    throw new Error("Invalid car data recieved");
                }
                console.log("Car Data: ", car); // Debug log

                document.getElementById('edit_id').value = car.idCars;
                document.getElementById('edit_nama_mobil').value = car.nama_mobil;
                document.getElementById('edit_merek').value = car.idMerek_fk;
                document.getElementById('edit_jenis').value = car.idJenis_fk;
                document.getElementById('edit_horse_power').value = car.horse_power;
            
                OpenEditModal(); // Show the modal after loading data
            })
            .catch(error => {
                console.error("Edit error:", error);
                alert("Failed to load car: " + error.message);
            });
    }
    
    // Your form submit handler...
    /*const fetchOptions = {
        credentials: 'include',  // This 'credentials' only makes the CORS more stricter, hard to serve DB to frontend.
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        }
    };*/

    // === 3 - DELETE CONTROL FUNCTIONS ===

    function deleteCar(id) {
        if(!confirm('Are you sure you want to delete this car?')) {
            return;
        }

        fetch(`${API_BASE}/api/cars/${id}`, {
            method: 'DELETE',
            headers: DEFAULT_HEADERS
        })
        .then(response => {
            if(!response.ok) throw new error(`HTTP Error! Status: ${response.status}`);
            return response.json();
        })
        .then(apiData => {
            showMessage('Car deleted successfully!', 'success');
            fetchCars(); // Refresh the list
        })
        .catch(error => {
            console.error('Delete error: ', error);
            showMessage('Failed to delete car, error message: ' + error.message, ' error');
        });
    }

    // Update fetchCars() to handle data.data
    function fetchCars() {
        fetch(`${API_BASE}/api/cars?_expand=merek&_expand=jenis`, { // Include related data
            headers: DEFAULT_HEADERS // No Credentials
        })
        /*.then(response => response.json())*/
        .then(response => {
            if(!response.ok) {
                throw new Error(`Network Response was not ok. HTTP ${response.status}`);
            }
            return response.json().then(data => {   // Properly parse JSON
                console.log("FULL API RESPONSE: ", data);   // Debug raw data
                return data;
            });
        })
        .then(apiData => {
            console.log("Raw API response:", apiData);
            const cars = apiData.data || apiData;
            const carList = document.getElementById('carList');
            //console.log('Car Data: ', apiData);
            carList.innerHTML = '';

            // Check for data.data structure
            /*if(data.length === 0) {*/
            if(!apiData.data || apiData.data.length === 0) {
                carList.innerHTML = '<p>No cars found.</p>';
                return;
            }

            // Use data.data instead of data
            /* data.forEach(car => { */
            cars.forEach(car => {
                // Safely extract brand and type names
                const brandDisplay = car.merek
                    ? car.merek.merek
                    : (car.idMerek_fk ? `[ID: ${car.idMerek_fk}]` : 'Not Specified');
                    //: `[ID:${car.merek}]`; // 👈 Show ID if name missing (prev changed idMerek_fk to merek)

                const typeDisplay = car.jenis
                    ? car.jenis.jenis
                    : (car.idJenis_fk ? `[ID: ${car.idJenis_fk}]` : 'Not Specified');
                    //: `[ID:${car.jenis}]`; // 👈 Show ID if name missing (prev changed idJenis_fk to jenis)

                const carItem = document.createElement('div');
                carItem.className = 'car-item';
                carItem.innerHTML = `
                    <h3>${car.nama_mobil}</h3>
                    <p><b>Brand: </b> ${brandDisplay}</p>
                    <p><b>Type: </b> ${typeDisplay}</p>
                    <p><b>Horse Power: </b> ${car.horse_power ?? 'N/A'}</p>
                    <!-- This is for Edit Button in Panel Update Menu -->
                    <button class="edit-btn" data-id="${car.idCars}">Edit</button>
                    <button class="delete-btn" data-id="${car.idCars}">Delete</button>
                `;
                carList.appendChild(carItem);
            });

            // Add Event Listeners for all Edit Buttons
            document.querySelectorAll('.edit-btn').forEach(btn => {
                btn.addEventListener('click', () => {
                    const carId = btn.getAttribute('data-id');
                    loadCarForEdit(carId); // Your existing edit function
                });
            });
        })
        .catch(error => {
            console.error('Error fetching cars. Car list Error: ', error);
            document.getElementById('carList').innerHTML = `
                <div class="error">Error loading cars: ${error.message}</div>
            `;
        });
    }

    // Update fetchBrands() to handle data.data
    function fetchBrands() {
        console.log("Fetching brands from:", `${API_BASE}/api/brands`);
        console.log("1. Starting brands fetch...") // Debug 1
        
        fetch(`${API_BASE}/api/brands`, {
            headers: DEFAULT_HEADERS // No Credentials!
        })
        /*.then(response => response.json())*/
        .then(response => {
            console.log("2. Brands response received. Status: ", response); // Debug 2
            if(!response.ok) {
                throw new Error(`Network Error, status: ${response.status}`);
            }
            return response.json().then(apiData => {
                console.log("3. Parsed JSON Data: ", apiData);
                return apiData;
            });
        })
        .then(apiData => {
            console.log('4. Processing Brands JSON Data: ', apiData);
            
            // Check if apiData.data exists, fallback to data if not
            //const brands = apiData.data || apiData;
            const brands = apiData.data || apiData;
            const addSelect = document.getElementById('merek');
            const editSelect = document.getElementById('edit_merek');
            
            // Clear existing options
            [addSelect, editSelect].forEach(select => {
                select.innerHTML = '<option value="">Select a brand</option>';
            });

            if(!brands || !Array.isArray(brands)) {
                throw new Error("Invalid brands data format");
            }
            
            if(brands.length === 0) {
                console.warn("No brands recieved from API");
                return;
            }

            // Add new options
            /*data.forEach(brand => { */
            brands.forEach(brand => {
                console.log("5. Adding brand: ", brand);
                
                const option = new Option(brand.merek, brand.idMerek);
                const editOption = new Option(brand.merek, brand.idMerek);

                addSelect.add(option);
                editSelect.add(editOption);
                
                // select.add(new Option(brand.merek, brand.idMerek));
                
                // Bisa pakai yang text, value, defaultSelected, dan selected
                /*select.add(new Option(
                    brand.merek,    // text
                    brand.idMerek,  // value
                    false,          // defaultSelected
                    false           // selected
                )); */

                // Atau pakai yang const option
                /*const option = document.createElement('option');
                option.value = brand.idMerek;
                option.textContent = brand.merek;
                select.appendChild(option);*/
            });
        })
        .catch(error => {
            console.error('6. Failed fetch brands: ', error);
            const selects = document.querySelectorAll('#merek', "#edit_merek");
            //const select = document.getElementById('merek');
            selects.forEach(select => {
                select.innerHTML = `
                    <option value="">Error loading brands (check console)</option>
                `;
            });
        });
    }

    // Update fetchTypes() to handle data.data
    function fetchTypes() {
        console.log("Fetching brands from:", `${API_BASE}/api/types`);
        console.log("1. Starting types fetch...") // Debug 1

        fetch(`${API_BASE}/api/types`, {
            headers: DEFAULT_HEADERS // No Credentials!
        })
        /*.then(response => response.json())*/
        .then(response => {
            console.log("2. Types response received. Status: ", response); // Debug 2
            if(!response.ok) {
                throw new Error(`Network Response was not ok`);
            }
            return response.json().then(apiData => {
                console.log("3. Parsed JSON Data: ", apiData);
                return apiData;
            });
        })
        .then(apiData => {
            console.log('4. Processing Types JSON Data: ', apiData);
            //const select = document.getElementById('jenis');
            const addSelect = document.getElementById('jenis');
            const editSelect = document.getElementById('edit_jenis');

            // Clear existing options
            [addSelect, editSelect].forEach(select => {
                select.innerHTML = '<option value="">Select a types</option>';
            });
            
            // Check if apiData.data exists, fallback to data if not
            const types = apiData.data || apiData;

            if(!types || !Array.isArray(types)) {
                throw new Error("Invalid types data format");
            }

            if(types.length === 0) {
                console.warn("No types recieved from API");
                return;
            }
            
            types.forEach(type => {
                console.log("5. Adding type: ", type);
                
                const option = new Option(type.jenis, type.idJenis);
                const editOption = new Option(type.jenis, type.idJenis);
                
                addSelect.add(option);
                editSelect.add(editOption);

                //select.add(new Option(type.jenis, type.idJenis));
            });
        })
        .catch(error => {
            console.error('6. Failed fetch types: ', error);
            const selects = document.querySelectorAll('#jenis', '#edit_jenis');
            //const select = document.getElementById('jenis');
            selects.forEach(select => {
                select.innerHTML = `
                    <option value="">Error loading jenis (check console)</option>
                `;
            })
        });
    }

    function showMessage(message, type) {
        const messageDiv = document.getElementById('message');
        messageDiv.textContent = message;
        messageDiv.className = type;

        // Auto-hide for 5 seconds
        setTimeout(() => {
            messageDiv.textContent = '';
            messageDiv.className = '';
        }, 5000);
    }

    //  === 1 - CREATE CAR EXECUTION ===
    // Update form submission to handle response structure (ADD CARS)
    document.getElementById('carForm').addEventListener('submit', function(e) {
        e.preventDefault();

        const formData = {  // changed idMerek_fk to merek, and idJenis_fk to jenis
            nama_mobil: document.getElementById('nama_mobil').value,
            idMerek_fk: parseInt(document.getElementById('merek').value), // Match Go struct
            idJenis_fk: parseInt(document.getElementById('jenis').value), // Match Go struct
            horse_power: parseInt(document.getElementById('horse_power').value) || 0 // Force number
        };

        console.log("Form Submission Data:", formData); // 👈 Add this

        fetch(`${API_BASE}/api/cars`, {
            method:'POST',
            headers: DEFAULT_HEADERS,
            body: JSON.stringify(formData)
        })
        /*.then(response => response.json())*/
        .then(response => {
            console.log('API Response: ', response);
            if(!response.ok) {
                throw new Error(`HTTP Error! Status ${response.status}`);
            }
            // Verify content type is JSON
            const contentType = response.headers.get('content-type');
            if(!contentType || !contentType.includes('application/json')) {
                throw new TypeError("Response isn't JSON");
            }
            return response.json();
        })
        .then(apiData => {
            // 👇 Add debug checks HERE (before processing cars)
            // (note: changed idMerek_fk to merek, and idJenis_fk to jenis)
            console.log("FULL RESPONSE: ", apiData);
            if(apiData.data?.merek === null && apiData.data?.idMerek_fk) {
                console.warn("Brand ID Exists but null: ", apiData.data.idMerek_fk);
            }
            if(apiData.data?.jenis === null && apiData.data?.idJenis_fk) {
                console.warn("Jenis ID Exists but null: ", apiData.data.idJenis_fk);
            }

            console.log("RAW CAR DATA: ", JSON.stringify(apiData, null, 2));
            const cars = apiData.data || apiData;

            document.getElementById('message').innerText = cars;

            if(apiData.status) {
                showMessage('Car Added Successfully!', 'success');
                document.getElementById('carForm').reset();
                fetchCars(); // Refresh the car list
            } else {
                showMessage('Failed to add car: ' + apiData.message, 'error');
            }
        })
        .catch(error => {
            showMessage('Error adding car: ' + Error.message, 'error');
        });
    });

    //  === 2 - HANDLE FORM SUBMISSION (UPDATE CARS) ===
    document.getElementById('editCarForm').addEventListener('submit', function(e) {
        e.preventDefault();

        const formData = {
            nama_mobil : document.getElementById('edit_nama_mobil').value,
            idMerek_fk : parseInt(document.getElementById('edit_merek').value),
            idJenis_fk : parseInt(document.getElementById('edit_jenis').value),
            horse_power : parseInt(document.getElementById('edit_horse_power').value)
        };

        const id = document.getElementById('edit_id').value;
        
        fetch(`${API_BASE}/api/cars/${id}`, {
            method : "PUT",
            headers : DEFAULT_HEADERS,
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            return response.json();
        })
        .then(updatedCar => { // updatedCar is now the full car object
            console.log("Update car: ", updatedCar);
            console.log("Updated car with relationships: ", updatedCar.merek);
            document.getElementById('editModal').style.display = 'none';
            fetchCars(); // Refresh the list
            showMessage('Car updated successfully!', 'success');
        })
        .catch(error => {
            console.log("Update Failed: ", error);
            showMessage('Update failed: ' + error, 'error');
        });
    });

    // === 3 - DELETE CARS EXECUTION ===
    document.addEventListener('click', function(e) {
        if(e.target.classList.contains('delete-btn')) {
            const carId = e.target.getAttribute('data-id');
            deleteCar(carId);
        }
    });
});

console.log("After DOMContentLoaded"); // Debug 3

