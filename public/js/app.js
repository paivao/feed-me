function clearFeed() {
    return {
        name: '',
        description: '',
        feed_type: "ip",
        is_public: true,
    }
}

function clearEntry() {
    return {
        value: '',
        description: '',
        valid_until: null,
    }
}

function clearLoginData() {
    return {
        username: '',
        password: '',
    }
}

const url_loc = window.location.pathname.split("/");
console.log(url_loc);
const last_path = url_loc[url_loc.length - 1];
if (last_path === "" || last_path.endsWith(".html"))
    url_loc.pop()
url_loc.pop()
url_loc.push("api")
const apiURL = url_loc.join("/");
console.log(apiURL);

const FEED_DATA_TYPE = {
    ip: "IP",
    url: "URL",
    domain: "Domain",
}

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

        feedFormCreate: true,
        entryFormCreate: true,

        entryPagination: {
            window: 20,
            page: 0,
            count: 0,
            lastPage() {
                return Math.floor((this.count -1) / (this.window))
            },
        },

        init() {
            axios.get(`${apiURL}/whoami`).then(response => {
                //console.log(response)
                this.user = response.data;
                this.fetchFeeds().then();
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

        feedActionName() {
            return this.feedFormCreate ? 'Adicionar' : 'Editar';
        },

        entryActionName() {
            return this.feedFormCreate ? 'Adicionar' : 'Editar';
        },

        async fetchFeeds() {
            try {
                const response = await axios.get(`${apiURL}/feed`)
                this.feeds = response.data;
            } catch (error) {
                console.error(error);
            }
        },

        async submitFeed() {
            const dafeed = Object.fromEntries(Object.entries(this.newFeed).filter(([_, v]) => v != null || v != ''));
            await (this.feedFormCreate ? this.createFeed(dafeed) : this.editFeed(dafeed));
            this.resetFeed()
        },

        resetFeed() {
            this.newFeed = clearFeed(),
            this.feedFormCreate = true;
        },

        feedMarkToEdit(feed) {
            this.newFeed = feed;
            this.feedFormCreate = false;
        },

        async submitEntry() {
            const daentry = Object.fromEntries(Object.entries(this.newEntry).filter(([_, v]) => v != null || v != ''));
            await (this.entryFormCreate ? this.createEntry(daentry) : this.editEntry(daentry));
            this.resetEntry()
        },

        resetEntry() {
            this.newEntry = clearEntry(),
            this.entryFormCreate = true;
        },

        entryMarkToEdit(entry) {
            this.newEntry = entry;
            this.entryFormCreate = false;
        },

        unselectFeed() {
            this.selectedFeed = null;
            this.entries = [];
        },

        // Category CRUD operations
        async createFeed(feed) {
            if (this.newFeed.name.trim() === '') return;
            try {
                const response = await axios.put(`${apiURL}/feed`, feed)
                this.feeds.push(response.data.feed);
            } catch (error) {
                console.error(error);
            }
        },

        async editFeed(feed) {
            try {
                const response = await axios.put(`${apiURL}/feed`, feed)
                const foundFeed = this.feeds.find(f => f.id === response.data.id)
                foundFeed.description = this.newFeed.description;
                foundFeed.is_public = this.newFeed.is_public;
            } catch (error) {
                console.error(error);
            }
        },

        async removeFeed(feed) {
            if (!confirm('Are you sure you want to delete this feed?'))
                return;
            try {
                const response = await axios.delete(`${apiURL}/feed/${feed.id}`)
                to_remove = this.feeds.findIndex(f => f.id !== response.data.id);
                if (to_remove != -1)
                    this.feeds.splice(to_remove, 1);
            } catch (error) {
                console.error(error);
            }
        },

        async showEntries(feed) {
            this.selectedFeed = feed;
            this.entryPagination.page = 0;
            await this.countEntries(feed);
            await this.fetchEntries(feed);
        },

        async countEntries(feed) {
            try {
                const response = await axios.get(`${apiURL}/entry/${feed.type}/${feed.id}/count`)
                this.entryPagination.count = response.data
            } catch (error) {
                console.error(error);
            }
        },

        async fetchEntries(feed) {
            try {
                const response = await axios.get(`${apiURL}/entry/${feed.type}/${feed.id}?window=${this.entryPagination.window}&page=${this.entryPagination.page}`)
                this.entries = response.data
            } catch (error) {
                console.error(error);
            }
        },

        // Item CRUD operations
        async createEntry(entry) {
            if (this.newEntry.value.trim() === '')
                return;
            const feed = this.selectedFeed;
            try {
                const response = axios.put(`${apiURL}/entry/${feed.type}/${feed.id}`, entry)
                this.entryPagination.page = Math.floor((this.entryPagination.count) / (this.entryPagination.window))
                this.entryPagination.count += 1
                await this.fetchEntries();

            } catch (error) {
                console.error(error);
            }
        },

        async editEntry(entry) {
            const feed = this.selectedFeed;
            try {
                const response = post(`${apiURL}/entry/${feed.type}/${feed.id}/${entry.id}`, entry)
                const foundEntry = this.entries.find(e => e.id == response.data.id)
                foundEntry.enabled = this.newEntry.enabled;
                foundEntry.description = this.newEntry.description;
                foundEntry.valid_until = this.newEntry.valid_until;
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
                this.entryPagination.count -= 1;
                await this.fetchEntries();
            } catch (error) {
                console.error(error);
            }
        },
    };
}

