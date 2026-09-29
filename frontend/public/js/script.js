import { renderCarList, getCars, refreshCarList } from "./ui.js";
import { STATUS } from "./config.js";

document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById("searchInput");
    const statusFilter = document.getElementById("statusFilter");
    const clearBtn = document.getElementById("clearBtn");

    console.log('DOM loaded - initializing filters'); // Debug log

    function filterData() {
        console.log('=== filterData called ==='); // Debug log

        const searchTerm = searchInput.value.toLowerCase().trim();
        const statusValue = statusFilter.value;

        console.log('Search term: ', `${searchTerm}`);  // Debug log
        console.log('Search term length: ', searchTerm.length);  // Debug log
        console.log('Status filter value: ', `"${statusValue}"`);  // Debug log
        console.log('STATUS.APPROVED from config: ', `"${STATUS.APPROVED}"`);
        console.log('STATUS.NEED_PREVIEW from config: ', `"${STATUS.NEED_PREVIEW}"`);
        console.log('Search input value: ', `${searchInput}`);  // Debug log
        console.log('Is search empty? ', searchTerm === '');   // Debug log
        console.log('Is status all or empty?', statusValue === '' || statusValue === 'all');  // Debug log

        // Check if the dropdown value matches our STATUS constants
        console.log('Does dropdown value match STATUS.APPROVED? ', statusValue === STATUS.APPROVED);
        console.log('Does dropdown value match STATUS.NEED_PREVIEW? ', statusValue === STATUS.NEED_PREVIEW);

        // const cars = getCars();  // Get the stored cars data (already flattened)
        // ganti ke allCars
        const allCars = getCars();
        console.log('Number of cars Available: ', allCars?.length); // Debug log

        // if(cars && cars.length > 0) {
        //     console.log('First car sample:', {
        //         name: cars[0].nama_mobil,
        //         status: cars[0].statusDisplay
        //     });
        // }

        if(!allCars || allCars.length === 0) {
            console.log('No cars data available'); // Debug log
            renderCarList([]);
            return;
        }

        // Log all car names to see what we're searching
        console.log('All car names: ', allCars.map(c => c.nama_mobil));

        // If search is empty AND status is "all" or empty, show all cars
        if (searchTerm === '' && (statusValue === '' || statusValue === 'all')) {
            console.log('Showing all cars (no filters) - count: ', allCars.length);
            renderCarList(allCars); // Pass the FULL dataset
            return;
        }

        // Otherwise, apply filters to the FULL dataset
        console.log('Applying filters to all', allCars.length, 'cars...');
        console.log('Searching for: ', `"${searchTerm}"`);

        const filteredResults = allCars.filter((car) => {
            // console.log('Checking car: ', car.nama_mobil, 'Status: ', car.statusDisplay); // Debug log
            const carNameLower = car.nama_mobil.toLowerCase();
            const carStatusLower = car.statusDisplay ? car.statusDisplay.toLowerCase() : '';

            // Search Filter - now using flattened statusDisplay
            let matchesSearch = true;
            if(searchTerm !== '') {
                matchesSearch = carNameLower.includes(searchTerm) ||
                                carStatusLower.includes(searchTerm);
                
                // Debug: Log cars that might match, because previously I use ' searchTerm === '' ' that cause input search that don't exist in the Car List kept showing LOL XD
                // if(searchTerm === 'zv') {  
                //     console.log(`Checking "${car.nama_mobil}":`, {
                //         name: carNameLower,
                //         includesZV: carNameLower.includes('zv'),
                //         status: carStatusLower,
                //         includesStatus: carStatusLower.includes('zv')
                //     });
                // }
            }
            
            // Status Filter - use flattened statusDisplay
            // empty or "all" means show all
            let matchesStatus = true;
            if(statusValue !== '' && statusValue !== 'all') {
                matchesStatus = carStatusLower === statusValue.toLowerCase();
            
                // Debug for first few cars
                if(allCars.indexOf(car) < 3) {
                    console.log(`Car "${car.nama_mobil}", Status: "${carStatusLower}"`);
                    console.log(`  Filter expecting: "${statusValue.toLowerCase()}"`);
                    console.log(`  Match? ${matchesStatus}`);
                }
            }

            return matchesSearch && matchesStatus;
            // const result = matchesSearch && matchesStatus;
            // const matchesStatus = statusValue === '' ||
            //     statusValue === 'all' ||
            //     carStatusLower === statusValue.toLowerCase();

            // console.log(`Car: ${car.nama_mobil}, matchesSearch: ${matchesSearch}, matchesStatus: ${matchesStatus}`);  // Debug log

            // Change Debug Log above to this below
            // Log only first 5 cars to avoid spam
            // if (cars.indexOf(car) < 5) {
            //     console.log(`Car ${cars.indexOf(car) + 1}:`, {
            //         nama: car.nama_mobil,
            //         status: car.statusDisplay,
            //         matchesSearch,
            //         matchesStatus,
            //         willInclude: matchesSearch && matchesStatus
            //     });
            // }

            // Log details for cars that might have "z" and "zv"
            if(searchTerm === 'z' || searchTerm === 'zv') {
                if(carNameLower.includes('z') || carNameLower.includes('zv')) {
                    console.log(`Car "${car.nama_mobil}" matches: `, {
                        matchesSearch,
                        matchesStatus,
                        result
                    });
                }
            }

            return result;
        });

        console.log('Filtered results count: ', filteredResults.length);  // Debug log

        if (filteredResults.length > 0) {
            console.log('First filtered car: ', {
                name: filteredResults[0].nama_mobil,
                status: filteredResults[0].statusDisplay
            });
        } else if(filteredResults === 0) {
            console.log('No cars matched the filters');
            // Log what we searched for
            console.log('  Search term was: ', `"${searchTerm}"`);
            // console.log('Avalable car names: ', cars.map(c => c.nama_mobil));
            console.log('  Status filter was: ', `"${statusValue}"`);
            console.log('  STATUS.APPROVED is: ', `"${STATUS.APPROVED}"`);
            console.log('  STATUS.NEED_PREVIEW is: ', `"${STATUS.NEED_PREVIEW}"`);
        } else {
            console.log(`Found ${filteredResults.length} cars matching the filters`);
            filteredResults.forEach((car, index) => {
                if (index < 5) {   // show first 5 matches
                    console.log(`  ${index + 1}. ${car.nama_mobil} - status: ${car.statusDisplay}`);
                }
            });
        }
        
        // console.log('Filtered data: ', filteredResults);  // Debug log, change with the debug log above.
        
        // Render the filtered results
        renderCarList(filteredResults); // Recently filterData
        console.log('=== FILTER DATA COMPLETE ===\n');
    }

    // Add event listeners
    if (searchInput) {
        console.log('Search input found - adding listener'); // Debug log
        // 'input' event triggers on every keystroke (including backspace/delete)
        
        searchInput.addEventListener('input', function(e) {
            console.log('Input event triggered.  Current value: ', `"${this.value}"`);
            filterData();
        }); // Real-time search on input

        // 'keypress' for Enter key
        searchInput.addEventListener('keypress', function(e) {
            if(e.key === 'enter') {
                e.preventDefault();
                console.log('Enter key pressed. Current value: ', `"${this.value}"`);
                filterData();
            }
        });

        searchInput.addEventListener('search', function(e) {
            console.log('Search event triggered (x button clicked).  Current value: ', `"${this.value}"`);
            if(this.value === '') {
                filterData();
            }
        });
    } else {
        console.log('Search input NOT found'); // Debug log
    }
    
    if (statusFilter) {
        console.log('Status filter found - adding listener'); // Debug log
        statusFilter.addEventListener('change', function(e) {
            console.log('Status changed to: ', this.value);
            filterData();
        });
    } else {
        console.log('Status filter NOT found'); // Debug log
    }

    // Clear button functionality
    if (clearBtn) {
        console.log('Clear button found - adding listener');
        clearBtn.addEventListener('click', function() {
            searchInput.value = '';
            statusFilter.value = ''; // or 'all' depending on your HTML
            console.log('Filters cleared. Search value: ', `"${searchInput.value}"`);
            filterData(); // Trigger filter to show all cars
        });
    }

    // Initial Load - use refreshCarList to fetch and flatten data
    console.log('Calling refreshCarList for initial load'); // Debug log
    refreshCarList();
});