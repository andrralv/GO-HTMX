function wait() {
  setTimeout(function() {
    console.log('Waiting for request...')
  }(), 2000);
}
document.addEventListener('htmx:afterRequest', function(evt) {
  wait();
  if (evt.srcElement.getAttribute('class') === "manga-get-button") {
    var _ul = document.getElementById('mangaList');
    _ul.innerHTML = '';
    JSON.parse(evt.detail.xhr.response).forEach(function(manga) {
      var _li = document.createElement('li');
      _li.textContent = manga.title + ' - ' + manga.author + ' - Stock: ' + manga.quantity;
      _ul.appendChild(_li);
    });
  }
  else if (evt.srcElement.localName === "body") {
    var _div = document.getElementById('mangaUnique');
    wait();
    const res = JSON.parse(evt.detail.xhr.response);
    _div.textContent = res.title + ' - ' + res.author + ' - Stock: ' + res.quantity;
    const buttongroup = document.getElementById('actionButtons');
    buttongroup.innerHTML = '';
    const _checkout = document.createElement('button');
    _checkout.textContent = "checkout";
    _checkout.setAttribute("hx-patch", "http://localhost:8080/checkout?id=" + res.id);
    _checkout.setAttribute("id", "button-checkout");
    _checkout.setAttribute("hx-swap", "none");
    if (parseInt(res.quantity) <= 0) {
      _checkout.disabled = true;
    }
    buttongroup.appendChild(_checkout);
    const _return = document.createElement('button');
    _return.textContent = "return";
    _return.setAttribute("hx-patch", "http://localhost:8080/return?id=" + res.id);
    _return.setAttribute("id", "button-return");
    _return.setAttribute("hx-swap", "none");
    if (parseInt(res.quantity) >= 3) {
      _return.disabled = true;
    }
    buttongroup.appendChild(_return);
    htmx.process(document.body);
  }
});
document.addEventListener('htmx:beforeSwap', function(evt) {
  if ((evt.target.getAttribute('id') === "button-return") || (evt.target.getAttribute('id') === "button-checkout")) {
    wait();
    var _res = JSON.parse(evt.detail.serverResponse);
    var _div = document.getElementById('mangaUnique');
    _div.textContent = _res.title + ' - ' + _res.author + ' - Stock: ' + _res.quantity;
    var _return  = document.getElementById('button-return');
    var _checkout  = document.getElementById('button-checkout');
    if (parseInt(_res.quantity) >= 3) {
      _return.disabled = true;
    } else {
      _return.disabled = false;
    }
    if (parseInt(_res.quantity) <= 0) {
      _checkout.disabled = true;
    } else {
      _checkout.disabled = false;
    }
  }
});
document.addEventListener('htmx:beforeRequest', function(evt) {
  if (evt.srcElement.getAttribute('class') === "manga-get-unique-button") {
    evt.detail.xhr.onloadstart = function(e){this.abort()}
    var _select = document.querySelector('#manga-list-select');
    var _selectValue = parseInt(_select.value);
    wait();
    htmx.ajax('GET', "http://localhost:8080/manga/"+_selectValue, '#mangaUnique')
  }
});