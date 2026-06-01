const csrfManager = (function() {
    let token = null;

    return {
        getToken: function() {
            return token;
        },
        setToken: function(newToken) {
            token = newToken;
        },
        updateFromResponse: function(response) {
            const newToken = response.headers.get('X-CSRF-Token');
            if (newToken) {
                token = newToken;
            }
        },
        reset: function() {
            token = null;
        }
    };
})();
