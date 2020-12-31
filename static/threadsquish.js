(function() {
    const $ = (selector) => document.querySelector(selector);
    const $$ = (selector) => Array.from(document.querySelectorAll(selector));

    let profile;

    document.addEventListener("DOMContentLoaded", init);
    const domParser = new DOMParser();

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
        renderForm();
    }

    function renderForm() {
        $('#game-name').innerText = profile.name;
        $('#game-description').innerText = profile.description;
        $('#format').innerText = profile.formatName;
        $('#example').innerText = profile.example;
        renderOptions(profile.options);

        $('#form-area').style.display = 'block';
    }

    function renderOptions(opts) {
        const menuOptions = $('#menu-options');
        menuOptions.innerHTML = '';
        let i = 2;

        for (let opt of opts) {
            const name = `menu${i}`;
            let innerText = `<div class="menu-option-container" id="${name}-container">\n`;
            innerText += `<label for="${name}">${opt.default}</label>\n`;
            if (opt.description !== "") {
                innerText += `<p class="explicatory">(${opt.description})</p>`
            }
            innerText += `<select name="${name}" id="${name}">\n`
            for (let [key, value] of Object.entries(opt.cases)) {
                innerText += `<option value="${key}">${value}</option>\n`
            }
            innerText += '</select>\n'
            innerText += '</div>'
            const doc = domParser.parseFromString(innerText, 'text/html');
            menuOptions.appendChild(doc.getElementsByClassName('menu-option-container')[0]);
            i++;
        }
    }
})();