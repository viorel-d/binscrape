const pushState = (state, url, title = '') => {
    history.pushState(state, title, url);
}

const changeHash = (hash) => {
    window.location.hash = hash;
}

const replaceState = (state, url, title = '') => {
    history.replaceState(state, title, url);
}

const newRouter = () => {
    const router = {
        init: function(routes) {
            this.routes = routes;
            this.nRoutes = Object.keys(routes).length;
            this.addEventListener('popstate', this.watchPopState);
            this.addEventListener('onhashchange', this.watchLocationHash);
        },
        routes: null,
        nRoutes: 0,
        eventListeners: null,
        addEventListener: function(eventName, callback, target = window) {
            if (!this.eventListeners) {
                this.eventListeners = {};
            }
            target.addEventListener(eventName, callback);
            this.eventListeners[eventName] = callback;
        },
        removeEventListener: function(eventName, target = window) {
            if (!this.eventListeners) {
                return;
            }
            const callback = this.eventListeners[eventName];
            if (typeof callback === 'function') {
                target.removeEventListener(eventName, callback);
                delete this.eventListeners[eventName];
            }
        },
        isLocationMatching: function(pathname) {
            const callback = this.eventListeners[pathname];
            if (typeof callback === 'function') {
                return [true, callback];
            }
            return [false, null];
        },
        watchPopState: function(e) {
            const { pathname } = e.state;
            const [isMatch, callback] = this.isLocationMatching(pathname);
            if (!isMatch) {
                return;
            }
            callback();
        },
        pushState: pushState,
        replaceState: replaceState,
        changeHash: changeHash,
        shutdown: function() {
            const boundEventListeners = Object.keys(this.eventListeners);
            const nEventListeners = boundEventListeners.length;
            if (nEventListeners === 0) {
                return;
            }
            for (let i = 0; i < nEventListeners; i++) {
                const eventName = boundEventListeners[i];
                this.removeEventListener(eventName);
            }
        }
    };

    return router;
};

const router = newRouter();

export default router;
