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
    domain: "Domain"
}

const FEED_ACT_NAME = {
    create: "Adicionar",
    edit: "Editar"
}

const ENTRY_ACT_NAME = {
    create: "Nova Entrada",
    edit: "Editar Entrada"
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
            await (this.feedFormCreate ? this.createFeed() : this.editFeed());
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
            await (this.entryFormCreate ? this.createEntry() : this.editEntry());
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
        async createFeed() {
            if (this.newFeed.name.trim() === '') return;
            try {
                const response = await axios.put(`${apiURL}/feed`, this.newFeed)
                this.feeds.push(response.data.feed);
            } catch (error) {
                console.error(error);
            }
        },

        async editFeed() {
            try {
                const response = await axios.put(`${apiURL}/feed`, this.newFeed)
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
            } catch (error) {
                console.error(error);
            }
        },

        async editEntry() {
            const feed = this.selectedFeed;
            try {
                const response = post(`${apiURL}/entry/${feed.type}/${feed.id}/${entry.id}`, this.newEntry)
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
                to_remove = this.entries.findIndex(e => e.id === response.data.id)
                if (to_remove != -1)
                    this.entries.splice(to_remove, 1)
            } catch (error) {
                console.error(error);
            }
        },
    };
}

