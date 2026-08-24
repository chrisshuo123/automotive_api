// modal.js
export function openEditModal() {
    document.getElementById('editModal').style.display = 'flex';
}

export function closeEditModal() {
    document.getElementById('editModal').style.display = 'none';
}

export function setupModalListeners() {
    // Close with X button
    document.querySelector('.close-modal').addEventListener('click', closeEditModal);
    
    // Close when clicking outside
    window.addEventListener('click', (e) => {
        const modal = document.getElementById('editModal');
        if (e.target === modal) {
            closeEditModal();
        }
    });
}