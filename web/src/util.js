import base64 from "base-64"

export async function postData(url = '', data = {}) {
    // Default options are marked with *
    let headers = new Headers();
    let authString = "mar-tina:password"
    headers.set('Authorization', 'Basic ' + btoa(authString))
    headers.set('Content-Type', 'application/json')

    const response = await fetch(url, {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: headers,
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data) // body data type must match "Content-Type" header
    });

    console.log(response)
    return response.json(); // parses JSON response into native JavaScript objects
}