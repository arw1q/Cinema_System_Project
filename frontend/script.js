const API_URL = 'http://localhost:2281';

let currentUser = null;
let currentToken = null;
let currentMovie = null;
let currentSession = null;
let selectedSeat = null;

function getAuthHeaders() {
    if (currentToken) {
        return {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${currentToken}`
        };
    }
    return {
        'Content-Type': 'application/json'
    };
}

function showNotification(message, type = 'success') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification ${type} show`;
    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}

function showPage(pageName) {
    const pages = document.querySelectorAll('.page');
    pages.forEach(page => page.classList.remove('active'));

    const targetPage = document.getElementById(`${pageName}Page`);
    if (targetPage) {
        targetPage.classList.add('active');
    }

    if (pageName === 'home') {
        loadMovies();
    } else if (pageName === 'bookings') {
        loadUserBookings();
    } else if (pageName === 'adminBookings') {
        loadAllBookings();
    }
}

function switchAuthTab(type) {
    const userForm = document.getElementById('userAuthForm');
    const adminForm = document.getElementById('adminAuthForm');
    const tabs = document.querySelectorAll('.tab-btn');

    tabs.forEach(tab => tab.classList.remove('active'));
    event.target.classList.add('active');

    if (type === 'user') {
        userForm.style.display = 'block';
        adminForm.style.display = 'none';
    } else {
        userForm.style.display = 'none';
        adminForm.style.display = 'block';
    }
}

async function handleRegister(event) {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;
    const email = username + '@cinema.com';

    try {
        const response = await fetch(`${API_URL}/api/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });

        const data = await response.json();

        if (data.success) {
            currentUser = data.data;
            currentToken = data.token;
            showNotification(`Welcome, ${username}!`);
            updateNavBar();
            showPage('home');
            document.getElementById('registerForm').reset();
        } else {
            showNotification(data.message || 'Registration failed', 'error');
        }
    } catch (error) {
        showNotification('Error connecting to server', 'error');
    }
}

async function handleLogin(event) {
    event.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${API_URL}/api/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (data.success) {
            currentUser = data.data;
            currentToken = data.token;
            showNotification(`Welcome back, ${username}!`);
            updateNavBar();
            showPage('home');
        } else {
            showNotification(data.message || 'Login failed', 'error');
        }
    } catch (error) {
        showNotification('Error connecting to server', 'error');
    }
}

async function handleAdminLogin(event) {
    event.preventDefault();
    const username = document.getElementById('adminUsername').value;
    const password = document.getElementById('adminPassword').value;

    try {
        const response = await fetch(`${API_URL}/api/auth/admin`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (data.success) {
            currentUser = { role: 'admin', username: 'admin', id: 'admin' };
            currentToken = data.token;
            showNotification('Admin access granted');
            updateNavBar();
            showPage('home');
        } else {
            showNotification('Invalid admin credentials', 'error');
        }
    } catch (error) {
        showNotification('Error connecting to server', 'error');
    }
}

function updateNavBar() {
    const navAuth = document.getElementById('navAuth');
    const adminControls = document.getElementById('adminControls');
    const adminBookingsLink = document.getElementById('adminBookingsLink');

    if (currentUser) {
        navAuth.innerHTML = `
            <span style="color: var(--silver); margin-right: 1rem;">${currentUser.username || 'Admin'}</span>
            <button onclick="handleLogout()" style="padding: 0.5rem 1rem; background: var(--crimson); color: white; border: none; border-radius: 8px; cursor: pointer;">Logout</button>
        `;

        if (currentUser.role === 'admin') {
            if (adminControls) adminControls.style.display = 'block';
            if (adminBookingsLink) adminBookingsLink.style.display = 'block';
        }
    }
}

function handleLogout() {
    currentUser = null;
    currentToken = null;
    showNotification('Logged out successfully');
    showPage('login');
    document.getElementById('navAuth').innerHTML = '';
    const adminControls = document.getElementById('adminControls');
    const adminBookingsLink = document.getElementById('adminBookingsLink');
    if (adminControls) adminControls.style.display = 'none';
    if (adminBookingsLink) adminBookingsLink.style.display = 'none';
}

async function loadMovies() {
    try {
        const response = await fetch(`${API_URL}/api/movies`);
        const data = await response.json();

        const moviesList = document.getElementById('moviesList');

        if (data.success && data.data && data.data.length > 0) {
            moviesList.innerHTML = data.data.map(movie => `
                <div class="movie-card">
                    <div onclick="showMovieDetail('${movie.id}')">
                        <img src="${movie.poster_url || 'https://via.placeholder.com/300x400?text=No+Poster'}" 
                             alt="${movie.title}" class="movie-poster">
                        <div class="movie-info">
                            <h3 class="movie-title">${movie.title}</h3>
                            <p class="movie-duration">${movie.duration} minutes</p>
                        </div>
                    </div>
                    ${currentUser && currentUser.role === 'admin' ? `
                        <div style="padding: 0 1.5rem 1.5rem;">
                            <button class="btn-delete" onclick="event.stopPropagation(); deleteMovie('${movie.id}')">Delete Movie</button>
                        </div>
                    ` : ''}
                </div>
            `).join('');
        } else {
            moviesList.innerHTML = `
                <div class="empty-state">
                    <h3>No movies available</h3>
                    <p>Check back soon for new releases!</p>
                </div>
            `;
        }
    } catch (error) {
        showNotification('Error loading movies', 'error');
    }
}

async function showMovieDetail(movieId) {
    try {
        const response = await fetch(`${API_URL}/api/movies/${movieId}`);
        const data = await response.json();

        if (data.success) {
            currentMovie = data.data;
            await loadSessions(movieId);
            showPage('movieDetail');
        }
    } catch (error) {
        showNotification('Error loading movie details', 'error');
    }
}

async function loadSessions(movieId) {
    try {
        const response = await fetch(`${API_URL}/api/sessions/movie/${movieId}`);
        const data = await response.json();

        const detailDiv = document.getElementById('movieDetail');
        const sessions = data.success && data.data ? data.data : [];

        detailDiv.innerHTML = `
            <div class="movie-detail-container">
                <img src="${currentMovie.poster_url || 'https://via.placeholder.com/400x600?text=No+Poster'}" 
                     alt="${currentMovie.title}" class="detail-poster">
                <div class="detail-info">
                    <h1>${currentMovie.title}</h1>
                    <div class="detail-meta">
                        <p>Duration: ${currentMovie.duration} minutes</p>
                    </div>
                    <div class="detail-description">
                        <p>${currentMovie.description || 'No description available.'}</p>
                    </div>
                    
                    <div class="sessions-section">
                        <h3>Select Showtime</h3>
                        ${currentUser && currentUser.role === 'admin' ? `
                            <button class="btn-admin" onclick="showAddSession()" style="margin-bottom: 1.5rem;">+ Add Session</button>
                        ` : ''}
                        ${sessions.length > 0 ? `
                            <div class="sessions-grid">
                                ${sessions.map(session => `
                                    <div class="session-card">
                                        <div onclick="selectSession('${session.id}')">
                                            <div class="session-time">${new Date(session.start_time).toLocaleString('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        })}</div>
                                            <div class="session-hall">Hall: ${session.hall_name}</div>
                                            <div class="session-price">$${session.price.toFixed(2)}</div>
                                        </div>
                                        ${currentUser && currentUser.role === 'admin' ? `
                                            <button class="btn-delete" style="margin-top: 1rem;" onclick="event.stopPropagation(); deleteSession('${session.id}')">Delete</button>
                                        ` : ''}
                                    </div>
                                `).join('')}
                            </div>
                        ` : '<p style="color: var(--silver);">No sessions available yet.</p>'}
                    </div>
                </div>
            </div>
        `;
    } catch (error) {
        showNotification('Error loading sessions', 'error');
    }
}

async function selectSession(sessionId) {
    if (!currentUser || currentUser.role === 'admin') {
        showNotification('Please log in as a user to book tickets', 'error');
        return;
    }

    currentSession = sessionId;
    await loadBookedSeats(sessionId);
}

async function loadBookedSeats(sessionId) {
    try {
        const response = await fetch(`${API_URL}/api/bookings/session/${sessionId}`);
        const data = await response.json();

        const bookedSeats = data.success && data.data ? data.data.map(b => b.seat_number) : [];
        showSeatSelection(bookedSeats);
    } catch (error) {
        showNotification('Error loading seats', 'error');
    }
}

function showSeatSelection(bookedSeats) {
    const modalBody = document.getElementById('modalBody');
    const rows = ['A', 'B', 'C', 'D', 'E'];
    const seatsPerRow = 10;

    let seatsHTML = '';
    for (let row of rows) {
        for (let i = 1; i <= seatsPerRow; i++) {
            const seatNumber = `${row}${i}`;
            const isBooked = bookedSeats.includes(seatNumber);
            const seatClass = isBooked ? 'seat booked' : 'seat';
            const clickHandler = isBooked ? '' : ` onclick="selectSeatClick('${seatNumber}')"`;

            seatsHTML += `<div class="${seatClass}" data-seat="${seatNumber}"${clickHandler}>${seatNumber}</div>`;
        }
    }

    modalBody.innerHTML = `
        <h2 style="font-family: 'Playfair Display', serif; margin-bottom: 1.5rem;">Select Your Seat</h2>
        <div class="seats-container">
            <div class="seats-grid">
                ${seatsHTML}
            </div>
            <div style="margin-top: 2rem; display: flex; gap: 2rem; justify-content: center; font-size: 0.9rem;">
                <div><span style="display: inline-block; width: 20px; height: 20px; background: rgba(255,255,255,0.1); border-radius: 4px; margin-right: 0.5rem;"></span> Available</div>
                <div><span style="display: inline-block; width: 20px; height: 20px; background: var(--gold); border-radius: 4px; margin-right: 0.5rem;"></span> Selected</div>
                <div><span style="display: inline-block; width: 20px; height: 20px; background: var(--crimson); border-radius: 4px; margin-right: 0.5rem;"></span> Booked</div>
            </div>
            <button class="btn-primary" onclick="confirmBooking()" style="margin-top: 2rem;">Confirm Booking</button>
        </div>
    `;

    document.getElementById('modal').classList.add('active');
}

function selectSeatClick(seatNumber) {
    const allSeats = document.querySelectorAll('.seat:not(.booked)');
    allSeats.forEach(seat => seat.classList.remove('selected'));

    const clickedSeat = document.querySelector(`[data-seat="${seatNumber}"]`);
    if (clickedSeat) {
        clickedSeat.classList.add('selected');
        selectedSeat = seatNumber;
    }
}

async function confirmBooking() {
    if (!selectedSeat) {
        showNotification('Please select a seat', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/bookings/create`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                user_id: currentUser.id,
                session_id: currentSession,
                seat_number: selectedSeat
            })
        });

        const data = await response.json();

        if (data.success) {
            showNotification('Booking confirmed!');
            closeModal();
            selectedSeat = null;
        } else {
            showNotification(data.message || 'Booking failed', 'error');
        }
    } catch (error) {
        showNotification('Error creating booking', 'error');
    }
}

async function loadUserBookings() {
    if (!currentUser || currentUser.role === 'admin') {
        document.getElementById('bookingsList').innerHTML = `
            <div class="empty-state">
                <h3>Please log in to view your bookings</h3>
            </div>
        `;
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/bookings/user/${currentUser.id}`);
        const data = await response.json();

        const bookingsList = document.getElementById('bookingsList');

        if (data.success && data.data && data.data.length > 0) {
            bookingsList.innerHTML = data.data.map(booking => `
                <div class="booking-card">
                    <div class="booking-header">
                        <h3 class="booking-title">Booking Confirmation</h3>
                    </div>
                    <div class="booking-details">
                        <p><strong>Seat:</strong> ${booking.seat_number}</p>
                        <p><strong>Booked on:</strong> ${new Date(booking.created_at).toLocaleString()}</p>
                    </div>
                </div>
            `).join('');
        } else {
            bookingsList.innerHTML = `
                <div class="empty-state">
                    <h3>No bookings yet</h3>
                    <p>Start by selecting a movie and booking a seat!</p>
                </div>
            `;
        }
    } catch (error) {
        showNotification('Error loading bookings', 'error');
    }
}

function showAddMovie() {
    const modalBody = document.getElementById('modalBody');
    modalBody.innerHTML = `
        <h2 style="font-family: 'Playfair Display', serif; margin-bottom: 1.5rem;">Add New Movie</h2>
        <form onsubmit="handleAddMovie(event)">
            <div class="form-group">
                <label>Movie Title</label>
                <input type="text" id="movieTitle" required>
            </div>
            <div class="form-group">
                <label>Description</label>
                <textarea id="movieDescription" rows="4" required></textarea>
            </div>
            <div class="form-group">
                <label>Duration (minutes)</label>
                <input type="number" id="movieDuration" required>
            </div>
            <div class="form-group">
                <label>Poster URL</label>
                <input type="url" id="moviePoster" placeholder="https://example.com/poster.jpg" required>
            </div>
            <button type="submit" class="btn-primary">Add Movie</button>
        </form>
    `;
    document.getElementById('modal').classList.add('active');
}

async function handleAddMovie(event) {
    event.preventDefault();

    const movieData = {
        title: document.getElementById('movieTitle').value,
        description: document.getElementById('movieDescription').value,
        duration: parseInt(document.getElementById('movieDuration').value),
        poster_url: document.getElementById('moviePoster').value
    };

    try {
        const response = await fetch(`${API_URL}/api/movies/create`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(movieData)
        });

        const data = await response.json();

        if (data.success) {
            showNotification('Movie added successfully!');
            closeModal();
            loadMovies();
        } else {
            showNotification(data.message || 'Failed to add movie', 'error');
        }
    } catch (error) {
        showNotification('Error adding movie', 'error');
    }
}

function showAddSession() {
    const modalBody = document.getElementById('modalBody');
    modalBody.innerHTML = `
        <h2 style="font-family: 'Playfair Display', serif; margin-bottom: 1.5rem;">Add New Session</h2>
        <form onsubmit="handleAddSession(event)">
            <div class="form-group">
                <label>Hall Name</label>
                <input type="text" id="sessionHall" placeholder="Hall A" required>
            </div>
            <div class="form-group">
                <label>Start Date & Time</label>
                <input type="datetime-local" id="sessionTime" required>
            </div>
            <div class="form-group">
                <label>Price ($)</label>
                <input type="number" step="0.01" id="sessionPrice" placeholder="15.50" required>
            </div>
            <button type="submit" class="btn-primary">Add Session</button>
        </form>
    `;
    document.getElementById('modal').classList.add('active');
}

async function handleAddSession(event) {
    event.preventDefault();

    const sessionData = {
        movie_id: currentMovie.id,
        hall_name: document.getElementById('sessionHall').value,
        start_time: new Date(document.getElementById('sessionTime').value).toISOString(),
        price: parseFloat(document.getElementById('sessionPrice').value)
    };

    try {
        const response = await fetch(`${API_URL}/api/sessions/create`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(sessionData)
        });

        const data = await response.json();

        if (data.success) {
            showNotification('Session added successfully!');
            closeModal();
            loadSessions(currentMovie.id);
        } else {
            showNotification(data.message || 'Failed to add session', 'error');
        }
    } catch (error) {
        showNotification('Error adding session', 'error');
    }
}

function closeModal() {
    document.getElementById('modal').classList.remove('active');
}

window.onclick = function(event) {
    const modal = document.getElementById('modal');
    if (event.target === modal) {
        closeModal();
    }
}

async function deleteMovie(movieId) {
    if (!confirm('Are you sure you want to delete this movie?')) {
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/movies/delete/${movieId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        const data = await response.json();

        if (data.success) {
            showNotification('Movie deleted successfully!');
            loadMovies();
        } else {
            showNotification(data.message || 'Failed to delete movie', 'error');
        }
    } catch (error) {
        showNotification('Error deleting movie', 'error');
    }
}

async function deleteSession(sessionId) {
    if (!confirm('Are you sure you want to delete this session?')) {
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/sessions/delete/${sessionId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        const data = await response.json();

        if (data.success) {
            showNotification('Session deleted successfully!');
            loadSessions(currentMovie.id);
        } else {
            showNotification(data.message || 'Failed to delete session', 'error');
        }
    } catch (error) {
        showNotification('Error deleting session', 'error');
    }
}

async function deleteBooking(bookingId) {
    if (!confirm('Are you sure you want to delete this booking?')) {
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/bookings/delete/${bookingId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        const data = await response.json();

        if (data.success) {
            showNotification('Booking deleted successfully!');
            loadAllBookings();
        } else {
            showNotification(data.message || 'Failed to delete booking', 'error');
        }
    } catch (error) {
        showNotification('Error deleting booking', 'error');
    }
}

async function loadAllBookings() {
    if (!currentUser || currentUser.role !== 'admin') {
        document.getElementById('adminBookingsList').innerHTML = `
            <div class="empty-state">
                <h3>Admin access required</h3>
            </div>
        `;
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/bookings/all`, {
            headers: getAuthHeaders()
        });
        const data = await response.json();

        const bookingsList = document.getElementById('adminBookingsList');

        if (data.success && data.data && data.data.length > 0) {
            bookingsList.innerHTML = data.data.map(booking => `
                <div class="booking-card">
                    <div class="booking-header">
                        <h3 class="booking-title">Booking #${booking.id.substring(0, 8)}</h3>
                        <button class="btn-delete" onclick="deleteBooking('${booking.id}')">Delete</button>
                    </div>
                    <div class="booking-details">
                        <p><strong>User ID:</strong> ${booking.user_id}</p>
                        <p><strong>Session ID:</strong> ${booking.session_id}</p>
                        <p><strong>Seat:</strong> ${booking.seat_number}</p>
                        <p><strong>Booked on:</strong> ${new Date(booking.created_at).toLocaleString()}</p>
                    </div>
                </div>
            `).join('');
        } else {
            bookingsList.innerHTML = `
                <div class="empty-state">
                    <h3>No bookings yet</h3>
                    <p>No tickets have been booked yet.</p>
                </div>
            `;
        }
    } catch (error) {
        showNotification('Error loading bookings', 'error');
    }
}