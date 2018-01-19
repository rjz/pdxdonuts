(function(){
  'use strict';
  var formEl = document.getElementById('forager');
  formEl.addEventListener('submit', function (e) {
    e.preventDefault();
    var opts = {
      location: e.location,
      keyword: e.keyword,
    };

    fetch(formEl.action, {
      method: 'POST', // or 'PUT'
      body: JSON.stringify(opts),
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
