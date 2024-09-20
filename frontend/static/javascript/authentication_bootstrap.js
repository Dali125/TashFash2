function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem('token');
    if (!token) {
        alert('No token found. Please log in again.');
        window.location.href = '/login';
        return;
    }

    const headers = {
        'Authorization': `Bearer ${token}`,
        ...options.headers
    };

    return fetch(url, { ...options, headers });
}