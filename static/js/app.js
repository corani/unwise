// Unwise UI Application
(function() {
    'use strict';

    // State
    let allBooks = [];
    let filteredBooks = [];
    let selectedBookId = null;

    // DOM Elements
    const elements = {
        // Books panel
        searchInput: null,
        booksLoading: null,
        booksError: null,
        booksErrorMessage: null,
        booksEmpty: null,
        booksList: null,
        
        // Highlights panel
        highlightsTitle: null,
        highlightsWelcome: null,
        highlightsLoading: null,
        highlightsError: null,
        highlightsErrorMessage: null,
        highlightsEmpty: null,
        highlightsList: null
    };

    // Initialize the application
    function init() {
        // Get DOM elements
        elements.searchInput = document.getElementById('searchInput');
        elements.booksLoading = document.getElementById('booksLoading');
        elements.booksError = document.getElementById('booksError');
        elements.booksErrorMessage = document.getElementById('booksErrorMessage');
        elements.booksEmpty = document.getElementById('booksEmpty');
        elements.booksList = document.getElementById('booksList');
        
        elements.highlightsTitle = document.getElementById('highlightsTitle');
        elements.highlightsWelcome = document.getElementById('highlightsWelcome');
        elements.highlightsLoading = document.getElementById('highlightsLoading');
        elements.highlightsError = document.getElementById('highlightsError');
        elements.highlightsErrorMessage = document.getElementById('highlightsErrorMessage');
        elements.highlightsEmpty = document.getElementById('highlightsEmpty');
        elements.highlightsList = document.getElementById('highlightsList');

        // Set up event listeners
        elements.searchInput.addEventListener('input', handleSearch);

        // Set up theme management
        initTheme();

        // Load books
        loadBooks();
    }

    // Theme management functions
    function initTheme() {
        const themeButtons = document.querySelectorAll('.theme-toggle');
        const storedTheme = localStorage.getItem('theme') || 'auto';

        // Update active button
        updateActiveThemeButton(storedTheme);

        // Add click handlers
        themeButtons.forEach(button => {
            button.addEventListener('click', () => {
                const theme = button.getAttribute('data-theme');
                setTheme(theme);
                localStorage.setItem('theme', theme);
                updateActiveThemeButton(theme);
            });
        });
    }

    function setTheme(theme) {
        if (theme === 'auto') {
            const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
            document.documentElement.setAttribute('data-bs-theme', systemTheme);
        } else {
            document.documentElement.setAttribute('data-bs-theme', theme);
        }
    }

    function updateActiveThemeButton(theme) {
        const themeButtons = document.querySelectorAll('.theme-toggle');
        themeButtons.forEach(button => {
            if (button.getAttribute('data-theme') === theme) {
                button.classList.add('active');
            } else {
                button.classList.remove('active');
            }
        });
    }

    // Load books from API
    async function loadBooks() {
        showBooksLoading();

        try {
            const response = await fetch('/ui/api/books');
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            allBooks = data.results || [];
            filteredBooks = [...allBooks];

            if (allBooks.length === 0) {
                showBooksEmpty();
            } else {
                renderBooks();
            }
        } catch (error) {
            console.error('Error loading books:', error);
            showBooksError(error.message);
        }
    }

    // Load highlights for a specific book
    async function loadHighlights(bookId) {
        showHighlightsLoading();

        try {
            const response = await fetch(`/ui/api/books/${bookId}/highlights`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            const highlights = data.results || [];

            if (highlights.length === 0) {
                showHighlightsEmpty();
            } else {
                renderHighlights(highlights);
            }
        } catch (error) {
            console.error('Error loading highlights:', error);
            showHighlightsError(error.message);
        }
    }

    // Render books list
    function renderBooks() {
        elements.booksList.innerHTML = '';

        filteredBooks.forEach(book => {
            const bookItem = createBookElement(book);
            elements.booksList.appendChild(bookItem);
        });

        hideBooksLoading();
        hideBooksError();
        hideBooksEmpty();
        elements.booksList.style.display = 'block';
    }

    // Create a book list item element
    function createBookElement(book) {
        const div = document.createElement('div');
        div.className = 'list-group-item book-item';
        if (book.id === selectedBookId) {
            div.classList.add('active');
        }

        const title = document.createElement('div');
        title.className = 'book-title';
        title.textContent = book.title || 'Untitled';

        const author = document.createElement('div');
        author.className = 'book-author';
        author.textContent = book.author || 'Unknown Author';

        const meta = document.createElement('div');
        meta.className = 'book-meta d-flex justify-content-between align-items-center mt-2';
        
        const highlightBadge = document.createElement('span');
        highlightBadge.className = 'badge bg-primary';
        highlightBadge.textContent = `${book.num_highlights} highlight${book.num_highlights !== 1 ? 's' : ''}`;
        
        const updatedDate = document.createElement('span');
        updatedDate.textContent = formatDate(book.updated);
        
        meta.appendChild(highlightBadge);
        meta.appendChild(updatedDate);

        div.appendChild(title);
        div.appendChild(author);
        div.appendChild(meta);

        div.addEventListener('click', () => selectBook(book));

        return div;
    }

    // Select a book and load its highlights
    function selectBook(book) {
        selectedBookId = book.id;
        
        // Update UI to show selected book
        document.querySelectorAll('.book-item').forEach(item => {
            item.classList.remove('active');
        });
        event.currentTarget.classList.add('active');

        // Update highlights title
        elements.highlightsTitle.textContent = `Highlights: ${book.title}`;

        // Load highlights
        loadHighlights(book.id);
    }

    // Render highlights list
    function renderHighlights(highlights) {
        elements.highlightsList.innerHTML = '';

        // Sort by location
        highlights.sort((a, b) => a.location - b.location);

        highlights.forEach(highlight => {
            const highlightCard = createHighlightElement(highlight);
            elements.highlightsList.appendChild(highlightCard);
        });

        hideHighlightsLoading();
        hideHighlightsError();
        hideHighlightsEmpty();
        hideHighlightsWelcome();
        elements.highlightsList.style.display = 'block';
    }

    // Create a highlight card element
    function createHighlightElement(highlight) {
        const card = document.createElement('div');
        card.className = 'card highlight-card';

        const cardBody = document.createElement('div');
        cardBody.className = 'card-body';

        // Highlight text
        const text = document.createElement('div');
        text.className = 'highlight-text';
        text.textContent = highlight.text;

        cardBody.appendChild(text);

        // Note (if present)
        if (highlight.note) {
            const note = document.createElement('div');
            note.className = 'highlight-note';
            note.textContent = `Note: ${highlight.note}`;
            cardBody.appendChild(note);
        }

        // Metadata
        const meta = document.createElement('div');
        meta.className = 'highlight-meta d-flex justify-content-between align-items-center';

        const metaInfo = document.createElement('div');
        const metaParts = [];
        
        if (highlight.chapter) {
            metaParts.push(`Chapter: ${highlight.chapter}`);
        }
        if (highlight.location) {
            metaParts.push(`Location: ${highlight.location}`);
        }
        metaParts.push(`Updated: ${formatDate(highlight.updated)}`);
        
        metaInfo.textContent = metaParts.join(' • ');

        // Copy button
        const copyBtn = document.createElement('button');
        copyBtn.className = 'btn btn-sm btn-outline-secondary copy-button';
        copyBtn.innerHTML = '📋 Copy';
        copyBtn.addEventListener('click', () => copyHighlight(highlight));

        meta.appendChild(metaInfo);
        meta.appendChild(copyBtn);

        cardBody.appendChild(meta);
        card.appendChild(cardBody);

        return card;
    }

    // Copy highlight text to clipboard
    function copyHighlight(highlight) {
        const text = highlight.text + (highlight.note ? `\n\nNote: ${highlight.note}` : '');
        
        navigator.clipboard.writeText(text).then(() => {
            // Show temporary success feedback
            const btn = event.currentTarget;
            const originalText = btn.innerHTML;
            btn.innerHTML = '✓ Copied!';
            btn.classList.remove('btn-outline-secondary');
            btn.classList.add('btn-success');
            
            setTimeout(() => {
                btn.innerHTML = originalText;
                btn.classList.remove('btn-success');
                btn.classList.add('btn-outline-secondary');
            }, 2000);
        }).catch(err => {
            console.error('Failed to copy:', err);
        });
    }

    // Handle search input
    function handleSearch(event) {
        const query = event.target.value.toLowerCase().trim();

        if (query === '') {
            filteredBooks = [...allBooks];
        } else {
            filteredBooks = allBooks.filter(book => {
                const titleMatch = book.title.toLowerCase().includes(query);
                const authorMatch = book.author.toLowerCase().includes(query);
                return titleMatch || authorMatch;
            });
        }

        if (filteredBooks.length === 0) {
            showBooksEmpty();
        } else {
            renderBooks();
        }
    }

    // Format date string
    function formatDate(dateString) {
        if (!dateString) return '';
        
        try {
            const date = new Date(dateString);
            const now = new Date();
            const diffTime = Math.abs(now - date);
            const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

            if (diffDays === 0) {
                return 'Today';
            } else if (diffDays === 1) {
                return 'Yesterday';
            } else if (diffDays < 7) {
                return `${diffDays} days ago`;
            } else if (diffDays < 30) {
                const weeks = Math.floor(diffDays / 7);
                return `${weeks} week${weeks !== 1 ? 's' : ''} ago`;
            } else if (diffDays < 365) {
                const months = Math.floor(diffDays / 30);
                return `${months} month${months !== 1 ? 's' : ''} ago`;
            } else {
                const years = Math.floor(diffDays / 365);
                return `${years} year${years !== 1 ? 's' : ''} ago`;
            }
        } catch (e) {
            return dateString;
        }
    }

    // UI State Management Functions
    function showBooksLoading() {
        elements.booksLoading.style.display = 'block';
        elements.booksError.style.display = 'none';
        elements.booksEmpty.style.display = 'none';
        elements.booksList.style.display = 'none';
    }

    function hideBooksLoading() {
        elements.booksLoading.style.display = 'none';
    }

    function showBooksError(message) {
        elements.booksLoading.style.display = 'none';
        elements.booksEmpty.style.display = 'none';
        elements.booksList.style.display = 'none';
        elements.booksError.style.display = 'block';
        elements.booksErrorMessage.textContent = message;
    }

    function hideBooksError() {
        elements.booksError.style.display = 'none';
    }

    function showBooksEmpty() {
        elements.booksLoading.style.display = 'none';
        elements.booksError.style.display = 'none';
        elements.booksList.style.display = 'none';
        elements.booksEmpty.style.display = 'block';
    }

    function hideBooksEmpty() {
        elements.booksEmpty.style.display = 'none';
    }

    function showHighlightsLoading() {
        elements.highlightsWelcome.style.display = 'none';
        elements.highlightsLoading.style.display = 'block';
        elements.highlightsError.style.display = 'none';
        elements.highlightsEmpty.style.display = 'none';
        elements.highlightsList.style.display = 'none';
    }

    function hideHighlightsLoading() {
        elements.highlightsLoading.style.display = 'none';
    }

    function showHighlightsError(message) {
        elements.highlightsWelcome.style.display = 'none';
        elements.highlightsLoading.style.display = 'none';
        elements.highlightsEmpty.style.display = 'none';
        elements.highlightsList.style.display = 'none';
        elements.highlightsError.style.display = 'block';
        elements.highlightsErrorMessage.textContent = message;
    }

    function hideHighlightsError() {
        elements.highlightsError.style.display = 'none';
    }

    function showHighlightsEmpty() {
        elements.highlightsWelcome.style.display = 'none';
        elements.highlightsLoading.style.display = 'none';
        elements.highlightsError.style.display = 'none';
        elements.highlightsList.style.display = 'none';
        elements.highlightsEmpty.style.display = 'block';
    }

    function hideHighlightsEmpty() {
        elements.highlightsEmpty.style.display = 'none';
    }

    function hideHighlightsWelcome() {
        elements.highlightsWelcome.style.display = 'none';
    }

    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
