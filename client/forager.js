(function(){
  'use strict';
  var formEl = document.getElementById('forager');
  var resultsEl = document.querySelector('.form__results');
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
    resultsEl.innerHTML = '';

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
      return res.json().then(json => {
        resultsEl.classList.add('form__results--error');
        resultsEl.innerHTML = json.error;
      });
    });
  });

  var mapEl = document.getElementById('google-map');

  var placesEl = document.createElement('div');
  placesEl.style.border = '1px solid red';
  document.body.appendChild(placesEl);

  window.initMap = function () {
    var pdx = {
      lat: 45.5231,
      lng: -122.6765,
    };
    var map = new google.maps.Map(mapEl, {
      zoom: 10,
      center: pdx,
    });

    var places = new google.maps.places.PlacesService(placesEl);

    var searchOpts = {
      bounds: map.getBounds(),
      query: 'donuts in Portland',
      type: 'restaurant|bakery',
    };

    // google.maps.event.addListener(map, 'bounds_changed', function () {
    //   console.log(searchOpts);
    // });

    places.textSearch(searchOpts, function (results, status, pagination) {
      alert(results.map(r => r.name).join('\n'));
      console.log({ results, status, pagination });
    });
  };
})();
