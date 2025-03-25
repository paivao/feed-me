function clearFeed() {
    return {
        name: '',
        description: '',
        feed_type: "ip",
        is_public: true
    }
}

function clearEntry() {
    return {
        name: '',
        description: '',
        valid_until: null
    }
}

function app() {
    return {
        // Authentication
        isAuthenticated: false,
        user: {},
        loginData: {
            username: '',
            password: ''
        },
        feed_types: ["ip", "url", "domain"],

        // CRUD Data
        feeds: [],
        newFeed: clearFeed(),

        entries: [],
        selectedFeed: null,
        newEntry: clearEntry(),

        // API URL
        apiUrl: '/api',

        login() {
            axios.post(`${this.apiUrl}/login`, this.loginData)
                .then(response => {
                    this.isAuthenticated = true;
                    this.user = response.data;
                    this.fetchFeeds();
                })
                .catch(error => {
                    alert('Login failed!');
                    console.error(error);
                });
        },

        logout() {
            axios.post(`${this.apiUrl}/logout`)
                .then(() => {
                    this.isAuthenticated = false;
                    this.user = {};
                    this.feeds = [];
                })
                .catch(error => {
                    alert('Logout failed!');
                    console.error(error);
                });
        },

        fetchFeeds() {
            axios.get(`${this.apiUrl}/feed`)
                .then(response => {
                    this.feeds = response.data;
                })
                .catch(error => {
                    console.error(error);
                });
        },

        // Category CRUD operations
        createFeed() {
            if (this.newFeed.name.trim() === '') return;

            axios.put(`${this.apiUrl}/feed`, this.newFeed)
                .then(response => {
                    this.feeds.push(response.data);
                    this.newFeed = clearFeed()
                })
                .catch(error => {
                    console.error(error);
                });
        },

        removeFeed(id) {
            if (confirm('Are you sure you want to delete this category?')) {
                axios.delete(`${this.apiUrl}/feed/${id}`)
                    .then(() => {
                        this.feeds = this.feeds.filter(feed => feed.id !== id);
                        if (this.selectedFeed.id === id) {
                            entries = [];
                            this.newEntry = clearEntry();
                            this.selectedFeed = null;
                        }
                    })
                    .catch(error => {
                        console.error(error);
                    });
            }
        },

        showEntries(feed) {
            this.selectedFeed = feed
            axios.get(`${this.apiUrl}/entry/${feed.type}/${feed.id}`)
                .then(response => {
                    this.entries = response.data
                })
                .catch(error => {
                    console.error(error);
                });
        },

        // Item CRUD operations
        createEntry(feed) {
            if (category.newItemName.trim() === '') return;

            axios.put(`${this.apiUrl}/items`, { name: category.newItemName, category_id: category.id })
                .then(response => {
                    category.items.push(response.data);
                    category.newItemName = ''; // Clear the input
                })
                .catch(error => {
                    console.error(error);
                });
        },

        editItem(item) {
            const newName = prompt('Edit Item Name:', item.name);
            if (newName && newName !== item.name) {
                axios.post(`${this.apiUrl}/items/${item.id}`, { name: newName })
                    .then(response => {
                        item.name = newName;
                    })
                    .catch(error => {
                        console.error(error);
                    });
            }
        },

        removeEntry(entry) {
            if (confirm('Are you sure you want to delete this item?')) {
                axios.delete(`${this.apiUrl}/items/${itemId}`)
                    .then(() => {
                        this.categories.forEach(category => {
                            category.items = category.items.filter(item => item.id !== itemId);
                        });
                    })
                    .catch(error => {
                        console.error(error);
                    });
            }
        }
    };
}