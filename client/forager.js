(function(){
  'use strict';
  var formEl = document.getElementById('forager');
  formEl.addEventListener('submit', function (e) {
    e.preventDefault();

    var inputs = formEl.querySelectorAll('input[type=text]');
    var body = [].slice.call(inputs).reduce(function (acc, el) {
      if (el.name.indexOf('pageData.') === 0) {
        acc.pageData[el.name.slice(9)] = el.value;
      } else {
        acc[el.name] = el.value;
      }
      return acc;
    }, {pageData: {}});

    body.types = body.types.split('|');

    fetch(formEl.action, {
      method: 'POST', // or 'PUT'
      body: JSON.stringify(body),
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    }).then((res) => {
      if (res.status === 200) {
        return res.blob().then(download);
      }
      return res.json().then(json => alert(JSON.stringify(json, null, 2)));
    });
  });
})();
