(function() {
    const $ = (selector) => document.querySelector(selector);
    const $$ = (selector) => Array.from(document.querySelectorAll(selector));

    let profile;

    document.addEventListener("DOMContentLoaded", init);

    function init() {
        $('#profile-load').onclick = handleProfileLoad;
    }

    async function handleProfileLoad(event) {
        event.preventDefault();
        const path = $('#profile').value;
        const params = new URLSearchParams({filePath: path})
        const url = "/profile?" + params.toString();
        const resp = await fetch(url);
        if (!resp.ok) {
            console.log(resp);
            alert(resp.statusText);
            return;
        }

        profile = await resp.json();
        console.log(profile);
    }
})();