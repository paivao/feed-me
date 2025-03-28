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
    return {
        // Authentication
        user: null,
        loginData: clearLoginData(),
        feed_types: ["ip", "url", "domain"],

        // CRUD Data
        feeds: [],
        newFeed: clearFeed(),

        entries: [],
        selectedFeed: null,
        newEntry: clearEntry(),

        init() {
            axios.get(`${apiURL}/whoami`).then(response => {
                //console.log(response)
                this.user = response.data;
            }).catch(() => {})
        },

        async login() {
            try {
                const response = await axios.post(`${apiURL}/login`, this.loginData)
                this.isAuthenticated = true;
                this.user = response.data;
                this.loginData = clearLoginData();
                await this.fetchFeeds();
            } catch (error) {
                alert('Login failed!');
                console.error(error);
            }
        },

        async logout() {
            try {
                const response = await axios.post(`${apiURL}/logout`);
                this.user = null;
                this.feeds = [];
            } catch (error) {
                alert('Logout failed!');
                console.error(error);
            }
        },

        async fetchFeeds() {
            try {
                const response = await axios.get(`${apiURL}/feed`)
                this.feeds = response.data;
            } catch (error) {
                console.error(error);
            }
        },

        // Category CRUD operations
        async createFeed() {
            if (this.newFeed.name.trim() === '') return;
            try {
                const response = await axios.put(`${apiURL}/feed`, this.newFeed)
                this.feeds.push(response.data.feed);
                this.newFeed = clearFeed()
            } catch (error) {
                console.error(error);
            }
        },

        async removeFeed(id) {
            if (!confirm('Are you sure you want to delete this feed?'))
                return;
            try {
                const response = await axios.delete(`${apiURL}/feed/${id}`)
                to_remove = this.feeds.findIndex(feed => feed.id !== response.data.id);
                if (to_remove != -1)
                    this.feeds.splice(to_remove, 1);
                if (this.selectedFeed.id === response.data.id) {
                    this.entries = [];
                    this.newEntry = clearEntry();
                    this.selectedFeed = null;
                }
            } catch (error) {
                console.error(error);
            }
        },

        async showEntries(feed) {
            this.selectedFeed = feed;
            try {
                const response = await axios.get(`${apiURL}/entry/${feed.type}/${feed.id}`)
                this.entries = response.data
            } catch (error) {
                console.error(error);
            }
        },

        // Item CRUD operations
        async createEntry() {
            if (this.newEntry.value.trim() === '')
                return;
            const feed = this.selectedFeed;
            try {
                const response = axios.put(`${apiURL}/entry/${feed.type}/${feed.id}`, this.newEntry)
                this.entries.push(response.data.entry);
                this.newEntry = clearEntry();
            } catch (error) {
                console.error(error);
            }
        },

        async editEntry(entry) {
            const newName = prompt('Edit Item Name:', entry.name);
            const feed = this.selectedFeed;
            if (newName === null || newName === entry.name) {
                return;
            }
            try {
                const response = post(`${apiURL}/entry/${feed.type}/${feed.id}/${entry.id}`, { name: newName })
                element = this.entries.find(e => e.id == response.data.id)
            } catch (error) {
                console.error(error);
            }
        },

        async removeEntry(entry) {
            if (!confirm('Are you sure you want to delete this item?'))
                return;
            const feed = this.selectedFeed;
            try {
                const response = await axios.delete(`${apiURL}/entry/${feed.type}/${feed.id}/${entry.id}`)
                to_remove = this.entries.findIndex(e => e.id === response.data.id)
                if (to_remove != -1)
                    this.entries.splice(to_remove, 1)
            } catch (error) {
                console.error(error);
            }
        },
    };
}

