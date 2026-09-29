// config.js
export const API_BASE = 'http://localhost:8000';
export const DEFAULT_HEADERS = {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
};

// use as Consistency to Standardize app's status values rather than plenty of if commands.
// Easy change whenever need to change "approved", "Approved", or "APPROVED" anytime
// Prevent typos
export const STATUS = {
    APPROVED: "approved",
    NEED_PREVIEW: "need preview"
}
console.log('Approved: ', STATUS.APPROVED);
console.log('Need Preview: ', STATUS.NEED_PREVIEW);

export const MESSAGE_TIMEOUT = 5000;