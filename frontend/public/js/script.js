document.addEventListener('DOMContentLoaded', function() {
    // Fetch Brand and Types when page loads
    fetchBrands();
    fetchTypes();
    fetchCars();

    document.getElementsByID('carForm').addEventListener('submit', function(e) {
        e.preventDefault();

        const formData = {
            nama_mobil: document.getElementById('nama_mobil').value,
            idMerek_fk: document.getElementById('merek').value,
            idJenis_fk: document.getElementById('jenis').value,
            horse_power: document.getElementById('horse_power').value
        };

        fetch('/api/cars', {
            method:'POST',
            headers: {
                'content-type': 'application/JSON',
            },
            body: JSON.stringify(formData)
        })
        .then(response => response.json())
        .then(data => {
            showMessage('Car Added Successfully!', success);
            document.getElementById('carForm').reset();
            fetchCars(); // Refresh the car list
        })
        .catch(error => {
            showMessage('Error adding car: ' + error, 'error');
        });
    });
});

function fetchBrands() {
    fetch('/api/merek')
        .then(response => response.json())
        .then(data => {
            const select = document.getElementById('merek');
            data.forEach(brand => {
                const option = document.createElement('option');
                option.value = brand.idMerek;
                option.textContent = brand.merek;
                select.appendChild(option);
            });
        })
        .catch(error => console.error('Error fetching brands: ', error));
}

function fetchTypes() {
    fetch('/api/jenis')
        .then(response => response.json())
        .then(data => {
            const select = document.getElementById('jenis');
            data.forEach(type => {
                const option = document.createElement('option');
                option.value = type.idJenis;
                option.value = type.jenis;
                select.appendChild(option);
            });
        })
        .catch(error => console.error('Error fetching types: ', error));
}

function fetchCars() {
    fetch('/api/cars')
        .then(response => response.json())
        .then(data => {
            const carList = document.getElementById('carList');
            carList.innerHTML = '';

            if(data.length === 0) {
                carList.innerHTML = '<p>No cars found.</p>';
                return;
            }

            data.forEach(car => {
                const carItem = document.createElement('div');
                carItem.className = 'car-item';
                carItem.innerHTML = `
                    <h3>${car.nama_mobil}</h3>
                    <p><b>Brand: </b> ${car.merek.merek}</p>
                    <p><b>Type: </b> ${car.jenis.jenis}</p>
                    <p><b>Horse Power: </b> ${car.horse_power}</p>
                `;
                carList.appendChild(carItem);
            })
        })
        .catch(error => console.error('Error fetching cars: ', error));
}

function showMessage(message, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.textContent = message;
    messageDiv.className = type;
}