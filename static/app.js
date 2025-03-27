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
        value: '',
        description: '',
        valid_until: null
    }
}

function clearLoginData() {
    return {
        username: '',
        password: ''
    }
}

const apiURL = '/api';

function app() {
    let user = null;
    axios.get(`${apiURL}/whoami`).then(
        response => {
            user = response
        }
    ).catch(
        _ => {
            user = null
        }
    )
    return {
        // Authentication
        user: user,
        loginData: clearLoginData(),
        feed_types: ["ip", "url", "domain"],

        // CRUD Data
        feeds: [],
        newFeed: clearFeed(),

        entries: [],
        selectedFeed: null,
        newEntry: clearEntry(),

        login() {
            axios.post(`${apiURL}/login`, this.loginData)
                .then(response => {
                    this.isAuthenticated = true;
                    this.user = response.data;
                    this.loginData = clearLoginData();
                    this.fetchFeeds();
                })
                .catch(error => {
                    alert('Login failed!');
                    console.error(error);
                });
        },

        logout() {
            axios.post(`${apiURL}/logout`)
                .then(() => {
                    this.user = null;
                    this.feeds = [];
                })
                .catch(error => {
                    alert('Logout failed!');
                    console.error(error);
                });
        },

        fetchFeeds() {
            axios.get(`${apiURL}/feed`)
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

            axios.put(`${apiURL}/feed`, this.newFeed)
                .then(response => {
                    this.feeds.push(response.data.feed);
                    this.newFeed = clearFeed()
                })
                .catch(error => {
                    console.error(error);
                });
        },

        removeFeed(id) {
            if (confirm('Are you sure you want to delete this feed?')) {
                axios.delete(`${apiURL}/feed/${id}`)
                    .then(response => {
                        to_remove = this.feeds.findIndex(feed => feed.id !== response.data.id);
                        if (to_remove != -1)
                            this.feeds.splice(to_remove, 1);
                        if (this.selectedFeed.id === response.data.id) {
                            this.entries = [];
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
            axios.get(`${apiURL}/entry/${feed.type}/${feed.id}`)
                .then(response => {
                    this.entries = response.data
                })
                .catch(error => {
                    console.error(error);
                });
        },

        // Item CRUD operations
        createEntry() {
            if (this.newEntry.value.trim() === '') return;
            const feed = this.selectedFeed;
            axios.put(`${apiURL}/entry/${feed.type}/${feed.id}`, this.newEntry)
                .then(response => {
                    this.entries.push(response.data.entry);
                    this.newEntry = clearEntry();
                })
                .catch(error => {
                    console.error(error);
                });
        },

        editEntry(entry) {
            const newName = prompt('Edit Item Name:', entry.name);
            const feed = this.selectedFeed;
            if (newName && newName !== entry.name) {
                axios.post(`${apiURL}/entry/${feed.type}/${feed.id}/${entry.id}`, { name: newName })
                    .then(response => {
                        element = this.entries.find(e => e.id == response.data.id)
                    })
                    .catch(error => {
                        console.error(error);
                    });
            }
        },

        removeEntry(entry) {
            if (confirm('Are you sure you want to delete this item?')) {
                const feed = this.selectedFeed;
                axios.delete(`${apiURL}/entry/${feed.type}/${feed.id}/${entry.id}`)
                    .then(response => {
                        to_remove = this.entries.findIndex(e => e.id === response.data.id)
                        if (to_remove != -1)
                            this.entries.splice(to_remove, 1)
                    })
                    .catch(error => {
                        console.error(error);
                    });
            }
        },
    };
}

