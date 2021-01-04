(function() {
    const $ = (selector) => document.querySelector(selector);
    const $$ = (selector) => Array.from(document.querySelectorAll(selector));

    let profile;

    document.addEventListener("DOMContentLoaded", init);
    const domParser = new DOMParser();

    async function init() {
        $('#profile-load').onclick = handleProfileLoad;
        $('#terminate-button').onclick = handleTerminate;
        $('#stop-button').onclick = handleStop;
        $('#start-button').onclick = handleStart;
        await getStatus();
        setInterval(getStatus, 2000)
    }

    async function getStatus() {
        const resp = await fetch('/status');
        if (!resp.ok) {
            console.log(resp);
            alert(resp.statusText);
            return;
        }

        const status = await resp.json();
        if (!status.started) {
            return;
        }

        // ((status.successful.length + status.failed.length + status.skipped.length) / status.enqueued.length)
        status.successful = status.successful || [];
        status.failed = status.failed || [];
        status.skipped = status.skipped || [];
        status.enqueued = status.enqueued || [];

        enterProgressMode();

        renderProgressSection(status);
    }

    async function handleStart(event) {
        event.preventDefault();
        const body = {
            degreeOfParallelism: parseInt($('#degree-of-parallelism').value, 10),
            script: $('#profile').value,
            directory: $('#directory').value,
            arguments: $$('#menu-options select').map(el => el.value),
        }
        const resp = await fetch('/start', {method: 'POST', body: JSON.stringify(body)})
        if (!resp.ok) {
            console.log(resp);
            alert(resp.statusText);
            return;
        }
        await getStatus();
    }

    async function handleStop(event) {
        event.preventDefault();
        const resp = await fetch('/stop', {method: 'POST'});
        if (!resp.ok) {
            console.log(resp);
            alert(resp.statusText);
            return;
        }
        await getStatus();
    }

    async function handleTerminate(event) {
        event.preventDefault();
        const resp = await fetch('/terminate', {method: 'POST'});
        if (!resp.ok) {
            console.log(resp);
            alert(resp.statusText);
            return;
        }
        alert('Application terminated.');
    }

    function renderProgressSection(status) {
        const progress = ((status.successful.length + status.failed.length + status.skipped.length) / status.enqueued.length) * 100;
        $('#stop-requested-msg').style.display = status.stopRequested ? 'block' : 'none';
        $('#progress-percentage').innerText = progress.toFixed(2);
        const successfulRecords = $('#successful-records');
        successfulRecords.innerHTML = '';
        for (let successful of status.successful) {
            successfulRecords.innerHTML += `<li> ${successful}`
        }
        const skippedRecords = $('#skipped-records');
        skippedRecords.innerHTML = '';
        for (let skipped of status.skipped) {
            skippedRecords.innerHTML += `<li> ${skipped}`
        }
        const failedRecords = $('#failed-records');
        failedRecords.innerHTML = '';
        for (let failed of status.failed) {
            failedRecords.innerHTML += `<li> ${failed}`
        }

        //$('#image-zone').innerHTML = `<img alt="preview" src="file:///${status.directory}\\preview.jpg">`
    }

    function enterProgressMode() {
        $('#not-running').style.display = 'none';
        $('#in-progress').style.display = 'block';
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
        $('#game-description').innerHTML = profile.description.replaceAll('\\n', '<br>');
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