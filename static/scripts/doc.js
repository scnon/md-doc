function toggleSearchBar() {
    var x = document.getElementById("search_bar");
    console.log(x)
    if (x.style.display === "none") {
        x.style.display = "block";
        document.getElementById("search_input").focus();
    } else {
        x.style.display = "none";
    }
}

function onSearchInput(e) {
    console.log('aaaa')
}

function onSearchOut(e) {
    document.getElementById("search_bar").style.display = "none";
}

document.addEventListener('keydown', (e) => {
    if (e.ctrlKey && e.key === "l") {
        toggleSearchBar();
    } else if(e.key === "Escape") {
        document.getElementById("search_bar").style.display = "none";
    }
})

document.addEventListener('input', (e) => {
    var input = document.getElementById("search_input").value;
    results = document.getElementById("search_result");

    $.ajax({
        type: "POST",
        url: "/api/doc/search",
        data: JSON.stringify({ "key": input }),
        dataType: "json",
        contentType: "text/plain",
        success: function (data, status) {
            console.log(data, status)
            if (data.code === 200 && data.data !== null && data.data !== undefined) {
                results.innerHTML = data.data;
            }
        }
    })
})